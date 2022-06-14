package database

import (
	"database/sql"
	"log"
)

var CreateTaskStatement *sql.Stmt

func InitTaskStatements() {
	CreateTaskStatement, err = Database.Prepare(`INSERT INTO public.task 
		(title,description,"user",start_date,due_date,status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING 
		title, description , "user",start_date ,due_date ,status ,main_task ,code `)

	if err != nil {
		log.Fatal("Couldn't initialize task statements ", err)
	}
}
