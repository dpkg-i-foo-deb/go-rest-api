package database

import (
	"database/sql"
	"log"
)

var CreateTaskStatement *sql.Stmt
var GetTaskStatement *sql.Stmt
var GetAllTasksStatement *sql.Stmt
var EditTaskStatement *sql.Stmt
var DeleteTaskStatement *sql.Stmt

func InitTaskStatements() {
	CreateTaskStatement, err = Database.Prepare(`INSERT INTO public.task 
		(title,description,"user",start_date,due_date,status,main_task)
		VALUES ($1, $2, $3, $4, $5, $6,$7)
		RETURNING 
		title, description , "user",start_date ,due_date ,status ,main_task ,code `)

	if err != nil {
		log.Fatal("Couldn't initialize task statements ", err)
	}

	GetTaskStatement, err = Database.Prepare(`SELECT title, description, code, main_task, "user", start_date, due_date, status
											FROM public.task t WHERE t.code=$1 AND t.user =$2`)

	if err != nil {
		log.Fatal("Couldn't initialize task statements ", err)
	}

	GetAllTasksStatement, err = Database.Prepare(`SELECT title, description, code, main_task, "user", start_date, due_date, status
													FROM public.task t WHERE t.user = $1`)

	if err != nil {
		log.Fatal("Couldn't initialize task statements ", err)
	}

	EditTaskStatement, err = Database.Prepare(`UPDATE public.task
												SET title=$1, description=$2, main_task=$3, start_date=$4, due_date=$5, status=$6
												WHERE code=$7 AND "user"=$8
												RETURNING title, description , main_task , start_date , due_date , status ,code 
											`)

	if err != nil {
		log.Fatal("Couldn't initialize task statements ", err)
	}

	DeleteTaskStatement, err = Database.Prepare(`DELETE FROM task WHERE code = $1 AND "user" = $2`)

	if err != nil {
		log.Fatal("Couldn't initialize task statements ", err)
	}
}
