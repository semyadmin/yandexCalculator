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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := upgradeMultiDivide(tt.args.exp); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("upgradeMultiDivide() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}
