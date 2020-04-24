package main

import (
	"go-sword/config"
	"go-sword/core"
)

func main() {

	c := core.Default()

	c.SetConfig(&config.Config{
		Database: &config.Db{
			Host:     "localhost",
			User:     "root",
			Password: "123456",
			Port:     3306,
			Database: "goadmin",
		},
		ModuleName: "go-sword",
		RootPath:   "go-sword-app",
	})

	c.Run()
}
