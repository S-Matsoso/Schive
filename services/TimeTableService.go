package services

import (
	"errors"
	"proj1/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TimetableService struct {
	db *gorm.DB
}

func NewTimetableService(db *gorm.DB) *TimetableService {
	return &TimetableService{db: db}
}

func (s *TimetableService) GetAllClassSlots() ([]models.ClassSlot, error) {
	var slots []models.ClassSlot
	err := s.db.Preload("Module").Find(&slots).Error
	return slots, err
}

func (s *TimetableService) CreateClassSlot(slot *models.ClassSlot) (*models.ClassSlot, error) {
	if slot.ModuleID == "" {
		return nil, errors.New("moduleId is required")
	}
	if slot.ID == "" {
		slot.ID = uuid.NewString()
	}
	if slot.DurationMins <= 0 {
		if slot.LengthHrs > 0 {
			slot.DurationMins = slot.LengthHrs * 60
		} else {
			slot.DurationMins = 60
		}
	}
	if slot.LengthHrs <= 0 {
		slot.LengthHrs = (slot.DurationMins + 59) / 60
	}
	err := s.db.Create(slot).Error
	if err != nil {
		return nil, err
	}
	_ = s.db.Preload("Module").First(slot, "id = ?", slot.ID).Error
	return slot, nil
}

func (s *TimetableService) DeleteClassSlot(id string) error {
	return s.db.Delete(&models.ClassSlot{}, "id = ?", id).Error
}
