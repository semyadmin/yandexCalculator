package responseexpression

import (
	"strconv"
	"time"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/entity"
)

func NewResponseExpression(id uint64,
	expression string,
	start time.Time,
	duration int64,
	isCalc bool,
	value float64,
	err error,
) entity.ResponseExpression {
	r := entity.ResponseExpression{}
	r.ID = strconv.FormatUint(id, 10)
	r.Expression = expression
	r.Start = start.Format("02.01.2006 15:04:05")
	r.End = start.Add(time.Duration(duration) * time.Second).Format("02.01.2006 15:04:05")
	r.Status = "progress"
	if isCalc {
		r.Expression = expression + "=" + strconv.FormatFloat(value, 'f', -1, 64)
		r.Status = "completed"
	}
	if err != nil {
		r.Expression = expression + "=error"
		r.Status = "error"
	}
	return r
}
