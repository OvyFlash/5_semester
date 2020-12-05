package config

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
)

var (
	configPath string
	c *Config = &Config{}
)

func init() {
	flag.StringVar(&configPath, "Basic dependencies of the server",
		"./internal/config/config.toml", "path to config file")

}

//Config struct with setting
type Config struct {
	BindAddr    string `toml:"bind_addr"` //Adress of APIServer start
	SessionKey  string `toml:"session_key"`
	DatabaseURL string `toml:"database_url"`
}

//MakeConfig ...
func MakeConfig() (*Config, error){
	flag.Parse() //write to configPath toml strings

	_, err := toml.DecodeFile(configPath, c) //parse info into config
	if err != nil {
		if len(os.Args) < 2 {
			fmt.Println("Додайте другим аргументом файл налаштувань")
			return nil, err
		}
		file, err := ioutil.ReadFile(os.Args[1])
		if err != nil {
			os.Exit(1)
		}
		_, err = toml.Decode(string(file), c)
		if err != nil {
			panic(err)
		}

	}

	return c, nil
}