package database

import (
	"database/sql"
	"log"
)

var err error
var LoginStatement *sql.Stmt

func InitLoginStatements() {
	LoginStatement, err = Database.Prepare("SELECT u.email, u.password FROM \"user\" u WHERE u.email = $1")

	if err != nil {
		log.Fatal("Couldn't initialize login statements ", err)
	}
}
