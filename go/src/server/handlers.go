package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

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
