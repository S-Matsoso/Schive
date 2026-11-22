package services

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type BackupService struct {
	DBPath string
}

func NewBackupService(dbPath string) *BackupService {
	return &BackupService{DBPath: dbPath}
}

func (s *BackupService) CreateBackup() (string, error) {
	if err := os.MkdirAll("backups", 0755); err != nil {
		return "", err
	}

	timestamp := time.Now().Format("2006-01-02_15-04-05")
	backupFileName := fmt.Sprintf("schive_backup_%s.zip", timestamp)
	backupPath := filepath.Join("backups", backupFileName)

	zipFile, err := os.Create(backupPath)
	if err != nil {
		return "", err
	}
	defer zipFile.Close()

	archive := zip.NewWriter(zipFile)
	defer archive.Close()

	// 1. Add database file if it exists
	if _, err := os.Stat(s.DBPath); err == nil {
		if err := addFileToZip(archive, s.DBPath, filepath.Base(s.DBPath)); err != nil {
			return "", err
		}
	}

	// 2. Add uploads directory if it exists
	if _, err := os.Stat("uploads"); err == nil {
		err = filepath.Walk("uploads", func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return err
			}
			relPath, err := filepath.Rel(".", path)
			if err != nil {
				relPath = path
			}
			return addFileToZip(archive, path, relPath)
		})
		if err != nil {
			return "", err
		}
	}

	absPath, err := filepath.Abs(backupPath)
	if err == nil {
		return absPath, nil
	}
	return backupPath, nil
}

func (s *BackupService) GetBackups() ([]string, error) {
	if _, err := os.Stat("backups"); os.IsNotExist(err) {
		return []string{}, nil
	}

	entries, err := os.ReadDir("backups")
	if err != nil {
		return nil, err
	}

	var backups []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".zip" {
			backups = append(backups, entry.Name())
		}
	}
	return backups, nil
}

func addFileToZip(archive *zip.Writer, filePath string, zipPath string) error {
	fileToZip, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	header.Name = filepath.ToSlash(zipPath)
	header.Method = zip.Deflate

	writer, err := archive.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)
	return err
}
