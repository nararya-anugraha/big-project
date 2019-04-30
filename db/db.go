package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DatabaseConfigType struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func GetDB(databaseConfig *DatabaseConfigType) (*sqlx.DB, error) {
	dataSourceName := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		databaseConfig.Host, databaseConfig.Port, databaseConfig.User, databaseConfig.Password, databaseConfig.DBName)
	return sqlx.Connect("postgres", dataSourceName)
}
