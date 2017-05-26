package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"github.com/golang/glog"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"database/sql"
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

//	default handler
func (env *Env) notFound(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("route not found... sorry! :("))
}

func (env *Env) createEntry(res http.ResponseWriter, req *http.Request) {

	// declare content type 'json'
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// capture and decode the payload from stream
	url := new(Url)
	decoder := json.NewDecoder(req.Body)
	error := decoder.Decode(&url)
	checkErr(error)

	// prepare insert statement
	query, _ := env.db.Prepare("INSERT IGNORE INTO nin SET url=?")

	// execute
	execute, err := query.Exec(url.Url)
	checkErr(err)

	// grab the id, encode, and suffix it to domain_name
	id, err2 := execute.LastInsertId()
	checkErr(err2)
	shorterUrl := domain_name + sep + "p" + sep + encode(int(id), base, a)

	// return the response
	res.WriteHeader(http.StatusCreated)
	if err3 := json.NewEncoder(res).Encode(&Url{
		shorterUrl,
	}); err3 != nil {
		log.Fatal(err)
	}
}

func (env *Env) retrieveEntry(res http.ResponseWriter, req *http.Request) {

	// declare content type 'json'
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vars := mux.Vars(req)
	shorterUrl := vars["encoded_value"]
	id := decode(shorterUrl, a)

	// query statement
	var decodedUri string
	err := env.db.QueryRow("SELECT URL FROM nin WHERE ID=?",
		strconv.Itoa(id)).Scan(&decodedUri)

	if len(decodedUri) == 0 || decodedUri == "" {
		decodedUri = domain_name + sep + "404" + sep
	}
	fmt.Println(err)
	//checkErr(err)

	http.Redirect(res, req, prefix + decodedUri, 301)
}

func main() {

	glog.Info("Service is warming up")
	// load config from file (remove this in the future, should use kub secret)
	config := loadConfig(configFile)

	dsn := config.DB_USER + ":" + config.DB_PASS + "@" + config.DB_HOST + "/" + config.DB_NAME
	db, err := NewDB(dsn)
	checkErr(err)
	defer db.Close()

	env := &Env{db : db}

	startService(env, db)
}

func startService(env *Env, db *sql.DB) {
	glog.Info("Current profile : ", configFile)
	var version string
	db.QueryRow("SELECT VERSION()").Scan(&version)
	glog.Info("Connected to: ", version)

	// CORS; re-review later
	headersOk := handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "Content-Length", "X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	r := mux.NewRouter()
	r.HandleFunc("/p/{encoded_value}", env.retrieveEntry)
	r.HandleFunc("/create", env.createEntry).Methods("POST")
	r.NotFoundHandler = http.Handler(http.StripPrefix("/404", http.FileServer(http.Dir("./static/404/"))))

	glog.Info("Service is ready")
	log.Fatal(http.ListenAndServe(":22222", handlers.CORS(originsOk, headersOk, methodsOk)(r)))
}
//
//curl -X POST -H "Content-Type: application/json" -d '{"url":"foxley.co/okay"}' http://127.0.0.1:8080/create -v
