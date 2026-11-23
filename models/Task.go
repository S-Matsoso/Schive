package models

import "time"

type Task struct {
	ID        string     `gorm:"primaryKey" json:"id"`
	ModuleID  string     `gorm:"index;not null" json:"moduleId"`
	Title     string     `gorm:"not null" json:"title"`
	Notes     string     `json:"notes"`
	DueDate   *time.Time `json:"dueDate,omitempty"`
	Done      bool       `gorm:"default:false" json:"done"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`

	Module *Module `gorm:"foreignKey:ModuleID" json:"module,omitempty"`
}
