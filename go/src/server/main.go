package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Route struct {
	name    string
	method  string
	pattern string
	handler http.HandlerFunc
}

type Response struct {
	Message string `json:"message"`
}

type Routes []Route

var clientMap map[string]Client

var routes = Routes{
	Route{"Test", "GET", "/test", http.HandlerFunc(TestHandler)},
	Route{"Register", "POST", "/register", http.HandlerFunc(RegisterClientHandler)},
	Route{"GetUser", "GET", "/client/{clientid}", http.HandlerFunc(GetUserHandler)},
}

func NewRegisteredRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.method).
			Path(route.pattern).
			Name(route.name).
			Handler(route.handler)
	}
	return router
}

func main() {

	clientMap = make(map[string]Client)
	router := NewRegisteredRouter()
	http.ListenAndServe(":8080", router)

}
