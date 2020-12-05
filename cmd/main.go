package main

import (
	"fmt"
	"log"
	"neckname/internal/config"
	"neckname/pkg/apiserver"
)

func main() {

	config, err := config.MakeConfig()
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting server ...")
	if err = apiserver.StartAPIServer(config); err != nil {
		log.Fatal(err)
	}

}

//migrate -path migrations/ -database "mysql://root:1_QwertY_2@/Social_Network" up
//migrate -path migrations/ -database "mysql://root:1_QwertY_2@/Social_Network" down
