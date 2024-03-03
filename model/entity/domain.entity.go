package entity

type Domain struct {
	ID          string `json:"id" gorm:"primary_key"`
	Domain      string `json:"domain"`
	ContainerID string `json:"container_id"`
	Type        string `json:"type"`
	Status      string `json:"status"`
	//CreatedAt   time.Time `json:"created_at"`
	//UpdatedAt   time.Time `json:"updated_at"`
}
