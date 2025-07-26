package models

type Todo struct {
	Text string `json:"text"`
	Done bool   `json:"done"`
}