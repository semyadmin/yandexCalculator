package config

import (
	"context"
	"testing"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/storage/postgresql"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/web_socket/manager"
)

func TestNew(t *testing.T) {
	type args struct {
		confFile string
	}
	tests := []struct {
		name string
		args args
		want *Config
	}{
		{
			name: "TestNew",
			args: args{
				confFile: "",
			},
			want: &Config{
				Host:      "localhost",
				HttpPort:  "8080",
				TCPPort:   "7777",
				Db:        postgresql.NewPostgresConnect("localhost", "5432", "postgres", "postgres", "orchestrator"),
				WSmanager: manager.NewManager(context.Background()),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.args.confFile)
			if got.Host != tt.want.Host {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
			if got.HttpPort != tt.want.HttpPort {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
			if got.TCPPort != tt.want.TCPPort {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
			if got.Db == nil {
				t.Errorf("DB is nil")
			}
			if got.WSmanager == nil {
				t.Errorf("WSmanager is nil")
			}
		})
	}
}
