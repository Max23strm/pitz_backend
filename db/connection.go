package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var DB *sql.DB

func DBconnection() error {
	errorVariables := godotenv.Load()
	if errorVariables != nil {
		panic(errorVariables)
	}

	conection, err := sql.Open("mysql", os.Getenv("DB_USER")+":"+os.Getenv("DB_PASSWORD")+"@tcp("+os.Getenv("DB_SERVER")+":"+os.Getenv("DB_PORT")+")/"+os.Getenv("DB_NAME")+"?parseTime=true")
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
