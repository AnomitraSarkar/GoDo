package ui

import (
	"path/filepath"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"todo-cli/types"
)

var app *tview.Application
var todoListView *tview.List
var fileListView *tview.List
var pages *tview.Pages
var currentFile string
var config types.Config
var activeList *tview.List

func StartApp(dir string, cfg types.Config) {
	config = cfg
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
		setFocus(todoListView)
	})

	flex := tview.NewFlex().
		AddItem(fileListView, 30, 1, true).
		AddItem(todoListView, 0, 2, false)

	pages.AddPage("main", flex, true, true)
	setFocus(fileListView)
	app.SetRoot(pages, true)

	registerKeybindings(dir)

	if err := app.Run(); err != nil {
		panic(err)
	}
}

func setFocus(primitive tview.Primitive) {
	if activeList != nil {
		activeList.SetBorderColor(tview.Styles.BorderColor)
	}

	app.SetFocus(primitive)

	if primitive == fileListView {
		activeList = fileListView
	} else if primitive == todoListView {
		activeList = todoListView
	}

	if activeList != nil {
		color := tcell.GetColor(config.ActiveWindowColor)
		activeList.SetBorderColor(color)
	}
}