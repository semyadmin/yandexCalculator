package entity

import (
	"testing"
)

func TestNewOperation(t *testing.T) {
	type args struct {
		id       string
		first    float64
		second   float64
		operator string
		duration uint64
	}
	tests := []struct {
		name string
		args args
		want *operation
	}{
		{
			name: "should create new operation",
			args: args{
				id:       "1",
				first:    10,
				second:   20,
				operator: "+",
				duration: 1,
			},
			want: &operation{
				id:        "1",
				first:     10,
				second:    20,
				operation: "+",
				duration:  1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewOperation(tt.args.id, tt.args.first, tt.args.second, tt.args.operator, tt.args.duration)
			if got.id != tt.want.id {
				t.Errorf("NewOperation() got.id = %v, want %v", got.id, tt.want.id)
			}
			if got.first != tt.want.first {
				t.Errorf("NewOperation() got.first = %v, want %v", got.first, tt.want.first)
			}
			if got.second != tt.want.second {
				t.Errorf("NewOperation() got.second = %v, want %v", got.second, tt.want.second)
			}
			if got.operation != tt.want.operation {
				t.Errorf("NewOperation() got.operation = %v, want %v", got.operation, tt.want.operation)
			}
			if got.duration != tt.want.duration {
				t.Errorf("NewOperation() got.duration = %v, want %v", got.duration, tt.want.duration)
			}
			if got.ch == nil {
				t.Errorf("NewOperation() got.ch = %v, want %v", got.ch, tt.want.ch)
			}
		})
	}
}

func TestOperationProperty(t *testing.T) {
	newOperator := NewOperation("1", 10, 20, "+", 1)
	if newOperator.Id() != "1" {
		t.Errorf("NewOperation() got.id = %v, want %v", newOperator.Id(), "1")
	}
	if newOperator.First() != 10 {
		t.Errorf("NewOperation() got.first = %v, want %v", newOperator.First(), 10)
	}
	if newOperator.Second() != 20 {
		t.Errorf("NewOperation() got.second = %v, want %v", newOperator.Second(), 20)
	}
	if newOperator.Operation() != "+" {
		t.Errorf("NewOperation() got.operation = %v, want %v", newOperator.Operation(), "+")
	}
	if newOperator.Duration() != 1 {
		t.Errorf("NewOperation() got.duration = %v, want %v", newOperator.Duration(), 1)
	}
	newOperator.Result(30)
	ch := newOperator.ResultChan()
	res := <-ch
	if res != 30 {
		t.Errorf("NewOperation() got.result = %v, want %v", res, 30)
	}
	newOperator.Result(30)
	res = newOperator.GetResult()
	if res != 30 {
		t.Errorf("NewOperation() got.result = %v, want %v", res, 30)
	}
	newOperator.Error("error")
	err := newOperator.GetError()
	if err.Error() != "error" {
		t.Errorf("NewOperation() got.error = %v, want %v", err, "error")
	}
}
