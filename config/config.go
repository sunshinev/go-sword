package config

type Config struct {
	Database *DbSet
	// The directory go-sword store new file
	RootPath string
	// Project go mod module name
	ModuleName string
	// Go-sword server port
	ServerPort string
}

type DbSet struct {
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
