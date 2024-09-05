package config

import (
	"flag"
)

type Config struct {
	ServerAddr string
	URLAddr    string
}

func MakeConfig() Config {
	con := Config{}
	flag.StringVar(&con.ServerAddr, "a", "localhost:8080", "server addr")
	flag.StringVar(&con.URLAddr, "b", "http://localhost:8080", "result server addr")
	flag.Parse()
	return con
}
