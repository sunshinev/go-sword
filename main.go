package main

import (
	"go-sword/config"
	"go-sword/core"
)

func main() {

	//f := untils.FileCopy{
	//	Dir:  make(chan *untils.Params, 100),
	//	File: make(chan *untils.Params, 100),
	//}
	//
	//err := f.Run("resource/", "go-sword-app/resource/")
	//
	//log.Fatalf("%v", err)

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
