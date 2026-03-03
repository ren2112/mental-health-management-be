package models

import (
	"gorm.io/gorm"
	"time"
)

type Teacher struct {
	ID           int            `gorm:"primaryKey;autoIncrement;comment:老师ID"`
	Name         string         `gorm:"type:varchar(50);not null;default:'';comment:姓名"`
	Sex          int8           `gorm:"type:tinyint(1);not null;default:0;comment:性别（0女 1男）"`
	Introduction string         `gorm:"type:text;not null;comment:简介"`
	Phone        string         `gorm:"type:varchar(20);not null;default:'';comment:联系电话"`
	Workspace    string         `gorm:"type:varchar(100);not null;default:'';comment:办公室"`
	Email        string         `gorm:"type:varchar(100);not null;default:'';uniqueIndex:uk_email;comment:邮箱"`
	Password     string         `gorm:"type:varchar(255);not null;default:'';comment:密码（加密存储）"`
	CreateTime   time.Time      `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:新建时间"`
	DeletedAt    gorm.DeletedAt // 软删除字段
}
