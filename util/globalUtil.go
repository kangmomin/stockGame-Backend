package util

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type res struct {
	Data any
	Err  bool
}

// ================= main db connection ====================

// main postgres db connection
var DB = connDB()

func connDB() *sql.DB {
	godotenv.Load(".env")
	dbinfo := `postgresql://` + os.Getenv("DB_USER") + `:` + os.Getenv("DB_PWD") + `@containers-us-west-72.railway.app:6455/railway`

	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		log.Fatalln(err)
	}

	return db
}

// ================= main db connection ====================
