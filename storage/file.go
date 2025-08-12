package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
	"todo-cli/models"
)

type TodoFileData struct {
	CreatedAt time.Time     `json:"createdAt"`
	Todos     []models.Todo `json:"todos"`
	Priority  bool          `json:"priority"`
}

func SaveTodos(path string, fileData TodoFileData) error {
	data, err := json.MarshalIndent(fileData, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func LoadTodos(path string) (TodoFileData, error) {
	var fileData TodoFileData
	data, err := os.ReadFile(path)
	if err != nil {
		return fileData, err
	}

	// First, try to unmarshal into the new format
	err = json.Unmarshal(data, &fileData)
	if err == nil {
		// Set default creation time if missing
		if fileData.CreatedAt.IsZero() {
			fileInfo, _ := os.Stat(path)
			fileData.CreatedAt = fileInfo.ModTime()
		}
		return fileData, nil
	}

	// If that fails, try the old format
	var oldTodos []models.Todo
	if err := json.Unmarshal(data, &oldTodos); err != nil {
		return fileData, err
	}

	// Get file mod time as creation time
	fileInfo, err := os.Stat(path)
	modTime := time.Now()
	if err == nil {
		modTime = fileInfo.ModTime()
	}

	// Set creation times for old todos
	for i := range oldTodos {
		if oldTodos[i].CreatedAt.IsZero() {
			oldTodos[i].CreatedAt = modTime
		}
	}

	return TodoFileData{
		CreatedAt: modTime,
		Todos:     oldTodos,
		Priority:  false,
	}, nil
}

func ListTodoFiles(dir string) ([]string, error) {
	var files []string
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".todo" {
			files = append(files, entry.Name())
		}
	}
	return files, nil
}