package controllers

import (
	"fmt"
	"proj1/models"
	"proj1/services"
)

type ModuleController struct {
	Service *services.ModuleService
}

func NewModuleController(service *services.ModuleService) *ModuleController {
	return &ModuleController{Service: service}
}

func (c *ModuleController) GetAllModules() ([]models.Module, error) {
	m, err := c.Service.GetAllModules()

	if err != nil {
		fmt.Println("Failed to fetch modules")
		return nil, err
	}
	return m, nil
}

func (c *ModuleController) GetModule(code string) (*models.Module, error) {
	m, err := c.Service.GetModule(code)
	if err != nil {
		fmt.Println("Failed to get module")
		return nil, err
	}
	return m, nil
}

func (c *ModuleController) CreateModule(module models.Module) (*models.Module, error) {
	m, err := c.Service.CreateModule(&module)
	if err != nil {
		fmt.Println("Failed to create a module:", err)
		return nil, err
	}
	return m, nil
}

func (c *ModuleController) UpdateModule(code string, module models.Module) (*models.Module, error) {
	m, err := c.Service.UpdateModule(code, &module)
	if err != nil {
		fmt.Printf("Failed to update module %s: %v\n", code, err)
		return nil, err
	}
	return m, nil
}

func (c *ModuleController) DeleteModule(code string) (string, error) {
	err := c.Service.DeleteModule(code)
	if err != nil {
		return fmt.Sprintf("Failed to delete module %v", code), err
	}
	return fmt.Sprintf("Successfully delete module %v", code), nil
}
