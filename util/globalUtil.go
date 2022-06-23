package util

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// ================= main db connection ====================

// main postgres db connection
var DB = connDB()

func connDB() *sql.DB {
	godotenv.Load(".env")
	dbinfo := fmt.Sprintf("host=containers-us-west-74.railway.app:7103 user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_USER"), os.Getenv("DB_PWD"), os.Getenv("DB_NAME"))

	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		log.Fatalln(err)
	}

	return db
}

// ================= main db connection ====================
