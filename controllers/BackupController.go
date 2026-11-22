package controllers

import (
	"fmt"
	"proj1/services"
)

type BackupController struct {
	Service *services.BackupService
}

func NewBackupController(service *services.BackupService) *BackupController {
	return &BackupController{Service: service}
}

func (c *BackupController) CreateBackup() (string, error) {
	path, err := c.Service.CreateBackup()
	if err != nil {
		fmt.Printf("Failed to create backup: %v\n", err)
		return "", err
	}
	return path, nil
}

func (c *BackupController) GetBackups() ([]string, error) {
	backups, err := c.Service.GetBackups()
	if err != nil {
		fmt.Printf("Failed to get backups: %v\n", err)
		return nil, err
	}
	return backups, nil
}
