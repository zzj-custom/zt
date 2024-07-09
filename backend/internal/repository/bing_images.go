package repository

import (
	"gorm.io/gorm"
	"log/slog"
	"sync"
	"time"
	"zt/backend/pkg/iMysql"
)

type BingImages struct {
	Id            uint      `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	Name          string    `gorm:"column:name;NOT NULL;comment:'图片名称'"`
	Copyright     string    `gorm:"column:copyright;default:;NOT NULL;comment:'版权'"`
	CopyrightLink string    `gorm:"column:copyright_link;default:;NOT NULL;comment:'版权链接'"`
	Url           string    `gorm:"column:url;NOT NULL;comment:'图片地址'"`
	Start         string    `gorm:"column:start;NOT NULL;comment:'开始时间'"`
	End           string    `gorm:"column:end;NOT NULL;comment:'结束时间'"`
	Location      string    `gorm:"column:location;default:zh-cn;NOT NULL;comment:'位置，中国:zh-CN'"`
	ClickCount    int       `gorm:"column:click_count;default:0;NOT NULL;comment:'点击次数'"`
	DownloadCount int       `gorm:"column:download_count;default:0;NOT NULL;comment:'下载次数'"`
	Hash          string    `gorm:"column:hash;NOT NULL;comment:'唯一hash值'"`
	CreatedAt     time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP;comment:'创建时间'"`
	UpdatedAt     time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP;comment:'更新时间'"`
}

func (b *BingImages) TableName() string {
	return "bing_images"
}

type BingImagesRepository struct {
	db *gorm.DB
}

var (
	bingImagesRepository *BingImagesRepository
	bingImagesRepoOnce   sync.Once
)

func NewBingImagesRepository() *BingImagesRepository {
	bingImagesRepoOnce.Do(func() {
		conn, err := iMysql.Conn("default")
		if err != nil {
			slog.With("err", err).Error("bingImagesRepository conn error")
			return
		}
		bingImagesRepository = new(BingImagesRepository)
		bingImagesRepository.db = conn
	})
	return bingImagesRepository
}

func (receiver BingImagesRepository) GetAll() ([]*BingImages, error) {
	var list []*BingImages
	result := receiver.db.Model(&BingImages{}).Limit(10).Find(&list)
	return list, result.Error

}
