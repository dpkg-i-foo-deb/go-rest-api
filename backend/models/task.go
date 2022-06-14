package models

type Task struct {
	Code        int    `json:"code"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
	MainTask    int    `json:"main_task"`
	StartDate   string `json:"start_date"`
	DueDate     string `json:"due_date"`
	User        string `json:"user"`
}
