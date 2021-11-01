package main

import (
	"APIGateway/pkg/api"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {

	// Read fiel configuration

	b, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatal(err)
	}
	var config api.Config
	err = json.Unmarshal(b, &config)
	if err != nil {
		log.Fatal(err)
	}

	a := api.New(config)

	fmt.Println("Файл конфига:", string(b))

	// Web server start
	err = http.ListenAndServe(":8080", a.Router())
	if err != nil {
		log.Fatal(err)
	}

}
