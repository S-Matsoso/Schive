package services

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBackupService(t *testing.T) {
	tempDir := t.TempDir()
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	_ = os.Chdir(tempDir)

	dbFile := filepath.Join(tempDir, "test.db")
	_ = os.WriteFile(dbFile, []byte("fake db data"), 0644)
	_ = os.MkdirAll("uploads/mod1", 0755)
	_ = os.WriteFile("uploads/mod1/file1.txt", []byte("hello notes"), 0644)

	backupSvc := NewBackupService(dbFile)
	backupPath, err := backupSvc.CreateBackup()
	if err != nil {
		t.Fatalf("failed to create backup: %v", err)
	}
	if _, err := os.Stat(backupPath); err != nil {
		t.Fatalf("backup file does not exist: %v", err)
	}

	backups, err := backupSvc.GetBackups()
	if err != nil || len(backups) != 1 {
		t.Fatalf("expected 1 backup, got %d (err: %v)", len(backups), err)
	}
}
