package services

import (
	"path/filepath"
	"proj1/models"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	dbPath := filepath.Join(t.TempDir(), "test.db")
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	err = db.AutoMigrate(&models.Module{}, &models.ModuleFolder{}, &models.File{}, &models.ClassSlot{}, &models.Assignment{}, &models.Folder{}, &models.Task{}, &models.Flashcard{})
	if err != nil {
		t.Fatalf("failed to auto migrate: %v", err)
	}
	return db
}

func TestCreateModuleValidation(t *testing.T) {
	db := setupTestDB(t)
	svc := NewModuleService(db)

	tests := []struct {
		name        string
		module      *models.Module
		wantErr     bool
		errContains string
	}{
		{
			name:        "nil module",
			module:      nil,
			wantErr:     true,
			errContains: "module data is required",
		},
		{
			name: "empty code",
			module: &models.Module{
				Code: "",
				Name: "Introduction to CS",
			},
			wantErr:     true,
			errContains: "module code is required",
		},
		{
			name: "whitespace-only code",
			module: &models.Module{
				Code: "   ",
				Name: "Introduction to CS",
			},
			wantErr:     true,
			errContains: "module code is required",
		},
		{
			name: "empty name",
			module: &models.Module{
				Code: "CS101",
				Name: "",
			},
			wantErr:     true,
			errContains: "module name is required",
		},
		{
			name: "whitespace-only name",
			module: &models.Module{
				Code: "CS101",
				Name: "   \t\n  ",
			},
			wantErr:     true,
			errContains: "module name is required",
		},
		{
			name: "valid module",
			module: &models.Module{
				Code:  "  CS101  ",
				Name:  "  Intro to Programming  ",
				Color: "#FF5722",
			},
			wantErr: false,
		},
		{
			name: "duplicate code (exact)",
			module: &models.Module{
				Code: "CS101",
				Name: "Another Course",
			},
			wantErr:     true,
			errContains: "module with this code already exists",
		},
		{
			name: "duplicate code (case insensitive)",
			module: &models.Module{
				Code: "cs101",
				Name: "Yet Another Course",
			},
			wantErr:     true,
			errContains: "module with this code already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			created, err := svc.CreateModule(tt.module)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateModule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if tt.errContains != "" && (err == nil || !contains(err.Error(), tt.errContains)) {
					t.Errorf("expected error containing %q, got %v", tt.errContains, err)
				}
			} else {
				if created == nil {
					t.Fatalf("expected non-nil created module")
				}
				if created.Code != "CS101" {
					t.Errorf("expected trimmed code 'CS101', got %q", created.Code)
				}
				if created.Name != "Intro to Programming" {
					t.Errorf("expected trimmed name 'Intro to Programming', got %q", created.Name)
				}
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || (len(s) > 0 && len(substr) > 0 && stringContains(s, substr)))
}

func stringContains(s, substr string) bool {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
