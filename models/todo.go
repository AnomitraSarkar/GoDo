package models

import "time"

type Todo struct {
	Text      string    `json:"text"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"createdAt"`
	Priority  bool      `json:"priority"`
}