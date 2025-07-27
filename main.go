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
		NewFile:  "n",
		EditFile: "e",
		DelFile:  "d",
		OpenFile: "o",
		AddTodo:  "a",
		EditTodo: "e",
		DelTodo:  "d",
		Toggle:   "t",
		Back:     "b",
		Quit:     "q",
	}

	defaultConfig := types.Config{
		UndoneColor:         "yellow",
		DoneColor:           "green",
		ActiveWindowColor:   "purple",
		UnactiveWindowColor: "aqua",
		Keymap:              defaultKeymap,
	}

	data, err := os.ReadFile("config.json")
	if err != nil {
		return defaultConfig
	}

	var config types.Config
	if err := json.Unmarshal(data, &config); err != nil {
		return defaultConfig
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