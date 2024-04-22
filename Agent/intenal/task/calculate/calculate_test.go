package calculate

import (
	"testing"
)

type ExpMock struct {
	id        string
	first     float64
	second    float64
	operation string
	result    float64
}

func (e *ExpMock) Id() string {
	return e.id
}

func (e *ExpMock) First() float64 {
	return e.first
}

func (e *ExpMock) Second() float64 {
	return e.second
}

func (e *ExpMock) Duration() uint64 {
	return 0
}

func (e *ExpMock) Operation() string {
	return e.operation
}

func (e *ExpMock) Result() float64 {
	return e.result
}

func (e *ExpMock) Error() string {
	return ""
}

func (e *ExpMock) SetResult(r float64) {
	e.result = r
}

func (e *ExpMock) SetError(string) {
}

func TestCalculateGRPC(t *testing.T) {
	type args struct {
		exp Expression
	}
	tests := []struct {
		name string
		args args
		want Expression
	}{
		{
			name: "plus",
			args: args{
				exp: &ExpMock{
					id:        "1",
					first:     1,
					second:    1,
					operation: "+",
				},
			},
			want: &ExpMock{
				id:        "1",
				first:     1,
				second:    1,
				operation: "+",
				result:    2,
			},
		},
		{
			name: "minus",
			args: args{
				exp: &ExpMock{
					id:        "1",
					first:     1,
					second:    1,
					operation: "-",
				},
			},
			want: &ExpMock{
				id:        "1",
				first:     1,
				second:    1,
				operation: "-",
				result:    0,
			},
		},
		{
			name: "multiply",
			args: args{
				exp: &ExpMock{
					id:        "1",
					first:     1,
					second:    1,
					operation: "*",
				},
			},
			want: &ExpMock{
				id:        "1",
				first:     1,
				second:    1,
				operation: "*",
				result:    1,
			},
		},
		{
			name: "divide",
			args: args{
				exp: &ExpMock{
					id:        "1",
					first:     2,
					second:    2,
					operation: "/",
				},
			},
			want: &ExpMock{
				id:        "1",
				first:     2,
				second:    2,
				operation: "/",
				result:    1,
			},
		},
		{
			name: "error divide",
			args: args{
				exp: &ExpMock{
					id:        "1",
					first:     2,
					second:    0,
					operation: "/",
				},
			},
			want: &ExpMock{
				id:        "1",
				first:     2,
				second:    0,
				operation: "/",
				result:    0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateGRPC(tt.args.exp)
			if got.Result() != tt.want.Result() {
				t.Errorf("CalculateGRPC() = %v, want %v", got, tt.want)
			}
		})
	}
}
