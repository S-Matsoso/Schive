package services

import "gorm.io/gorm"

// "errors"

// "proj1/models"

// "github.com/google/uuid"

type TimetableService struct {
	db *gorm.DB
}

func NewTimetableService(db *gorm.DB) *TimetableService {
	return &TimetableService{db: db}
}
