package arithmetic

import (
	"encoding/json"
	"errors"
	"go/ast"
	"go/parser"
	"log/slog"
	"sync"
	"time"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/config"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/entity"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/tasks/queue"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/tasks/upgrade"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/web_socket/client"
)

type ASTTree struct {
	expression *entity.Expression
	X          *ASTTree
	Y          *ASTTree
	Operator   string
	Value      string
	IsCalc     bool
	IsParent   bool
	queue      *queue.MapQueue
	config     *config.Config
	Start      time.Time
	Duration   int64
	Err        error
	sync.Mutex
}

type result struct {
	err error
	res string
}

// Создаем AST дерево из выражения
func NewASTTree(expression *entity.Expression,
	config *config.Config,
	queue *queue.MapQueue,
) (*ASTTree, error) {
	if expression.Err != nil {
		return nil, expression.Err
	}
	// Добавляем кавычки, где только возможно
	upgradeExp := upgrade.Upgrade([]byte(expression.Expression))
	slog.Info("Усовершенствованное выражение", "выражение:", string(upgradeExp))
	// Создаем AST дерево
	tr, err := parser.ParseExpr(string(upgradeExp))
	if err != nil {
		return nil, err
	}
	// Преобразуем AST дерево в нашу структуру ASTTree
	a := create(tr)
	a.queue = queue
	a.config = config
	// Вычисляем примерное время выполнения выражения
	a.Duration = duration(a, config)
	a.Start = time.Now()
	return a, nil
}

// Создаем AST дерево из базы данных
func NewASTTreeDB(
	id uint64,
	expression *entity.Expression,
	value string,
	isErr bool,
	currentResult string,
	config *config.Config,
	queue *queue.MapQueue,
) (*ASTTree, error) {
	tr, err := parser.ParseExpr(currentResult)
	if err != nil {
		return nil, err
	}
	a := create(tr)
	a.expression = expression
	a.Value = value
	if isErr {
		a.Err = errors.New("error")
	}
	a.IsCalc = true
	a.queue = queue
	a.config = config
	a.Duration = duration(a, config)
	a.Start = time.Now()
	return a, nil
}

func create(tr ast.Expr) *ASTTree {
	a := new(ASTTree)
	switch nod := tr.(type) {
	case *ast.BasicLit:
		a.Value = nod.Value
		a.IsCalc = true
	case *ast.ParenExpr:
		a.X = create(nod.X)
		a.IsParent = true
	case *ast.BinaryExpr:
		a.X = create(nod.X)
		a.Y = create(nod.Y)
		a.Operator = nod.Op.String()
	case *ast.UnaryExpr:
		a.Value = nod.Op.String() + create(nod.X).Value
		a.IsCalc = true
	}
	return a
}

func duration(a *ASTTree, config *config.Config) int64 {
	if a == nil {
		return 0
	}
	res := int64(0)
	if a.Operator == "+" {
		res += config.Plus
	}
	if a.Operator == "-" {
		res += config.Minus
	}
	if a.Operator == "*" {
		res += config.Multiply
	}
	if a.Operator == "*" {
		res += config.Divide
	}
	res += duration(a.X, config)
	res += duration(a.Y, config)
	return res
}

// Вычисляем выражение
func Calculate(a *ASTTree, c *config.Config) {
	if a.IsCalc || a.Err != nil || a == nil {
		return
	}
	ch := make(chan result)
	go getResult(a, ch, a, "P")
	res := <-ch
	a.Lock()
	if res.err != nil {
		a.Err = res.err
	} else {
		a.Value = res.res
		a.IsCalc = true
	}
	a.Unlock()
	resp := entity.NewResponseExpression(a.expression.ID, a.expression.Expression, a.Start, a.Duration, "progress", a.IsCalc, a.expression.Result, a.Err)
	answer, err := json.Marshal(resp)
	if err != nil {
		slog.Error("Проблема с формированием ответа", "ошибка:", err)
		return
	}
	go func() {
		c.WSmanager.MessageCh <- &client.Message{
			Payload: answer,
			Type:    client.ClientExpression,
		}
	}()
	slog.Info("Выражение вычислено", "выражение:", a.expression.Expression, "результат:", a.Value)
}

// Печатаем полученное выражение, вычисленное в процессе, что бы
// не считать все выражение заново
func PrintExpression(a *ASTTree) string {
	if a.IsCalc {
		return a.Value
	}
	if a.IsParent {
		return "(" + PrintExpression(a.X) + ")"
	}

	return PrintExpression(a.X) + a.Operator + PrintExpression(a.Y)
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
func calculate(resX, operator, resY string, parent *ASTTree, level string) result {
	resultCh := make(chan string)
	deadline := int64(0)
	switch operator {
	case "+":
		deadline = parent.config.Plus
	case "-":
		deadline = parent.config.Minus
	case "*":
		deadline = parent.config.Multiply
	case "/":
		deadline = parent.config.Divide
	}
	send := &queue.SendInfo{
		Id:         parent.expression.Expression + "-" + level,
		Expression: resX + operator + resY,
		Result:     resultCh,
		Deadline:   uint64(deadline),
		IdExp:      parent.expression.Expression,
	}
	parent.queue.Enqueue(send)
	res := result{}
	resExp := <-send.Result
	if resExp == "error" {
		res.err = errors.New("error")
		return res
	}
	res.res = resExp
	return res
}
