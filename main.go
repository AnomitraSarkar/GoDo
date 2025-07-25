package main

import (
	"os"
	"path/filepath"
	"todo-cli/ui"
)

func main() {
	dir := "todos"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, os.ModePerm); err != nil {
			panic(err)
		}
	}
	ui.StartApp(filepath.Join(".", dir))
}