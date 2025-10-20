package config

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDatabase() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL belum diset di env")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("gagal sql.Open:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("gagal ping db:", err)
	}

	DB = db
	log.Println("✅ connected to DB")
}
