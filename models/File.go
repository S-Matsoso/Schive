package models

import "time"

type File struct {
	ID           string    `gorm:"primaryKey" json:"id"`
	ModuleID     string    `gorm:"index;not null" json:"moduleId"`
	FolderID     *string   `gorm:"index" json:"folderId,omitempty"`
	AssignmentID *string   `gorm:"index" json:"assignmentId,omitempty"`
	Name         string    `json:"name"`
	Path         string    `json:"path"`
	FileType     string    `json:"fileType"` // "pdf" | "docx" etc.
	FileSize     int64     `json:"fileSize"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	Data         []byte    `gorm:"type:bytes" json:"-"`
}
