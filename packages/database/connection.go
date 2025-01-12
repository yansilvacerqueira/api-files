package database

import (
	"database/sql"
	"fmt"
	"os"
)

func NewConnection() (*sql.DB, error) {
	psqlconn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("HOST_DB"),
		os.Getenv("PORT_DB"),
		os.Getenv("USER_DB"),
		os.Getenv("PASSWORD_DB"),
		os.Getenv("NAME_DB"),
	)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
