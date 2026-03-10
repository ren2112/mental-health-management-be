package vo

type StudentVO struct {
	StudentID int    `json:"studentID"`
	Name      string `json:"name"`
	Sex       int8   `json:"sex"`
	StudentNo string `json:"studentNo"`
	Email     string `json:"email"`
}
