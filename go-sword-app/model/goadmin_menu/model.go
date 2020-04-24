package goadmin_menu

import "time"

type GoadminMenu struct {
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	Header    string    `gorm:"column:header" json:"header"`
	Icon      string    `gorm:"column:icon" json:"icon"`
	ID        int       `gorm:"column:id;primary_key" json:"id"`
	Order     int       `gorm:"column:order" json:"order"`
	ParentID  int       `gorm:"column:parent_id" json:"parent_id"`
	Title     string    `gorm:"column:title" json:"title"`
	Type      int       `gorm:"column:type" json:"type"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
	URI       string    `gorm:"column:uri" json:"uri"`
}

// TableName sets the insert table name for this struct type
func (g *GoadminMenu) TableName() string {
	return "goadmin_menu"
}
