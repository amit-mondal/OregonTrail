package main

import (
	"github.com/gorilla/mux"
	"math/rand"
	"net/http"
	"time"
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

var clientMap map[string]*Client
var lastEventLocation, currAvgLocation Location
var state State
var pendingEvent Event
var eventClientId string

var routes = Routes{
	Route{"Test", "GET", "/test", http.HandlerFunc(TestHandler)},
	Route{"Register", "POST", "/register", http.HandlerFunc(RegisterClientHandler)},
	Route{"GetUser", "GET", "/client/{clientid}", http.HandlerFunc(GetUserHandler)},
	Route{"StartGame", "GET", "/start", http.HandlerFunc(StartGameHandler)},
	Route{"CheckIn", "POST", "/checkin", http.HandlerFunc(CheckInHandler)},
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

	// We want a nice random seed
	rand.Seed(time.Now().UTC().UnixNano())
	pendingEvent = None
	clientMap = make(map[string]*Client)
	state = WaitForGameStart
	router := NewRegisteredRouter()
	http.ListenAndServe(":8080", router)

}
