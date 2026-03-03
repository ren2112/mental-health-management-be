package converter

import (
	"mental-health-management-be/models"
	"mental-health-management-be/vo"
)

func ToTeacherVO(t models.Teacher) vo.TeacherVO {
	return vo.TeacherVO{
		TeacherID:    t.ID,
		Name:         t.Name,
		Sex:          t.Sex,
		Introduction: t.Introduction,
		Email:        t.Email,
		Workspace:    t.Workspace,
		Phone:        t.Phone,
	}
}

func ToTeacherVOList(list []models.Teacher) []vo.TeacherVO {
	res := make([]vo.TeacherVO, 0, len(list))

	for _, t := range list {
		res = append(res, ToTeacherVO(t))
	}

	return res
}
