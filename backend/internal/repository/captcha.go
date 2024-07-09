package repository

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"log/slog"
	"reflect"
	"sync"
	"time"
	"zt/backend/pkg/iMysql"
)

const (
	CaptchaStatusNoExpired = 1
	CaptchaStatusExpired   = 2
)

// CaptchaLog 验证码表
type CaptchaLog struct {
	Id        int       `gorm:"column:id;type:bigint(20);AUTO_INCREMENT;primary_key;comment:自增id" json:"id"`
	Email     string    `gorm:"column:email;type:varchar(128);comment:邮箱地址;NOT NULL" json:"email"`
	Captcha   int       `gorm:"column:captcha;type:int(11);default:0;comment:验证码;NOT NULL" json:"captcha"`
	Status    int       `gorm:"column:status;type:int(11);default:1;comment:状态，1-未过期，2-已过期;NOT NULL" json:"status"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`
}

func (m *CaptchaLog) TableName() string {
	return "captcha_log"
}

type CaptchaRepository struct {
	db *gorm.DB
}

var (
	captchaRepository *CaptchaRepository
	captchaOnce       sync.Once
)

func NewCaptchaRepository() *CaptchaRepository {
	captchaOnce.Do(func() {
		conn, err := iMysql.Conn("default")
		if err != nil {
			slog.With("err", err).Error("captcha连接失败")
			return
		}
		captchaRepository = new(CaptchaRepository)
		captchaRepository.db = conn
	})
	return captchaRepository
}

func (receiver CaptchaRepository) StatByEmailAndCode(email string, captcha int) (*CaptchaLog, error) {
	var record CaptchaLog
	result := receiver.db.Where("email=? AND captcha = ?", email, captcha).Order("created_at desc").First(&record)
	return &record, result.Error
}

func (receiver CaptchaRepository) CreateByModel(model *CaptchaLog) error {
	return receiver.db.Model(&CaptchaLog{}).Create(model).Error
}

func (receiver CaptchaRepository) UpdateByModel(id int, model *CaptchaLog) error {
	if id <= 0 {
		return errors.New("id不能为空")
	}

	if reflect.DeepEqual(model, &CaptchaLog{}) {
		return errors.New("model不能为空")
	}

	return receiver.db.Model(&CaptchaLog{}).Where("id=?", id).Updates(model).Error
}
