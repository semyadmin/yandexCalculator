package arithmetic

import (
	"errors"
	"testing"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/config"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/tasks/queue"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/validator"
)

func TestNewASTTree(t *testing.T) {
	type args struct {
		expression string
		config     *config.Config
		queue      *queue.MapQueue
		validator  func(string) bool
	}
	conf := config.New()
	q := queue.NewMapQueue(queue.NewLockFreeQueue(), conf)
	tests := []struct {
		name    string
		args    args
		want    *ASTTree
		wantErr bool
	}{
		{
			name: "simple",
			args: args{
				expression: "-1+(1+2)*3",
				config:     conf,
				queue:      q,
				validator:  validator.Validator,
			},
			want: &ASTTree{
				Expression: "-1+(1+2)*3",
				Value:      "",
				IsCalc:     false,
				IsParent:   true,
				Err:        nil,
			},
			wantErr: false,
		},
		{
			name: "with error",
			args: args{
				expression: "-1+(1+2)*3+",
				config:     conf,
				queue:      q,
				validator:  validator.Validator,
			},
			want: &ASTTree{
				Expression: "-1+(1+2)*3+",
				Value:      "",
				IsCalc:     false,
				IsParent:   true,
				Err:        errors.New("invalid expression"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewASTTree(tt.args.expression, tt.args.config, tt.args.queue, tt.args.validator)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewASTTree() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want.Expression != got.Expression &&
				tt.want.Value != got.Value &&
				tt.want.IsCalc != got.IsCalc &&
				tt.want.IsParent != got.IsParent &&
				tt.want.Err == got.Err {
				t.Errorf("NewASTTree() = %v, want %v", got, tt.want)
			}
		})
	}
}
