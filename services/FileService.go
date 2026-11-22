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

func (s *FileService) CreateFolder(moduleID string, name string) (*models.Folder, error) {
	cleanModuleID := strings.TrimSpace(moduleID)
	cleanName := strings.TrimSpace(name)

	if cleanModuleID == "" {
		return nil, errors.New("moduleID is required")
	}
	if cleanName == "" {
		return nil, errors.New("folder name is required")
	}

	var existing []models.Folder
	if err := s.DB.Where("module_id = ?", cleanModuleID).Find(&existing).Error; err != nil {
		return nil, err
	}
	for _, f := range existing {
		if strings.EqualFold(strings.TrimSpace(f.Name), cleanName) {
			return nil, errors.New("folder with this name already exists in this module")
		}
	}

	folder := models.Folder{
		ID:        uuid.NewString(),
		ModuleID:  cleanModuleID,
		Name:      cleanName,
		CreatedAt: time.Now(),
		FileCount: 0,
	}

	if err := s.DB.Create(&folder).Error; err != nil {
		return nil, err
	}

	return &folder, nil
}

func (s *FileService) GetFoldersByModule(moduleID string) ([]models.Folder, error) {
	cleanModuleID := strings.TrimSpace(moduleID)
	if cleanModuleID == "" {
		return nil, errors.New("moduleID is required")
	}

	var folders []models.Folder
	if err := s.DB.Where("module_id = ?", cleanModuleID).Order("created_at asc").Find(&folders).Error; err != nil {
		return nil, err
	}

	for i := range folders {
		var count int64
		s.DB.Model(&models.File{}).Where("folder_id = ?", folders[i].ID).Count(&count)
		folders[i].FileCount = int(count)
	}

	return folders, nil
}

func (s *FileService) DeleteFolder(folderID string) error {
	cleanFolderID := strings.TrimSpace(folderID)
	if cleanFolderID == "" {
		return errors.New("folderID is required")
	}

	var files []models.File
	if err := s.DB.Where("folder_id = ?", cleanFolderID).Find(&files).Error; err != nil {
		return err
	}

	for _, f := range files {
		if f.Path != "" {
			_ = os.Remove(f.Path)
		}
	}

	if err := s.DB.Where("folder_id = ?", cleanFolderID).Delete(&models.File{}).Error; err != nil {
		return err
	}

	return s.DB.Where("id = ?", cleanFolderID).Delete(&models.Folder{}).Error
}

func (s *FileService) MoveFile(fileID string, targetFolderID string) error {
	cleanFileID := strings.TrimSpace(fileID)
	if cleanFileID == "" {
		return errors.New("fileID is required")
	}

	var file models.File
	if err := s.DB.First(&file, "id = ?", cleanFileID).Error; err != nil {
		return err
	}

	cleanTarget := strings.TrimSpace(targetFolderID)
	if cleanTarget == "" {
		file.FolderID = nil
	} else {
		var folder models.Folder
		if err := s.DB.First(&folder, "id = ?", cleanTarget).Error; err != nil {
			return errors.New("target folder not found")
		}
		file.FolderID = &cleanTarget
	}

	return s.DB.Save(&file).Error
}

func (s *FileService) RenameFile(id string, newName string) error {
	cleanID := strings.TrimSpace(id)
	cleanName := strings.TrimSpace(newName)
	if cleanID == "" || cleanName == "" {
		return errors.New("file id and new name are required")
	}

	var file models.File
	if err := s.DB.First(&file, "id = ?", cleanID).Error; err != nil {
		return err
	}

	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(cleanName), "."))
	file.Name = cleanName
	if ext != "" {
		file.FileType = ext
	}

	return s.DB.Save(&file).Error
}

func (s *FileService) GetFileContent(id string) (string, error) {
	cleanID := strings.TrimSpace(id)
	if cleanID == "" {
		return "", errors.New("file id is required")
	}

	var file models.File
	if err := s.DB.First(&file, "id = ?", cleanID).Error; err != nil {
		return "", err
	}

	if file.Path == "" {
		return "", errors.New("file path not found")
	}

	data, err := os.ReadFile(file.Path)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(data), nil
}

func (s *FileService) UploadFile(moduleID string, folderID string, name string, base64Data string) (*models.File, error) {
	cleanModuleID := strings.TrimSpace(moduleID)
	cleanName := strings.TrimSpace(name)
	cleanFolderID := strings.TrimSpace(folderID)

	if cleanModuleID == "" {
		return nil, errors.New("moduleID is required")
	}
	if cleanName == "" {
		return nil, errors.New("file name is required")
	}

	var folderIDPtr *string
	if cleanFolderID != "" {
		var folder models.Folder
		if err := s.DB.First(&folder, "id = ? AND module_id = ?", cleanFolderID, cleanModuleID).Error; err != nil {
			return nil, errors.New("specified folder does not exist")
		}
		folderIDPtr = &cleanFolderID
	}

	data, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return nil, errors.New("invalid base64 data: " + err.Error())
	}

	storageDir := filepath.Join("uploads", cleanModuleID)
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return nil, errors.New("failed to create upload directory: " + err.Error())
	}

	fileID := uuid.NewString()
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(cleanName), "."))
	safeName := filepath.Base(cleanName)
	targetPath := filepath.Join(storageDir, fileID+"_"+safeName)

	if err := os.WriteFile(targetPath, data, 0644); err != nil {
		return nil, errors.New("failed to save file on disk: " + err.Error())
	}

	now := time.Now()
	file := models.File{
		ID:        fileID,
		ModuleID:  cleanModuleID,
		FolderID:  folderIDPtr,
		Name:      safeName,
		Path:      targetPath,
		FileType:  ext,
		FileSize:  int64(len(data)),
		CreatedAt: now,
		UpdatedAt: now,
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
	if err != nil {
		return nil, err
	}

	// Update live file size and last modification date from disk
	for i := range files {
		if files[i].Path != "" {
			if stat, err := os.Stat(files[i].Path); err == nil {
				files[i].UpdatedAt = stat.ModTime()
				files[i].FileSize = stat.Size()
			}
		}
		if files[i].UpdatedAt.IsZero() {
			files[i].UpdatedAt = files[i].CreatedAt
		}
	}

	return files, nil
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
