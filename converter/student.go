package converter

import (
	"mental-health-management-be/models"
	"mental-health-management-be/vo"
)

// 单个转换
func ToStudentVO(s models.Student) vo.StudentVO {
	return vo.StudentVO{
		StudentID: s.ID,
		Name:      s.Name,
		Sex:       s.Sex,
		StudentNo: s.StudentNo,
		Email:     s.Email,
	}
}

// 列表转换 ⭐⭐⭐（最常用）
func ToStudentVOList(list []models.Student) []vo.StudentVO {
	res := make([]vo.StudentVO, 0, len(list))

	for _, s := range list {
		res = append(res, ToStudentVO(s))
	}

	return res
}
