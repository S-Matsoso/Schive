package models

type ClassSlot struct {
	ID        string `gorm:"primaryKey" json:"id"`
	ModuleID  string `gorm:"index" json:"moduleId"`
	Day       string `json:"day"` // "mon", "tue", "wed", "thu", "fri"
	StartHour int    `json:"startHour"`
	LengthHrs int    `json:"lengthHrs"`
	Room      string `json:"room"`
}
