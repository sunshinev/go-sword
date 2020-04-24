package main

import (
	"go-sword/gosword/core"
	"log"
)

func main() {

	e := core.Engine{}

	err := e.Run()
	if err != nil {
		log.Fatalf("%v", err)
	}
}
