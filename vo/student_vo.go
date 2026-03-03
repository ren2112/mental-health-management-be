package vo

type StudentVO struct {
	StudentID   int    `json:"studentID"`
	StudentName string `json:"studentName"`
	Sex         int8   `json:"sex"`
	StudentNum  string `json:"studentNum"`
	Email       string `json:"email"`
}
