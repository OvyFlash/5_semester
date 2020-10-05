package main

import (
	"flag"
	"log"
	"neckname/internal/app/apiserver"

	"github.com/BurntSushi/toml"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "Basic dependencies of the server",
		"configs/apiserver.toml", "path to config file")
}

func main() {
	flag.Parse() //write to configPath toml strings

	config := apiserver.NewConfig()//new config to write all toml info
	_, err := toml.DecodeFile(configPath, config)//parse info into config
	if err != nil {
		log.Fatal(err)
	}

	if err := apiserver.StartAPIServer(config); err != nil {
		log.Fatal(err)
	}
	
}

//migrate -path migrations/ -database "mysql://root:1_QwertY_2@/Social_Network" up
//migrate -path migrations/ -database "mysql://root:1_QwertY_2@/Social_Network" down