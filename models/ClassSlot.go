package models

type ClassSlot struct {
	ID           string  `gorm:"primaryKey" json:"id"`
	ModuleID     string  `gorm:"index;not null" json:"moduleId"`
	Day          string  `gorm:"not null" json:"day"` // "mon", "tue", "wed", "thu", "fri"
	StartHour    int     `json:"startHour"`
	LengthHrs    int     `json:"lengthHrs"`
	DurationMins int     `json:"durationMins"` // 30, 40, 50, 60, 120
	Room         string  `json:"room"`
	Module       *Module `gorm:"foreignKey:ModuleID;references:ID" json:"module,omitempty"`
}
