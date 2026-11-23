package main

import (
	"context"
	"proj1/controllers"
	"proj1/services"
	"proj1/store"
)

// App struct
type App struct {
	ctx                  context.Context
	ModuleController     *controllers.ModuleController
	FileController       *controllers.FileController
	TimetableController  *controllers.TimetableController
	AssignmentController *controllers.AssignmentController
	TaskController       *controllers.TaskController
	FlashcardController  *controllers.FlashcardController
	BackupController     *controllers.BackupController
}

func NewApp() *App {
	dbPath := "schiveApplicationDB.db"
	db, err := store.InitDB(dbPath)
	if err != nil {
		panic(err)
	}

	moduleService := services.NewModuleService(db)
	fileService := services.NewFileService(db)
	timeTableService := services.NewTimetableService(db)
	assignmentService := services.NewAssignmentService(db)
	taskService := services.NewTaskService(db)
	flashcardService := services.NewFlashcardService(db)
	backupService := services.NewBackupService(dbPath)

	return &App{
		ModuleController:     controllers.NewModuleController(moduleService),
		FileController:       controllers.NewFileController(fileService),
		TimetableController:  controllers.NewTimetableController(timeTableService),
		AssignmentController: controllers.NewAssignmentController(assignmentService),
		TaskController:       controllers.NewTaskController(taskService),
		FlashcardController:  controllers.NewFlashcardController(flashcardService),
		BackupController:     controllers.NewBackupController(backupService),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.FileController.SetContext(ctx)
}

func (a *App) shutdown(ctx context.Context) {
}
