package config

type Config struct {
	Database *Db
}

type Db struct {
	Host     string
	User     string
	Password string
	Port     int
	Database string
}
