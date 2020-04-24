package config

type Config struct {
	Database *Db
	// The directory go-sword store new file
	RootPath string
	// Project go mod module name
	ModuleName string
}

type Db struct {
	Host     string
	User     string
	Password string
	Port     int
	Database string
}

func (c *Config) getRootPath() string {
	if c.RootPath == "" {
		c.RootPath = "go-sword-app"
	}

	return c.RootPath
}
