package entity

import (
	"database/sql"
	"time"
)

type Container struct {
	ID        string         `json:"id" gorm:"primary_key"`
	UserId    string         `json:"user_id" gorm:"index:containers_user_id_foreign"`
	ServerID  string         `json:"server_id" gorm:"index:containers_server_id_foreign"`
	Name      string         `json:"name"`
	Status    string         `json:"status"`
	Config    string         `json:"config"`
	IPAddress string         `json:"ip_address"`
	Domain    sql.NullString `json:"domain"`
	IsActive  uint8          `json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}
