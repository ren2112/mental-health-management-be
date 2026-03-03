package models

import (
	"gorm.io/gorm"
	"time"
)

type Manager struct {
	ID         int            `gorm:"primaryKey;autoIncrement;comment:管理员ID"`
	Name       string         `gorm:"type:varchar(50);not null;default:'';comment:姓名"`
	Email      string         `gorm:"type:varchar(100);not null;default:'';uniqueIndex:uk_email;comment:邮箱"`
	Phone      string         `gorm:"type:varchar(20);not null;default:'';comment:电话"`
	Password   string         `gorm:"type:varchar(255);not null;default:'';comment:密码（加密存储）"`
	CreateTime time.Time      `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:新建时间"`
	DeletedAt  gorm.DeletedAt // 软删除字段
}
