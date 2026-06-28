package models

type File struct {
	ID           string  `gorm:"primaryKey" json:"id"`
	ModuleID     string  `gorm:"index" json:"moduleId"`
	AssignmentID *string `gorm:"index" json:"assignmentId,omitempty"`
	Name         string  `json:"name"`
	Path         string  `json:"path"`
	FileType     string  `json:"fileType"` // "pdf" | "docx" etc.
}
