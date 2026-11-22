package models

import "time"

type Flashcard struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	ModuleID  string    `gorm:"index;not null" json:"moduleId"`
	Question  string    `gorm:"not null" json:"question"`
	Answer    string    `gorm:"not null" json:"answer"`
	Mastered  bool      `gorm:"default:false" json:"mastered"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
