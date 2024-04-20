package expression

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/config"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/entity"
	jwttoken "github.com/adminsemy/yandexCalculator/Orchestrator/internal/services/jwt_token"
	responseexpression "github.com/adminsemy/yandexCalculator/Orchestrator/internal/services/response_expression"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/storage/memory"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/tasks/queue"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/tasks/upgrade"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/validator"
)

var rightStorage = &entity.Expression{
	Expression: "2+2",
	Result:     4,
	Start:      time.Now(),
}

type StorageMock struct{}

func (s *StorageMock) GeByExpression(string, string) (*entity.Expression, error) {
	return rightStorage, nil
}

func (s *StorageMock) Set(*entity.Expression) error {
	return nil
}

func (s *StorageMock) GetById(uint64, string) (*entity.Expression, error) {
	return nil, nil
}

type StorageEmptyMock struct{}

func (s *StorageEmptyMock) GeByExpression(string, string) (*entity.Expression, error) {
	return nil, nil
}

func (s *StorageEmptyMock) Set(*entity.Expression) error {
	return nil
}

func (s *StorageEmptyMock) GetById(uint64, string) (*entity.Expression, error) {
	return nil, nil
}

type StorageErrExpressionNotExistsMock struct{}

func (s *StorageErrExpressionNotExistsMock) GeByExpression(string, string) (*entity.Expression, error) {
	return nil, memory.ErrExpressionNotExists
}

func (s *StorageErrExpressionNotExistsMock) Set(*entity.Expression) error {
	return nil
}

func (s *StorageErrExpressionNotExistsMock) GetById(uint64, string) (*entity.Expression, error) {
	return nil, nil
}

type UserStorageMock struct{}

func (s *UserStorageMock) SetConfig(string, *entity.Config) error {
	return nil
}

func (s *UserStorageMock) GetConfig(string) (*entity.Config, error) {
	return &entity.Config{
		Plus:     0,
		Minus:    0,
		Multiply: 0,
		Divide:   0,
	}, nil
}

type UserStorageWithErrorMock struct{}

func (s *UserStorageWithErrorMock) SetConfig(string, *entity.Config) error {
	return nil
}

func (s *UserStorageWithErrorMock) GetConfig(string) (*entity.Config, error) {
	return nil, errors.New("error")
}

type QueueMock struct{}

func (q *QueueMock) Enqueue(queue.Expression) {}

func TestNewExpression(t *testing.T) {
	type args struct {
		conf        *config.Config
		storage     Storage
		queue       Queue
		expression  string
		token       string
		userStorage UserStorage
		now         time.Time
	}
	token, _ := jwttoken.GenerateToken("test", 15)
	rightStorageMarshal, _ := json.Marshal(responseexpression.NewResponseExpression(
		rightStorage.ID,
		rightStorage.Expression,
		rightStorage.Start,
		rightStorage.Duration,
		rightStorage.IsCalc,
		rightStorage.Result,
		rightStorage.Err,
	))
	conf := config.New("../../../config/.env")
	now := time.Now()
	newExp := entity.NewExpression("2+2--2*3/4", "2+2--2*3/4", validator.Validator, "test", now, upgrade.Upgrade)
	newExp.ID = 1
	newExp.Duration = duration(newExp.Expression, &entity.Config{Plus: 0, Minus: 0, Multiply: 0, Divide: 0})
	newExpMarshal, _ := json.Marshal(responseexpression.NewResponseExpression(
		newExp.ID,
		newExp.Expression,
		newExp.Start,
		newExp.Duration,
		newExp.IsCalc,
		newExp.Result,
		newExp.Err,
	))
	errorDb := entity.NewExpression("2+2++", "2+2++", validator.Validator, "test", now, upgrade.Upgrade)
	errorDb.ID = 2
	errorDb.Err = errors.New("db error")
	errorDbMarshal, _ := json.Marshal(responseexpression.NewResponseExpression(
		errorDb.ID,
		errorDb.Expression,
		errorDb.Start,
		errorDb.Duration,
		errorDb.IsCalc,
		errorDb.Result,
		errorDb.Err,
	))
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "find expression right",
			args: args{
				conf:        conf,
				storage:     &StorageMock{},
				queue:       &QueueMock{},
				expression:  "2+2",
				token:       token,
				userStorage: &UserStorageMock{},
				now:         now,
			},
			want:    rightStorageMarshal,
			wantErr: false,
		},
		{
			name: "error token",
			args: args{
				conf:        conf,
				storage:     &StorageMock{},
				queue:       &QueueMock{},
				expression:  "2+2",
				token:       "wrong token",
				userStorage: &UserStorageMock{},
				now:         now,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "new expression",
			args: args{
				conf:        conf,
				storage:     &StorageErrExpressionNotExistsMock{},
				queue:       &QueueMock{},
				expression:  "2+2--2*3/4",
				token:       token,
				userStorage: &UserStorageMock{},
				now:         now,
			},
			want:    newExpMarshal,
			wantErr: false,
		},
		{
			name: "db errror true",
			args: args{
				conf:        conf,
				storage:     &StorageErrExpressionNotExistsMock{},
				queue:       &QueueMock{},
				expression:  "2+2++",
				token:       token,
				userStorage: &UserStorageMock{},
				now:         now,
			},
			want:    errorDbMarshal,
			wantErr: false,
		},
		{
			name: "wrong expression",
			args: args{
				conf:        conf,
				storage:     &StorageErrExpressionNotExistsMock{},
				queue:       &QueueMock{},
				expression:  "2+2",
				token:       token,
				userStorage: &UserStorageWithErrorMock{},
				now:         now,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "nil expression",
			args: args{
				conf:        conf,
				storage:     &StorageEmptyMock{},
				queue:       &QueueMock{},
				expression:  "2+2",
				token:       token,
				userStorage: &UserStorageMock{},
				now:         now,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewExpression(
				tt.args.conf,
				tt.args.storage,
				tt.args.queue,
				tt.args.expression,
				tt.args.token,
				tt.args.userStorage,
				tt.args.now)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewExpression() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewExpression() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}
