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
	MusicPTypeSrc         = 1
	MusicPTypeDst         = 2
	MusicStatusIncomplete = 1
	MusicStatusComplete   = 2
)

// Music 网易云转换音乐表
type Music struct {
	Id         int             `gorm:"column:id;type:bigint(20) unsigned;AUTO_INCREMENT;comment:自增id;primary_key" json:"id"`
	AccountId  int             `gorm:"column:account_id;type:bigint(20);default:0;comment:账号id;NOT NULL" json:"account_id"`
	Flag       string          `gorm:"column:flag;type:varchar(32);comment:文件标识;NOT NULL" json:"flag"`
	MusicId    int             `gorm:"column:music_id;type:bigint(20);default:0;comment:音乐id;NOT NULL" json:"music_id"`
	MusicName  string          `gorm:"column:music_name;type:varchar(255);comment:音乐名称;NOT NULL" json:"music_name"`
	Artist     [][]interface{} `gorm:"serializer:json;column:artist;type:varchar(128);comment:艺术家;NOT NULL" json:"artist"`
	AlbumPic   string          `gorm:"column:album_pic;type:varchar(500);comment:专辑图片;NOT NULL" json:"album_pic"`
	MusicDocId string          `gorm:"column:music_doc_id;type:varchar(128);comment:音乐文件id;NOT NULL" json:"music_doc_id"`
	Duration   int             `gorm:"column:duration;type:int(11);default:0;comment:音乐时长;NOT NULL" json:"duration"`
	MvId       int             `gorm:"column:mv_id;type:int(11);default:0;comment:视频id;NOT NULL" json:"mv_id"`
	Format     string          `gorm:"column:format;type:varchar(64);comment:音乐格式;NOT NULL" json:"format"`
	SrcPath    string          `gorm:"column:src_path;type:varchar(255);comment:文件原始路径;NOT NULL" json:"src_path"`
	DstPath    string          `gorm:"column:dst_path;type:varchar(255);comment:文件上传后的路径;NOT NULL" json:"dst_path"`
	ParsePath  string          `gorm:"column:parse_path;type:varchar(255);comment:文件解析后的路径;NOT NULL" json:"parse_path"`
	PType      int             `gorm:"column:p_type;type:int(11);default:0;comment:地址选择类型，1 - 原文件地址，2 - 自定义文件地址;NOT NULL" json:"p_type"`
	Status     int             `gorm:"column:status;type:int(11);default:0;comment:状态，1-未处理，2-已处理;NOT NULL" json:"status"`
	CreatedAt  time.Time       `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	UpdatedAt  time.Time       `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`
}

func (m *Music) TableName() string {
	return "music"
}

type MusicRepository struct {
	db *gorm.DB
}

var (
	musicRepository *MusicRepository
	musicRepoOnce   sync.Once
)

func NewMusicRepository() *MusicRepository {
	musicRepoOnce.Do(func() {
		conn, err := iMysql.Conn("default")
		if err != nil {
			slog.With("err", err).Error("musicRepository conn error")
			return
		}
		musicRepository = new(MusicRepository)
		musicRepository.db = conn
	})
	return musicRepository
}

func (receiver *MusicRepository) CreateByModel(model *Music) error {
	if model == nil || reflect.DeepEqual(model, &Music{}) {
		return errors.New("model is nil")
	}
	result := receiver.db.Model(&Music{}).Create(model)
	return result.Error
}

func (receiver *MusicRepository) BulkCreate(model []*Music) error {
	if model == nil || reflect.DeepEqual(model, []*Music{}) {
		return errors.New("model is nil")
	}
	return receiver.db.Model(&Music{}).CreateInBatches(model, 100).Error
}

func (receiver *MusicRepository) StatByFlag(flag string) (*Music, error) {
	var record Music
	result := receiver.db.Model(&Music{}).Where("flag = ?", flag).First(&record)
	return &record, result.Error
}

func (receiver *MusicRepository) UpdateByModel(id int, model *Music) error {
	if id <= 0 {
		return errors.New("id is invalid")
	}

	if model == nil || reflect.DeepEqual(model, &Music{}) {
		return errors.New("model is nil")
	}
	result := receiver.db.Model(&Music{}).Where("id = ?", id).Updates(model)
	return result.Error
}
