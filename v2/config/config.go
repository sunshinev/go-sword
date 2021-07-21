package config

import (
	"database/sql"
	"errors"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

// App global config
var GlobalConfig *Config

type Config struct {
	DatabaseSet DbSet  `yaml:"db"`        // MySQL config
	RootPath    string `yaml:"root_path"` // The directory go-sword store new file
	ModuleName  string // Project go mod module name
	ServerPort  string `yaml:"tool_port"` // Go-sword server port
	DbConn      *sql.DB
}

type DbSet struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Port     int    `yaml:"port"`
	Database string `yaml:"database"`
}

func LoadConfig(configPath string) error {
	modName, err := readGoMod()
	if err != nil {
		log.Fatalf("read go mod err %v", err)
	}

	body, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}

	conf := Config{}
	err = yaml.Unmarshal(body, &conf)
	if err != nil {
		return err
	}

	conf.RootPath = strings.TrimRight(conf.RootPath, "/")
	conf.ModuleName = modName
	GlobalConfig = &conf

	initDbConnect()

	return nil
}

func readGoMod() (string, error) {
	// 获取go.mod文件中的module定义
	modBody, err := ioutil.ReadFile("go.mod")
	if err != nil {
		return "", err
	}

	log.Printf("%s", modBody)

	r := regexp.MustCompile(`module (.*)\n`)
	x := r.FindStringSubmatch(string(modBody))
	log.Printf("%v", x)

	if len(x) == 2 {
		return x[1], nil
	}

	return "", errors.New("parse `module` from go.mod error")
}

func initDbConnect() {
	dbc := GlobalConfig.DatabaseSet
	// user:password@(localhost)/dbname?charset=utf8&parseTime=True&loc=Local
	db, err := sql.Open("mysql", dbc.User+":"+dbc.Password+"@tcp("+dbc.Host+":"+strconv.Itoa(dbc.Port)+")/"+dbc.Database+"?&parseTime=True")
	if err != nil {
		log.Fatalf("%v", err)
	}

	GlobalConfig.DbConn = db
}
