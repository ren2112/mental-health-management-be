package models

import (
	"time"
)

type AppointmentSlot struct {
	ID            int       `gorm:"primaryKey"`
	TeacherID     int       `gorm:"not null;uniqueIndex:idx_slot_unique"`
	SlotTime      time.Time `gorm:"not null;uniqueIndex:idx_slot_unique"`
	AppointmentID int       `gorm:"not null;index"`

	CreateTime time.Time `gorm:"autoCreateTime"`
}
