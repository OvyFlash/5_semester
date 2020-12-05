package main

import (
	"flag"
	"io/ioutil"
	"log"
	"neckname/internal/app/apiserver"
	"os"
	"fmt"
	
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

	config := apiserver.NewConfig()               //new config to write all toml info
	_, err := toml.DecodeFile(configPath, config) //parse info into config
	if err != nil {
		if len(os.Args) < 2 {
			fmt.Println("Додайте другим аргументом файл налаштувань")
			return
		}
		file, err := ioutil.ReadFile(os.Args[1])
		if err != nil {
			os.Exit(1)
		}
		_, err = toml.Decode(string(file), config)
		if err != nil {
			panic(err)
		}

	}

	if err := apiserver.StartAPIServer(config); err != nil {
		log.Fatal(err)
	}

}

//migrate -path migrations/ -database "mysql://root:1_QwertY_2@/Social_Network" up
//migrate -path migrations/ -database "mysql://root:1_QwertY_2@/Social_Network" down
