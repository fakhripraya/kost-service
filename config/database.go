package config

import (
	"fmt"
	"os"

	"github.com/fakhripraya/kost-service/entities"

	"gorm.io/gorm"
)

// DB is an ORM for MYSQL database
var DB *gorm.DB

// DBConfig represents db configuration
type DBConfig struct {
	Host     string
	Port     int
	User     string
	DBName   string
	Password string
}

// BuildDBConfig is a function that builds the database config based on DBConfig structure
func BuildDBConfig(db *entities.DatabaseConfiguration) *DBConfig {
	dbConfig := DBConfig{
		Host:     db.Host,
		Port:     db.Port,
		User:     db.User,
		Password: db.Password,
		DBName:   db.Dbname,
	}
	return &dbConfig
}

// DbURL is a function that returns the connected db DSN
func DbURL(dbConfig *DBConfig, config *entities.Configuration) string {

	var DbURLString string
	dev := "development"

	if config.API.Environment == dev {
		DbURLString = fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?tls=true&charset=utf8&parseTime=True&loc=Local",
			dbConfig.User,
			dbConfig.Password,
			dbConfig.Host,
			dbConfig.Port,
			dbConfig.DBName,
		)
	} else {
		DbURLString = os.Getenv("DB_CONNECTION_STRING")
	}

	return DbURLString
}
