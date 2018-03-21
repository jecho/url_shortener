package main

import (
	"encoding/json"
	"net/http"
	"time"
	"github.com/golang/glog"
	"os"
)

func main() {

	// config := &Config{ "http://localhost:11124", mariadbService }
	// configMap
	config := Config{
		MARIADB_SERVICE : os.Getenv("MARIADB_SERVICE"),
		REDIS_SERVICE : os.Getenv("REDIS_SERVICE"),
	}

	quickErr := map[string]string{"errorMessage": "your request could not be completed",}
	jsonBlob, _ := json.Marshal(quickErr)

	glog.Info("Initializing Gateway Service")

	env := Env{&config}
	r := NewRouter(&env)

	glog.Info("Adding Middleware for Timeout Protocals")
	muxWithMiddlewares := http.TimeoutHandler(r, time.Second*5, string(jsonBlob))
	glog.Info("Service is ONLINE")
	glog.Fatal(http.ListenAndServe(":8080", muxWithMiddlewares))
}