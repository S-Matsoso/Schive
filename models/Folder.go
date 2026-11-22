package models

import "time"

type Folder struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	ModuleID  string    `gorm:"index;not null" json:"moduleId"`
	Name      string    `gorm:"not null" json:"name"`
	Color     string    `json:"color,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	FileCount int       `gorm:"-" json:"fileCount"`
	Files     []File    `gorm:"foreignKey:FolderID;constraint:OnDelete:CASCADE" json:"files,omitempty"`
}
