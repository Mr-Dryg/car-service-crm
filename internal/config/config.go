package config

import (
	"fmt"
	"net/url"
	"os"
)

var ErrIncompleteDatabaseEnv = fmt.Errorf("database environment variables are incomplete")

type Config struct {
	DatabaseURL string
}

func Load() (Config, error) {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")
	if sslmode == "" {
		sslmode = "disable"
	}

	if user == "" || password == "" || host == "" || port == "" || name == "" {
		return Config{}, ErrIncompleteDatabaseEnv
	}

	databaseURL := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(user, password),
		Host:   host + ":" + port,
		Path:   name,
	}

	query := databaseURL.Query()
	query.Set("sslmode", sslmode)
	databaseURL.RawQuery = query.Encode()

	return Config{
		DatabaseURL: databaseURL.String(),
	}, nil
}
