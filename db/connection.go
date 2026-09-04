package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func DBconnection() error {

	if os.Getenv("RAILWAY_ENVIRONMENT") == "" {
		fmt.Println("Is not railway")
		errorVariables := godotenv.Load()
		if errorVariables != nil {
			fmt.Println(errorVariables)
			panic(errorVariables)
		}
	}

	sslmode := os.Getenv("DB_SSLMODE")
	if sslmode == "" {
		sslmode = "disable"
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_SERVER"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		sslmode,
	)

	conection, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}

	if err := conection.Ping(); err != nil {
		return err
	}

	DB = conection
	fmt.Println("✅ Database connection established")
	return nil
}
func CerrarConexion() {
	if DB != nil {
		_ = DB.Close()
	}
}
