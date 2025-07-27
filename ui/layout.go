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

func StartApp(dir string, cfg types.Config) {
	config = cfg
	app = tview.NewApplication()
	pages = tview.NewPages()

	fileListView = tview.NewList()
	updateFileListTitle()

	todoListView = tview.NewList()
	updateTodoListTitle()

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

func updateFileListTitle() {
	key := config.Keymap
	title := " Todo Files "
	title += "(" + key.NewFile + ":new "
	title += key.EditFile + ":edit "
	title += key.DelFile + ":del "
	title += key.OpenFile + ":open "
	title += key.Quit + ":quit) "
	fileListView.SetTitle(title).SetBorder(true)
}

func updateTodoListTitle() {
	key := config.Keymap
	title := " Todos "
	title += "(" + key.AddTodo + ":add "
	title += key.EditTodo + ":edit "
	title += key.DelTodo + ":del "
	title += key.Toggle + ":toggle "
	title += key.Back + ":back "
	title += key.Quit + ":quit) "
	todoListView.SetTitle(title).SetBorder(true)
}

func setFocus(primitive tview.Primitive) {
	activeColor := tcell.GetColor(config.ActiveWindowColor)
	unactiveColor := tcell.GetColor(config.UnactiveWindowColor)

	if primitive == fileListView {
		fileListView.SetBorderColor(activeColor)
		todoListView.SetBorderColor(unactiveColor)
	} else if primitive == todoListView {
		todoListView.SetBorderColor(activeColor)
		fileListView.SetBorderColor(unactiveColor)
	}

	app.SetFocus(primitive)
}