package arithmetic

import (
	"strconv"
	"time"
)

type Expression struct {
	ID         string
	Expression string
	Start      string
	End        string
	Status     string
}

// Конвертируем выражение в структуру для отправки
func NewExpression(a *ASTTree) Expression {
	r := Expression{}
	a.Lock()
	r.ID = strconv.FormatUint(a.ID, 10)
	r.Expression = a.Expression
	r.Start = a.Start.Format("02.01.2006 15:04:05")
	r.End = a.Start.Add(time.Duration(a.Duration) * time.Second).Format("02.01.2006 15:04:05")
	r.Status = "progress"
	if a.IsCalc {
		r.Expression = a.Expression + "=" + a.Value
		r.Status = "completed"
	}
	if a.Err != nil {
		r.Expression = a.Expression + "=error"
		r.Status = "error"
	}
	a.Unlock()
	return r
}
