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

func RandomClient() string {
	var randNum int
	if len(clientMap) == 1 {
		randNum = 0
	} else {
		// Intn panics when you give it 0 as an arg
		randNum = rand.Intn(len(clientMap) - 1)
	}
	ctr := 0
	for key, _ := range clientMap {
		if ctr == randNum {
			return key
		}
		ctr = ctr + 1
	}
	fmt.Println("ERR: failed to find random client")
	return ""
}
