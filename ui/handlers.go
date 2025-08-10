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
		color := config.UndoneColor
		if t.Done {
			prefix = "✅"
			color = config.DoneColor
		}
		
		text := fmt.Sprintf("[%s]%s %s", color, prefix, t.Text)
		todoListView.AddItem(text, "", 0, nil)
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

	inputField.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			app.SetFocus(form.GetButton(0))
		case tcell.KeyEsc:
			pages.RemovePage("dialog")
		case tcell.KeyTab:
			app.SetFocus(form.GetButton(0))
		}
	})

	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			pages.RemovePage("dialog")
			return nil
		case tcell.KeyTab:
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
	key := config.Keymap

	fileListView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case rune(key.NewFile[0]): // Create new file
			showInputDialog("New File", "File name (without .todo):", "", func(name string) {
				if name != "" {
					name = strings.TrimSuffix(name, ".todo")
					path := filepath.Join(dir, name+".todo")
					_ = storage.SaveTodos(path, []models.Todo{})
					refreshFileList(dir)
				}
			})

		case rune(key.DelFile[0]): // Delete selected file
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

		case rune(key.OpenFile[0]): // Open selected file
			index := fileListView.GetCurrentItem()
			if index >= 0 && index < fileListView.GetItemCount() {
				name, _ := fileListView.GetItemText(index)
				currentFile = filepath.Join(dir, name)
				refreshTodoList()
				setFocus(todoListView)
			}

		case rune(key.EditFile[0]): // Edit file name
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
			
		case rune(key.Quit[0]): // Quit application
			app.Stop()
		}
		return event
	})

	todoListView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case rune(key.AddTodo[0]): // Add new todo
			showInputDialog("New Todo", "Enter todo text:", "", func(text string) {
				if text != "" {
					todos, _ := storage.LoadTodos(currentFile)
					todos = append(todos, models.Todo{Text: text, Done: false})
					_ = storage.SaveTodos(currentFile, todos)
					refreshTodoList()
				}
			})

		case rune(key.DelTodo[0]): // Delete selected todo
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

		case rune(key.Toggle[0]): // Toggle done
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
			
		case rune(key.Back[0]): // Go back to file list
			setFocus(fileListView)
			
		case rune(key.EditTodo[0]): // Edit todo text
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
			
		case rune(key.Quit[0]): // Quit application
			app.Stop()
		}
		return event
	})
}