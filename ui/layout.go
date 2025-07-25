package ui

import (
	"path/filepath"

	"github.com/rivo/tview"
)

var app *tview.Application
var todoListView *tview.List
var fileListView *tview.List
var currentFile string

func StartApp(dir string) {
	app = tview.NewApplication()

	fileListView = tview.NewList()
	fileListView.SetTitle(" Todo Files (n:new d:del o:open) ").SetBorder(true)

	todoListView = tview.NewList()
	todoListView.SetTitle(" Todos (a:add d:del t:toggle b:back) ").SetBorder(true)

	refreshFileList(dir)

	fileListView.SetSelectedFunc(func(index int, name string, secondary string, shortcut rune) {
		currentFile = filepath.Join(dir, name)
		refreshTodoList()
		app.SetFocus(todoListView)
	})

	flex := tview.NewFlex().
		AddItem(fileListView, 30, 1, true).
		AddItem(todoListView, 0, 2, false)

	app.SetRoot(flex, true).SetFocus(fileListView)

	registerKeybindings(dir)

	if err := app.Run(); err != nil {
		panic(err)
	}
}