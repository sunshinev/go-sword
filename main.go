package main

import (
	"go-sword/config"
	"go-sword/core"
)

func main() {

	eng := core.Default()

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
