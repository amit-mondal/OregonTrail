package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func LogBadRequest(w http.ResponseWriter, s string) {
	fmt.Printf("ERR: %s\n", s)
	w.WriteHeader(http.StatusBadRequest)
}

func TestHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Message: "This is a test endpoint for Oregon Trail Go.",
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func RegisterClientHandler(w http.ResponseWriter, r *http.Request) {
	if state == WaitForGameStart {
		var tmpClient struct {
			Id       string   `json:"id"`
			Location Location `json:"location"`
		}
		if err := json.NewDecoder(r.Body).Decode(&tmpClient); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var newClient Client
		newClient.Location = tmpClient.Location
		newClient.Id = tmpClient.Id
		newClient.IsAlive = true
		/*Give random starting values, can be changed as seen fit*/
		newClient.Water = 2
		newClient.Food = 5
		newClient.Bullets = 10
		newClient.Supplies = 5
		clientMap[newClient.Id] = newClient
		response := Response{
			Message: "Success",
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Printf("Client %s registered\n", newClient.Id)
	} else {
		LogBadRequest(w, "Attempted to add during game start")
	}
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if client, exists := clientMap[vars["clientid"]]; exists {
		if err := json.NewEncoder(w).Encode(client); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func StartGameHandler(w http.ResponseWriter, r *http.Request) {
	if state == WaitForGameStart {
		state = WaitForCheckIn
		fmt.Println("Started Game")
	} else {
		LogBadRequest(w, "Game already started")
	}
}
