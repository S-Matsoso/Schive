package services

import (
	"path/filepath"
	"proj1/models"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestAssignmentDB(t *testing.T) *gorm.DB {
	dbPath := filepath.Join(t.TempDir(), "test_assignment.db")
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

func TestAssignmentService(t *testing.T) {
	db := setupTestAssignmentDB(t)
	modSvc := NewModuleService(db)
	assignSvc := NewAssignmentService(db)

	mod, err := modSvc.CreateModule(&models.Module{
		Code: "CS201",
		Name: "Data Structures",
	})
	if err != nil {
		t.Fatalf("failed to create module: %v", err)
	}

	// 1. Create assignment
	dueDate := time.Now().Add(48 * time.Hour).Format(time.RFC3339)
	assign, err := assignSvc.CreateAssignment(mod.ID, "Homework 1", "Binary Trees", 1, "pending", "high", dueDate)
	if err != nil {
		t.Fatalf("failed to create assignment: %v", err)
	}
	if assign.Title != "Homework 1" || assign.Priority != "high" {
		t.Fatalf("unexpected assignment values: %+v", assign)
	}

	// 2. Get by module
	list, err := assignSvc.GetAssignmentsByModule(mod.ID)
	if err != nil || len(list) != 1 {
		t.Fatalf("expected 1 assignment, got %d (err: %v)", len(list), err)
	}

	// 3. Update progress (should also move the status)
	progressed, err := assignSvc.UpdateAssignmentProgress(assign.ID, 40)
	if err != nil {
		t.Fatalf("failed to update progress: %v", err)
	}
	if progressed.Progress != 40 || progressed.Status != "in_progress" {
		t.Fatalf("unexpected progress update: %+v", progressed)
	}

	// 4. Update status
	err = assignSvc.UpdateAssignmentStatus(assign.ID, "completed")
	if err != nil {
		t.Fatalf("failed to update status: %v", err)
	}
	var completed models.Assignment
	if err := db.First(&completed, "id = ?", assign.ID).Error; err != nil {
		t.Fatalf("failed to reload assignment: %v", err)
	}
	if completed.Progress != 100 {
		t.Fatalf("expected progress 100 after completing, got %d", completed.Progress)
	}

	// 5. Update assignment details
	updated, err := assignSvc.UpdateAssignment(assign.ID, "Homework 1 - Revised", "AVL Trees", "completed", "medium", "")
	if err != nil || updated.Title != "Homework 1 - Revised" {
		t.Fatalf("failed to update assignment: %v", err)
	}

	// 6. Delete assignment
	err = assignSvc.DeleteAssignment(assign.ID)
	if err != nil {
		t.Fatalf("failed to delete assignment: %v", err)
	}
	list, _ = assignSvc.GetAssignmentsByModule(mod.ID)
	if len(list) != 0 {
		t.Fatalf("expected 0 assignments after delete, got %d", len(list))
	}
}
