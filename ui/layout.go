package ui

import (
	"path/filepath"

	"github.com/rivo/tview"
)

var app *tview.Application
var todoListView *tview.List
var fileListView *tview.List
var pages *tview.Pages
var currentFile string

func StartApp(dir string) {
	app = tview.NewApplication()
	pages = tview.NewPages()

	fileListView = tview.NewList()
	fileListView.SetTitle(" Todo Files (n:new e:edit d:del o:open q:quit) ").SetBorder(true)

	todoListView = tview.NewList()
	todoListView.SetTitle(" Todos (a:add e:edit d:del t:toggle b:back q:quit) ").SetBorder(true)

	refreshFileList(dir)

	fileListView.SetSelectedFunc(func(index int, name string, secondary string, shortcut rune) {
		currentFile = filepath.Join(dir, name)
		refreshTodoList()
		app.SetFocus(todoListView)
	})

	flex := tview.NewFlex().
		AddItem(fileListView, 30, 1, true).
		AddItem(todoListView, 0, 2, false)

	pages.AddPage("main", flex, true, true)
	app.SetRoot(pages, true).SetFocus(fileListView)

	registerKeybindings(dir)

	if err := app.Run(); err != nil {
		panic(err)
	}
}