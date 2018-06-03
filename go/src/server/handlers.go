package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

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
		/*Give random starting values between 3-6*/
		newClient.Water = rand.Intn(4) + 3
		newClient.Food = rand.Intn(4) + 3
		newClient.Bullets = rand.Intn(4) + 3
		newClient.Supplies = rand.Intn(4) + 3
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
		state = WaitForCheckIn
		fmt.Println("Started Game")
	}
	WriteMessage(w, "Success")
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
		GameOver        bool    `json:"game_over"`
		PercentComplete float64 `json:"percent_complete"`
		Event           Event   `json:"event"`
		EventClientId   string  `json:"event_client"`
		Client          Client  `json:"client"`
	}
	// Defaults to 0, so make sure to set to None (-1)
	checkInResponse.Event = None
	checkInResponse.PercentComplete = distanceTravelled / GameDistance
	checkInResponse.GameOver = NumLivingClients() == 0

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
				//DoEvent(pendingEvent, eventClientId)
			} else {
				SetAllClientState(WillCheckIn)
				pendingEvent = None
			}
		} else {
			// Not every client has checked in yet, so just record its location.
			if client.State != HasCheckedIn {
				client.State = HasCheckedIn
				currAvgLocation.Lat += clientInfo.Location.Lat
				currAvgLocation.Lon += clientInfo.Location.Lon
				fmt.Printf("Client %s has checked in (Lat: %d, Lon: %d)\n", client.Id, clientInfo.Location.Lat, clientInfo.Location.Lon)
			}
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
		fmt.Printf("Client %s has received the event\n", client.Id)
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
		vars := mux.Vars(r)
		respondingClient := clientMap[vars["clientid"]]
		action := vars["action"]
		//Do something here to handle client response.
		if action == "true" {
			valid := DoEvent(w, pendingEvent, vars["clientid"])
			if !valid {
				//Don't change the state, so just return
				return
			}
			state = WaitForCheckIn
			SetAllClientState(WillCheckIn)
			fmt.Println("Client responded 'true' to event")
			pendingEvent = None
		} else {
			respondingClient.State = HasMadeDecision
			fmt.Println("Client responded 'false' to event")
			if AllLivingClientState(HasMadeDecision) {
				IgnoreEvent(w, pendingEvent, vars["clientid"])
				state = WaitForCheckIn
				SetAllClientState(WillCheckIn)
				fmt.Println("All clients responded 'false' to event")
				pendingEvent = None
			} else {
				WriteMessage(w, "Waiting for party to arrive at a decision")
			}
		}
	} else {
		LogBadRequest(w, "Client attempted to respond while there was no open event - another client may have already responded, or not every client has received the event")
		WriteMessage(w, "No open event to respond to - another client may have already responded, or not every client has received the event")
	}
}

func EventHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	action := vars["action"]
	eventNum, err := strconv.Atoi(vars["eventNum"])
	if err != nil {
		fmt.Printf("Atoi error")
	}
	pendingEvent = Event(eventNum)
	eventClientId = RandomClient()
	if action == "true" {
		valid := DoEvent(w, pendingEvent, vars["clientid"])
		if !valid {
			//Don't change the state, so just return
			return
		}
	} else {
		IgnoreEvent(w, pendingEvent, vars["clientid"])
	}
}
