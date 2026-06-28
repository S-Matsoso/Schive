package main

import (
	"context"
	"proj1/controllers"
	"proj1/services"
	"proj1/store"
)

// App struct
type App struct {
	ctx                 context.Context
	ModuleController    *controllers.ModuleController
	FileController      *controllers.FileController
	TimetableController *controllers.TimetableController
}

func NewApp() *App {
	db, err := store.InitDB("schiveApplicationDB.db")
	if err != nil {
		panic(err)
	}

	moduleService := services.NewModuleService(db)
	fileService := services.NewFileService(db)
	timeTableService := services.NewTimetableService(db)

	return &App{
		ModuleController:    controllers.NewModuleController(moduleService),
		FileController:      controllers.NewFileController(fileService),
		TimetableController: controllers.NewTimetableController(timeTableService),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) shutdown(ctx context.Context) {
}
