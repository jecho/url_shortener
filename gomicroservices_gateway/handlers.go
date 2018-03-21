package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
)

//	default handler
func (env *Env) notFound(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("blackswan."))
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
	jsonBlob, err := json.Marshal(urlPayload)
	glog.Info("Payload contents: ", jsonBlob)
	checkErr(err)

	// do something to urlPayload
	// ignored for now

	// create new request
	reqForward, err := http.NewRequest("POST", "http://" + env.config.MARIADB_SERVICE + "/create", bytes.NewBuffer(jsonBlob))
	reqForward.Header.Set("Content-Type", "application/json")
	checkErr(err)

	// communicate with other service
	glog.Info("Communicating with ", env.config.MARIADB_SERVICE)
	client := &http.Client{}
	resp, err := client.Do(reqForward)
	//checkErr(err)

	defer resp.Body.Close()
	// capture response form other service, and deal with any fatals
	if resp.Status != "201" || resp.Status != "200" {
		//glog.Fatal()
	}
	// read and validate the response
	body, err := ioutil.ReadAll(resp.Body)
	checkErr(err)
	glog.Info("Received Response from ", env.config.MARIADB_SERVICE)

	// may be unnesessary
	var dat Url
	err = json.Unmarshal(body, &dat)
	checkErr(err)

	// communicate with the redis service
	// combined the original request with the returned request from mariadb_service
	// -- wrap with context to deal with timeout, retry, etc
	combinedStruct := &Entry{Url(jsonBlob).Url, dat.Url}
	glog.Info("Communicating with ", env.config.REDIS_SERVICE)
	reqForward2, err := http.NewRequest("POST", "http://" + env.config.REDIS_SERVICE + "/create", bytes.NewBuffer([]byte(combinedStruct)))
	reqForward2.Header.Set("Content-Type", "application/json")
	//checkErr(err)
	if err != nil {
		// failed
	} else {
		glog.Info("Received Response from ", env.config.REDIS_SERVICE)
	}

	// return the response
	res.WriteHeader(http.StatusCreated)
	if err3 := json.NewEncoder(res).Encode(&dat); err3 != nil {
		glog.Fatal(err)
	}
}

func (env *Env) retrieveEntry(res http.ResponseWriter, req *http.Request) {

	glog.Info("Outbound Traffic: ")
	// declare content type 'json'
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// capture the payload from stream
	vars := mux.Vars(req)
	shorterUrl := vars["encoded_value"]

	// do something to urlPayload
	// ignored for now

	// check if it exists in the in memory database
	reqForwardRedis, err := http.NewRequest("GET", env.config.REDIS_SERVICE + "/p/" + shorterUrl, nil)
	reqForwardRedis.Header.Set("Content-Type", "application/json")

	// make the call to the standard database if it did not exist prior
	if err != nil {
		// create new request
		reqForward, err := http.NewRequest("GET", env.config.MARIADB_SERVICE + "/p/" + shorterUrl, nil)
		reqForward.Header.Set("Content-Type", "application/json")
		checkErr(err)

		// communicate with other service
		glog.Info("Communicating with ", env.config.MARIADB_SERVICE)
		client := &http.Client{}
		resp, err := client.Do(reqForward)
		checkErr(err)

		defer resp.Body.Close()
		if resp.Status != "201" || resp.Status != "200" {
			//glog.Fatal()
		}

		body, err := ioutil.ReadAll(resp.Body)
		checkErr(err)

		var dat Url
		err = json.Unmarshal(body, &dat)
		checkErr(err)

		// redirect request
		glog.Info("Redirecting: ", &dat.Url)
		http.Redirect(res, req, dat.Url, 301)
	}

	// communicate with other service
	glog.Info("Communicating with ", env.config.REDIS_SERVICE)
	client := &http.Client{}
	resp, err := client.Do(reqForwardRedis)
	checkErr(err)

	defer resp.Body.Close()
	if resp.Status != "201" || resp.Status != "200" {
		//glog.Fatal()
	}

	body, err := ioutil.ReadAll(resp.Body)
	checkErr(err)

	var dat Url
	err = json.Unmarshal(body, &dat)
	checkErr(err)

	// redirect request
	glog.Info("Redirecting: ", &dat.Url)
	http.Redirect(res, req, dat.Url, 301)
}