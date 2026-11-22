package models

import "time"

type Assignment struct {
	ID               string     `gorm:"primaryKey" json:"id"`
	ModuleID         string     `gorm:"index; not null" json:"moduleId"`
	Title            string     `gorm:"not null" json:"title"`
	Description      string     `json:"description"`
	AssignmentNumber int        `json:"assignmentNumber"`
	Status           string     `gorm:"default:'pending'" json:"status"`   // "pending" | "in_progress" | "completed"
	Priority         string     `gorm:"default:'medium'" json:"priority"` // "low" | "medium" | "high"
	Progress         int        `gorm:"default:0" json:"progress"`       // 0-100
	DueDate          *time.Time `json:"dueDate,omitempty"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        time.Time  `json:"updatedAt"`

	Module *Module `gorm:"foreignKey:ModuleID" json:"module,omitempty"`
	Files  []File  `gorm:"foreignKey:AssignmentID;constraint:OnDelete:CASCADE" json:"files,omitempty"`
}
