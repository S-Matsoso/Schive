package controllers

import "proj1/services"

type TimetableController struct {
	Service *services.TimetableService
}

func NewTimetableController(service *services.TimetableService) *TimetableController {
	return &TimetableController{Service: service}
}
