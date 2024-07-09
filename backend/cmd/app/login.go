package app

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"log/slog"
	"zt/backend/cmd/bootstrapper/config"
	"zt/backend/internal/code"
	"zt/backend/internal/repository"
	"zt/backend/internal/response"
	"zt/backend/internal/utils"
	"zt/backend/pkg/email"
)

const defaultAccountName = "游客"

func (a *App) Captcha(to string) *response.Reply {
	if to == "" {
		return response.FailReply(response.QueryParamsError)
	}

	// 验证码是否存在
	c, nc := utils.GenerateRandomNumber(6), code.NewCode(to)
	ok, err := nc.Validate()
	if err != nil {
		slog.With(
			slog.String("to", to),
		).With("err", err).Error("验证验证码错误")
		return response.FailReply(response.ValidateCaptchaFail)
	}

	if ok {
		return response.FailReply(response.CaptchaRepeat)
	}

	go func() {
		// 判断账号是否存在
		accountRepo := repository.NewAccountRepository()
		account, err := accountRepo.StatByEmail(to)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			slog.With(
				slog.String("to", to),
			).With("err", err).Error("查询账号错误")
			return
		}

		var (
			name    = defaultAccountName
			subject = "注册新账号"
		)

		if account != nil {
			name = account.Name
			subject = "登录"
		}

		// 发送验证码
		e := email.NewEmail(config.Get().Email)
		if err = e.Send(
			to,
			c,
			email.WithOptionsAccount(name),
			email.WithOptionsWeb("zt"),
			email.WithOptionsSubject(subject),
		); err != nil {
			slog.With(
				slog.String("to", to),
				slog.Int("code", c),
			).With("err", err).Error("send code error")
			return
		}

		// 创建记录
		if err = repository.NewCaptchaRepository().CreateByModel(&repository.CaptchaLog{
			Email:   to,
			Captcha: c,
			Status:  repository.CaptchaStatusNoExpired,
		}); err != nil {
			slog.With(
				slog.String("to", to),
				slog.Int("code", c),
			).With("err", err).Error("create code error")
			return
		}

		// redis保存验证码
		if err = nc.Set(c); err != nil {
			slog.With(
				slog.Int("code", c),
			).With("err", err).Error("set code error")
			return
		}
	}()

	slog.With(
		slog.String("to", to),
		slog.Int("code", c),
	).Info("发送验证码成功")

	return response.OkReply(nil)
}

type LoginRequest struct {
	Email   string `json:"email"`
	Captcha int    `json:"captcha"`
}

func (a *App) Login(req LoginRequest) *response.Reply {
	// 判断验证码是否可用
	captchaRepo := repository.NewCaptchaRepository()
	captcha, err := captchaRepo.StatByEmailAndCode(req.Email, req.Captcha)
	if err != nil {
		slog.With(
			slog.String("email", req.Email),
			slog.Int("code", req.Captcha),
		).With("err", err).Error("查询验证码错误")
		return response.FailReply(response.ValidateCaptchaFail)
	}

	if captcha.Status != repository.CaptchaStatusNoExpired {
		return response.FailReply(response.CaptchaExpired)
	}

	// 如果Account里面没有数据，就直接创建一个
	accountRepository := repository.NewAccountRepository()
	account, err := accountRepository.StatByEmail(req.Email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		insert := &repository.Account{
			Name:  defaultAccountName,
			Email: req.Email,
		}
		id, err := accountRepository.CreateByModel(insert)
		if err != nil {
			slog.With("email", req.Email).With("err", err).Error("创建账号错误")
			return response.FailReply(response.LoginFail)
		}
		insert.Id = id
		account = insert
	} else if err != nil {
		slog.With("email", req.Email).With("err", err).Error("查询账号错误")
		return response.FailReply(response.LoginFail)
	}

	go func() {
		// 更新验证码状态
		if err = captchaRepo.UpdateByModel(captcha.Id, &repository.CaptchaLog{
			Status: repository.CaptchaStatusExpired,
		}); err != nil {
			slog.With(
				slog.String("email", req.Email),
				slog.Int("code", req.Captcha),
			).With("err", err).Error("更新验证码错误")
			return
		}
	}()

	return response.OkReply(account)
}
