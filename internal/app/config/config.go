// Package config с основными настройками приложения.
package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

// Config структура основных настроек.
type Config struct {
	ServerAddr  string `json:"server_address"`    // Адрес сервера.
	BaseURL     string `json:"base_url"`          // Адрес сервера для перенаправления.
	LogLevel    string `json:"info"`              // Уровень логирования.
	StoragePath string `json:"file_storage_path"` // Путь к файлу для хранения в файле.
	DSN         string `json:"database_dsn"`      // Строка подключения к БД.
	Secret      string `json:"secret"`            // Строка с секрета для JWT.
	EnableHTTPS bool   `json:"enable_https"`      // Использовать https.
}

// MakeConfig возвращает структуру с настройками.
func MakeConfig() Config {
	con := getDefaultConfig()
	flag.String("c", "", "config file path")
	flag.String("сonfig", "", "config file path")
	flag.StringVar(&con.ServerAddr, "a", con.ServerAddr, "server address")
	flag.StringVar(&con.BaseURL, "b", con.BaseURL, "result base url")
	flag.StringVar(&con.LogLevel, "l", con.LogLevel, "log level")
	flag.StringVar(&con.StoragePath, "f", con.StoragePath, "storage file path")
	flag.StringVar(&con.DSN, "d", con.DSN, "DSN")
	flag.StringVar(&con.Secret, "k", con.Secret, "secret key")
	flag.BoolVar(&con.EnableHTTPS, "s", con.EnableHTTPS, "Enable HTTPS")
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

	if envHTTPS := os.Getenv("ENABLE_HTTPS"); envHTTPS != "" {
		con.EnableHTTPS = true
	}

	return con
}

func getDefaultConfig() Config {

	conf := Config{
		ServerAddr: "localhost:8080",
		BaseURL:    "http://localhost:8080",
		LogLevel:   "info",
		Secret:     "secret_key",
	}
	confPath := ""
	confPathIndex := 0
	for i, v := range os.Args[1:] {
		if v == "-c" || v == "-config" {
			confPathIndex = i + 1
		}
		if i == confPathIndex {
			confPath = v
		}
	}

	if envPath := os.Getenv("CONFIG"); envPath != "" {
		confPath = envPath
	}

	if confPath == "" {
		return conf
	}

	file, err := os.OpenFile(confPath, os.O_RDONLY, 0666)

	if err != nil {
		fmt.Println("Ошибка открытия файла конфигурации")
		return conf
	}

	dec := json.NewDecoder(file)
	err = dec.Decode(&conf)
	if err != nil {
		fmt.Println("Ошибка декодирования файла конфигурации")
	}
	return conf
}
