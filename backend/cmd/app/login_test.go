package app

import (
	"context"
	"reflect"
	"testing"
	"zt/backend/internal/response"
)

func TestApp_Captcha(t *testing.T) {
	type fields struct {
		ctx context.Context
	}
	type args struct {
		to string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *response.Reply
	}{
		{
			name: "test",
			fields: fields{
				ctx: context.TODO(),
			},
			args: args{
				to: "1844066417@qq.com",
			},
			want: &response.Reply{
				Code:   0,
				Msg:    "success",
				Result: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &App{
				ctx: tt.fields.ctx,
			}
			if got := a.Captcha(tt.args.to); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Captcha() = %v, want %v", got, tt.want)
			}
		})
	}
}
