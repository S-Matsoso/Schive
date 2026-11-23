package services

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

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
	err := s.DB.Preload("Files").Preload("Folders").Preload("ClassSlots").
		Preload("Assignments").Preload("Flashcards").Find(&modules).Error
	return modules, err
}

func (s *ModuleService) GetModule(code string) (*models.Module, error) {
	cleanCode := strings.TrimSpace(code)
	if cleanCode == "" {
		return nil, errors.New("module code is required")
	}
	var module models.Module
	err := s.DB.Preload("Files").Preload("Folders").Preload("ClassSlots").
		Preload("Assignments").Preload("Flashcards").Where("code = ?", cleanCode).First(&module).Error
	if err != nil {
		return nil, err
	}
	return &module, nil
}

func (s *ModuleService) CreateModule(newModule *models.Module) (*models.Module, error) {
	if newModule == nil {
		return nil, errors.New("module data is required")
	}

	code := strings.TrimSpace(newModule.Code)
	name := strings.TrimSpace(newModule.Name)
	color := strings.TrimSpace(newModule.Color)

	if code == "" {
		return nil, errors.New("module code is required")
	}
	if name == "" {
		return nil, errors.New("module name is required")
	}
	if color == "" {
		color = "#1976D2"
	}

	existing, err := s.GetAllModules()
	if err != nil {
		return nil, err
	}
	for _, m := range existing {
		if strings.EqualFold(strings.TrimSpace(m.Code), code) {
			return nil, errors.New("module with this code already exists")
		}
	}

	folderID := newModule.FolderID
	if folderID != nil && strings.TrimSpace(*folderID) == "" {
		folderID = nil
	}

	module := models.Module{
		ID:       uuid.NewString(),
		Code:     code,
		Name:     name,
		Color:    color,
		FolderID: folderID,
	}

	res := s.DB.Create(&module)
	if res.Error != nil {
		return nil, res.Error
	}

	return &module, nil
}

func (s *ModuleService) UpdateModule(code string, updated *models.Module) (*models.Module, error) {
	cleanCode := strings.TrimSpace(code)
	if cleanCode == "" {
		return nil, errors.New("module code is required")
	}
	if updated == nil {
		return nil, errors.New("updated module data is required")
	}

	var module models.Module
	if err := s.DB.Where("code = ?", cleanCode).First(&module).Error; err != nil {
		return nil, errors.New("module not found")
	}

	newName := strings.TrimSpace(updated.Name)
	newColor := strings.TrimSpace(updated.Color)
	newCode := strings.TrimSpace(updated.Code)

	if newName == "" {
		return nil, errors.New("module name is required")
	}
	if newCode == "" {
		return nil, errors.New("module code is required")
	}

	if !strings.EqualFold(module.Code, newCode) {
		var count int64
		s.DB.Model(&models.Module{}).Where("code = ? AND id != ?", newCode, module.ID).Count(&count)
		if count > 0 {
			return nil, errors.New("another module already has this code")
		}
		module.Code = newCode
	}

	module.Name = newName
	if newColor != "" {
		module.Color = newColor
	}

	if updated.FolderID != nil && strings.TrimSpace(*updated.FolderID) == "" {
		updated.FolderID = nil
	}
	module.FolderID = updated.FolderID

	if err := s.DB.Save(&module).Error; err != nil {
		return nil, err
	}
	return &module, nil
}

func (s *ModuleService) DeleteModule(code string) error {
	cleanCode := strings.TrimSpace(code)
	if cleanCode == "" {
		return errors.New("module code is required")
	}

	var module models.Module
	if err := s.DB.Where("code = ?", cleanCode).First(&module).Error; err != nil {
		return err
	}

	// Delete uploads directory for this module
	uploadDir := filepath.Join("uploads", module.ID)
	_ = os.RemoveAll(uploadDir)

	_ = s.DB.Where("module_id = ?", module.ID).Delete(&models.File{}).Error
	_ = s.DB.Where("module_id = ?", module.ID).Delete(&models.Folder{}).Error
	_ = s.DB.Where("module_id = ?", module.ID).Delete(&models.ClassSlot{}).Error
	_ = s.DB.Where("module_id = ?", module.ID).Delete(&models.Assignment{}).Error
	_ = s.DB.Where("module_id = ?", module.ID).Delete(&models.Flashcard{}).Error

	return s.DB.Delete(&module).Error
}

// --- Module folders ---

func (s *ModuleService) GetAllModuleFolders() ([]models.ModuleFolder, error) {
	var folders []models.ModuleFolder
	err := s.DB.Order("created_at asc").Find(&folders).Error
	return folders, err
}

func (s *ModuleService) CreateModuleFolder(name string) (*models.ModuleFolder, error) {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return nil, errors.New("folder name is required")
	}

	folder := models.ModuleFolder{
		ID:   uuid.NewString(),
		Name: trimmed,
	}
	if err := s.DB.Create(&folder).Error; err != nil {
		return nil, err
	}
	return &folder, nil
}

func (s *ModuleService) UpdateModuleFolder(id string, name string) (*models.ModuleFolder, error) {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return nil, errors.New("folder name is required")
	}

	var folder models.ModuleFolder
	if err := s.DB.First(&folder, "id = ?", id).Error; err != nil {
		return nil, errors.New("folder not found")
	}
	folder.Name = trimmed
	if err := s.DB.Save(&folder).Error; err != nil {
		return nil, err
	}
	return &folder, nil
}

// DeleteModuleFolder removes the folder and leaves its modules ungrouped.
func (s *ModuleService) DeleteModuleFolder(id string) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("folder id is required")
	}
	if err := s.DB.Model(&models.Module{}).Where("folder_id = ?", id).Update("folder_id", nil).Error; err != nil {
		return err
	}
	return s.DB.Delete(&models.ModuleFolder{}, "id = ?", id).Error
}
