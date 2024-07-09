package email

// QQ 邮箱：
// SMTP 服务器地址：smtp.qq.com（SSL协议端口：465/994 | 非SSL协议端口：25）
// 163 邮箱：
// SMTP 服务器地址：smtp.163.com（端口：25）

type Config struct {
	Host     string `json:"host" toml:"host"`
	Port     int    `json:"port" toml:"port"`
	UserName string `json:"userName" toml:"user_name" mapstructure:"user_name"`
	Password string `json:"password" toml:"password"`
}

type Email struct {
	cfg    *Config
	Extend *Options `json:"extend"`
}

type Options struct {
	Cc      []string `json:"cc"`      // 抄送
	Bcc     []string `json:"bcc"`     // 密送
	Attach  string   `json:"attach"`  // 附件
	Video   string   `json:"video"`   // 视频
	Image   string   `json:"image"`   // 图片
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

func WithEmailOptionsBcc(bcc []string) Option {
	return func(opts *Options) {
		opts.Bcc = bcc
	}
}

func WithOptionsAttach(attach string) Option {
	return func(opts *Options) {
		opts.Attach = attach
	}
}

func WithOptionsVideo(video string) Option {
	return func(opts *Options) {
		opts.Video = video
	}
}

func WithOptionsImage(image string) Option {
	return func(opts *Options) {
		opts.Image = image
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

func WithOptionsSubject(subject string) Option {
	return func(opts *Options) {
		opts.Subject = subject
	}
}
