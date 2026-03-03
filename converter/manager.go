package converter

import (
	"mental-health-management-be/models"
	"mental-health-management-be/vo"
)

// 单个转换
func ToManagerVO(m models.Manager) vo.ManagerVO {
	return vo.ManagerVO{
		ManagerID: m.ID,
		Name:      m.Name,
		Email:     m.Email,
		Phone:     m.Phone,
	}
}

// 列表转换
func ToManagerVOList(list []models.Manager) []vo.ManagerVO {
	res := make([]vo.ManagerVO, 0, len(list))

	for _, m := range list {
		res = append(res, ToManagerVO(m))
	}

	return res
}
