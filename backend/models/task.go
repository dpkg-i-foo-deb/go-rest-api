package models

type Task struct {
	Code        *int    `json:"code"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	MainTask    *int    `json:"main_task"`
	StartDate   string  `json:"start_date"`
	DueDate     string  `json:"due_date"`
	Status      bool    `json:"status"`
	UserEmail   *string `json:"user_email"`
}
