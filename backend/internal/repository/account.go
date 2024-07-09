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

// Account 账号表
type Account struct {
	Id        int       `gorm:"column:id;type:bigint(20) unsigned;AUTO_INCREMENT;comment:自增id;primary_key" json:"id"`
	Name      string    `gorm:"column:name;type:varchar(64);comment:账号名称;NOT NULL" json:"name"`
	Email     string    `gorm:"column:email;type:varchar(128);comment:邮箱地址;NOT NULL" json:"email"`
	Avatar    string    `gorm:"column:avatar;type:varchar(255);comment:头像地址;NOT NULL" json:"avatar"`
	Mobile    int       `gorm:"column:mobile;type:bigint(20);default:0;comment:手机号;NOT NULL" json:"mobile"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`
}

func (m *Account) TableName() string {
	return "account"
}

type AccountRepository struct {
	db *gorm.DB
}

var (
	accountRepository *AccountRepository
	accountOnce       sync.Once
)

func NewAccountRepository() *AccountRepository {
	accountOnce.Do(func() {
		conn, err := iMysql.Conn("default")
		if err != nil {
			slog.With(slog.String("database", "default")).With("err", err).Error("account数据库连接失败")
			return
		}
		accountRepository = new(AccountRepository)
		accountRepository.db = conn
	})
	return accountRepository
}

func (r *AccountRepository) StatByEmail(email string) (*Account, error) {
	var account Account
	result := r.db.Model(&Account{}).Where("email = ?", email).First(&account)
	return &account, result.Error
}

func (r *AccountRepository) CreateByModel(model *Account) (int, error) {
	if model == nil || reflect.DeepEqual(model, &Account{}) {
		return 0, errors.New("model is nil")
	}
	result := r.db.Model(&Account{}).Create(model).Error
	return model.Id, result
}
