package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

//BadGrass count for events
var BadGrass int

func LogBadRequest(w http.ResponseWriter, s string) {
	fmt.Printf("ERR: %s\n", s)
	w.WriteHeader(http.StatusBadRequest)
}

func WriteMessage(w http.ResponseWriter, s string) {
	response := Response{
		Message: s,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func TestHandler(w http.ResponseWriter, r *http.Request) {
	WriteMessage(w, "This is a test endpoint for Oregon Trail Go")
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
		newClient.State = WillCheckIn
		clientMap[newClient.Id] = &newClient
		WriteMessage(w, "Success")
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
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func StartGameHandler(w http.ResponseWriter, r *http.Request) {
	if state == WaitForGameStart {
		WriteMessage(w, "Success")
		state = WaitForCheckIn
		fmt.Println("Started Game")
		//Also set the BadGrass to 0
		BadGrass = 0
	} else {
		LogBadRequest(w, "Game already started")
	}
}

func CheckInHandler(w http.ResponseWriter, r *http.Request) {

	// First decode into a temp struct
	var clientInfo struct {
		Id       string   `json:"id"`
		Location Location `json:"location"`
	}

	var client *Client

	// Try to decode
	if err := json.NewDecoder(r.Body).Decode(&clientInfo); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Try to get client
	if cli, exists := clientMap[clientInfo.Id]; exists {
		client = cli

	} else {
		LogBadRequest(w, "Invalid client ID")
	}

	// Struct defining response format
	var checkInResponse struct {
		PercentComplete float64 `json:"percent_complete"`
		Event           Event   `json:"event"`
		EventClientId   string  `json:"event_client"`
		Client          Client  `json:"client"`
	}
	// Defaults to 0, so make sure to set to None (-1)
	checkInResponse.Event = None
	checkInResponse.PercentComplete = distanceTravelled / GameDistance

	switch state {
	case WaitForGameStart:
		LogBadRequest(w, "Cannot check in until game is started")
	case WaitForCheckIn:

		if AllClientState(HasCheckedIn) {
			fmt.Println("All clients checked in")
			if UpdateLocation() {
				SetAllClientState(WillReceive)
				state = WaitForReceive
				pendingEvent = RandomEvent()
				eventClientId = RandomClient()
				fmt.Printf("Event %d selected\n", pendingEvent)
				DoEvent(pendingEvent, eventClientId)
			} else {
				SetAllClientState(WillCheckIn)
				pendingEvent = None
			}
		} else {
			// Not every client has checked in yet, so just record its location.
			client.State = HasCheckedIn
			currAvgLocation.Lat += clientInfo.Location.Lat
			currAvgLocation.Lon += clientInfo.Location.Lon
			fmt.Printf("Client %s has checked in\n", client.Id)
		}

		// Whenever a client checks in during this state, we simply echo back the client
		checkInResponse.Client = *client
		if err := json.NewEncoder(w).Encode(checkInResponse); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return

	case WaitForReceive:
		if pendingEvent != None {
			checkInResponse.Event = pendingEvent
			checkInResponse.EventClientId = eventClientId
		}
		checkInResponse.Client = *client
		if err := json.NewEncoder(w).Encode(checkInResponse); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		client.State = HasReceived
		// If that was the last client to receive, move to the next state
		if AllClientState(HasReceived) {
			fmt.Println("All clients received event")
			SetAllClientState(WillMakeDecision)
			state = WaitForDecision
			return
		}

	case WaitForDecision:
		WriteMessage(w, "Waiting for response to pending event")
	}
}

func RespondHandler(w http.ResponseWriter, r *http.Request) {
	if state == WaitForDecision {
		vars := mux.Vars()
		//respondingClient := clientMap[vars["clientid"]]
		action := vars["action"]
		// TODO: Do something else here to handle client response.
		if action == "true" {
			valid := DoEvent(pendingEvent, vars["clientid"])
			if !valid {
				//Don't change the state, so just return
				return
			}
		} else {
			IgnoreEvent(pendingEvent, vars["clientid"])
		}
		state = WaitForCheckIn
		SetAllClientState(WillCheckIn)
		fmt.Println("Client responded to event")
		pendingEvent = None
	} else {
		LogBadRequest(w, "Client attempted to respond while there was no open event - another client may have already responded, or not every client has received the event")
		WriteMessage(w, "No open event to respond to - another client may have already responded, or no every client has received the event")
	}
}
