package elastic

import (
	"github.com/elastic/go-elasticsearch/v8"
	"log"
)

func NewClient(host, username, password string, secure bool) *elasticsearch.Client {
	if host == "" {
		log.Fatalf("Failed to create Elasticsearch client: host is required")
	}

	cfg := elasticsearch.Config{
		Addresses: []string{host},
	}

	if username != "" && password != "" {
		cfg.Username = username
		cfg.Password = password
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to create Elasticsearch client: %v", err)
	}
	return client
}
