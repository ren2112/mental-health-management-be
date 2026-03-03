package converter

import (
	"mental-health-management-be/models"
	"mental-health-management-be/vo"
)

func ToPostVO(p models.Post) vo.PostVO {
	return vo.PostVO{
		PostID:     p.ID,
		Title:      p.Title,
		AuthorID:   p.AuthorID,
		AuthorName: p.Author.Name, // 来自 Preload("Author")
		Cover:      p.Cover,
		Content:    p.Content,
		CreateTime: p.CreateTime.Unix(),
		EditTime:   p.UpdateTime.Unix(),
	}
}

// Post 列表转换
func ToPostVOList(list []models.Post) []vo.PostVO {
	res := make([]vo.PostVO, 0, len(list))

	for _, p := range list {
		res = append(res, ToPostVO(p))
	}

	return res
}
