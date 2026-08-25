package clickhouse

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/ClickHouse/clickhouse-go/v2"
)

type Config struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string
	Secure   bool
}

func Connect(cfg Config) (*sql.DB, error) {
	if cfg.Host == "" {
		return nil, errors.New("clickhouse host required")
	}
	if cfg.Port == 0 {
		cfg.Port = 9000
	}
	if cfg.Database == "" {
		return nil, errors.New("clickhouse database required")
	}
	scheme := "tcp"
	if cfg.Secure {
		scheme = "https"
	}
	dsn := fmt.Sprintf("clickhouse://%s:%s@%s:%d/%s?dial_timeout=10s&secure=%t", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.Secure)
	if scheme == "https" {
		dsn += "&skip_verify=true"
	}
	db, err := sql.Open("clickhouse", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}
