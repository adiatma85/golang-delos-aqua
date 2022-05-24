package db

import (
	"fmt"
	"time"

	"github.com/adiatma85/golang-rest-template-api/internal/pkg/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB  *gorm.DB
	err error
)

// Database instance
type Database struct {
	*gorm.DB
}

// SetupDB is a function to open connection to database
func SetupDB() {
	var db = DB

	configuration := config.GetConfig()

	// Viper Config
	driver := configuration.Database.Driver
	database := configuration.Database.Dbname
	username := configuration.Database.Username
	password := configuration.Database.Password
	host := configuration.Database.Host
	port := configuration.Database.Port

	// Gorm config
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}

	switch driver {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, database)
		db, err = gorm.Open(mysql.Open(dsn), gormConfig)
		if err != nil {
			fmt.Println("db err:", err)
		}
	case "postgres":
		dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", host, port, username, database, password)
		db, err = gorm.Open(postgres.Open(dsn), gormConfig)
		if err != nil {
			fmt.Println("db err:", err)
		}
	}
	// Set up the connection pools
	sqlDb, _ := db.DB()
	sqlDb.SetMaxIdleConns(configuration.Database.MaxIdleConns)
	sqlDb.SetMaxOpenConns(configuration.Database.MaxOpenConns)
	sqlDb.SetConnMaxLifetime(time.Duration(configuration.Database.MaxLifetime))

	DB = db
	migration()
}

// Setup for testing database
func SetupTestingDb(host, username, password, port, database string) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("db err for testing :", err)
		panic(err.Error())
	}

	DB = db
	migration()
}

// AutoMigrate project models
func migration() {
}

func GetDB() *gorm.DB {
	return DB
}
