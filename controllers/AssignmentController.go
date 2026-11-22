package controllers

import (
	"fmt"
	"proj1/models"
	"proj1/services"
)

type AssignmentController struct {
	Service *services.AssignmentService
}

func NewAssignmentController(service *services.AssignmentService) *AssignmentController {
	return &AssignmentController{Service: service}
}

func (c *AssignmentController) CreateAssignment(moduleID string, title string, description string, assignmentNumber int, status string, priority string, dueDateStr string) (*models.Assignment, error) {
	assignment, err := c.Service.CreateAssignment(moduleID, title, description, assignmentNumber, status, priority, dueDateStr)
	if err != nil {
		fmt.Printf("Failed to create assignment: %v\n", err)
		return nil, err
	}
	return assignment, nil
}

func (c *AssignmentController) GetAssignmentsByModule(moduleID string) ([]models.Assignment, error) {
	assignments, err := c.Service.GetAssignmentsByModule(moduleID)
	if err != nil {
		fmt.Printf("Failed to get assignments for module %s: %v\n", moduleID, err)
		return nil, err
	}
	return assignments, nil
}

func (c *AssignmentController) GetAllAssignments() ([]models.Assignment, error) {
	assignments, err := c.Service.GetAllAssignments()
	if err != nil {
		fmt.Printf("Failed to get all assignments: %v\n", err)
		return nil, err
	}
	return assignments, nil
}

func (c *AssignmentController) UpdateAssignment(id string, title string, description string, status string, priority string, dueDateStr string) (*models.Assignment, error) {
	assignment, err := c.Service.UpdateAssignment(id, title, description, status, priority, dueDateStr)
	if err != nil {
		fmt.Printf("Failed to update assignment: %v\n", err)
		return nil, err
	}
	return assignment, nil
}

func (c *AssignmentController) UpdateAssignmentStatus(id string, status string) error {
	err := c.Service.UpdateAssignmentStatus(id, status)
	if err != nil {
		fmt.Printf("Failed to update assignment status: %v\n", err)
		return err
	}
	return nil
}

func (c *AssignmentController) UpdateAssignmentProgress(id string, progress int) (*models.Assignment, error) {
	assignment, err := c.Service.UpdateAssignmentProgress(id, progress)
	if err != nil {
		fmt.Printf("Failed to update assignment progress: %v\n", err)
		return nil, err
	}
	return assignment, nil
}

func (c *AssignmentController) DeleteAssignment(id string) error {
	err := c.Service.DeleteAssignment(id)
	if err != nil {
		fmt.Printf("Failed to delete assignment: %v\n", err)
		return err
	}
	return nil
}
