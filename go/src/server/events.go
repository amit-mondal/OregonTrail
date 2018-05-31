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
	eventClient := clientMap[clientId]
	switch event {
	case Dysentery:
		eventClient.IsAlive = false
		fmt.Printf("Client %s died of dysentery.\n", clientId)
		return true
	case InadequateGrass:
		fmt.Printf("Inadequate Grass event, counter = %s.\n", BadGrass)
		BadGrass++
		if BadGrass >= 8 {
			for k, clients := range clientMap {
				clients.isAlive = false
			}
		}
	case Food:
		fmt.Printf("Client %s got food.\n", clientId)
		if eventClient.Bullets >= 1 {
			eventClient.Bullets--
			eventClient.Food += 2
			return true
		} else {
			//Client doesn't have enough bullets, so return false and do nothing
			return false
		}
	case BadWater:
		fmt.Printf("Client %s drank bad water!\n", clientId)
		if eventClient.Water >= 1 {
			eventClient.Water--
			return true
		} else {
			//Not enough water
			return false
		}
	case BrokenTongue:
		fmt.Printf("Client %s got a broken tongue!\n", clientId)
		if eventClient.Supplies >= 1 {
			eventClient.Supplies--
			return true
		} else {
			//Not enough supplies
			return false
		}
	case Starvation:
		fmt.Printf("Client %s is starving!\n", clientId)
		if eventClient.Food >= 1 {
			eventClient.Food--
			return true
		} else {
			//Not enough food
			return false
		}
	case Town:
		fmt.Printf("Client %s has reached a town!\n", clientId)
		//Give them a random amount of items from the twon
		eventClient.Food += rand.Intn(2)
		eventClient.Water += rand.Intn(2)
		eventClient.Supplies += rand.Intn(2)
		eventClient.Bullets += rand.Intn(2)
		return true
	case SnakeBite:
		fmt.Printf("Client %s has been bitten by a snake!\n", clientId)
		if eventClient.Supplies >= 1 {
			eventClient.Supplies--
			return true
		} else {
			//Not enough supplies
			return false
		}
	case Bandits:
		fmt.Printf("Client %s has been attacked by bandits!\n", clientId)
		if eventClient.Bullets >= 1 {
			eventClient.Bullets--
			return true
		} else {
			//Not enough bullets
			return false
		}
	default:
		fmt.Printf("Event %d not handled\n", event)
		return false
	}
}

func IgnoreEvent(event Event, clientId string) {
	eventClient := clientMap[clientId]
	switch event {
	case Dysentery:
		eventClient.IsAlive = false
		fmt.Printf("Client %s died of dysentery.\n", clientId)
	case InadequateGrass:
		fmt.Printf("Inadequate Grass event, counter = %s.\n", BadGrass)
		BadGrass++
		if BadGrass >= 8 {
			for k, clients := range clientMap {
				clients.isAlive = false
			}
		}
	case Food:
		fmt.Printf("Client %s ignored the food.\n", clientId)
	case BadWater:
		fmt.Printf("Client %s died from drinking bad water!\n", clientId)
		eventClient.IsAlive = false
	case BrokenTongue:
		fmt.Printf("The party is left stranded on the trail...!\n")
		//Everyone dies
		for k, clients := range clientMap {
			clients.IsAlive = false
		}
	case Starvation:
		fmt.Printf("Client %s starved to death!\n", clientId)
		eventClient.IsAlive = false
	case Town:
		fmt.Printf("Client %s has reached a town!\n", clientId)
		//Give them a random amount of items from the town
		eventClient.Food += rand.Intn(2)
		eventClient.Water += rand.Intn(2)
		eventClient.Supplies += rand.Intn(2)
		eventClient.Bullets += rand.Intn(2)
	case SnakeBite:
		fmt.Printf("Client %s has died by a snake bite!\n", clientId)
		eventClient.IsAlive = false
	case Bandits:
		fmt.Printf("Client %s has been attacked by banditts!\n", clientId)
		//Lose a random amount of items from the town
		if eventClient.Food >= 1 {
			eventClient.Food -= rand.Intn(1)
		}
		if eventClient.Water >= 1 {
			eventClient.Water -= rand.Intn(1)
		}
		if eventClient.Supplies >= 1 {
			eventClient.Supplies -= rand.Intn(1)
		}
		if eventClient.Bullets >= 1 {
			eventClient.Bullets -= rand.Intn(1)
		}
		//Small chance of 1/50 to lose your life (mean bandits)
		chance := rand.Intn(49)
		if chance == 0 {
			eventClient.IsAlive = false
		}
	default:
		fmt.Printf("Event %d not handled\n", event)
	}
}

func CheckForEvent() bool {
	return (rand.Intn(EventChance-1) == 0)
}
