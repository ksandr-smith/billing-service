package store

type Config struct {
	Host     string `toml:"db_host"`
	Port     string `toml:"db_port"`
	Database string `toml:"db_name"`
	Username string `toml:"db_user"`
	Password string `toml:"db_password"`
}

func NewConfig() *Config {
	return &Config{
		Host:     "localhost",
		Port:     "5432",
		Database: "postgres",
		Username: "postgres",
		Password: "1234",
	}
}
