package goadmin_permissions

import "time"

type GoadminPermissions struct {
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
	HTTPMethod string    `gorm:"column:http_method" json:"http_method"`
	HTTPPath   string    `gorm:"column:http_path" json:"http_path"`
	ID         int       `gorm:"column:id;primary_key" json:"id"`
	Name       string    `gorm:"column:name" json:"name"`
	Slug       string    `gorm:"column:slug" json:"slug"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName sets the insert table name for this struct type
func (g *GoadminPermissions) TableName() string {
	return "goadmin_permissions"
}
