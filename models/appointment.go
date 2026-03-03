package models

import "time"

type Appointment struct {
	ID         int       `gorm:"primaryKey;autoIncrement;comment:预约ID"`
	Title      string    `gorm:"type:varchar(100);not null;default:'';comment:预约标题"`
	Detail     string    `gorm:"type:text;not null;comment:预约详情"`
	StudentID  int       `gorm:"type:int(11);not null;default:0;comment:发起预约学生ID"`
	TeacherID  int       `gorm:"type:int(11);not null;default:0;comment:预约老师ID"`
	Status     int8      `gorm:"type:tinyint(1);not null;default:0;comment:预约状态（0待审核 1通过 2拒绝）"`
	StartTime  time.Time `gorm:"type:datetime;not null;comment:开始时间"`
	EndTime    time.Time `gorm:"type:datetime;not null;comment:结束时间"`
	CreateTime time.Time `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:新建时间"`

	// 可选关联（推荐加，毕业设计加分）
	Student Student `gorm:"foreignKey:StudentID"`
	Teacher Teacher `gorm:"foreignKey:TeacherID"`
}
