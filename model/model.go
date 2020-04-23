package model

import (
	"go-sword/config"
	"log"

	"github.com/Shelnutt2/db2struct"
	_ "github.com/go-sql-driver/mysql"
)

func Create(c *config.Db, table string) {

	// Use db2struct (https://github.com/Shelnutt2/db2struct)
	columnDataTypes, err := db2struct.GetColumnsFromMysqlTable(c.User, c.Password, c.Host, c.Port, c.Database, table)
	if err != nil {
		panic(err.Error())
	}

	struc, err := db2struct.Generate(*columnDataTypes, table, table, table, true, true, false)
	if err != nil {
		panic(err.Error())
	}

	log.Printf("%s", struc)
}
