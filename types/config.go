package types

type Keymap struct {
	NewFile  string `json:"newFile"`
	EditFile string `json:"editFile"`
	DelFile  string `json:"delFile"`
	OpenFile string `json:"openFile"`
	AddTodo  string `json:"addTodo"`
	EditTodo string `json:"editTodo"`
	DelTodo  string `json:"delTodo"`
	Toggle   string `json:"toggle"`
	Back     string `json:"back"`
	Quit     string `json:"quit"`
}

type Config struct {
	UndoneColor         string `json:"undoneColor"`
	DoneColor           string `json:"doneColor"`
	ActiveWindowColor   string `json:"activeWindowColor"`
	UnactiveWindowColor string `json:"unactiveWindowColor"`
	Keymap              Keymap `json:"keymap"`
}