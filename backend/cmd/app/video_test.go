package app

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"zt/backend/internal/response"
	_ "zt/backend/pkg/extractor/bilibili"
)

func TestApp_List(t *testing.T) {
	type fields struct {
		ctx context.Context
	}
	type args struct {
		u string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *response.Reply
	}{
		{
			name: "测试视频列表",
			fields: fields{
				ctx: context.TODO()},
			args: args{
				u: "https://www.bilibili.com/video/BV1VH4y1w7De",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &App{
				ctx: tt.fields.ctx,
			}
			if got := a.List(tt.args.u); !reflect.DeepEqual(got, tt.want) {
				r, _ := json.Marshal(got.Result)
				fmt.Println(string(r))
				t.Errorf("List() = %v, want %v", got, tt.want)
			}
		})
	}
}
