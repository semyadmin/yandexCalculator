package expression

import (
	"reflect"
	"testing"
	"time"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/entity"
	jwttoken "github.com/adminsemy/yandexCalculator/Orchestrator/internal/services/jwt_token"
)

var start = time.Now()

type StorageGetAllMock struct{}

func (s StorageGetAllMock) GetAll(token string) []*entity.Expression {
	return []*entity.Expression{
		{
			ID:         1,
			Expression: "2+2",
			Start:      start,
		},
	}
}

func TestGetAllExpressions(t *testing.T) {
	type args struct {
		storage StorageGetAll
		token   string
	}
	token, _ := jwttoken.GenerateToken("test", 15)
	tests := []struct {
		name string
		args args
		want []entity.ResponseExpression
	}{
		{
			name: "positive test",
			args: args{
				storage: &StorageGetAllMock{},
				token:   token,
			},
			want: []entity.ResponseExpression{
				{
					ID:         "1",
					Expression: "2+2",
					Start:      start.Format("02.01.2006 15:04:05"),
					End:        start.Format("02.01.2006 15:04:05"),
					Status:     "progress",
				},
			},
		},
		{
			name: "invalid token",
			args: args{
				storage: &StorageGetAllMock{},
				token:   "token",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAllExpressions(tt.args.storage, tt.args.token); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllExpressions() = %v, want %v", got, tt.want)
			}
		})
	}
}
