package goadmin_user_permissions

import "time"

type GoadminUserPermissions struct {
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
	PermissionID int       `gorm:"column:permission_id;primary_key" json:"permission_id"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`
	UserID       int       `gorm:"column:user_id;primary_key" json:"user_id"`
}

// TableName sets the insert table name for this struct type
func (g *GoadminUserPermissions) TableName() string {
	return "goadmin_user_permissions"
}
