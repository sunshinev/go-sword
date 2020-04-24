package model

import "time"

type GoadminRoleUsers struct {
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	RoleID    int       `gorm:"column:role_id;primary_key" json:"role_id"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
	UserID    int       `gorm:"column:user_id;primary_key" json:"user_id"`
}

// TableName sets the insert table name for this struct type
func (g *GoadminRoleUsers) TableName() string {
	return "goadmin_role_users"
}
