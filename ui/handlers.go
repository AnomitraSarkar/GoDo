package ui

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"todo-cli/models"
	"todo-cli/storage"

	"github.com/gdamore/tcell/v2"
)

func refreshFileList(dir string) {
	fileListView.Clear()
	todoFiles, err := storage.ListTodoFiles(dir)
	if err != nil {
		return
	}
	for _, file := range todoFiles {
		fileListView.AddItem(file, "", 0, nil)
	}
}

func refreshTodoList() {
	todoListView.Clear()
	if currentFile == "" {
		return
	}
	todos, err := storage.LoadTodos(currentFile)
	if err != nil {
		return
	}
	for _, t := range todos {
		prefix := "❌"
		if t.Done {
			prefix = "✅"
		}
		todoListView.AddItem(fmt.Sprintf("%s %s", prefix, t.Text), "", 0, nil)
	}
}

func registerKeybindings(dir string) {
	fileListView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'n': // Create new file
			app.Suspend(func() {
				fmt.Print("Enter new file name (without .todo): ")
				reader := bufio.NewReader(os.Stdin)
				name, _ := reader.ReadString('\n')
				name = strings.TrimSpace(name)
				if name != "" {
					name = strings.TrimSuffix(name, ".todo")
					path := filepath.Join(dir, name+".todo")
					_ = storage.SaveTodos(path, []models.Todo{})
				}
			})
			refreshFileList(dir)

		case 'd': // Delete selected file
			index := fileListView.GetCurrentItem()
			if index >= 0 && index < fileListView.GetItemCount() {
				name, _ := fileListView.GetItemText(index)
				_ = os.Remove(filepath.Join(dir, name))
				if currentFile == filepath.Join(dir, name) {
					currentFile = ""
					todoListView.Clear()
				}
				refreshFileList(dir)
			}

		case 'o': // Open selected file
			index := fileListView.GetCurrentItem()
			if index >= 0 && index < fileListView.GetItemCount() {
				name, _ := fileListView.GetItemText(index)
				currentFile = filepath.Join(dir, name)
				refreshTodoList()
				app.SetFocus(todoListView)
			}
		}
		return event
	})

	todoListView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'a': // Add new todo
			app.Suspend(func() {
				fmt.Print("Enter new todo: ")
				reader := bufio.NewReader(os.Stdin)
				text, _ := reader.ReadString('\n')
				text = strings.TrimSpace(text)
				if text != "" {
					todos, _ := storage.LoadTodos(currentFile)
					todos = append(todos, models.Todo{Text: text, Done: false})
					_ = storage.SaveTodos(currentFile, todos)
				}
			})
			refreshTodoList()

		case 'd': // Delete selected todo
			index := todoListView.GetCurrentItem()
			todos, err := storage.LoadTodos(currentFile)
			if err != nil {
				return event
			}
			if index >= 0 && index < len(todos) {
				todos = append(todos[:index], todos[index+1:]...)
				_ = storage.SaveTodos(currentFile, todos)
			}
			refreshTodoList()

		case 't': // Toggle done
			index := todoListView.GetCurrentItem()
			todos, err := storage.LoadTodos(currentFile)
			if err != nil {
				return event
			}
			if index >= 0 && index < len(todos) {
				todos[index].Done = !todos[index].Done
				_ = storage.SaveTodos(currentFile, todos)
			}
			refreshTodoList()
			
		case 'b': // Go back to file list
			app.SetFocus(fileListView)
		}
		return event
	})
}