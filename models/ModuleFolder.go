package models

import "time"

type ModuleFolder struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	CreatedAt time.Time `json:"createdAt"`

	Modules []Module `gorm:"foreignKey:FolderID" json:"modules,omitempty"`
}
