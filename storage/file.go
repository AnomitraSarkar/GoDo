package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"todo-cli/models"
)

func SaveTodos(path string, todos []models.Todo) error {
	data, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func LoadTodos(path string) ([]models.Todo, error) {
	var todos []models.Todo
	data, err := os.ReadFile(path)
	if err != nil {
		return todos, err
	}
	if len(data) == 0 {
		return todos, nil
	}
	err = json.Unmarshal(data, &todos)
	return todos, err
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