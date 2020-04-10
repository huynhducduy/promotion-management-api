package db

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

var db *sql.DB

func openConnection() {

	conn, err := sql.Open("mysql", config.DB_USER+":"+config.DB_PASS+"@tcp("+config.DB_HOST+":"+config.DB_PORT+")/"+config.DB_NAME+"?parseTime=true")

	if err != nil {
		log.Fatalf("Cannot open connection, %s", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Cannot ping connection, %s", err)
	}

	conn.SetMaxOpenConns(25)
	conn.SetMaxIdleConns(25)
	// The sum is the maximum number of concurrent connections
	conn.SetConnMaxLifetime(5 * time.Minute)

	db = conn
}

func getConnection() *sql.DB {
	return db
}
