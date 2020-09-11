package main

import (
	"fmt"

	"github.com/sunshinev/go-sword/core"
)

func main() {

	fmt.Print(`
+---------------------------------------------------+
|                                                   |
|            Welcome to use Go-Sword                |
|                                                   |
|                Visualized tool                    |
|        Fastest to create CRUD background          |
|      https://github.com/sunshinev/go-sword        |
|                                                   |
+---------------------------------------------------+
`)

	core.Init().Run()
}
