package config

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"

	"github.com/jinzhu/gorm"
)

var GlobalEngineConfig = &Config{}
var GlobalDbConnect = &gorm.DB{}

type Config struct {
	Database   *DatabaseConfig `yaml:"db"`
	ServerPort string          `yaml:"server_port"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Port     int    `yaml:"port"`
	Db       string `yaml:"database"`
}

func LoadConfig(configPath string) error {
	body, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(body, GlobalEngineConfig)
	if err != nil {
		return err
	}

	return nil
}

func InitDB() {
	d := GlobalEngineConfig.Database
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", d.User,
		d.Password, d.Host, d.Port, d.Db))
	if err != nil {
		log.Fatal("%v", err)
	}

	GlobalDbConnect = db
}
