package upgrade

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
		{
			name: "multy with plus and brackets",
			args: args{
				exp: []byte("2+2-111+(888*89*2*2)+599*33*555+1"),
			},
			want: []byte("2+2-111+((888*89)*(2*2))+(599*33)*555+1"),
		},
		{
			name: "first minus",
			args: args{
				exp: []byte("-1+(1+2)*3"),
			},
			want: []byte("-1+(1+2)*3"),
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

func Test_upgradePlusMinus(t *testing.T) {
	type args struct {
		exp []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "simple plus",
			args: args{
				exp: []byte("1+2"),
			},
			want: []byte("(1+2)"),
		},
		{
			name: "simple minus",
			args: args{
				exp: []byte("1-2"),
			},
			want: []byte("(1-2)"),
		},
		{
			name: "4 plus",
			args: args{
				exp: []byte("2+2+2+2"),
			},
			want: []byte("(2+2)+(2+2)"),
		},
		{
			name: "4 minus",
			args: args{
				exp: []byte("2-2-2-2"),
			},
			want: []byte("(2-2)+(-2-2)"),
		},
		{
			name: "2 plus 1 minus",
			args: args{
				exp: []byte("2+2-2+2"),
			},
			want: []byte("(2+2)+(-2+2)"),
		},
		{
			name: "1 plus 1 minus",
			args: args{
				exp: []byte("2+2-2"),
			},
			want: []byte("(2+2)+-2"),
		},
		{
			name: "2 plus 1 minus",
			args: args{
				exp: []byte("2+2+2-2"),
			},
			want: []byte("(2+2)+(2-2)"),
		},
		{
			name: "1 plus 2 minus",
			args: args{
				exp: []byte("2-2+2-2"),
			},
			want: []byte("(2-2)+(2-2)"),
		},
		{
			name: "no brackets",
			args: args{
				exp: []byte("1+1/1+1+1*1+2"),
			},
			want: []byte("1+1/1+1+1*1+2"),
		},
		{
			name: "negative numbers",
			args: args{
				exp: []byte("-1+1+-1+2"),
			},
			want: []byte("(-1+1)+(-1+2)"),
		},
		{
			name: "with brackets",
			args: args{
				exp: []byte("-1+1+-1+(2*2)+(2+2+1)"),
			},
			want: []byte("(-1+1)+-1+(2*2)+((2+2)+1)"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := upgradePlusMinus(tt.args.exp); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("upgradeMultiDivide() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

func Test_upgradeDoubleMinus(t *testing.T) {
	type args struct {
		exp []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "simple plus",
			args: args{
				exp: []byte("-1--2+2+2-2+-2+2--2"),
			},
			want: []byte("-1+2+2+2-2+-2+2+2"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := updateDoubleMinus(tt.args.exp); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("upgradeMultiDivide() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

func TestUpgrade(t *testing.T) {
	type args struct {
		exp []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "mix",
			args: args{
				exp: []byte("5*9+8*6-9+7-11+333"),
			},
			want: []byte("(5*9)+(8*6)+(-9+7)+(-11+333)"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Upgrade(tt.args.exp); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Upgrade() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}
