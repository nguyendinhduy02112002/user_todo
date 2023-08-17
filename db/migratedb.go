package db

import (
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/v7/esutil"
	"github.com/nguyendinhduy02112002/user_todo/config"
)

func MigrateDB() {
	_, err := config.MI.DB.Exec("USE todos")
	if err != nil {
		panic(err.Error())
	}

	rows, err := config.MI.DB.Query("SELECT id, username, password FROM users")
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var username string
		var password string
		err := rows.Scan(&id, &username, &password)
		if err != nil {
			log.Println(err)
			continue
		}

		// Send data to Elasticsearch
		doc := map[string]interface{}{
			"id":       id,
			"username": username,
			"password": password,
		}

		_, err = config.EI.Client.Index("users", esutil.NewJSONReader(&doc), config.EI.Client.Index.WithDocumentID(fmt.Sprintf("%d", id)))
		if err != nil {
			log.Println(err)
		}
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
}
