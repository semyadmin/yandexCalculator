package duration

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/config"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/entity"
	jwttoken "github.com/adminsemy/yandexCalculator/Orchestrator/internal/services/jwt_token"
)

type UserStorageMock struct{}

func (u UserStorageMock) SetConfig(user string, conf *entity.Config) error {
	return nil
}

func (u UserStorageMock) GetConfig(user string) (*entity.Config, error) {
	return &entity.Config{}, nil
}

type UserStorageMockWrong struct{}

func (u UserStorageMockWrong) SetConfig(user string, conf *entity.Config) error {
	return errors.New("error")
}

func (u UserStorageMockWrong) GetConfig(user string) (*entity.Config, error) {
	return &entity.Config{}, nil
}

func TestChangeDuration(t *testing.T) {
	type args struct {
		config      *config.Config
		data        []byte
		token       string
		userStorage UserStorage
	}
	conf := config.New("../../../config/.env")
	token, _ := jwttoken.GenerateToken("test", conf.TokenLimit)

	testDuration := entity.Config{
		Plus:     10,
		Minus:    10,
		Multiply: 10,
		Divide:   10,
	}
	data, _ := json.Marshal(testDuration)
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "no errors",
			args: args{
				config:      conf,
				data:        data,
				token:       token,
				userStorage: UserStorageMock{},
			},
			want:    data,
			wantErr: false,
		},
		{
			name: "error data",
			args: args{
				config:      conf,
				data:        nil,
				token:       token,
				userStorage: UserStorageMock{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error save data",
			args: args{
				config:      conf,
				data:        data,
				token:       token,
				userStorage: UserStorageMockWrong{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error token",
			args: args{
				config:      conf,
				data:        data,
				token:       "test",
				userStorage: UserStorageMock{},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ChangeDuration(tt.args.config, tt.args.data, tt.args.token, tt.args.userStorage)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChangeDuration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChangeDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}
