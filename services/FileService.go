package services

import (
	"encoding/base64"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"proj1/models"
	"runtime"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FileService struct {
	DB *gorm.DB
}

func NewFileService(db *gorm.DB) *FileService {
	return &FileService{DB: db}
}

func (s *FileService) UploadFile(moduleID string, name string, base64Data string) (*models.File, error) {
	if moduleID == "" {
		return nil, errors.New("moduleID is required")
	}
	if name == "" {
		return nil, errors.New("file name is required")
	}

	data, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return nil, errors.New("invalid base64 data: " + err.Error())
	}

	storageDir := filepath.Join("uploads", moduleID)
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return nil, errors.New("failed to create upload directory: " + err.Error())
	}

	fileID := uuid.NewString()
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(name), "."))
	safeName := filepath.Base(name)
	targetPath := filepath.Join(storageDir, fileID+"_"+safeName)

	if err := os.WriteFile(targetPath, data, 0644); err != nil {
		return nil, errors.New("failed to save file on disk: " + err.Error())
	}

	file := models.File{
		ID:        fileID,
		ModuleID:  moduleID,
		Name:      safeName,
		Path:      targetPath,
		FileType:  ext,
		FileSize:  int64(len(data)),
		CreatedAt: time.Now(),
	}

	if err := s.DB.Create(&file).Error; err != nil {
		_ = os.Remove(targetPath)
		return nil, err
	}

	return &file, nil
}

func (s *FileService) GetFilesByModule(moduleID string) ([]models.File, error) {
	var files []models.File
	err := s.DB.Where("module_id = ?", moduleID).Order("created_at desc").Find(&files).Error
	return files, err
}

func (s *FileService) DeleteFile(id string) error {
	var file models.File
	if err := s.DB.First(&file, "id = ?", id).Error; err != nil {
		return err
	}

	if file.Path != "" {
		_ = os.Remove(file.Path)
	}

	return s.DB.Delete(&models.File{}, "id = ?", id).Error
}

func (s *FileService) OpenFile(id string) error {
	var file models.File
	if err := s.DB.First(&file, "id = ?", id).Error; err != nil {
		return err
	}

	absPath, err := filepath.Abs(file.Path)
	if err != nil {
		absPath = file.Path
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/C", "start", "", absPath)
	case "darwin":
		cmd = exec.Command("open", absPath)
	default: // linux, etc.
		cmd = exec.Command("xdg-open", absPath)
	}

	return cmd.Start()
}
