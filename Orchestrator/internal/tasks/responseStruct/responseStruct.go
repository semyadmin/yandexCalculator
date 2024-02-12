package responseStruct

import (
	"strconv"
	"time"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/tasks/arithmetic"
)

type Expression struct {
	ID         string
	Expression string
	Start      string
	End        string
	Status     string
}

func NewExpression(a *arithmetic.ASTTree) Expression {
	r := Expression{}
	a.Lock()
	r.ID = strconv.FormatUint(a.ID, 10)
	r.Expression = a.Expression
	r.Start = a.Start.Format("02.01.2006 15:04:05")
	r.End = a.Start.Add(time.Duration(a.Duration) * time.Second).Format("02.01.2006 15:04:05")
	r.Status = "progress"
	if a.IsCalc {
		r.Status = "completed"
	}
	a.Unlock()
	return r
}
