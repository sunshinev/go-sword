package main

import (
	"flag"

	"github.com/sunshinev/go-sword/config"
	"github.com/sunshinev/go-sword/core"
)

var dbHost = flag.String("host", "localhost", "MySQL Host")
var dbUser = flag.String("user", "", "MySQL user")
var dbPassword = flag.String("password", "", "MySQL password")
var dbDatabase = flag.String("db", "", "MySQL database ")
var dbPort = flag.Int("port", 3306, "MySQL port")
var serverPort = flag.String("p", "8080", "Go-sword Server port")

func main() {

	flag.Parse()
	c := core.Default()

	c.SetConfig(&config.Config{
		ServerPort: *serverPort,
		Database: &config.DbSet{
			Host:     *dbHost,
			User:     *dbUser,
			Password: *dbPassword,
			Port:     *dbPort,
			Database: *dbDatabase,
		},
		ModuleName: "github.com/sunshinev/go-sword",
	})

	c.Run()
}
