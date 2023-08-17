package config

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
	_ "go.elastic.co/apm"
	"go.elastic.co/apm/module/apmsql"
)

const (
	username = "root"
	password = "root"
	hostname = "localhost:3306"
	dbname   = "todos"
)

type MySQLInstance struct {
	DB *sql.DB
}

var MI MySQLInstance

func dsn(dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbName)
}

func ConnectDB() {
	apmDriver := apmsql.Wrap(&mysql.MySQLDriver{})
	apmsql.Register("mysql", apmDriver)
	db, err := apmsql.Open("mysql", dsn(""))
	if err != nil {
		fmt.Printf("error set agent: %s", err)
		panic(err.Error())
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MySQL")

	MI = MySQLInstance{
		DB: db,
	}
}
