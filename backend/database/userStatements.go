package database

import (
	"database/sql"
	"log"
)

var err error
var LoginStatement *sql.Stmt
var SignUpStatement *sql.Stmt

func InitUserStatements() {
	LoginStatement, err = Database.Prepare("SELECT u.email, u.password FROM \"user\" u WHERE u.email = $1")

	if err != nil {
		log.Fatal("Couldn't initialize login statements ", err)
	}

	SignUpStatement, err = Database.Prepare(
		"INSERT INTO \"user\" (email,\"password\", first_name, last_name) VALUES ($1,$2,$3,$4) RETURNING email")

	if err != nil {
		log.Fatal("Couldn't initialize sign up statements ", err)
	}
}
