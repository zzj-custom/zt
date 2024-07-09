package utils

import "testing"

func TestGenerateRandomNumber(t *testing.T) {
	type args struct {
		length int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "测试生成随机数字",
			args: args{
				length: 6,
			},
			want: 123456,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateRandomNumber(tt.args.length); got != tt.want {
				t.Errorf("GenerateRandomNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}
