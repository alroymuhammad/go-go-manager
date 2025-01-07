package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type DBConfig struct {
	Driver   string
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

var DB *sql.DB

func ConnectDB() *sql.DB {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}


	dbConfig := DBConfig{
		Driver:   os.Getenv("DB_DRIVER"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSL_MODE"),
	}

	if dbConfig.Driver == "" || dbConfig.Host == "" || dbConfig.Port == "" || dbConfig.User == "" ||
		dbConfig.Password == "" || dbConfig.DBName == "" || dbConfig.SSLMode == "" {
		log.Fatal("Missing required database environment variables")
	}

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.User,
		dbConfig.Password,
		dbConfig.DBName,
		dbConfig.SSLMode,
	)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Error pinging database:", err)
	}

	fmt.Println("Successfully connected to database!")

	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(25)
	DB.SetConnMaxLifetime(time.Minute * 5)

	return DB
}
