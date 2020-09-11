package config

import (
	"database/sql"
	"flag"
	"log"
	"strconv"
	"strings"
)

// App global config
var GlobalConfig *Config

type Config struct {
	DatabaseSet DbSet  // MySQL config
	RootPath    string // The directory go-sword store new file
	ModuleName  string // Project go mod module name
	ServerPort  string // Go-sword server port
	DbConn      *sql.DB
}

type DbSet struct {
	Host     string
	User     string
	Password string
	Port     int
	Database string
}

func (c Config) InitConfig() {

	var dbHost = flag.String("host", "localhost", "MySQL Host")
	var dbUser = flag.String("user", "", "MySQL user")
	var dbPassword = flag.String("password", "", "MySQL password")
	var dbDatabase = flag.String("db", "", "MySQL database ")
	var dbPort = flag.Int("port", 3306, "MySQL port")
	var serverPort = flag.String("p", "8080", "Go-sword Server port")
	var rootPath = flag.String("module", "go-sword-app/", "New project module, the same as  'module' in go.mod file.  ")

	flag.Parse()

	GlobalConfig = &Config{
		ServerPort: *serverPort,
		DatabaseSet: DbSet{
			Host:     *dbHost,
			User:     *dbUser,
			Password: *dbPassword,
			Port:     *dbPort,
			Database: *dbDatabase,
		},
		ModuleName: "github.com/sunshinev/go-sword",
		RootPath:   strings.TrimRight(*rootPath, "/"),
	}

	// Init MySQL connection
	c.initDbConnect()
}

func (c Config) initDbConnect() {
	dbc := GlobalConfig.DatabaseSet
	// user:password@(localhost)/dbname?charset=utf8&parseTime=True&loc=Local
	db, err := sql.Open("mysql", dbc.User+":"+dbc.Password+"@tcp("+dbc.Host+":"+strconv.Itoa(dbc.Port)+")/"+dbc.Database+"?&parseTime=True")
	if err != nil {
		log.Fatalf("%v", err)
	}

	GlobalConfig.DbConn = db
}
