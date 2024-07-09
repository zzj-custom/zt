package sms

import (
	"sync"
	"zt/backend/pkg/email"
)

const (
	mobileRegex = `^(13[0-9]|14[5-9]|15[0-3]|15[5-9]|16[2-7]|17[0-8]|18[0-9]|19[0-3]|19[5-9])[0-9]{8}$`
)

var (
	sms     *Sms
	smsOnce sync.Once
)

func NewSms(cfg *Config) *Sms {
	smsOnce.Do(func() {
		sms = new(Sms)
		sms.cfg = cfg
	})
	return sms
}

func (s *Sms) Send(to string, code int, opts ...email.Option) error {
	// TODO 阿里云短信发送
	return nil
}
