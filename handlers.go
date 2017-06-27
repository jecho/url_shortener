package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	s "strings"
	"strconv"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
)

//	default handler
func (env *Env) notFound(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("route not found... sorry! :("))
}

func (env *Env) createEntry(res http.ResponseWriter, req *http.Request) {

	glog.Info("Inbound Traffic: ")
	// declare content type 'json'
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// capture and decode the payload from stream
	urlPayload := new(Url)
	decoder := json.NewDecoder(req.Body)
	error := decoder.Decode(&urlPayload)
	checkErr(error)

	// check for http, https, and www, fail other cases, should subsitiute as regex later
	updatedUrl := urlPayload.Url
	updatedUrl = s.Replace(updatedUrl, "http://www", "", 1)
	updatedUrl = s.Replace(updatedUrl, "https://www", "", 1)
	updatedUrl = s.Replace(updatedUrl, "https://", "", 1)
	updatedUrl = s.Replace(updatedUrl, "http://", "", 1)

	// check if it really exists?
	//_, err := url.ParseRequestURI(updatedUrl)
	//if err != nil {
	//	glog.Fatal(err)
	//}

	// prepare insert statement
	query, _ := env.db.Prepare("INSERT IGNORE INTO nin SET url=?")

	// execute
	execute, err := query.Exec(updatedUrl)
	checkErr(err)

	// grab the id, encode, and suffix it to domain_name
	id, err2 := execute.LastInsertId()
	checkErr(err2)
	shorterUrl := domain_name + sep + "p" + sep + encode(int(id), base, a)

	// return the response
	glog.Info(shorterUrl)
	res.WriteHeader(http.StatusCreated)
	if err3 := json.NewEncoder(res).Encode(&Url{
		shorterUrl,
	}); err3 != nil {
		glog.Fatal(err)
	}
}

func (env *Env) retrieveEntry(res http.ResponseWriter, req *http.Request) {

	glog.Info("Outbound Traffic: ")
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
	fmt.Print(err)
	//checkErr(err)

	glog.Info("Redirecting: ", prefix + decodedUri)
	http.Redirect(res, req, prefix + decodedUri, 301)
}