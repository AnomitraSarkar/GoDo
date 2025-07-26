package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"todo-cli/models"
	"todo-cli/storage"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
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

func showInputDialog(title, label, initialText string, callback func(string)) {
	inputField := tview.NewInputField().
		SetLabel(label).
		SetText(initialText)

	form := tview.NewForm().
		AddFormItem(inputField).
		AddButton("OK", func() {
			text := inputField.GetText()
			pages.RemovePage("dialog")
			callback(text)
		}).
		AddButton("Cancel", func() {
			pages.RemovePage("dialog")
		})

	// Set up keyboard navigation
	inputField.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			app.SetFocus(form.GetButton(0)) // Focus OK button
		case tcell.KeyEsc:
			pages.RemovePage("dialog")
		case tcell.KeyTab:
			app.SetFocus(form.GetButton(0)) // Focus OK button
		}
	})

	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			pages.RemovePage("dialog")
			return nil
		case tcell.KeyTab:
			// Cycle through buttons
			currentFocus := app.GetFocus()
			if currentFocus == form.GetButton(0) {
				app.SetFocus(form.GetButton(1))
			} else {
				app.SetFocus(form.GetButton(0))
			}
			return nil
		}
		return event
	})

	form.SetTitle(title).SetBorder(true)
	flex := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(form, 7, 1, false).
			AddItem(nil, 0, 1, false), 60, 1, false).
		AddItem(nil, 0, 1, false)

	pages.AddPage("dialog", flex, true, true)
	app.SetFocus(inputField)
}

func registerKeybindings(dir string) {
	fileListView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'n': // Create new file
			showInputDialog("New File", "File name (without .todo):", "", func(name string) {
				if name != "" {
					name = strings.TrimSuffix(name, ".todo")
					path := filepath.Join(dir, name+".todo")
					_ = storage.SaveTodos(path, []models.Todo{})
					refreshFileList(dir)
				}
			})

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
			
		case 'e': // Edit file name
			index := fileListView.GetCurrentItem()
			if index >= 0 && index < fileListView.GetItemCount() {
				oldName, _ := fileListView.GetItemText(index)
				oldName = strings.TrimSuffix(oldName, ".todo")
				showInputDialog("Edit File", "New name:", oldName, func(newName string) {
					if newName != "" {
						newName = strings.TrimSuffix(newName, ".todo")
						newPath := filepath.Join(dir, newName+".todo")
						oldPath := filepath.Join(dir, oldName+".todo")
						_ = os.Rename(oldPath, newPath)
						if currentFile == oldPath {
							currentFile = newPath
						}
						refreshFileList(dir)
					}
				})
			}
			
		case 'q': // Quit application
			app.Stop()
		}
		return event
	})

	todoListView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'a': // Add new todo
			showInputDialog("New Todo", "Enter todo text:", "", func(text string) {
				if text != "" {
					todos, _ := storage.LoadTodos(currentFile)
					todos = append(todos, models.Todo{Text: text, Done: false})
					_ = storage.SaveTodos(currentFile, todos)
					refreshTodoList()
				}
			})

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
			
		case 'e': // Edit todo text
			index := todoListView.GetCurrentItem()
			todos, err := storage.LoadTodos(currentFile)
			if err != nil {
				return event
			}
			if index >= 0 && index < len(todos) {
				showInputDialog("Edit Todo", "New text:", todos[index].Text, func(newText string) {
					if newText != "" {
						todos[index].Text = newText
						_ = storage.SaveTodos(currentFile, todos)
						refreshTodoList()
					}
				})
			}
			
		case 'q': // Quit application
			app.Stop()
		}
		return event
	})
}