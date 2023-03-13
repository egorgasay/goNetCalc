package config

import (
	"flag"
	"os"
)

const (
	defaultHost = ":80"
)

type Flag struct {
	host *string
}

var f Flag

func init() {
	f.host = flag.String("a", defaultHost, "-a=host")
}

type Config struct {
	Host string
}

func New() *Config {
	flag.Parse()

	if addr, ok := os.LookupEnv("RUN_ADDRESS"); ok {
		f.host = &addr
	}

	return &Config{
		Host: *f.host,
	}
}
