package entity

type Container struct {
	ID     string `json:"id" gorm:"primary_key"`
	UserId string `json:"user_id"`
	Name   string `json:"name"`
	Status string `json:"status"`
	Config
}
