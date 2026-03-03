package vo

type ManagerVO struct {
	ManagerID int    `json:"managerID"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}
