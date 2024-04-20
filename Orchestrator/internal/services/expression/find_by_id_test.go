package expression

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/entity"
	jwttoken "github.com/adminsemy/yandexCalculator/Orchestrator/internal/services/jwt_token"
	responseexpression "github.com/adminsemy/yandexCalculator/Orchestrator/internal/services/response_expression"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/tasks/upgrade"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/validator"
)

type StorageMockGetById struct{}

func (s StorageMockGetById) GetById(number uint64, token string) (*entity.Expression, error) {
	ent := entity.NewExpression("2+2", "2+2", validator.Validator, "test", time.Now(), upgrade.Upgrade)
	return ent, nil
}

func (s StorageMockGetById) Set(*entity.Expression) error {
	return nil
}

func (s *StorageMockGetById) GeByExpression(string, string) (*entity.Expression, error) {
	return nil, nil
}

type StorageMockGetByIdNotFound struct{}

func (s StorageMockGetByIdNotFound) GetById(number uint64, token string) (*entity.Expression, error) {
	return nil, errors.New("error")
}

func (s StorageMockGetByIdNotFound) Set(*entity.Expression) error {
	return nil
}

func (s *StorageMockGetByIdNotFound) GeByExpression(string, string) (*entity.Expression, error) {
	return nil, nil
}

func TestGetById(t *testing.T) {
	token, _ := jwttoken.GenerateToken("test", 15)
	validEnt := entity.NewExpression("2+2", "2+2", validator.Validator, "test", time.Now(), upgrade.Upgrade)
	validEntMarshal, _ := json.Marshal(responseexpression.NewResponseExpression(
		validEnt.ID,
		validEnt.Expression,
		validEnt.Start,
		validEnt.Duration,
		validEnt.IsCalc,
		validEnt.Result,
		validEnt.Err,
	))
	type args struct {
		storage Storage
		number  string
		token   string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "positive test",
			args: args{
				storage: &StorageMockGetById{},
				number:  "2",
				token:   token,
			},
			want:    validEntMarshal,
			wantErr: false,
		},
		{
			name: "error token",
			args: args{
				storage: &StorageMockGetById{},
				number:  "2",
				token:   "123",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error parse number",
			args: args{
				storage: &StorageMockGetById{},
				number:  "2+1",
				token:   token,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error found expression",
			args: args{
				storage: &StorageMockGetByIdNotFound{},
				number:  "2",
				token:   token,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetById(tt.args.storage, tt.args.number, tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetById() = %v, want %v", got, tt.want)
			}
		})
	}
}
