package types

type Keymap struct {
	NewFile      string `json:"newFile"`
	EditFile     string `json:"editFile"`
	DelFile      string `json:"delFile"`
	OpenFile     string `json:"openFile"`
	AddTodo      string `json:"addTodo"`
	EditTodo     string `json:"editTodo"`
	DelTodo      string `json:"delTodo"`
	Toggle       string `json:"toggle"`
	Back         string `json:"back"`
	Quit         string `json:"quit"`
	MoveToTop    string `json:"moveToTop"`    // gg
	MoveToBottom string `json:"moveToBottom"` // G
	SetPriority  string `json:"setPriority"`  // p
}

type Config struct {
	UndoneColor         string `json:"undoneColor"`
	DoneColor           string `json:"doneColor"`
	ActiveWindowColor   string `json:"activeWindowColor"`
	UnactiveWindowColor string `json:"unactiveWindowColor"`
	PriorityColor       string `json:"priorityColor"` // New
	RelativeTime        bool   `json:"relativeTime"`  // New
	Keymap              Keymap `json:"keymap"`
}