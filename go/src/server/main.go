package main

import (
	"encoding/json"
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

var routes = Routes{
	Route{"Test", "GET", "/test", http.HandlerFunc(TestRoute)},
}

func TestRoute(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Message: "This is a test endpoint for Oregon Trail Go.",
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
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

	router := NewRegisteredRouter()

	http.ListenAndServe(":8080", router)

}
