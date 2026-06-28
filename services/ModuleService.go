package services

import (
	"errors"

	"proj1/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ModuleService struct {
	DB *gorm.DB
}

func NewModuleService(db *gorm.DB) *ModuleService {
	return &ModuleService{DB: db}
}

func (s *ModuleService) GetAllModules() ([]models.Module, error) {
	var modules []models.Module
	err := s.DB.Preload("Files").Preload("ClassSlots").
		Preload("Assignments").Find(&modules).Error
	return modules, err
}

func (s *ModuleService) GetModule(code string) (*models.Module, error) {
	var module models.Module
	err := s.DB.Preload("Files").Preload("ClassSlots").
		Preload("Assignments").First(&models.Module{}, code).Error
	return &module, err
}

func (s *ModuleService) CreateModule(newModule *models.Module) (*models.Module, error) {
	existing, err := s.GetAllModules()
	if err != nil {
		return &models.Module{}, err
	}
	for _, m := range existing {
		if m.Code == newModule.Code {
			return &models.Module{}, errors.New("module with this code already exists")
		}
	}

	module := models.Module{
		ID:    uuid.NewString(),
		Code:  newModule.Code,
		Name:  newModule.Name,
		Color: newModule.Color,
	}

	res := s.DB.Create(module)

	return &module, res.Error
}

func (s *ModuleService) DeleteModule(code string) error {
	return s.DB.Delete(&models.Module{}, code).Error
}
