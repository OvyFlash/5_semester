package apiserver

//Config contains all server configs
type Config struct {
	BindAddr    string `toml:"bind_addr"` //Adress of APIServer start
	LogLever    string `toml:"log_level"`
	//DatabaseULR string `toml:"database_url"`
	SessionKey  string `toml:"session_key"`
	DatabaseURL string `toml:"database_url"`
}

//NewConfig returns new instance of config
func NewConfig() *Config {
	return &Config{
		BindAddr: ":8080",
		LogLever: "debug",
	}
}
