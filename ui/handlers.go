package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"todo-cli/models"
	"todo-cli/storage"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func firstRune(s string) rune {
	if s == "" {
		return 0
	}
	for _, r := range s {
		return r
	}
	return 0
}

func refreshFileList(dir string) {
	fileListView.Clear()
	todoFiles, err := storage.ListTodoFiles(dir)
	if err != nil {
		return
	}
	for _, file := range todoFiles {
		path := filepath.Join(dir, file)
		fileData, err := storage.LoadTodos(path)
		if err != nil {
			continue
		}
		
		text := file
		if fileData.Priority {
			text = fmt.Sprintf("[%s]★ %s[-]", config.PriorityColor, file)
		}
		timeStr := formatTime(fileData.CreatedAt, config.RelativeTime)
		fileListView.AddItem(text, timeStr, 0, nil)
	}
}

func refreshTodoList() {
	todoListView.Clear()
	if currentFile == "" {
		return
	}
	fileData, err := storage.LoadTodos(currentFile)
	if err != nil {
		return
	}
	for _, t := range fileData.Todos {
		prefix := "❌"
		color := config.UndoneColor
		if t.Done {
			prefix = "✅"
			color = config.DoneColor
		}
		if t.Priority {
			color = config.PriorityColor
		}
		
		timeStr := formatTime(t.CreatedAt, config.RelativeTime)
		text := fmt.Sprintf("[%s]%s %s[-] (%s)", color, prefix, t.Text, timeStr)
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
		case firstRune(key.NewFile): // Create new file
			showInputDialog("New File", "File name (without .todo):", "", func(name string) {
				if name != "" {
					name = strings.TrimSuffix(name, ".todo")
					path := filepath.Join(dir, name+".todo")
					fileData := storage.TodoFileData{
						CreatedAt: time.Now(),
						Todos:     []models.Todo{},
						Priority:  false,
					}
					_ = storage.SaveTodos(path, fileData)
					refreshFileList(dir)
				}
			})

		case firstRune(key.DelFile): // Delete selected file
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

		case firstRune(key.OpenFile): // Open selected file
			index := fileListView.GetCurrentItem()
			if index >= 0 && index < fileListView.GetItemCount() {
				name, _ := fileListView.GetItemText(index)
				currentFile = filepath.Join(dir, name)
				refreshTodoList()
				setFocus(todoListView)
			}

		case firstRune(key.EditFile): // Edit file name
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
			
		case firstRune(key.Quit): // Quit application
			app.Stop()

		case firstRune(key.SetPriority): // Toggle file priority
			index := fileListView.GetCurrentItem()
			if index >= 0 && index < fileListView.GetItemCount() {
				name, _ := fileListView.GetItemText(index)
				path := filepath.Join(dir, name)
				fileData, err := storage.LoadTodos(path)
				if err != nil {
					return event
				}
				fileData.Priority = !fileData.Priority
				_ = storage.SaveTodos(path, fileData)
				refreshFileList(dir)
			}
			
		case firstRune(key.MoveToTop): // gg - move to top
			fileListView.SetCurrentItem(0)
		case firstRune(key.MoveToBottom): // G - move to bottom
			fileListView.SetCurrentItem(fileListView.GetItemCount() - 1)
		}
		return event
	})

	todoListView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case firstRune(key.AddTodo): // Add new todo
			showInputDialog("New Todo", "Enter todo text:", "", func(text string) {
				if text != "" {
					fileData, _ := storage.LoadTodos(currentFile)
					fileData.Todos = append(fileData.Todos, models.Todo{
						Text:      text,
						Done:      false,
						CreatedAt: time.Now(),
						Priority:  false,
					})
					_ = storage.SaveTodos(currentFile, fileData)
					refreshTodoList()
				}
			})

		case firstRune(key.DelTodo): // Delete selected todo
			index := todoListView.GetCurrentItem()
			fileData, err := storage.LoadTodos(currentFile)
			if err != nil {
				return event
			}
			if index >= 0 && index < len(fileData.Todos) {
				fileData.Todos = append(fileData.Todos[:index], fileData.Todos[index+1:]...)
				_ = storage.SaveTodos(currentFile, fileData)
			}
			refreshTodoList()

		case firstRune(key.Toggle): // Toggle done
			index := todoListView.GetCurrentItem()
			fileData, err := storage.LoadTodos(currentFile)
			if err != nil {
				return event
			}
			if index >= 0 && index < len(fileData.Todos) {
				fileData.Todos[index].Done = !fileData.Todos[index].Done
				_ = storage.SaveTodos(currentFile, fileData)
			}
			refreshTodoList()
			
		case firstRune(key.Back): // Go back to file list
			setFocus(fileListView)
			
		case firstRune(key.EditTodo): // Edit todo text
			index := todoListView.GetCurrentItem()
			fileData, err := storage.LoadTodos(currentFile)
			if err != nil {
				return event
			}
			if index >= 0 && index < len(fileData.Todos) {
				showInputDialog("Edit Todo", "New text:", fileData.Todos[index].Text, func(newText string) {
					if newText != "" {
						fileData.Todos[index].Text = newText
						_ = storage.SaveTodos(currentFile, fileData)
						refreshTodoList()
					}
				})
			}
			
		case firstRune(key.Quit): // Quit application
			app.Stop()

		case firstRune(key.SetPriority): // Toggle todo priority
			index := todoListView.GetCurrentItem()
			fileData, err := storage.LoadTodos(currentFile)
			if err != nil {
				return event
			}
			if index >= 0 && index < len(fileData.Todos) {
				fileData.Todos[index].Priority = !fileData.Todos[index].Priority
				_ = storage.SaveTodos(currentFile, fileData)
				refreshTodoList()
			}
			
		case firstRune(key.MoveToTop): // gg - move to top
			todoListView.SetCurrentItem(0)
		case firstRune(key.MoveToBottom): // G - move to bottom
			todoListView.SetCurrentItem(todoListView.GetItemCount() - 1)
		}
		return event
	})
}