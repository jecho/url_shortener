package main

import (
	"os"
	"encoding/json"
	"log"
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
	return
}

func loadConfig(s string) Config {
	file, _ := os.Open(s)
	decoder := json.NewDecoder(file)
	config := Config{}
	err := decoder.Decode(&config)
	checkErr(err)
	return config
}
