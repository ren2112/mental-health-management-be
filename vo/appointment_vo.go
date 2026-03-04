package vo

import "time"

type AppointmentVO struct {
	Title       string    `json:"title"`
	Detail      string    `json:"detail"`
	TeacherID   int       `json:"teacherID"`
	TeacherName string    `json:"teacherName"`
	Status      int8      `json:"status"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
}
