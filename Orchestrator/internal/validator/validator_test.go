package validator

import "testing"

func TestValidator(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 bool
	}{
		{
			name: "simple",
			args: args{
				str: "1+2",
			},
			want:  "1+2",
			want1: true,
		},
		{
			name: "simple error",
			args: args{
				str: "-1",
			},
			want:  "",
			want1: false,
		},
		{
			name: "first minus",
			args: args{
				str: "-1+22222222*3",
			},
			want:  "-1+22222222*3",
			want1: true,
		},
		{
			name: "first plus error",
			args: args{
				str: "+1+22222222*3",
			},
			want:  "",
			want1: false,
		},
		{
			name: "last plus error",
			args: args{
				str: "1+22222222*3+",
			},
			want:  "",
			want1: false,
		},
		{
			name: "invalid negative number",
			args: args{
				str: "1+22222222*--3",
			},
			want:  "",
			want1: false,
		},
		{
			name: "valid negative number",
			args: args{
				str: "1+22222222*-3",
			},
			want:  "1+22222222*-3",
			want1: true,
		},
		{
			name: "valid negative number with brackets",
			args: args{
				str: "1+22222222/(-3--4)",
			},
			want:  "1+22222222/(-3--4)",
			want1: true,
		},
		{
			name: "random string",
			args: args{
				str: "1+22222222*p-4",
			},
			want:  "",
			want1: false,
		},
		{
			name: "invalid with opened bracket",
			args: args{
				str: "1+22222222*(-4",
			},
			want:  "",
			want1: false,
		},
		{
			name: "invalid with closed bracket",
			args: args{
				str: "1+22222222*-4)",
			},
			want:  "",
			want1: false,
		},
		{
			name: "invalid without operators before opened bracket",
			args: args{
				str: "1+22222222(*4-4)",
			},
			want:  "",
			want1: false,
		},
		{
			name: "invalid with opened bracket with operators",
			args: args{
				str: "1+22222222-(*4-4)",
			},
			want:  "",
			want1: false,
		},
		{
			name: "invalid with opened bracket without operators",
			args: args{
				str: "1+22222222(4-4)-5",
			},
			want:  "",
			want1: false,
		},
		{
			name: "invalid with closed bracket with operators",
			args: args{
				str: "1+22222222*(4-4)5",
			},
			want:  "",
			want1: false,
		},
		{
			name: "invalid with closed bracket with operators",
			args: args{
				str: "0.1+22222222*(4-4.)+5",
			},
			want:  "",
			want1: false,
		},
		{
			name: "invalid with first closed bracket",
			args: args{
				str: "0.1+22222222)-4-4.(+5",
			},
			want:  "",
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := Validator(tt.args.str)
			if got != tt.want {
				t.Errorf("Validator() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Validator() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
