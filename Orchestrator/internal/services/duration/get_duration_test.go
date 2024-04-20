package duration

import (
	"errors"
	"reflect"
	"testing"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/entity"
	jwttoken "github.com/adminsemy/yandexCalculator/Orchestrator/internal/services/jwt_token"
)

type UserStorageMockCorrect struct{}

func (u UserStorageMockCorrect) SetConfig(login string, conf *entity.Config) error {
	return nil
}

func (u UserStorageMockCorrect) GetConfig(login string) (*entity.Config, error) {
	return &entity.Config{
		Plus:     1,
		Minus:    2,
		Multiply: 3,
		Divide:   4,
	}, nil
}

type UserStorageMockIncorrect struct{}

func (u UserStorageMockIncorrect) SetConfig(login string, conf *entity.Config) error {
	return nil
}

func (u UserStorageMockIncorrect) GetConfig(login string) (*entity.Config, error) {
	return nil, errors.New("error")
}

type UserStorageMockIncorrectMarshalling struct{}

func (u UserStorageMockIncorrectMarshalling) SetConfig(login string, conf *entity.Config) error {
	return nil
}

func (u UserStorageMockIncorrectMarshalling) GetConfig(login string) (*entity.Config, error) {
	return nil, nil
}

func TestGetDuration(t *testing.T) {
	type args struct {
		token       string
		userStorage UserStorage
	}
	token, _ := jwttoken.GenerateToken("test", 15)
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "correct",
			args: args{
				token:       token,
				userStorage: UserStorageMockCorrect{},
			},
			want:    []byte(`{"plus":1,"minus":2,"multiply":3,"divide":4}`),
			wantErr: false,
		},
		{
			name: "incorrect token",
			args: args{
				token:       "incorrect",
				userStorage: UserStorageMockCorrect{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "incorrect user storage",
			args: args{
				token:       token,
				userStorage: UserStorageMockIncorrect{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error marshalling",
			args: args{
				token:       token,
				userStorage: UserStorageMockIncorrectMarshalling{},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetDuration(tt.args.token, tt.args.userStorage)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDuration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}
