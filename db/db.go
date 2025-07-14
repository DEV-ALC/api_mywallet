package db

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Connect() {
	var err error
	dsn := "root:h7bd@tcp(127.0.0.1:3306)/carteira?charset=utf8&parseTime=true&loc=Local"
	// dsn := "user:senha@tcp(127.0.0.1:3306)/nome"
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("✅ Banco conectado com sucesso!")
}
