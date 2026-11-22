package models

type Module struct {
	ID    string `gorm:"primaryKey" json:"id"`
	Code  string `gorm:"uniqueIndex; not null" json:"code"`
	Name  string `gorm:"not null" json:"name"`
	Color string `json:"color"`

	Files       []File       `gorm:"foreignKey:ModuleID;constraint:OnDelete:CASCADE" json:"files,omitempty"`
	Folders     []Folder     `gorm:"foreignKey:ModuleID;constraint:OnDelete:CASCADE" json:"folders,omitempty"`
	ClassSlots  []ClassSlot  `gorm:"foreignKey:ModuleID;constraint:OnDelete:CASCADE" json:"classSlots,omitempty"`
	Assignments []Assignment `gorm:"foreignKey:ModuleID;constraint:OnDelete:CASCADE" json:"assignments,omitempty"`
	Flashcards  []Flashcard  `gorm:"foreignKey:ModuleID;constraint:OnDelete:CASCADE" json:"flashcards,omitempty"`
}
