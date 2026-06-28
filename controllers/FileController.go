package controllers

import "proj1/services"

type FileController struct {
	Service *services.FileService
}

func NewFileController(service *services.FileService) *FileController {
	return &FileController{Service: service}
}
