package converter

import (
	"mental-health-management-be/models"
	"mental-health-management-be/vo"
)

func ToAppointmentVO(a models.Appointment) vo.AppointmentVO {
	return vo.AppointmentVO{
		Title:       a.Title,
		Detail:      a.Detail,
		TeacherID:   a.TeacherID,
		TeacherName: a.Teacher.Name,
		Status:      a.Status,
		StartTime:   a.StartTime,
		EndTime:     a.EndTime,
	}
}

// ===== 列表转换 =====
func ToAppointmentVOList(list []models.Appointment) []vo.AppointmentVO {
	res := make([]vo.AppointmentVO, 0, len(list))

	for _, a := range list {
		res = append(res, ToAppointmentVO(a))
	}

	return res
}
