package controllers

import (
	"fmt"
	"proj1/models"
	"proj1/services"
)

type FlashcardController struct {
	Service *services.FlashcardService
}

func NewFlashcardController(service *services.FlashcardService) *FlashcardController {
	return &FlashcardController{Service: service}
}

func (c *FlashcardController) CreateFlashcard(moduleID string, question string, answer string) (*models.Flashcard, error) {
	flashcard, err := c.Service.CreateFlashcard(moduleID, question, answer)
	if err != nil {
		fmt.Printf("Failed to create flashcard: %v\n", err)
		return nil, err
	}
	return flashcard, nil
}

func (c *FlashcardController) GetFlashcardsByModule(moduleID string) ([]models.Flashcard, error) {
	cards, err := c.Service.GetFlashcardsByModule(moduleID)
	if err != nil {
		fmt.Printf("Failed to get flashcards: %v\n", err)
		return nil, err
	}
	return cards, nil
}

func (c *FlashcardController) ToggleMastered(id string) (*models.Flashcard, error) {
	card, err := c.Service.ToggleMastered(id)
	if err != nil {
		fmt.Printf("Failed to toggle flashcard mastered: %v\n", err)
		return nil, err
	}
	return card, nil
}

func (c *FlashcardController) UpdateFlashcard(id string, question string, answer string) (*models.Flashcard, error) {
	card, err := c.Service.UpdateFlashcard(id, question, answer)
	if err != nil {
		fmt.Printf("Failed to update flashcard: %v\n", err)
		return nil, err
	}
	return card, nil
}

func (c *FlashcardController) DeleteFlashcard(id string) error {
	err := c.Service.DeleteFlashcard(id)
	if err != nil {
		fmt.Printf("Failed to delete flashcard: %v\n", err)
		return err
	}
	return nil
}
