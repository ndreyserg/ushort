package config

import (
	"flag"
	"os"
)

type Config struct {
	ServerAddr string
	BaseURL    string
	LogLevel   string
}

func MakeConfig() Config {
	con := Config{}
	flag.StringVar(&con.ServerAddr, "a", "localhost:8080", "server address")
	flag.StringVar(&con.BaseURL, "b", "http://localhost:8080", "result base url")
	flag.StringVar(&con.LogLevel, "l", "info", "log level")
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

	return con
}
