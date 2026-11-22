package controllers

import (
	"fmt"
	"proj1/models"
	"proj1/services"
)

type FileController struct {
	Service *services.FileService
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
