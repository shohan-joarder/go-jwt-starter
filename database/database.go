package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var DB *sql.DB

func InitDB() {
	godotenv.Load()
	db_user := os.Getenv("DB_USER_NAME")
	db_pass := os.Getenv("DB_PASSWORD")
	db_host := os.Getenv("DB_HOST")
	db_port := os.Getenv("DB_PORT")
	db_name := os.Getenv("DB_DATABASE_NAME")
	db_driver := os.Getenv("DB_DRIVER")
	var err error
	DB, err = sql.Open(db_driver, db_user+":"+db_pass+"@tcp("+db_host+":"+db_port+")/"+db_name)
	if err != nil {
		log.Fatal("Failed to connect 1")
	}
	err = DB.Ping()
	if err != nil {
		log.Fatal("Failed to connect 2")
	}
	// log.Fatal("connect ")
	fmt.Println("Connected to the database!")
}
