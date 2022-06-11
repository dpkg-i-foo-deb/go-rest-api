package database

import (
	"database/sql"
	"log"
)

var SignUpStatement *sql.Stmt

func InitSignUpStatements() {
	SignUpStatement, err = Database.Prepare(
		"INSERT INTO \"user\" (email,\"password\") VALUES ($1,$2) RETURNING email")

	if err != nil {
		log.Fatal("Couldn't initialize sign up statements ", err)
	}
}
