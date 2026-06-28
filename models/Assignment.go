package models

import "time"

type Assignment struct {
	ID               string     `gorm:"primaryKey" json:"id"`
	ModuleID         string     `gorm:"index; not null" json:"moduleId"`
	AssignmentNumber int        `json:"assignmentNumber"`
	DueDate          *time.Time `json:"dueDate,omitempty"`

	Files []File `gorm:"foreignKey:AssignmentID;constraint:OnDelete:CASCADE" json:"files,omitempty"`
}
