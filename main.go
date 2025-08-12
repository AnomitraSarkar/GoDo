package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"todo-cli/types"
	"todo-cli/ui"
)

func loadConfig() types.Config {
	defaultKeymap := types.Keymap{
		NewFile:      "n",
		EditFile:     "e",
		DelFile:      "d",
		OpenFile:     "o",
		AddTodo:      "a",
		EditTodo:     "e",
		DelTodo:      "d",
		Toggle:       "t",
		Back:         "b",
		Quit:         "q",
		MoveToTop:    "g",
		MoveToBottom: "G",
		SetPriority:  "p",
	}

	defaultConfig := types.Config{
		UndoneColor:         "yellow",
		DoneColor:           "green",
		ActiveWindowColor:   "purple",
		UnactiveWindowColor: "aqua",
		PriorityColor:       "red",
		RelativeTime:        true,
		Keymap:              defaultKeymap,
	}

	data, err := os.ReadFile("config.json")
	if err != nil {
		return defaultConfig
	}

	// Create config with defaults
	config := defaultConfig

	// Unmarshal into temporary config
	var tempConfig types.Config
	if err := json.Unmarshal(data, &tempConfig); err != nil {
		return defaultConfig
	}

	// Merge non-keymap fields
	if tempConfig.UndoneColor != "" {
		config.UndoneColor = tempConfig.UndoneColor
	}
	if tempConfig.DoneColor != "" {
		config.DoneColor = tempConfig.DoneColor
	}
	if tempConfig.ActiveWindowColor != "" {
		config.ActiveWindowColor = tempConfig.ActiveWindowColor
	}
	if tempConfig.UnactiveWindowColor != "" {
		config.UnactiveWindowColor = tempConfig.UnactiveWindowColor
	}
	if tempConfig.PriorityColor != "" {
		config.PriorityColor = tempConfig.PriorityColor
	}
	config.RelativeTime = tempConfig.RelativeTime

	// Merge keymap ensuring no empty strings
	if tempConfig.Keymap.NewFile != "" {
		config.Keymap.NewFile = tempConfig.Keymap.NewFile
	}
	if tempConfig.Keymap.EditFile != "" {
		config.Keymap.EditFile = tempConfig.Keymap.EditFile
	}
	if tempConfig.Keymap.DelFile != "" {
		config.Keymap.DelFile = tempConfig.Keymap.DelFile
	}
	if tempConfig.Keymap.OpenFile != "" {
		config.Keymap.OpenFile = tempConfig.Keymap.OpenFile
	}
	if tempConfig.Keymap.AddTodo != "" {
		config.Keymap.AddTodo = tempConfig.Keymap.AddTodo
	}
	if tempConfig.Keymap.EditTodo != "" {
		config.Keymap.EditTodo = tempConfig.Keymap.EditTodo
	}
	if tempConfig.Keymap.DelTodo != "" {
		config.Keymap.DelTodo = tempConfig.Keymap.DelTodo
	}
	if tempConfig.Keymap.Toggle != "" {
		config.Keymap.Toggle = tempConfig.Keymap.Toggle
	}
	if tempConfig.Keymap.Back != "" {
		config.Keymap.Back = tempConfig.Keymap.Back
	}
	if tempConfig.Keymap.Quit != "" {
		config.Keymap.Quit = tempConfig.Keymap.Quit
	}
	if tempConfig.Keymap.MoveToTop != "" {
		config.Keymap.MoveToTop = tempConfig.Keymap.MoveToTop
	}
	if tempConfig.Keymap.MoveToBottom != "" {
		config.Keymap.MoveToBottom = tempConfig.Keymap.MoveToBottom
	}
	if tempConfig.Keymap.SetPriority != "" {
		config.Keymap.SetPriority = tempConfig.Keymap.SetPriority
	}

	return config
}

func main() {
	dir := "todos"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, os.ModePerm); err != nil {
			panic(err)
		}
	}

	config := loadConfig()
	ui.StartApp(filepath.Join(".", dir), config)
}