package config

import (
	"flag"
	"os"
)

type Config struct {
	ServerAddr  string
	BaseURL     string
	LogLevel    string
	StoragePath string
	DSN         string
	Secret      string
}

func MakeConfig() Config {
	con := Config{}
	flag.StringVar(&con.ServerAddr, "a", "localhost:8080", "server address")
	flag.StringVar(&con.BaseURL, "b", "http://localhost:8080", "result base url")
	flag.StringVar(&con.LogLevel, "l", "info", "log level")
	flag.StringVar(&con.StoragePath, "f", "", "storage file path")
	flag.StringVar(&con.DSN, "d", "", "DSN")
	flag.StringVar(&con.Secret, "s", "secret_key", "secret key")
	flag.Parse()

	if envServerAddr := os.Getenv("SERVER_ADDRESS"); envServerAddr != "" {
		con.ServerAddr = envServerAddr
	}

	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		con.BaseURL = envBaseURL
	}

	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		con.LogLevel = envLogLevel
	}

	if envStoragePath := os.Getenv("FILE_STORAGE_PATH"); envStoragePath != "" {
		con.StoragePath = envStoragePath
	}

	if envDSN := os.Getenv("DATABASE_DSN"); envDSN != "" {
		con.DSN = envDSN
	}

	if envSecret := os.Getenv("SECRET_KEY"); envSecret != "" {
		con.Secret = envSecret
	}

	return con
}
