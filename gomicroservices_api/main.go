package main

import (
	"net/http"
	"github.com/golang/glog"
	"github.com/gorilla/handlers"
	"os"
)

const (
	// declarations
	a string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	base int = 62	// character's known in 'a'
	domain_name = "nenley.co" // remove port when kube manifest are ready
	prefix = "http://"
	sep = "/"
	configFile = ".env/foxley_mock.json"
)

func main() {

	glog.Info("Service is warming up")
	// load config from file (remove this in the future, should use kub secret)
	//config := loadConfig(configFile) // default left, but remove when doing kube

	// for now override, if doing docker, omit; will put flag later
	config := Config{
		DB_HOST : os.Getenv("DB_HOST"),
		DB_NAME : os.Getenv("DB_NAME"),
		DB_USER : os.Getenv("DB_USER"),
		DB_PASS : os.Getenv("DB_PASS"),
	}

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
//curl -X POST -H "Content-Type: application/json" -d '{"url":"foxley.co/okay"}' http://http://a0c28d4b07e3c11e7b68402d49427e8f-634873964.us-west-2.elb.amazonaws.com/create -v
//curl -X POST -H "Content-Type: application/json" -d '{ "unfiltered":"foxley.co/okay","filtered":"foxley.co/okay2"}' http://localhost:12345/create -v