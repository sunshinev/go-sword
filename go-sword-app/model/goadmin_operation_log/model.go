package goadmin_operation_log

import "time"

type GoadminOperationLog struct {
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	ID        int       `gorm:"column:id;primary_key" json:"id"`
	Input     string    `gorm:"column:input" json:"input"`
	IP        string    `gorm:"column:ip" json:"ip"`
	Method    string    `gorm:"column:method" json:"method"`
	Path      string    `gorm:"column:path" json:"path"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
	UserID    int       `gorm:"column:user_id" json:"user_id"`
}

// TableName sets the insert table name for this struct type
func (g *GoadminOperationLog) TableName() string {
	return "goadmin_operation_log"
}
