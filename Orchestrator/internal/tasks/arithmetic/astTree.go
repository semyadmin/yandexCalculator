package arithmetic

import (
	"encoding/json"
	"go/ast"
	"go/parser"
	"log/slog"
	"strconv"
	"sync"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/config"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/entity"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/storage/memory"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/storage/postgresql/postgresql_expression"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/tasks/queue"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/web_socket/client"
)

type ASTTree struct {
	expression  *entity.Expression
	X           *ASTTree
	Y           *ASTTree
	Operator    string
	Value       float64
	IsCalc      bool
	IsParent    bool
	queue       *queue.MapQueue
	config      *config.Config
	userStorage *memory.UserStorage
	Err         error
	sync.Mutex
}

type result struct {
	err error
	res float64
}

// Создаем AST дерево из выражения
func NewASTTree(expression *entity.Expression,
	config *config.Config,
	queue *queue.MapQueue,
	userStorage *memory.UserStorage,
) (*ASTTree, error) {
	if expression.Err != nil {
		return nil, nil
	}
	// Создаем AST дерево
	tr, err := parser.ParseExpr(string(expression.CalculatedExpression))
	if err != nil {
		return nil, err
	}
	// Преобразуем AST дерево в нашу структуру ASTTree
	a := create(tr)
	a.expression = expression
	a.queue = queue
	a.config = config
	a.IsCalc = expression.IsCalc
	a.userStorage = userStorage
	go a.calc()
	return a, nil
}

// Создаем AST дерево из базы данных
func NewASTTreeDB(
	expression *entity.Expression,
	config *config.Config,
	queue *queue.MapQueue,
) (*ASTTree, error) {
	tr, err := parser.ParseExpr(expression.CalculatedExpression)
	if err != nil {
		return nil, err
	}
	a := create(tr)
	a.expression = expression
	a.Value = expression.Result
	a.IsCalc = expression.IsCalc
	a.queue = queue
	a.config = config
	go a.calc()
	return a, nil
}

func create(tr ast.Expr) *ASTTree {
	a := new(ASTTree)
	switch nod := tr.(type) {
	case *ast.BasicLit:
		a.Value, _ = strconv.ParseFloat(nod.Value, 64)
		a.IsCalc = true
	case *ast.ParenExpr:
		a.X = create(nod.X)
		a.IsParent = true
	case *ast.BinaryExpr:
		a.X = create(nod.X)
		a.Y = create(nod.Y)
		a.Operator = nod.Op.String()
	case *ast.UnaryExpr:
		v := create(nod.X).Value
		if nod.Op.String() == "-" {
			v = -v
		}
		a.Value = v
		a.IsCalc = true
	}
	return a
}

// Вычисляем выражение
func (a *ASTTree) calc() {
	if a.IsCalc || a.Err != nil || a == nil {
		return
	}
	var err error
	ch := make(chan result)
	go getResult(a, ch, a, "P")
	res := <-ch
	a.Lock()
	if res.err != nil {
		a.Err = res.err
		a.expression.Err = res.err
	} else {
		a.Value = res.res
		a.IsCalc = true
		a.expression.Result = res.res
		if res.err != nil {
			a.expression.Err = err
		}
		a.expression.IsCalc = true
	}
	a.Unlock()
	resp := entity.NewResponseExpression(a.expression.ID, a.expression.Expression, a.expression.Start, a.expression.Duration, a.IsCalc, a.expression.Result, a.Err)
	answer, err := json.Marshal(resp)
	if err != nil {
		slog.Error("Проблема с формированием ответа", "ошибка:", err)
		return
	}
	// Отправляем ответ клиенту веб сокета
	go func() {
		a.config.WSmanager.MessageCh <- &client.Message{
			Payload: answer,
			Type:    client.ClientExpression,
		}
	}()
	slog.Info("Выражение вычислено", "выражение:", a.expression.Expression, "результат:", a.expression.Result)
	// Сохраняем выражение в базу данных
	newExp := postgresql_expression.Expression{
		BaseID:        a.expression.ID,
		Expression:    a.expression.Expression,
		Value:         a.expression.Result,
		User:          a.expression.User,
		CurrentResult: a.PrintExpression(),
	}
	if res.err != nil {
		newExp.Err = true
	}
	a.config.Db.Expression.Update(newExp)
}

// Печатаем полученное выражение, вычисленное в процессе, что бы
// не считать все выражение заново
func (a *ASTTree) PrintExpression() string {
	a.Lock()
	defer a.Unlock()
	return stringASTTree(a)
}

func stringASTTree(a *ASTTree) string {
	if a.IsCalc {
		return strconv.FormatFloat(a.Value, 'f', -1, 64)
	}
	if a.IsParent {
		return "(" + stringASTTree(a.X) + ")"
	}

	return stringASTTree(a.X) + a.Operator + stringASTTree(a.Y)
}

// Вычисляем каждую операцию. Если уже вычислено
// то возвращаем результат
func getResult(a *ASTTree, ch chan result, parent *ASTTree, level string) {
	if a == nil {
		return
	}
	res := result{}
	if a.IsCalc {
		res.res = a.Value
		ch <- res
		return
	}
	resChX := make(chan result)
	go getResult(a.X, resChX, parent, level+"X")
	if a.Y == nil {
		res = <-resChX
		if res.err != nil {
			ch <- res
			return
		}
		a.Lock()
		a.Value = res.res
		a.IsCalc = true
		a.Unlock()
		ch <- res
		return
	}
	resChY := make(chan result)
	go getResult(a.Y, resChY, parent, level+"Y")
	resX := <-resChX
	if resX.err != nil {
		ch <- resX
		return
	}
	resY := <-resChY
	if resY.err != nil {
		ch <- resY
		return
	}
	res = calculate(resX.res, a.Operator, resY.res, parent, level)
	ch <- res
	if res.err != nil {
		return
	}
	a.Lock()
	a.Value = res.res
	a.IsCalc = true
	a.Unlock()
}

// Вычисляем операцию в зависимости от оператора
func calculate(resX float64, operator string, resY float64, parent *ASTTree, level string) result {
	deadline := int64(0)
	config, err := parent.userStorage.GetConfig(parent.expression.User)
	if err != nil {
		config = &entity.Config{}
	}
	switch operator {
	case "+":
		deadline = config.Plus
	case "-":
		deadline = config.Minus
	case "*":
		deadline = config.Multiply
	case "/":
		deadline = config.Divide

	}
	// Отправляем выражение в очередь для вычисления агентом
	send := entity.NewOperation(parent.expression.Expression+"-"+parent.expression.User+"-"+level, resX, resY, operator, uint64(deadline))
	parent.queue.Enqueue(send)
	res := result{}
	resExp := <-send.ResultChan()
	if send.GetError() != nil {
		res.err = send.GetError()
		return res
	}
	res.res = resExp
	// Обновляем выражение в базе данных
	newExp := postgresql_expression.Expression{
		BaseID:        parent.expression.ID,
		Expression:    parent.expression.Expression,
		Value:         parent.expression.Result,
		User:          parent.expression.User,
		CurrentResult: parent.PrintExpression(),
	}
	if res.err != nil {
		newExp.Err = true
	}
	parent.config.Db.Expression.Update(newExp)
	return res
}
