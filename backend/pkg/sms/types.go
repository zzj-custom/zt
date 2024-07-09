package sms

type Config struct {
}

type Sms struct {
	cfg    *Config
	Extend *Options `json:"extend"`
}

type Options struct {
	Cc      []string `json:"cc"`      // 抄送
	Subject string   `json:"subject"` // 主题
	Account string   `json:"account"` // 账号名称
	Web     string   `json:"web"`     // app名称
}

type Option func(*Options)

func WithEmailOptionsCc(cc []string) Option {
	return func(opts *Options) {
		opts.Cc = cc
	}
}

func WithOptionsAccount(account string) Option {
	return func(opts *Options) {
		opts.Account = account
	}
}

func WithOptionsWeb(web string) Option {
	return func(opts *Options) {
		opts.Web = web
	}
}
