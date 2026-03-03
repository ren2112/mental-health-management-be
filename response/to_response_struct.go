package response

//
//import (
//	"partyus_app/models"
//	"time"
//)
//
//type RespUser struct {
//	ID         int    `json:"id"`
//	Username   string `json:"username"`
//	Avatar     string `json:"avatar"`
//	Email      string `json:"email"` // 邮箱
//	Phone      string `json:"phone"`
//	Department int    `json:"department"` //0表示区团委，1表示社区团组织，2表示高校团组织，3表示企业团组织
//}
//
//func ToResponseUser(u models.User) RespUser {
//	res := RespUser{
//		u.ID,
//		u.Username,
//		u.Avatar,
//		u.Email,
//		u.Phone,
//		u.Department,
//	}
//	return res
//}
//
//type RespPost struct {
//	ID        int        `json:"id" `
//	UserID    int        `json:"user_id" `
//	User      RespUser   `json:"user" `
//	Type      int        `json:"type" `       // 文件类型: 0 图文帖子, 1 视频帖子
//	Title     string     `json:"title" `      // 标题
//	Content   string     `json:"content" `    // 内容
//	Cover     string     `json:"cover" `      // 封面地址
//	Video     string     `json:"video" `      // 视频地址
//	Part      int        `json:"part" `       // 分类: 0 理论学习, 1 走进高新等
//	IsAudit   int        `json:"is_audit"`    // 审核状态: 0 未审核, 1 审核通过, 2 审核不通过
//	CreatedAt time.Time  `json:"created_at" ` // 创建时间
//	DeletedAt *time.Time `json:"deleted_at"`  // 删除时间
//}
//
//func ToResponsePost(p models.Post) RespPost {
//	return RespPost{
//		ID:        p.ID,
//		UserID:    p.UserID,
//		User:      ToResponseUser(p.User), // 调用 ToResponseUser 转换用户信息
//		Type:      p.Type,
//		Title:     p.Title,
//		Content:   p.Content,
//		Cover:     p.Cover,
//		Video:     p.Video,
//		Part:      p.Part,
//		IsAudit:   p.IsAudit,
//		CreatedAt: p.CreatedAt,
//		DeletedAt: func() *time.Time {
//			if p.DeletedAt.Valid {
//				return &p.DeletedAt.Time
//			}
//			return nil
//		}(),
//	}
//}
