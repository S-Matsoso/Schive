package services

import (
	"errors"
	"strings"
	"time"

	"proj1/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TaskService struct {
	DB *gorm.DB
}

func NewTaskService(db *gorm.DB) *TaskService {
	return &TaskService{DB: db}
}

func parseTaskDueDate(dueDateStr string) *time.Time {
	clean := strings.TrimSpace(dueDateStr)
	if clean == "" {
		return nil
	}
	if parsed, err := time.Parse(time.RFC3339, clean); err == nil {
		return &parsed
	}
	if parsed, err := time.Parse("2006-01-02", clean); err == nil {
		return &parsed
	}
	return nil
}

func (s *TaskService) CreateTask(moduleID string, title string, notes string, dueDateStr string) (*models.Task, error) {
	cleanModuleID := strings.TrimSpace(moduleID)
	cleanTitle := strings.TrimSpace(title)

	if cleanModuleID == "" {
		return nil, errors.New("moduleID is required")
	}
	if cleanTitle == "" {
		return nil, errors.New("title is required")
	}

	var module models.Module
	if err := s.DB.First(&module, "id = ?", cleanModuleID).Error; err != nil {
		return nil, errors.New("module not found")
	}

	now := time.Now()
	task := models.Task{
		ID:        uuid.NewString(),
		ModuleID:  cleanModuleID,
		Title:     cleanTitle,
		Notes:     strings.TrimSpace(notes),
		DueDate:   parseTaskDueDate(dueDateStr),
		Done:      false,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.DB.Create(&task).Error; err != nil {
		return nil, err
	}

	return &task, nil
}

func (s *TaskService) GetTasksByModule(moduleID string) ([]models.Task, error) {
	cleanModuleID := strings.TrimSpace(moduleID)
	if cleanModuleID == "" {
		return nil, errors.New("moduleID is required")
	}

	var tasks []models.Task
	err := s.DB.Where("module_id = ?", cleanModuleID).
		Order("done asc, due_date asc, created_at desc").
		Find(&tasks).Error
	return tasks, err
}

func (s *TaskService) UpdateTask(id string, title string, notes string, dueDateStr string) (*models.Task, error) {
	cleanID := strings.TrimSpace(id)
	cleanTitle := strings.TrimSpace(title)
	if cleanID == "" || cleanTitle == "" {
		return nil, errors.New("id and title are required")
	}

	var task models.Task
	if err := s.DB.First(&task, "id = ?", cleanID).Error; err != nil {
		return nil, err
	}

	task.Title = cleanTitle
	task.Notes = strings.TrimSpace(notes)
	task.DueDate = parseTaskDueDate(dueDateStr)
	task.UpdatedAt = time.Now()

	if err := s.DB.Save(&task).Error; err != nil {
		return nil, err
	}

	return &task, nil
}

func (s *TaskService) SetTaskDone(id string, done bool) (*models.Task, error) {
	cleanID := strings.TrimSpace(id)
	if cleanID == "" {
		return nil, errors.New("id is required")
	}

	var task models.Task
	if err := s.DB.First(&task, "id = ?", cleanID).Error; err != nil {
		return nil, err
	}

	task.Done = done
	task.UpdatedAt = time.Now()

	if err := s.DB.Save(&task).Error; err != nil {
		return nil, err
	}

	return &task, nil
}

func (s *TaskService) DeleteTask(id string) error {
	cleanID := strings.TrimSpace(id)
	if cleanID == "" {
		return errors.New("id is required")
	}

	return s.DB.Delete(&models.Task{}, "id = ?", cleanID).Error
}
