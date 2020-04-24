package goadmin_site

import "time"

type GoadminSite struct {
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	Description string    `gorm:"column:description" json:"description"`
	ID          int       `gorm:"column:id;primary_key" json:"id"`
	Key         string    `gorm:"column:key" json:"key"`
	State       int       `gorm:"column:state" json:"state"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
	Value       string    `gorm:"column:value" json:"value"`
}

// TableName sets the insert table name for this struct type
func (g *GoadminSite) TableName() string {
	return "goadmin_site"
}
