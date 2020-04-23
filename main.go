package main

import (
	"go-sword/config"
	"go-sword/engine"
)

func main() {

	eng := engine.Default()

	eng.SetConfig(&config.Config{
		Database: &config.Db{
			Host:     "localhost",
			User:     "root",
			Password: "123456",
			Port:     3306,
			Database: "goadmin",
		},
	})

	eng.Run()
}
