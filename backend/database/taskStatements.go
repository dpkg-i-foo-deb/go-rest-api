package database

import (
	"database/sql"
	"log"
)

var createTaskStatement *sql.Stmt

func InitTaskStatements() {
	createTaskStatement, err = Database.Prepare(`INSERT INTO public.task 
		(title,description,"user",start_date,due_date,status)
		VALUES ($1, $2, $3, $4, $5, $6)`)

	if err != nil {
		log.Fatal("Couldn't initialize task statements ", err)
	}
}
