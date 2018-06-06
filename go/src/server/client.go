package main

import (
	"fmt"
	"math/rand"
)

type ClientState int

const (
	WillCheckIn      ClientState = 0
	HasCheckedIn     ClientState = 1
	WillReceive      ClientState = 2
	HasReceived      ClientState = 3
	WillMakeDecision ClientState = 4
	HasMadeDecision  ClientState = 5
)

type Client struct {
	Id       string      `json:"id"`
	Location Location    `json:"location"`
	IsAlive  bool        `json:"is_alive"`
	Food     int         `json:"food"`
	Water    int         `json:"water"`
	Bullets  int         `json:"bullets"`
	Supplies int         `json:"supplies"`
	State    ClientState `json:"-"`
}

func SetAllClientState(clientState ClientState) {
	for key, _ := range clientMap {
		clientMap[key].State = clientState
	}
}

func AllClientState(clientState ClientState) bool {
	for _, value := range clientMap {
		if value.State != clientState {
			return false
		}
	}
	return true
}

func AllLivingClientState(clientState ClientState) bool {
	for _, value := range clientMap {
		if value.IsAlive && value.State != clientState {
			return false
		}
	}
	return true
}

func NumLivingClients() int {
	var numAlive int
	for _, client := range clientMap {
		if client.IsAlive {
			numAlive += 1
		}
	}
	return numAlive
}

func RandomClient() string {
	var randNum int
	numAlive := NumLivingClients()
	if numAlive < 2 {
		randNum = 0
	} else {
		// Intn panics when you give it 0 as an arg
		randNum = rand.Intn(numAlive)
	}
	ctr := 0
	for key, client := range clientMap {
		// Ignore clients that are not alive.
		if client.IsAlive {
			if ctr == randNum {
				return key
			}
			ctr = ctr + 1
		}
	}
	fmt.Println("ERR: failed to find random client")
	return ""
}
