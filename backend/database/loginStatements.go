package database

import (
	"database/sql"
	"log"
)

var err error
var LoginStatement *sql.Stmt

func InitLoginStatements() {
	LoginStatement, err = Database.Prepare("SELECT u.email FROM \"user\" u WHERE u.password = $1")

	if err != nil {
		log.Fatal("Couldn't initialize login statements ", err)
	}
}
