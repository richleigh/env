package main

import (
	"fmt"
	"os"

	"github.com/richleigh/env"
)

type config struct {
	Home         string `env:"HOME"`
	Port         int    `env:"PORT"`
	IsProduction bool   `env:"PRODUCTION,optional"`
}

func main() {
	os.Setenv("HOME", "/tmp/fakehome")
	cfg := config{}
	err := env.Parse(&cfg)
	if err == nil {
		fmt.Println(cfg)
	} else {
		fmt.Println(err)
	}
}
