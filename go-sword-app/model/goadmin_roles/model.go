package goadmin_roles

import "time"

type GoadminRoles struct {
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	ID        int       `gorm:"column:id;primary_key" json:"id"`
	Name      string    `gorm:"column:name" json:"name"`
	Slug      string    `gorm:"column:slug" json:"slug"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName sets the insert table name for this struct type
func (g *GoadminRoles) TableName() string {
	return "goadmin_roles"
}
