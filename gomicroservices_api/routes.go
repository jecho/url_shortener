package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/golang/glog"
)

type Route struct {
	Name string
	Method string
	Pattern string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter(env *Env) *mux.Router {

	glog.Info("Initializing routes and handlers")

	var routes = Routes{
		Route {
			"retrieveEntry",
			"GET",
			"/p/{encoded_value}",
			env.retrieveEntry,
		},
		Route {
			"createEntry",
			"POST",
			"/create",
			env.createEntry,
		},
	}

	router := mux.NewRouter()
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	router.NotFoundHandler = http.Handler(http.StripPrefix("/404/", http.FileServer(http.Dir("./static/404/"))))

	glog.Info("Routes and Handlers loaded")

	return router
}
