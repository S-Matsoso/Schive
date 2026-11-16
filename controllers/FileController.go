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

func (c *FileController) UploadFile(moduleID string, name string, base64Data string) (*models.File, error) {
	file, err := c.Service.UploadFile(moduleID, name, base64Data)
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
