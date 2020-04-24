package model

import "time"

type GoadminUsers struct {
	Avatar        string    `gorm:"column:avatar" json:"avatar"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at"`
	ID            int       `gorm:"column:id;primary_key" json:"id"`
	Name          string    `gorm:"column:name" json:"name"`
	Password      string    `gorm:"column:password" json:"password"`
	RememberToken string    `gorm:"column:remember_token" json:"remember_token"`
	UpdatedAt     time.Time `gorm:"column:updated_at" json:"updated_at"`
	Username      string    `gorm:"column:username" json:"username"`
}

// TableName sets the insert table name for this struct type
func (g *GoadminUsers) TableName() string {
	return "goadmin_users"
}
