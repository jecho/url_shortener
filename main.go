package main

import (
	"github.com/gorilla/mux"
	"fmt"
	"log"
	"net/http"
	"encoding/json"
	"strconv"
)

const (
	// default config values
	// seperate and gitignore later
	DB_HOST = "tcp(127.0.0.1:3306)"
	DB_NAME = "nintendo"
	DB_USER = "root"
	DB_PASS = "nintendo"

	// declarations
	a string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	base int = 62	// character's known in 'a'
	domain_name = "localhost:22222"
	prefix = "http://"
	sep = "/"
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
	shorterUrl := domain_name + sep + encode(int(id), base, a)

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
		decodedUri = domain_name + sep + "error"
	}
	fmt.Println(err)
	//checkErr(err)

	http.Redirect(res, req, prefix + decodedUri, 301)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
	return
}

func main() {

	dsn := DB_USER + ":" + DB_PASS + "@" + DB_HOST + "/" + DB_NAME
	//db, err := sql.Open("mysql", dsn)
	db, err := NewDB(dsn)
	checkErr(err)
	defer db.Close()

	env := &Env{db}
	var version string
	db.QueryRow("SELECT VERSION()").Scan(&version)
	fmt.Println("Connected to:", version)

	r := mux.NewRouter()
	r.HandleFunc("/p/{encoded_value}", env.retrieveEntry)
	r.HandleFunc("/create", env.createEntry).Methods("POST")
	r.NotFoundHandler = http.HandlerFunc(env.notFound)

	log.Fatal(http.ListenAndServe(":22222", r))
}
//
//curl -X POST -H "Content-Type: application/json" -d '{"url":"foxley.co/okay"}' http://127.0.0.1:8080/create -v
