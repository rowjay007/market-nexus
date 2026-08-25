package elasticsearch

import (
	"errors"
	es8 "github.com/elastic/go-elasticsearch/v8"
)

type Config struct {
	Addresses []string
	Username  string
	Password  string
}

func Connect(cfg Config) (*es8.Client, error) {
	if len(cfg.Addresses) == 0 {
		return nil, errors.New("elasticsearch addresses required")
	}
	return es8.NewClient(es8.Config{
		Addresses: cfg.Addresses,
		Username:  cfg.Username,
		Password:  cfg.Password,
	})
}
