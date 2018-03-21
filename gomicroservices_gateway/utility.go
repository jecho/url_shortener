package main

import (
	"os"
	"encoding/json"
	"github.com/golang/glog"
)

func checkErr(err error) {
	if err != nil {
		glog.Fatal(err)
	}
	//return
}

func loadConfig(s string) Config {
	file, err := os.Open(s)
	checkErr(err)
	decoder := json.NewDecoder(file)
	config := Config{}
	err = decoder.Decode(&config)
	checkErr(err)
	return config
}
