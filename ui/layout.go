package ui

import (
	"path/filepath"
	"time"

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

func splashScreen(dir string) *tview.Flex {
	asciiArt := 
`
     ,o888888o.        ,o888888o.     8 888888888o.          ,o888888o.     
    8888     '88.   . 8888     '88.   8 8888    '^888.    . 8888     '88.   
 ,8 8888       '8. ,8 8888       '8b  8 8888        '88. ,8 8888       '8b  
 88 8888           88 8888        '8b 8 8888         '88 88 8888        '8b 
 88 8888           88 8888         88 8 8888          88 88 8888         88 
 88 8888           88 8888         88 8 8888          88 88 8888         88 
 88 8888   8888888 88 8888        ,8P 8 8888         ,88 88 8888        ,8P 
 '8 8888       .8' '8 8888       ,8P  8 8888        ,88' '8 8888       ,8P  
    8888     ,88'   ' 8888     ,88'   8 8888    ,o88P'    ' 8888     ,88'   
     '8888888P'        '8888888P'     8 888888888P'          '8888888P'     
`

	splashText := tview.NewTextView().
		SetText(asciiArt).
		SetTextAlign(tview.AlignCenter)

	message := tview.NewTextView().
		SetText("Press Enter to continue").
		SetTextAlign(tview.AlignCenter)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(splashText, 0, 1, false).
		AddItem(message, 1, 1, false).
		AddItem(nil, 0, 1, false)

	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			initializeMainUI(dir)
			pages.RemovePage("splash")
			return nil
		}
		return nil
	})

	return flex
}

func showErrorDialog(msg string) {
	modal := tview.NewModal().
		SetText(msg).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(_ int, _ string) {
			pages.RemovePage("errorDialog")
		})
	pages.AddPage("errorDialog", modal, true, true)
}

func showInfoDialogAndExit(msg string) {
	modal := tview.NewModal().
		SetText(msg).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(_ int, _ string) {
			app.Stop()
		})
	pages.AddPage("infoDialog", modal, true, true)
}

func showUpdateDialog() {
	modal := tview.NewModal().
		SetText("An update is available. Do you want to update now?").
		AddButtons([]string{"Yes", "Later"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			pages.RemovePage("updateDialog")
			if buttonLabel == "Yes" {
				updatingModal := tview.NewModal().
					SetText("Updating... Please wait.")
				pages.AddPage("updating", updatingModal, true, true)

				go func() {
					err := updateApplication()
					app.QueueUpdateDraw(func() {
						pages.RemovePage("updating")
						if err != nil {
							showErrorDialog("Update failed: " + err.Error())
						} else {
							showInfoDialogAndExit("✓ Update successful. Restart the application.")
						}
					})
				}()
			}
		})
	pages.AddPage("updateDialog", modal, true, true)
}

func StartApp(dir string, cfg types.Config) {
	config = cfg
	app = tview.NewApplication()
	pages = tview.NewPages()

	// Check for updates in background
	go func() {
		// Wait for UI to initialize
		time.Sleep(500 * time.Millisecond)
		
		updateAvailable, err := checkForUpdate()
		if err != nil {
			return // Silently ignore errors
		}
		if updateAvailable {
			app.QueueUpdateDraw(func() {
				showUpdateDialog()
			})
		}
	}()

	// Show splash
	pages.AddPage("splash", splashScreen(dir), true, true)
	app.SetRoot(pages, true)

	if err := app.Run(); err != nil {
		panic(err)
	}
}

func initializeMainUI(dir string) {
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

	pages.AddPage("main", flex, true, false)
	pages.SwitchToPage("main")
	setFocus(fileListView)

	registerKeybindings(dir)
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