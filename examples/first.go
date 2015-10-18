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
	Password     string `env:"PASSWORD,sensitive,optional"`
}

func main() {
	os.Setenv("HOME", "/tmp/fakehome")
	fmt.Printf("PASSWORD before: '%s'\n", os.Getenv("PASSWORD"))
	cfg := config{}
	err := env.Parse(&cfg)
	if err == nil {
		fmt.Println(cfg)
	} else {
		fmt.Println(err)
	}
	fmt.Printf("PASSWORD after: '%s'\n", os.Getenv("PASSWORD"))
}
