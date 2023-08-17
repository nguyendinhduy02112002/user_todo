package config

import (
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/v7"
)

type ElasticInstance struct {
	Client *elasticsearch.Client
}

var EI ElasticInstance

func NewElasticsearchClient() {
	client, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	res, err := client.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()
	fmt.Println("connected elasticsearc")
	EI = ElasticInstance{
		Client: client,
	}
}
