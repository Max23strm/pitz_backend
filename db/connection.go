package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var DSN = "host=localhost user=max password=secretpassword dbname=gorm port=5001"

// var DSN = "host=localhost user=postgres password=mysecretpassword dbname=gorm port=5432"

var DB *sql.DB

func DBconnection() {
	errorVariables := godotenv.Load()
	if errorVariables != nil {
		panic(errorVariables)
	}

	conection, err := sql.Open("mysql", os.Getenv("DB_USER")+":"+"@tcp("+os.Getenv("DB_SERVER")+":"+os.Getenv("DB_PORT")+")/"+os.Getenv("DB_NAME"))
	if err != nil {
		log.Fatal(err)
	}

	DB = conection
}
func CerrarConexion() {
	DB.Close()
}

//Concectarnos a la BD
