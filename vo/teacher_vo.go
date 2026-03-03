package vo

type TeacherVO struct {
	TeacherID    int    `json:"teacherID"`
	Name         string `json:"name"`
	Sex          int8   `json:"sex"`
	Introduction string `json:"introduction"`
	Email        string `json:"email"`
	Workspace    string `json:"workspace"`
	Phone        string `json:"phone"`
}
