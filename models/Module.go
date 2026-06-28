package models

type Module struct {
	ID    string `gorm:"primaryKey" json:"id"`
	Code  string `gorm:"uniqueIndex" json:"code"`
	Name  string `json:"name"`
	Color string `json:"color"`

	Files       []File       `gorm:"foreignKey:ModuleID;constraint:OnDelete:CASCADE" json:"files,omitempty"`
	ClassSlots  []ClassSlot  `gorm:"foreignKey:ModuleID;constraint:OnDelete:CASCADE" json:"classSlots,omitempty"`
	Assignments []Assignment `gorm:"foreignKey:ModuleID;constraint:OnDelete:CASCADE" json:"assigments,omitempty"`
}
