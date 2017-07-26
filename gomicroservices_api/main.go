package main

import (
	"net/http"
	"github.com/golang/glog"
	"github.com/gorilla/handlers"
)

const (
	// declarations
	a string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	base int = 62	// character's known in 'a'
	domain_name = "foxley.co:22222" // remove port when kube manifest are ready
	prefix = "http://"
	sep = "/"
	configFile = ".env/foxley_mock.json"
)

func main() {

	glog.Info("Service is warming up")
	// load config from file (remove this in the future, should use kub secret)
	config := loadConfig(configFile)

	dsn := config.DB_USER + ":" + config.DB_PASS + "@" + config.DB_HOST + "/" + config.DB_NAME
	db, err := NewDB(dsn)
	checkErr(err)
	defer db.Close()
	env := &Env{db : db}

	startService(env)
}

func startService(env *Env) {
	glog.Info("Current config: ", configFile)
	var version string
	env.db.QueryRow("SELECT VERSION()").Scan(&version)
	glog.Info("Connected to: ", version)

	// CORS; re-review later
	headersOk := handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "Content-Length", "X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	r := NewRouter(env)

	glog.Info("Service is ready")
	glog.Fatal(http.ListenAndServe(":22222", handlers.CORS(originsOk, headersOk, methodsOk)(r)))
}
//
//curl -X POST -H "Content-Type: application/json" -d '{"url":"foxley.co/okay"}' http://a6b3114d5723f11e78ea302ddb01ab83-239777298.us-west-2.elb.amazonaws.com/create -v
