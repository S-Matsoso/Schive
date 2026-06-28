package services

import (
	// "errors"
	// "path/filepath"
	// "strings"

	// "proj1/models"

	// "github.com/google/uuid"
	"gorm.io/gorm"
)

type FileService struct {
	DB *gorm.DB
}

func NewFileService(db *gorm.DB) *FileService {
	return &FileService{DB: db}
}
