package arithmetic

import (
	"testing"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/config"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/entity"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/tasks/queue"
)

func TestNewASTTree(t *testing.T) {
	type args struct {
		expression *entity.Expression
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
				expression: &entity.Expression{
					Expression: "1+2",
				},
				config: conf,
				queue:  q,
				validator: func(s string) bool {
					return true
				},
			},
			want: &ASTTree{
				IsCalc: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewASTTree(tt.args.expression, tt.args.config, tt.args.queue)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewASTTree() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want.Value != got.Value &&
				tt.want.IsCalc != got.IsCalc &&
				tt.want.IsParent != got.IsParent &&
				tt.want.Err == got.Err {
				t.Errorf("NewASTTree() = %v, want %v", got, tt.want)
			}
		})
	}
}
