package main

import (
	"fmt"
	"math/rand"
)

type Event int

const (
	None            Event = -1
	Dysentery       Event = 0
	InadequateGrass Event = 1
	Food            Event = 2
	BadWater        Event = 3
	BrokenTongue    Event = 4
	Starvation      Event = 5
	Town            Event = 6
	SnakeBite       Event = 7
	Bandits         Event = 8
)

const (
	NumEvents   int = 9 // We don't count the "None" event
	EventChance int = 2 // Probability of an event is 1 / EventChance
)

func RandomEvent() Event {
	return Event(rand.Intn(NumEvents - 1)) // Should never give us "None"
}

func DoEvent(event Event, clientId string) {
	switch event {
	case Dysentery:
		clientMap[clientId].IsAlive = false
		fmt.Printf("Client %s died of dysentery\n", clientId)
	default:
		fmt.Printf("Event %d not handled\n", event)
	}
}

func CheckForEvent() bool {
	return (rand.Intn(EventChance-1) == 0)
}
