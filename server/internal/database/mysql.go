package database

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/lepoai/lepo/server/internal/config"
)

func BootstrapMySQL() (*sql.DB, error) {
	appConfig, err := config.GetAppConfig()
	if err != nil {
		return nil, err
	}

	mysqlConfig := mysql.Config{
		User:                 appConfig.PlanetScaleConfig.Username,
		Passwd:               appConfig.PlanetScaleConfig.Password,
		DBName:               appConfig.PlanetScaleConfig.DatabaseName,
		Net:                  appConfig.PlanetScaleConfig.Protocol,
		Addr:                 appConfig.PlanetScaleConfig.Server,
		AllowNativePasswords: true,
		TLSConfig:            "true",
	}

	db, err := sql.Open("mysql", mysqlConfig.FormatDSN())
	if err != nil {
		return nil, err
	}
	log.Println("info: MySQL connection established")
	return db, nil
}

func CloseMySQLConnection(db *sql.DB) error {
	log.Println("Closing MySQL Connection")
	return db.Close()
}
