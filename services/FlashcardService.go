package services

import (
	"errors"
	"strings"
	"time"

	"proj1/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FlashcardService struct {
	DB *gorm.DB
}

func NewFlashcardService(db *gorm.DB) *FlashcardService {
	return &FlashcardService{DB: db}
}

func (s *FlashcardService) CreateFlashcard(moduleID string, question string, answer string) (*models.Flashcard, error) {
	cleanModuleID := strings.TrimSpace(moduleID)
	cleanQ := strings.TrimSpace(question)
	cleanA := strings.TrimSpace(answer)

	if cleanModuleID == "" {
		return nil, errors.New("moduleID is required")
	}
	if cleanQ == "" || cleanA == "" {
		return nil, errors.New("question and answer are required")
	}

	now := time.Now()
	flashcard := models.Flashcard{
		ID:        uuid.NewString(),
		ModuleID:  cleanModuleID,
		Question:  cleanQ,
		Answer:    cleanA,
		Mastered:  false,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.DB.Create(&flashcard).Error; err != nil {
		return nil, err
	}

	return &flashcard, nil
}

func (s *FlashcardService) GetFlashcardsByModule(moduleID string) ([]models.Flashcard, error) {
	cleanModuleID := strings.TrimSpace(moduleID)
	var flashcards []models.Flashcard
	err := s.DB.Where("module_id = ?", cleanModuleID).Order("created_at desc").Find(&flashcards).Error
	return flashcards, err
}

func (s *FlashcardService) ToggleMastered(id string) (*models.Flashcard, error) {
	cleanID := strings.TrimSpace(id)
	if cleanID == "" {
		return nil, errors.New("id is required")
	}

	var flashcard models.Flashcard
	if err := s.DB.First(&flashcard, "id = ?", cleanID).Error; err != nil {
		return nil, err
	}

	flashcard.Mastered = !flashcard.Mastered
	flashcard.UpdatedAt = time.Now()

	if err := s.DB.Save(&flashcard).Error; err != nil {
		return nil, err
	}

	return &flashcard, nil
}

func (s *FlashcardService) UpdateFlashcard(id string, question string, answer string) (*models.Flashcard, error) {
	cleanID := strings.TrimSpace(id)
	cleanQ := strings.TrimSpace(question)
	cleanA := strings.TrimSpace(answer)

	if cleanID == "" || cleanQ == "" || cleanA == "" {
		return nil, errors.New("id, question, and answer are required")
	}

	var flashcard models.Flashcard
	if err := s.DB.First(&flashcard, "id = ?", cleanID).Error; err != nil {
		return nil, err
	}

	flashcard.Question = cleanQ
	flashcard.Answer = cleanA
	flashcard.UpdatedAt = time.Now()

	if err := s.DB.Save(&flashcard).Error; err != nil {
		return nil, err
	}

	return &flashcard, nil
}

func (s *FlashcardService) DeleteFlashcard(id string) error {
	cleanID := strings.TrimSpace(id)
	if cleanID == "" {
		return errors.New("id is required")
	}

	return s.DB.Delete(&models.Flashcard{}, "id = ?", cleanID).Error
}
