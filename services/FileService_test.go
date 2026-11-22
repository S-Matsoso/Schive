package services

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"proj1/models"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestFileDB(t *testing.T) (*gorm.DB, string) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_file.db")
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	err = db.AutoMigrate(&models.Module{}, &models.File{}, &models.ClassSlot{}, &models.Assignment{}, &models.Folder{})
	if err != nil {
		t.Fatalf("failed to auto migrate: %v", err)
	}
	return db, tempDir
}

func TestFoldersAndFileService(t *testing.T) {
	db, tempDir := setupTestFileDB(t)
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	_ = os.Chdir(tempDir)

	fileSvc := NewFileService(db)
	moduleID := "mod-123"

	// 1. Create folder
	folder, err := fileSvc.CreateFolder(moduleID, "Notes")
	if err != nil {
		t.Fatalf("unexpected error creating folder: %v", err)
	}
	if folder.Name != "Notes" {
		t.Fatalf("expected folder name 'Notes', got %q", folder.Name)
	}

	// 2. Duplicate folder validation
	_, err = fileSvc.CreateFolder(moduleID, "notes")
	if err == nil {
		t.Fatalf("expected error creating duplicate folder, got nil")
	}

	// 3. Upload file to folder
	sampleContent := base64.StdEncoding.EncodeToString([]byte("Hello world notes"))
	file, err := fileSvc.UploadFile(moduleID, folder.ID, "lecture1.odt", sampleContent)
	if err != nil {
		t.Fatalf("unexpected error uploading file: %v", err)
	}
	if file.FolderID == nil || *file.FolderID != folder.ID {
		t.Fatalf("expected file to be in folder %s, got %v", folder.ID, file.FolderID)
	}

	// 4. Upload file to root
	rootFile, err := fileSvc.UploadFile(moduleID, "", "syllabus.pdf", sampleContent)
	if err != nil {
		t.Fatalf("unexpected error uploading root file: %v", err)
	}
	if rootFile.FolderID != nil {
		t.Fatalf("expected root file to have nil folderId, got %v", rootFile.FolderID)
	}

	// 5. Check folder file counts
	folders, err := fileSvc.GetFoldersByModule(moduleID)
	if err != nil {
		t.Fatalf("unexpected error getting folders: %v", err)
	}
	if len(folders) != 1 || folders[0].FileCount != 1 {
		t.Fatalf("expected 1 folder with FileCount=1, got %+v", folders)
	}

	// 6. Test file modification time and stat
	time.Sleep(10 * time.Millisecond)
	// Append to file on disk to simulate user editing in LibreOffice
	err = os.WriteFile(file.Path, []byte("Hello world notes - edited"), 0644)
	if err != nil {
		t.Fatalf("failed to write to file: %v", err)
	}

	files, err := fileSvc.GetFilesByModule(moduleID)
	if err != nil {
		t.Fatalf("unexpected error getting files: %v", err)
	}
	if len(files) != 2 {
		t.Fatalf("expected 2 files, got %d", len(files))
	}
	var loadedOdt models.File
	for _, f := range files {
		if f.ID == file.ID {
			loadedOdt = f
		}
	}
	if loadedOdt.FileSize != int64(len("Hello world notes - edited")) {
		t.Errorf("expected updated file size, got %d", loadedOdt.FileSize)
	}
	if loadedOdt.UpdatedAt.IsZero() {
		t.Errorf("expected non-zero UpdatedAt")
	}

	// 7. Move rootFile into folder
	err = fileSvc.MoveFile(rootFile.ID, folder.ID)
	if err != nil {
		t.Fatalf("unexpected error moving file: %v", err)
	}
	folders, _ = fileSvc.GetFoldersByModule(moduleID)
	if folders[0].FileCount != 2 {
		t.Fatalf("expected folder FileCount=2 after move, got %d", folders[0].FileCount)
	}

	// 8. Delete folder cascades
	err = fileSvc.DeleteFolder(folder.ID)
	if err != nil {
		t.Fatalf("unexpected error deleting folder: %v", err)
	}
	folders, _ = fileSvc.GetFoldersByModule(moduleID)
	if len(folders) != 0 {
		t.Fatalf("expected 0 folders after delete, got %d", len(folders))
	}
	files, _ = fileSvc.GetFilesByModule(moduleID)
	if len(files) != 0 {
		t.Fatalf("expected 0 files after cascade delete, got %d", len(files))
	}
}
