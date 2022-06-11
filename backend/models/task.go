package models

type task struct {
	Code        int    `json:"code"`
	Name        string `json:"name"`
	Tile        string `json:"title"`
	Description string `json:"description"`
	Subtask     int    `json:"subtask"`
	User        string `json:"user"`
}
