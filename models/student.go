package models

import (
	"gorm.io/gorm"
	"time"
)

type Student struct {
	ID         int            `gorm:"primaryKey;autoIncrement;comment:学生ID"`
	Name       string         `gorm:"type:varchar(50);not null;comment:姓名"`
	Sex        int8           `gorm:"type:tinyint(1);not null;comment:性别（0女 1男）"`
	StudentNo  string         `gorm:"type:varchar(20);not null;uniqueIndex:uk_student_no;comment:学号"`
	Email      string         `gorm:"type:varchar(100);not null;uniqueIndex:uk_email;comment:邮箱"`
	Password   string         `gorm:"type:varchar(255);not null;comment:密码（加密存储）"`
	CreateTime time.Time      `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:新建时间"`
	DeletedAt  gorm.DeletedAt // 软删除字段
}
