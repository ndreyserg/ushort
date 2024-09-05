package config

import (
	"flag"
)

type Config struct {
	ServerAddr string
	UrlAddr    string
}

func MakeConfig() Config {
	con := Config{}
	flag.StringVar(&con.ServerAddr, "h", "localhost:8080", "server addr")
	flag.StringVar(&con.UrlAddr, "b", "localhost:8080", "result server addr")
	flag.Parse()
	return con
}
