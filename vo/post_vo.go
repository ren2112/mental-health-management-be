package vo

type PostVO struct {
	PostID     int    `json:"postID"`
	Title      string `json:"title"`
	AuthorID   int    `json:"authorID"`
	AuthorName string `json:"authorName"`
	Cover      string `json:"cover"`
	Content    string `json:"content"`
	CreateTime int64  `json:"createTime"`
	EditTime   int64  `json:"editTime"`
}
