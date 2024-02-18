package arithmetic

import (
	"errors"
	"go/ast"
	"go/parser"
	"log/slog"
	"sync"
	"time"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/config"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/tasks/queue"
)

type ASTTree struct {
	ID         uint64
	Expression string
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

func NewASTTree(expression string, config *config.Config, queue *queue.MapQueue) (*ASTTree, error) {
	upgradeExp := Upgrade([]byte(expression))
	slog.Info("Усовершенствованное выражение", "выражение:", string(upgradeExp))
	tr, err := parser.ParseExpr(string(upgradeExp))
	if err != nil {
		return nil, err
	}
	a := create(tr)
	if err != nil {
		return nil, err
	}
	a.Expression = expression
	a.queue = queue
	a.config = config
	a.Duration = duration(a, config)
	a.Start = time.Now()
	return a, nil
}

func NewASTTreeDB(
	id uint64,
	expression string,
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
	if err != nil {
		return nil, err
	}
	a.ID = id
	a.Expression = expression
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

func (a *ASTTree) Calculate() {
	if a.IsCalc {
		return
	}
	ch := make(chan result)
	go getResult(a, ch, a, "P")
	res := <-ch
	if res.err != nil {
		a.Lock()
		a.Err = res.err
		a.Unlock()
		return
	}
	a.Lock()
	a.Value = res.res
	a.IsCalc = true
	a.Unlock()
	slog.Info("Выражение вычислено", "выражение:", a.Expression, "результат:", a.Value)
}

func (a *ASTTree) GetExpression() string {
	return a.Expression
}

func (a *ASTTree) SetID(id uint64) {
	a.ID = id
}

func PrintExpression(a *ASTTree) string {
	if a.IsCalc {
		return a.Value
	}
	if a.IsParent {
		return "(" + PrintExpression(a.X) + ")"
	}

	return PrintExpression(a.X) + a.Operator + PrintExpression(a.Y)
}

func getResult(a *ASTTree, ch chan result, parent *ASTTree, level string) {
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
		Id:         parent.Expression + "-" + level,
		Expression: resX + operator + resY,
		Result:     resultCh,
		Deadline:   uint64(deadline),
		IdExp:      parent.Expression,
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

func Upgrade(exp []byte) []byte {
	operatorIndex := 0
	result := make([]byte, 0)
	left := 0
	for i := 0; i < len(exp); i++ {
		if exp[i] == '+' || exp[i] == '-' {
			if exp[operatorIndex] == '+' || exp[operatorIndex] == '-' {
				array := append([]byte{'('}, exp[left:i]...)
				array = append(array, ')')
				array = append(array, exp[i])
				result = append(result, array...)
				left = i + 1
				operatorIndex = 0
				continue
			}
			operatorIndex = i
			continue
		}
		if exp[i] == '*' || exp[i] == '/' {
			if operatorIndex > left {
				result = append(result, exp[left:operatorIndex+1]...)
				left = operatorIndex + 1
			}
			array, index := upgradeMultiDivide(exp, i+1, left)
			result = append(result, array...)
			if index == len(exp) {
				return result
			}
			operatorIndex = 0
			if exp[index] == ')' {
				return result
			}
			if exp[index] == '(' {
				left = index
				i = index - 1
				continue
			}
			i = index
			left = i + 1
			if index == len(exp) {
				break
			}
			result = append(result, exp[index])
			continue
		}
		if exp[i] == '(' {
			result = append(result, exp[left:i+1]...)
			array := Upgrade(exp[i+1:])

			result = append(result, array...)
			closeBoarder := -1
			for j := i + 1; j < len(exp); j++ {
				if exp[j] == ')' {
					closeBoarder++
				}
				if exp[j] == '(' {
					closeBoarder--
				}
				if closeBoarder == 0 {
					i = j
					break
				}
			}
			result = append(result, exp[i])
			if i == len(exp)-1 {
				return result
			}
			if exp[i+1] == '+' || exp[i+1] == '-' ||
				exp[i+1] == '*' || exp[i+1] == '/' {
				i = i + 1
				result = append(result, exp[i])
			}
			left = i + 1
			operatorIndex = 0
			continue
		}
		if exp[i] == ')' {
			result = append(result, exp[left:i]...)
			return result
		}
	}
	if operatorIndex != 0 {
		array := append([]byte{'('}, exp[left:]...)
		array = append(array, ')')
		result = append(result, array...)
	} else {
		result = append(result, exp[left:]...)
	}

	return result
}

func upgradeMultiDivide(exp []byte, index int, left int) ([]byte, int) {
	result := make([]byte, 0)
	openBorder := true
	for i := index; i < len(exp); i++ {
		if exp[i] == '*' || exp[i] == '/' {
			if openBorder {
				array := append([]byte{'('}, exp[left:i]...)
				array = append(array, ')')
				array = append(array, exp[i])
				result = append(result, array...)
				left = i + 1
				openBorder = false
				continue
			}
			openBorder = true
		}
		if exp[i] == '+' || exp[i] == '-' ||
			exp[i] == ')' {
			if openBorder {
				array := append([]byte{'('}, exp[left:i]...)
				array = append(array, ')')
				result = append(result, array...)
			} else {
				result = append(result, exp[left:i]...)
			}
			return result, i
		}
		if exp[i] == '(' {
			result = append(result, exp[left:i]...)
			return result, i
		}
	}
	if openBorder {
		array := append([]byte{'('}, exp[left:]...)
		array = append(array, ')')
		result = append(result, array...)
	} else {
		result = append(result, exp[left:]...)
	}
	return result, len(exp)
}
