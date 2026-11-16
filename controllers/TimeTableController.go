package controllers

import (
	"fmt"
	"proj1/models"
	"proj1/services"
)

type TimetableController struct {
	Service *services.TimetableService
}

func NewTimetableController(service *services.TimetableService) *TimetableController {
	return &TimetableController{Service: service}
}

func (c *TimetableController) GetAllClassSlots() ([]models.ClassSlot, error) {
	slots, err := c.Service.GetAllClassSlots()
	if err != nil {
		fmt.Printf("Failed to get class slots: %v\n", err)
		return nil, err
	}
	return slots, nil
}

func (c *TimetableController) CreateClassSlot(slot models.ClassSlot) (*models.ClassSlot, error) {
	newSlot, err := c.Service.CreateClassSlot(&slot)
	if err != nil {
		fmt.Printf("Failed to create class slot: %v\n", err)
		return nil, err
	}
	return newSlot, nil
}

func (c *TimetableController) DeleteClassSlot(id string) error {
	err := c.Service.DeleteClassSlot(id)
	if err != nil {
		fmt.Printf("Failed to delete class slot: %v\n", err)
		return err
	}
	return nil
}
