package controllers

import (
	"fmt"
	"proj1/models"
	"proj1/services"
)

type TaskController struct {
	Service *services.TaskService
}

func NewTaskController(service *services.TaskService) *TaskController {
	return &TaskController{Service: service}
}

func (c *TaskController) CreateTask(moduleID string, title string, notes string, dueDateStr string) (*models.Task, error) {
	task, err := c.Service.CreateTask(moduleID, title, notes, dueDateStr)
	if err != nil {
		fmt.Printf("Failed to create task: %v\n", err)
		return nil, err
	}
	return task, nil
}

func (c *TaskController) GetTasksByModule(moduleID string) ([]models.Task, error) {
	tasks, err := c.Service.GetTasksByModule(moduleID)
	if err != nil {
		fmt.Printf("Failed to get tasks for module %s: %v\n", moduleID, err)
		return nil, err
	}
	return tasks, nil
}

func (c *TaskController) UpdateTask(id string, title string, notes string, dueDateStr string) (*models.Task, error) {
	task, err := c.Service.UpdateTask(id, title, notes, dueDateStr)
	if err != nil {
		fmt.Printf("Failed to update task: %v\n", err)
		return nil, err
	}
	return task, nil
}

func (c *TaskController) SetTaskDone(id string, done bool) (*models.Task, error) {
	task, err := c.Service.SetTaskDone(id, done)
	if err != nil {
		fmt.Printf("Failed to update task status: %v\n", err)
		return nil, err
	}
	return task, nil
}

func (c *TaskController) DeleteTask(id string) error {
	if err := c.Service.DeleteTask(id); err != nil {
		fmt.Printf("Failed to delete task: %v\n", err)
		return err
	}
	return nil
}
