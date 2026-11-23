package controllers

import (
	"context"
	"errors"
	"fmt"
	"proj1/models"
	"proj1/services"
	"strings"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type FileController struct {
	Service *services.FileService
	ctx     context.Context
}

// SetContext stores the Wails runtime context so native dialogs can be shown.
func (c *FileController) SetContext(ctx context.Context) {
	c.ctx = ctx
}

// ExportFile asks the user where to save a copy of a Schive-owned file and
// writes it there. Returns an empty string when the dialog is cancelled.
func (c *FileController) ExportFile(fileID string) (string, error) {
	file, err := c.Service.GetFile(fileID)
	if err != nil {
		fmt.Printf("Failed to export file: %v\n", err)
		return "", err
	}
	if c.ctx == nil {
		return "", errors.New("file dialog is not available")
	}

	filterPattern := "*"
	if ext := strings.TrimSpace(file.FileType); ext != "" {
		filterPattern = "*." + ext
	}

	destination, err := wailsRuntime.SaveFileDialog(c.ctx, wailsRuntime.SaveDialogOptions{
		Title:           "Export " + file.Name,
		DefaultFilename: file.Name,
		Filters: []wailsRuntime.FileFilter{
			{DisplayName: "Schive file (" + filterPattern + ")", Pattern: filterPattern},
			{DisplayName: "All files", Pattern: "*.*"},
		},
	})
	if err != nil {
		fmt.Printf("Failed to open save dialog: %v\n", err)
		return "", err
	}
	if destination == "" {
		return "", nil
	}

	exportedPath, err := c.Service.ExportFile(fileID, destination)
	if err != nil {
		fmt.Printf("Failed to export file: %v\n", err)
		return "", err
	}
	return exportedPath, nil
}

func NewFileController(service *services.FileService) *FileController {
	return &FileController{Service: service}
}

func (c *FileController) CreateFolder(moduleID string, name string) (*models.Folder, error) {
	folder, err := c.Service.CreateFolder(moduleID, name)
	if err != nil {
		fmt.Printf("Failed to create folder: %v\n", err)
		return nil, err
	}
	return folder, nil
}

func (c *FileController) GetFoldersByModule(moduleID string) ([]models.Folder, error) {
	folders, err := c.Service.GetFoldersByModule(moduleID)
	if err != nil {
		fmt.Printf("Failed to get folders: %v\n", err)
		return nil, err
	}
	return folders, nil
}

func (c *FileController) DeleteFolder(id string) error {
	err := c.Service.DeleteFolder(id)
	if err != nil {
		fmt.Printf("Failed to delete folder: %v\n", err)
		return err
	}
	return nil
}

func (c *FileController) MoveFile(fileID string, targetFolderID string) error {
	err := c.Service.MoveFile(fileID, targetFolderID)
	if err != nil {
		fmt.Printf("Failed to move file: %v\n", err)
		return err
	}
	return nil
}

func (c *FileController) RenameFile(id string, newName string) error {
	err := c.Service.RenameFile(id, newName)
	if err != nil {
		fmt.Printf("Failed to rename file: %v\n", err)
		return err
	}
	return nil
}

func (c *FileController) GetFileContent(id string) (string, error) {
	content, err := c.Service.GetFileContent(id)
	if err != nil {
		fmt.Printf("Failed to get file content: %v\n", err)
		return "", err
	}
	return content, nil
}

func (c *FileController) UploadFile(moduleID string, folderID string, name string, base64Data string) (*models.File, error) {
	file, err := c.Service.UploadFile(moduleID, folderID, name, base64Data)
	if err != nil {
		fmt.Printf("Failed to upload file: %v\n", err)
		return nil, err
	}
	return file, nil
}

func (c *FileController) UploadAssignmentInstructions(moduleID string, assignmentID string, name string, base64Data string) (*models.File, error) {
	file, err := c.Service.UploadAssignmentInstructions(moduleID, assignmentID, name, base64Data)
	if err != nil {
		fmt.Printf("Failed to upload assignment instructions: %v\n", err)
		return nil, err
	}
	return file, nil
}

func (c *FileController) CreateAssignmentFile(moduleID string, assignmentID string, name string) (*models.File, error) {
	file, err := c.Service.CreateAssignmentFile(moduleID, assignmentID, name)
	if err != nil {
		fmt.Printf("Failed to create assignment file: %v\n", err)
		return nil, err
	}
	return file, nil
}

func (c *FileController) GetFilesByModule(moduleID string) ([]models.File, error) {
	files, err := c.Service.GetFilesByModule(moduleID)
	if err != nil {
		fmt.Printf("Failed to get files: %v\n", err)
		return nil, err
	}
	return files, nil
}

func (c *FileController) DeleteFile(id string) error {
	err := c.Service.DeleteFile(id)
	if err != nil {
		fmt.Printf("Failed to delete file: %v\n", err)
		return err
	}
	return nil
}

func (c *FileController) OpenFile(id string) error {
	err := c.Service.OpenFile(id)
	if err != nil {
		fmt.Printf("Failed to open file: %v\n", err)
		return err
	}
	return nil
}
