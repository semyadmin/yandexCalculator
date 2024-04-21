package arithmetic

import (
	"errors"
	"testing"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/config"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/entity"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/tasks/queue"
)

type MapQueueMock struct{}

func (m *MapQueueMock) Enqueue(expression queue.Expression) {
	first := expression.First()
	second := expression.Second()
	operator := expression.Operation()
	result := float64(0)
	switch operator {
	case "+":
		result = first + second
	case "-":
		result = first - second
	case "*":
		result = first * second
	case "/":
		if second == 0 {
			expression.Error("division by zero")
			expression.Result(0)
		}
		result = first / second
	}
	expression.Result(result)
}

type UserStorageMock struct{}

func (u *UserStorageMock) SetConfig(login string, conf *entity.Config) error {
	return nil
}

func (u *UserStorageMock) GetConfig(login string) (*entity.Config, error) {
	return &entity.Config{}, nil
}

func TestNewASTTree(t *testing.T) {
	type args struct {
		expression  *entity.Expression
		config      *config.Config
		queue       *MapQueueMock
		userStorage *UserStorageMock
	}
	conf := config.New("../../../config/.env")
	q := &MapQueueMock{}
	u := &UserStorageMock{}
	tests := []struct {
		name    string
		args    args
		want    *ASTTree
		wantErr bool
	}{
		{
			name: "maximum true test",
			args: args{
				expression: &entity.Expression{
					Expression:           "(2+2)+(2*2/2)-((3+3)+(3+-3))",
					CalculatedExpression: "(2+2)+(2*2/2)-((3+3)+(3+-3))",
				},
				config:      conf,
				queue:       q,
				userStorage: u,
			},
			want: &ASTTree{
				IsCalc: true,
			},
			wantErr: false,
		},
		{
			name: "error divide by zero",
			args: args{
				expression: &entity.Expression{
					Expression:           "2/0",
					CalculatedExpression: "2/0",
				},
				config:      conf,
				queue:       q,
				userStorage: u,
			},
			want: &ASTTree{
				IsCalc: true,
			},
			wantErr: false,
		},
		{
			name: "error AST",
			args: args{
				expression: &entity.Expression{
					Expression:           "wqer",
					CalculatedExpression: "+-*/",
				},
				config:      conf,
				queue:       q,
				userStorage: u,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "expression is calc",
			args: args{
				expression: &entity.Expression{
					Expression:           "(2+2)+(2*2/2)-((3+3)+(3+-3))",
					CalculatedExpression: "(2+2)+(2*2/2)-((3+3)+(3+-3))",
					IsCalc:               true,
				},
				config:      conf,
				queue:       q,
				userStorage: u,
			},
			want: &ASTTree{
				IsCalc: true,
			},
			wantErr: false,
		},
		{
			name: "expression is err",
			args: args{
				expression: &entity.Expression{
					Expression:           "(2+2)+(2*2/2)-((3+3)+(3+-3))",
					CalculatedExpression: "(2+2)+(2*2/2)-((3+3)+(3+-3))",
					Err:                  errors.New("error"),
				},
				config:      conf,
				queue:       q,
				userStorage: u,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewASTTree(tt.args.expression, tt.args.config, tt.args.queue, tt.args.userStorage)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewASTTree() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil &&
				tt.want.Value != got.Value &&
				tt.want.IsCalc != got.IsCalc &&
				tt.want.IsParent != got.IsParent &&
				tt.want.Err == got.Err {
				t.Errorf("NewASTTree() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewASTTreeDB(t *testing.T) {
	conf := config.New("../../../config/.env")
	q := &MapQueueMock{}
	u := &UserStorageMock{}
	type args struct {
		expression  *entity.Expression
		config      *config.Config
		queue       Queue
		userStorage UserStorage
	}
	tests := []struct {
		name    string
		args    args
		want    *ASTTree
		wantErr bool
	}{
		{
			name: "test 1",
			args: args{
				expression: &entity.Expression{
					Expression:           "(2+2)+(2*2/2)-((3+3)+(3+-3))",
					CalculatedExpression: "100",
					IsCalc:               true,
					Result:               float64(100),
				},
				config:      conf,
				queue:       q,
				userStorage: u,
			},
			want: &ASTTree{
				IsCalc: true,
				Value:  100,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewASTTreeDB(tt.args.expression, tt.args.config, tt.args.queue, tt.args.userStorage)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewASTTree() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil &&
				tt.want.Value != got.Value &&
				tt.want.IsCalc != got.IsCalc &&
				tt.want.IsParent != got.IsParent &&
				tt.want.Err == got.Err {
				t.Errorf("NewASTTree() = %v, want %v", got, tt.want)
			}
		})
	}
}
