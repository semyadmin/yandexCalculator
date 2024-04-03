package validator

import "testing"

func TestValidator(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "simple",
			args: args{
				str: "1+2",
			},
			want: true,
		},
		{
			name: "simple error",
			args: args{
				str: "-1",
			},
			want: false,
		},
		{
			name: "first minus",
			args: args{
				str: "-1+22222222*3",
			},
			want: true,
		},
		{
			name: "first plus error",
			args: args{
				str: "+1+22222222*3",
			},
			want: false,
		},
		{
			name: "last plus error",
			args: args{
				str: "1+22222222*3+",
			},
			want: false,
		},
		{
			name: "invalid negative number",
			args: args{
				str: "1+22222222*--3",
			},
			want: false,
		},
		{
			name: "valid negative number",
			args: args{
				str: "1+22222222*-3",
			},
			want: true,
		},
		{
			name: "valid negative number with brackets",
			args: args{
				str: "1+22222222/(-3--4)",
			},
			want: true,
		},
		{
			name: "random string",
			args: args{
				str: "1+22222222*p-4",
			},
			want: false,
		},
		{
			name: "invalid with opened bracket",
			args: args{
				str: "1+22222222*(-4",
			},
			want: false,
		},
		{
			name: "invalid with closed bracket",
			args: args{
				str: "1+22222222*-4)",
			},
			want: false,
		},
		{
			name: "invalid without operators before opened bracket",
			args: args{
				str: "1+22222222(*4-4)",
			},
			want: false,
		},
		{
			name: "invalid with opened bracket with operators",
			args: args{
				str: "1+22222222-(*4-4)",
			},
			want: false,
		},
		{
			name: "invalid with opened bracket without operators",
			args: args{
				str: "1+22222222(4-4)-5",
			},
			want: false,
		},
		{
			name: "invalid with closed bracket with operators",
			args: args{
				str: "1+22222222*(4-4)5",
			},
			want: false,
		},
		{
			name: "invalid with closed bracket with operators",
			args: args{
				str: "0.1+22222222*(4-4.)+5",
			},
			want: false,
		},
		{
			name: "invalid with first closed bracket",
			args: args{
				str: "0.1+22222222)-4-4.(+5",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Validator(tt.args.str)
			if got != tt.want {
				t.Errorf("Validator() got = %v, want %v", got, tt.want)
			}
		})
	}
}
