package models

import (
	"gorm.io/gorm"
	"time"
)

type Post struct {
	ID         int            `gorm:"primaryKey;autoIncrement;comment:文章ID"`
	Title      string         `gorm:"type:varchar(200);not null;default:'';comment:标题"`
	Cover      string         `gorm:"type:varchar(255);not null;default:'';comment:封面图片地址"`
	Content    string         `gorm:"type:text;not null;comment:正文"`
	AuthorID   int            `gorm:"type:int(11);not null;default:0;comment:作者ID（老师）"`
	CreateTime time.Time      `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:新建时间"`
	UpdateTime time.Time      `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;autoUpdateTime;comment:编辑时间"`
	DeletedAt  gorm.DeletedAt // 软删除字段
	Author     Teacher        `gorm:"foreignKey:AuthorID;references:ID"`
}
