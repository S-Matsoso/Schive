package services

import (
	"path/filepath"
	"proj1/models"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestFlashcardDB(t *testing.T) *gorm.DB {
	dbPath := filepath.Join(t.TempDir(), "test_flashcard.db")
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	err = db.AutoMigrate(&models.Module{}, &models.File{}, &models.ClassSlot{}, &models.Folder{}, &models.Assignment{}, &models.Flashcard{})
	if err != nil {
		t.Fatalf("failed to auto migrate: %v", err)
	}
	return db
}

func TestFlashcardService(t *testing.T) {
	db := setupTestFlashcardDB(t)
	cardSvc := NewFlashcardService(db)
	moduleID := "mod-999"

	// 1. Create flashcard
	card, err := cardSvc.CreateFlashcard(moduleID, "What is O(1)?", "Constant time complexity")
	if err != nil {
		t.Fatalf("failed to create flashcard: %v", err)
	}
	if card.Question != "What is O(1)?" || card.Mastered != false {
		t.Fatalf("unexpected flashcard: %+v", card)
	}

	// 2. Toggle mastered
	card, err = cardSvc.ToggleMastered(card.ID)
	if err != nil || !card.Mastered {
		t.Fatalf("expected card to be mastered, got %v", card.Mastered)
	}

	// 3. Update flashcard
	card, err = cardSvc.UpdateFlashcard(card.ID, "What is O(log n)?", "Logarithmic time complexity")
	if err != nil || card.Question != "What is O(log n)?" {
		t.Fatalf("failed to update flashcard: %v", err)
	}

	// 4. Get by module
	cards, err := cardSvc.GetFlashcardsByModule(moduleID)
	if err != nil || len(cards) != 1 {
		t.Fatalf("expected 1 card, got %d", len(cards))
	}

	// 5. Delete flashcard
	err = cardSvc.DeleteFlashcard(card.ID)
	if err != nil {
		t.Fatalf("failed to delete flashcard: %v", err)
	}
	cards, _ = cardSvc.GetFlashcardsByModule(moduleID)
	if len(cards) != 0 {
		t.Fatalf("expected 0 cards after delete, got %d", len(cards))
	}
}
