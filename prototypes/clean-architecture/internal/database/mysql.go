package database

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/humanbeeng/lepo/prototypes/clean-architecture/internal/config"
)

func BootstrapMySQL() (*sql.DB, error) {
	appConfig, err := config.GetAppConfig()
	if err != nil {
		return nil, err
	}

	mysqlConfig := mysql.Config{
		User:                 appConfig.DBConfig.Username,
		Passwd:               appConfig.DBConfig.Password,
		DBName:               appConfig.DBConfig.DatabaseName,
		Net:                  appConfig.DBConfig.Protocol,
		Addr:                 appConfig.DBConfig.Server,
		AllowNativePasswords: true,
		TLSConfig:            "true",
	}

	db, err := sql.Open("mysql", mysqlConfig.FormatDSN())
	if err != nil {
		return nil, err
	}
	return db, nil
}

func CloseMySQLConnection(db *sql.DB) error {
	log.Println("Closing MySQL Connection")
	return db.Close()
}
