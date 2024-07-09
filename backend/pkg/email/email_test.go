package email

import "testing"

func TestEmail_Send(t *testing.T) {
	type fields struct {
		cfg    *Config
		Extend *Options
	}
	type args struct {
		to   string
		code int
		opts []Option
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "测试发送",
			fields: fields{
				cfg: &Config{
					Host:     "smtp.qq.com",
					Port:     25,
					UserName: "1844066417@qq.com",
					Password: "jthlffwdvnfmbedi",
				},
			},
			args: args{
				to:   "1844066417@qq.com",
				code: 123456,
				opts: []Option{
					WithOptionsWeb("zt"),
					WithOptionsAccount("zzj"),
					WithOptionsSubject("测试发送"),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Email{
				cfg:    tt.fields.cfg,
				Extend: tt.fields.Extend,
			}
			if err := e.Send(tt.args.to, tt.args.code, tt.args.opts...); (err != nil) != tt.wantErr {
				t.Errorf("Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
