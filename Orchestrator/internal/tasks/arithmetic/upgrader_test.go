package arithmetic

import (
	"reflect"
	"testing"
)

func Test_upgradeMultiDivide(t *testing.T) {
	type args struct {
		exp []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "simple1",
			args: args{
				exp: []byte("1*2"),
			},
			want: []byte("(1*2)"),
		},
		{
			name: "simple2",
			args: args{
				exp: []byte("1/2"),
			},
			want: []byte("1/2"),
		},
		{
			name: "first plus 1",
			args: args{
				exp: []byte("1+2*2"),
			},
			want: []byte("1+(2*2)"),
		},
		{
			name: "first plus 2",
			args: args{
				exp: []byte("1+2/2"),
			},
			want: []byte("1+2/2"),
		},
		{
			name: "multi",
			args: args{
				exp: []byte("1*2*2*2"),
			},
			want: []byte("(1*2)*(2*2)"),
		},
		{
			name: "div",
			args: args{
				exp: []byte("1/2/2/2"),
			},
			want: []byte("1/2/2/2"),
		},
		{
			name: "multy+div",
			args: args{
				exp: []byte("1*2/2*2"),
			},
			want: []byte("1*2/2*2"),
		},
		{
			name: "multy with negative number 1",
			args: args{
				exp: []byte("-11*2*2*2"),
			},
			want: []byte("(-11*2)*(2*2)"),
		},
		{
			name: "multy with negative number 1",
			args: args{
				exp: []byte("-11*2*-2*-2"),
			},
			want: []byte("(-11*2)*(-2*-2)"),
		},
		{
			name: "multy with brackets 1",
			args: args{
				exp: []byte("2*2*2*(2*2*2)*2"),
			},
			want: []byte("(2*2)*2*((2*2)*2)*2"),
		},
		{
			name: "multy with brackets 2",
			args: args{
				exp: []byte("2*2*2*(2*2*2)"),
			},
			want: []byte("(2*2)*2*((2*2)*2)"),
		},
		{
			name: "many brackets",
			args: args{
				exp: []byte("(2*2)*2*((2*2)*2*2)"),
			},
			want: []byte("((2*2))*2*(((2*2))*(2*2))"),
		},
		{
			name: "multy with plus minus",
			args: args{
				exp: []byte("2+2-111+888*8/9*2*2+599*33*555+1"),
			},
			want: []byte("2+2-111+888*8/9*(2*2)+(599*33)*555+1"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := upgradeMultiDivide(tt.args.exp); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("upgradeMultiDivide() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}
