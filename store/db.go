package store

import (
	"fmt" // Added missing fmt import
	"log"
	"proj1/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(dbPath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		// It's usually a good idea to log the connection error before returning nil
		log.Printf("Failed to connect to the database: %v", err)
		return nil, err
	}

	err = db.AutoMigrate(
		&models.Module{}, &models.ModuleFolder{}, &models.File{}, &models.ClassSlot{}, &models.Folder{}, &models.Assignment{}, &models.Task{}, &models.Flashcard{})
	if err != nil {
		fmt.Printf("Failed to migrate database: %v\n", err)
		return nil, err
	}

	fmt.Println("\033[32mDatabase successfully migrated.\033[0m")

	return db, nil
}
