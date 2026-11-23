package services

import (
	"errors"
	"strings"
	"time"

	"proj1/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AssignmentService struct {
	DB *gorm.DB
}

func NewAssignmentService(db *gorm.DB) *AssignmentService {
	return &AssignmentService{DB: db}
}

func (s *AssignmentService) CreateAssignment(moduleID string, title string, description string, assignmentNumber int, status string, priority string, dueDateStr string) (*models.Assignment, error) {
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

	cleanStatus := strings.ToLower(strings.TrimSpace(status))
	if cleanStatus == "" {
		cleanStatus = "pending"
	}

	cleanPriority := strings.ToLower(strings.TrimSpace(priority))
	if cleanPriority == "" {
		cleanPriority = "medium"
	}

	var dueDatePtr *time.Time
	if strings.TrimSpace(dueDateStr) != "" {
		parsedDate, err := time.Parse(time.RFC3339, strings.TrimSpace(dueDateStr))
		if err == nil {
			dueDatePtr = &parsedDate
		} else {
			// fallback YYYY-MM-DD
			if parsedDate2, err2 := time.Parse("2006-01-02", strings.TrimSpace(dueDateStr)); err2 == nil {
				dueDatePtr = &parsedDate2
			}
		}
	}

	progress := 0
	if cleanStatus == "completed" {
		progress = 100
	}

	now := time.Now()
	assignment := models.Assignment{
		ID:               uuid.NewString(),
		ModuleID:         cleanModuleID,
		Title:            cleanTitle,
		Description:      strings.TrimSpace(description),
		AssignmentNumber: assignmentNumber,
		Status:           cleanStatus,
		Priority:         cleanPriority,
		Progress:         progress,
		DueDate:          dueDatePtr,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	if err := s.DB.Create(&assignment).Error; err != nil {
		return nil, err
	}

	return &assignment, nil
}

func (s *AssignmentService) GetAssignmentsByModule(moduleID string) ([]models.Assignment, error) {
	cleanModuleID := strings.TrimSpace(moduleID)
	var assignments []models.Assignment
	err := s.DB.Preload("Module").Preload("Files").Where("module_id = ?", cleanModuleID).Order("due_date asc, created_at desc").Find(&assignments).Error
	return assignments, err
}

func (s *AssignmentService) GetAllAssignments() ([]models.Assignment, error) {
	var assignments []models.Assignment
	err := s.DB.Preload("Module").Preload("Files").Order("due_date asc, created_at desc").Find(&assignments).Error
	return assignments, err
}

func (s *AssignmentService) UpdateAssignment(id string, title string, description string, status string, priority string, dueDateStr string) (*models.Assignment, error) {
	cleanID := strings.TrimSpace(id)
	cleanTitle := strings.TrimSpace(title)

	if cleanID == "" || cleanTitle == "" {
		return nil, errors.New("id and title are required")
	}

	var assignment models.Assignment
	if err := s.DB.First(&assignment, "id = ?", cleanID).Error; err != nil {
		return nil, err
	}

	assignment.Title = cleanTitle
	assignment.Description = strings.TrimSpace(description)
	if strings.TrimSpace(status) != "" {
		assignment.Status = strings.ToLower(strings.TrimSpace(status))
		if assignment.Status == "completed" {
			assignment.Progress = 100
		} else if assignment.Status == "pending" {
			assignment.Progress = 0
		}
	}
	if strings.TrimSpace(priority) != "" {
		assignment.Priority = strings.ToLower(strings.TrimSpace(priority))
	}

	if strings.TrimSpace(dueDateStr) != "" {
		parsedDate, err := time.Parse(time.RFC3339, strings.TrimSpace(dueDateStr))
		if err == nil {
			assignment.DueDate = &parsedDate
		} else if parsedDate2, err2 := time.Parse("2006-01-02", strings.TrimSpace(dueDateStr)); err2 == nil {
			assignment.DueDate = &parsedDate2
		}
	} else if dueDateStr == "" {
		assignment.DueDate = nil
	}

	assignment.UpdatedAt = time.Now()

	if err := s.DB.Save(&assignment).Error; err != nil {
		return nil, err
	}

	return &assignment, nil
}

func (s *AssignmentService) UpdateAssignmentStatus(id string, status string) error {
	cleanID := strings.TrimSpace(id)
	cleanStatus := strings.ToLower(strings.TrimSpace(status))
	if cleanID == "" || cleanStatus == "" {
		return errors.New("id and status are required")
	}

	updates := map[string]interface{}{
		"status":     cleanStatus,
		"updated_at": time.Now(),
	}
	switch cleanStatus {
	case "completed":
		updates["progress"] = 100
	case "pending":
		updates["progress"] = 0
	}

	return s.DB.Model(&models.Assignment{}).Where("id = ?", cleanID).Updates(updates).Error
}

func (s *AssignmentService) UpdateAssignmentProgress(id string, progress int) (*models.Assignment, error) {
	cleanID := strings.TrimSpace(id)
	if cleanID == "" {
		return nil, errors.New("id is required")
	}

	if progress < 0 {
		progress = 0
	}
	if progress > 100 {
		progress = 100
	}

	var assignment models.Assignment
	if err := s.DB.First(&assignment, "id = ?", cleanID).Error; err != nil {
		return nil, err
	}

	assignment.Progress = progress
	switch {
	case progress >= 100:
		assignment.Status = "completed"
	case progress > 0:
		assignment.Status = "in_progress"
	default:
		assignment.Status = "pending"
	}
	assignment.UpdatedAt = time.Now()

	if err := s.DB.Save(&assignment).Error; err != nil {
		return nil, err
	}

	return &assignment, nil
}

func (s *AssignmentService) DeleteAssignment(id string) error {
	cleanID := strings.TrimSpace(id)
	if cleanID == "" {
		return errors.New("id is required")
	}

	return s.DB.Delete(&models.Assignment{}, "id = ?", cleanID).Error
}
