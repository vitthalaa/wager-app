package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"github.com/vitthalaa/wager-app/internal/config"
)

func OpenConnection(conf *config.DataBaseConfig) (dbConn *sql.DB, err error) {
	connStr := fmt.Sprintf(
		"port=%d host=%s user=%s password=%s dbname=%s sslmode=disable",
		conf.DBPort, conf.DBHost, conf.DBUser, conf.DBPass, conf.DBName)

	dbConn, err = sql.Open("postgres", connStr)
	if err == nil {
		log.Print("Connection opened to DB: " + conf.DBName)
	} else {
		log.Printf("DB connection failed: %v", err)
		return
	}

	dbConn.SetMaxOpenConns(conf.DBMaxOpenConn)
	dbConn.SetMaxIdleConns(conf.DBMaxIdleConn)

	return
}
