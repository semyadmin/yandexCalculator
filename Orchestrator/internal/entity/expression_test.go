package entity

import (
	"errors"
	"testing"
	"time"
)

func TestNewExpression(t *testing.T) {
	type args struct {
		exp       string
		calcExp   string
		validator func(string) bool
		user      string
		start     time.Time
		updater   func([]byte) []byte
	}
	start := time.Now()
	tests := []struct {
		name string
		args args
		want *Expression
	}{
		{
			name: "valid expression",
			args: args{
				exp:       "exp",
				calcExp:   "calcExp",
				validator: func(s string) bool { return true },
				user:      "user",
				start:     start,
				updater:   func(b []byte) []byte { return b },
			},
			want: &Expression{
				Expression:           "exp",
				CalculatedExpression: "calcExp",
				User:                 "user",
				Start:                start,
				Err:                  nil,
			},
		},
		{
			name: "invalid expression",
			args: args{
				exp:       "exp",
				calcExp:   "calcExp",
				validator: func(s string) bool { return false },
				user:      "user",
				start:     start,
				updater:   func(b []byte) []byte { return b },
			},
			want: &Expression{
				Expression:           "exp",
				CalculatedExpression: "calcExp",
				User:                 "user",
				Start:                start,
				Err:                  errors.New(errInvalidExpression),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewExpression(tt.args.exp, tt.args.calcExp, tt.args.validator, tt.args.user, tt.args.start, tt.args.updater)
			if got.Expression != tt.want.Expression {
				t.Errorf("NewExpression() = %v, want %v", got, tt.want)
			}
			if got.CalculatedExpression != tt.want.CalculatedExpression {
				t.Errorf("NewExpression() = %v, want %v", got, tt.want)
			}
			if got.User != tt.want.User {
				t.Errorf("NewExpression() = %v, want %v", got, tt.want)
			}
			if got.Start != tt.want.Start {
				t.Errorf("NewExpression() = %v, want %v", got, tt.want)
			}
			if tt.want.Err != nil && got.Err.Error() != tt.want.Err.Error() {
				t.Errorf("NewExpression() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExpression_SetResult(t *testing.T) {
	type fields struct {
		ID                   uint64
		Start                time.Time
		Duration             int64
		Expression           string
		CalculatedExpression string
		Result               float64
		IsCalc               bool
		Err                  error
		User                 string
	}
	type args struct {
		r   float64
		err error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "valid expression",
			fields: fields{
				ID:                   1,
				Start:                time.Now(),
				Duration:             0,
				Expression:           "exp",
				CalculatedExpression: "calcExp",
				Result:               1,
				IsCalc:               true,
				Err:                  nil,
				User:                 "user",
			},
			args: args{
				r:   1,
				err: nil,
			},
		},
		{
			name: "invalid expression",
			fields: fields{
				ID:                   1,
				Start:                time.Now(),
				Duration:             0,
				Expression:           "exp",
				CalculatedExpression: "calcExp",
				Result:               0,
				IsCalc:               false,
				Err:                  errors.New(errInvalidExpression),
				User:                 "user",
			},
			args: args{
				r:   1,
				err: errors.New(errInvalidExpression),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Expression{
				ID:                   tt.fields.ID,
				Start:                tt.fields.Start,
				Duration:             tt.fields.Duration,
				Expression:           tt.fields.Expression,
				CalculatedExpression: tt.fields.CalculatedExpression,
				Result:               tt.fields.Result,
				IsCalc:               tt.fields.IsCalc,
				Err:                  tt.fields.Err,
				User:                 tt.fields.User,
			}
			e.SetResult(tt.args.r, tt.args.err)

			if e.Result != tt.fields.Result {
				t.Errorf("SetResult() = %v, want %v", e.Result, tt.args.r)
			}
			if e.Err != nil && e.Err.Error() != tt.args.err.Error() {
				t.Errorf("SetResult() = %v, want %v", e.Err, tt.args.err)
			}
		})
	}
}

func TestExpression_SetId(t *testing.T) {
	type fields struct {
		ID                   uint64
		Start                time.Time
		Duration             int64
		Expression           string
		CalculatedExpression string
		Result               float64
		IsCalc               bool
		Err                  error
		User                 string
	}
	type args struct {
		id uint64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "valid expression",
			fields: fields{
				ID:                   1,
				Start:                time.Now(),
				Duration:             0,
				Expression:           "exp",
				CalculatedExpression: "calcExp",
				Result:               1,
				IsCalc:               true,
				Err:                  nil,
				User:                 "user",
			},
			args: args{
				id: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Expression{
				ID:                   tt.fields.ID,
				Start:                tt.fields.Start,
				Duration:             tt.fields.Duration,
				Expression:           tt.fields.Expression,
				CalculatedExpression: tt.fields.CalculatedExpression,
				Result:               tt.fields.Result,
				IsCalc:               tt.fields.IsCalc,
				Err:                  tt.fields.Err,
				User:                 tt.fields.User,
			}
			e.SetId(tt.args.id)

			if e.ID != tt.fields.ID {
				t.Errorf("SetId() = %v, want %v", e.ID, tt.args.id)
			}
		})
	}
}
