package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var DB *sql.DB

func DBconnection() {
	errorVariables := godotenv.Load()
	if errorVariables != nil {
		panic(errorVariables)
	}

	conection, err := sql.Open("mysql", os.Getenv("DB_USER")+":"+os.Getenv("DB_PASSWORD")+"@tcp("+os.Getenv("DB_SERVER")+":"+os.Getenv("DB_PORT")+")/"+os.Getenv("DB_NAME")+"?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}

	DB = conection
}
func CerrarConexion() {
	DB.Close()
}
