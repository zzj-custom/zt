package email

import (
	"crypto/tls"
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/gomail.v2"
	"log/slog"
	"os"
	"regexp"
	"sync"
)

const (
	emailRegex = "^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]+$"
	template   = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>邮箱验证码</title>
    <style>
        table {
            width: 700px;
            margin: 0 auto;
        }
        #top {
            width: 700px;
            border-bottom: 1px solid #ccc;
            margin: 0 auto 30px;
        }
        #top table {
            font: 12px Tahoma, Arial, 宋体;
            height: 40px;
        }
        #content {
            width: 680px;
            padding: 0 10px;
            margin: 0 auto;
        }
        #content_top {
            line-height: 1.5;
            font-size: 14px;
            margin-bottom: 25px;
            color: #4d4d4d;
        }
        #content_top strong {
            display: block;
            margin-bottom: 15px;
        }
        #content_top strong span {
            color: #f60;
            font-size: 16px;
        }
        #verificationCode {
            color: #f60;
            font-size: 24px;
        }
        #content_bottom {
            margin-bottom: 30px;
        }
        #content_bottom small {
            display: block;
            margin-bottom: 20px;
            font-size: 12px;
            color: #747474;
        }
        #bottom {
            width: 700px;
            margin: 0 auto;
        }
        #bottom div {
            padding: 10px 10px 0;
            border-top: 1px solid #ccc;
            color: #747474;
            margin-bottom: 20px;
            line-height: 1.3em;
            font-size: 12px;
        }
        #content_top strong span {
            font-size: 18px;
            color: #FE4F70;
        }
        #sign {
            text-align: right;
            font-size: 18px;
            color: #FE4F70;
            font-weight: bold;
        }
        #verificationCode {
            height: 100px;
            width: 680px;
            text-align: center;
            margin: 30px 0;
        }
        #verificationCode div {
            height: 100px;
            width: 680px;
        }
        #verificationCode .button {
            margin-left: 10px;
            height: 80px;
            resize: none;
            border: none;
            outline: none;
            padding: 10px 15px;
            background: #ededed;
            border-radius: 17px;
            box-shadow: 6px 6px 12px #cccccc,
            -6px -6px 12px #ffffff;
        }
        #verificationCode .button:hover {
            box-shadow: inset 6px 6px 4px #d1d1d1,
            inset -6px -6px 4px #ffffff;
        }

        .code{
            color: #FE4F70;
            font-weight: bold;
            font-size: 42px;
            text-align: center;
        }
    </style>
</head>
<body>
<table>
    <tbody>
    <tr>
        <td>
            <div id="top">
                <table>
                    <tbody><tr><td></td></tr></tbody>
                </table>
            </div>
            <div id="content">
                <div id="content_top">
                    <strong>尊敬的%s用户：您好！</strong>
                    <strong>
                        您正在进行<span>登录</span>操作，请在登录验证码栏中输入以下验证码完成操作：
                    </strong>
                    <div id="verificationCode">
                        <button class="button"><span class="code">%d</span></button>
                    </div>
                </div>
                <div id="content_bottom">
                    <small>
                        注意：此操作可能会修改您的密码、登录邮箱或绑定手机。如非本人操作，请及时登录并修改密码以保证帐户安全
                        <br>（工作人员不会向你索取此验证码，请勿泄漏！)
                    </small>
                </div>
            </div>
            <div id="bottom">
                <div>
                    <p>此为系统邮件，请勿回复<br>
                        请保管好您的邮箱，避免账号被他人盗用
                    </p>
                    <p id="sign">—— 您的%s网站</p>
                </div>
            </div>
        </td>
    </tr>
    </tbody>
</table>
</body>`
)

var (
	email     *Email
	emailOnce sync.Once
)

func NewEmail(cfg *Config) *Email {
	emailOnce.Do(func() {
		email = new(Email)
		email.cfg = cfg
	})
	return email
}

func (e *Email) ValidateEmail(email []string) bool {
	if len(email) == 0 {
		slog.With("email", email).Error("邮箱不能为空")
		return false
	}

	// 编译正则表达式
	re := regexp.MustCompile(emailRegex)
	for _, s := range email {
		if !re.MatchString(s) {
			slog.With("email", s).Error("邮箱格式不正确")
			return false
		}
	}
	return true
}

func (e *Email) validate(str string, opts ...Option) bool {
	if !e.ValidateEmail([]string{str}) {
		return false
	}

	options := &Options{}
	for _, opt := range opts {
		opt(options)
	}

	if options.Cc != nil {
		if !e.ValidateEmail(options.Cc) {
			return false
		}
	}

	if options.Bcc != nil {
		if !e.ValidateEmail(options.Bcc) {
			return false
		}
	}

	if options.Attach != "" {
		if !e.fileExists(options.Attach) {
			return false
		}
	}

	if options.Video != "" {
		if !e.fileExists(options.Video) {
			return false
		}
	}

	if options.Image != "" {
		if !e.fileExists(options.Image) {
			return false
		}
	}
	e.Extend = options

	return true
}

// 检查文件是否存在
func (e *Email) fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			slog.With("filePath", filePath).Error("文件不存在")
			return false
		}
	}
	return true
}

func (e *Email) Send(to string, code int, opts ...Option) error {
	if !e.validate(to, opts...) {
		return errors.New("验证码发送失败")
	}

	var (
		m   = gomail.NewMessage()
		err error
	)

	// 设置请求header
	m = e.SetHeader(m, to)
	// 设置body
	m, err = e.SetBody(m, code)
	if err != nil {
		return errors.Wrap(err, "设置body失败")
	}
	// 设置附件
	m = e.Attach(m)
	// 发送邮件
	d := gomail.NewDialer(
		e.cfg.Host,
		e.cfg.Port,
		e.cfg.UserName,
		e.cfg.Password,
	)

	// 关闭SSL协议认证
	d.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
	}

	if err = d.DialAndSend(m); err != nil {
		return errors.Wrap(err, "发送邮件失败")
	}
	return nil
}

func (e *Email) SetHeader(m *gomail.Message, to string) *gomail.Message {
	m.SetHeader("From", e.cfg.UserName)
	//m.SetHeader("From", "alias"+"<zzjlovetl>") // 增加发件人别名

	m.SetHeader("To", to) // 收件人，可以多个收件人，但必须使用相同的 SMTP 连接
	if e.Extend.Cc != nil {
		m.SetHeader("Cc", e.Extend.Cc...) // 抄送，可以多个
	}

	if e.Extend.Bcc != nil {
		m.SetHeader("Bcc", e.Extend.Bcc...) // 密送，可以多个
	}

	if e.Extend.Subject != "" {
		m.SetHeader("Subject", e.Extend.Subject) // 邮件主题
	}
	return m
}

func (e *Email) SetBody(m *gomail.Message, code int) (*gomail.Message, error) {
	var (
		account = e.Extend.Account
		web     = e.Extend.Web
	)

	if e.Extend.Account == "" {
		account = "系统用户"
	}

	if e.Extend.Web == "" {
		web = "公主的智能小工具"
	}

	m.SetBody("text/html", fmt.Sprintf(template, account, code, web))
	return m, nil
}

func (e *Email) Attach(m *gomail.Message) *gomail.Message {
	if e.Extend.Attach != "" {
		m.Attach(e.Extend.Attach) // 附件文件，可以是文件，照片，视频等等
	}

	if e.Extend.Video != "" {
		m.Attach(e.Extend.Video) // 视频
	}

	if e.Extend.Image != "" {
		m.Attach(e.Extend.Image) // 图片
	}
	return m
}
