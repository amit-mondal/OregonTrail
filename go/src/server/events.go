package main

import (
	"fmt"
	"math/rand"
	"net/http"
)

type Event int

const (
	None        Event = -1
	Dysentery   Event = 0
	Bandits     Event = 1
	Food        Event = 2
	BadWater    Event = 3
	BrokenWheel Event = 4
	Starvation  Event = 5
	Town        Event = 6
	SnakeBite   Event = 7
)

const (
	NumEvents   int = 8 // We don't count the "None" event
	EventChance int = 2 // Probability of an event is 1 / EventChance
)

func RandomEvent() Event {
	return Event(rand.Intn(NumEvents)) // Should never give us "None"
}

func DoEvent(w http.ResponseWriter, event Event, clientId string) bool {
	respondingClient := clientMap[clientId]
	eventClient := clientMap[eventClientId]
	switch event {
	case Dysentery:
		eventClient.IsAlive = false
		fmt.Printf("Client %s died of dysentery.\n", eventClientId)
		message := fmt.Sprintf("Alas, no one could save %s from dysentery.", eventClientId)
		WriteMessage(w, message)
		return true
	case Food:
		fmt.Printf("Client %s got food.\n", clientId)
		if respondingClient.Bullets >= 1 {
			respondingClient.Bullets--
			chance := rand.Intn(20)
			if chance == 0 {
				//1/20 chance of missing
				WriteMessage(w, "Too eager! You pulled the trigger too soon and missed!\n -1 Bullet")
			} else {
				respondingClient.Food += 2
				WriteMessage(w, "With ease you land a bulleyes!\n -1 Bullet, +2 Food")
			}
			return true
		} else {
			//Client doesn't have enough bullets, so return false and do nothing
			w.WriteHeader(http.StatusBadRequest)
			WriteMessage(w, "You do not have enough bullets.")
			return false
		}
	case BadWater:
		fmt.Printf("Client %s drank clean water.\n", eventClientId)
		if respondingClient.Water >= 1 {
			if clientId == eventClientId {
				WriteMessage(w, "You quickly chug your clean water to cleanse yourself.\n -1 Water")
			} else {
				message := fmt.Sprintf("You throw your flask to %s to save your pal.\n -1 Water", eventClientId)
				WriteMessage(w, message)
			}
			respondingClient.Water--
			return true
		} else {
			//Not enough water
			w.WriteHeader(http.StatusBadRequest)
			WriteMessage(w, "You do not have enough water.")
			return false
		}
	case BrokenWheel:
		fmt.Printf("Client %s got a broken wheel!\n", eventClientId)
		if respondingClient.Supplies >= 1 {
			WriteMessage(w, "You used your supplies to fix the wheel. Your journey continues.\n -1 Supplies")
			respondingClient.Supplies--
			return true
		} else {
			//Not enough supplies
			w.WriteHeader(http.StatusBadRequest)
			WriteMessage(w, "You do not have enough supplies.")
			return false
		}
	case Starvation:
		fmt.Printf("Client %s was fed.\n", eventClientId)
		if respondingClient.Food >= 1 {
			if clientId == eventClientId {
				WriteMessage(w, "In mere seconds, you engulf your meal to fill the void of your stomach.\n -1 Food")
			} else {
				message := fmt.Sprintf("A feast is nothing without your friend! You give your food to %s.\n -1 Food", eventClientId)
				WriteMessage(w, message)
			}
			respondingClient.Food--
			return true
		} else {
			//Not enough food
			w.WriteHeader(http.StatusBadRequest)
			WriteMessage(w, "You do not have enough food.")
			return false
		}
	case Town:
		fmt.Printf("Client %s has reached a town!\n", eventClientId)
		//Give them a random amount of items from the town
		foodNum := rand.Intn(3)
		waterNum := rand.Intn(3)
		suppNum := rand.Intn(3)
		bullNum := rand.Intn(3)
		eventClient.Food += foodNum
		eventClient.Water += waterNum
		eventClient.Supplies += suppNum
		eventClient.Bullets += bullNum
		message := fmt.Sprintf("%s enters a town and is showered in kindess. | +%d Food, +%d Water, +%d Supplies, +%d Bullets", eventClientId, foodNum, waterNum, suppNum, bullNum)
		WriteMessage(w, message)
		return true
	case SnakeBite:
		fmt.Printf("Client %s survived the snake bite\n", eventClientId)
		if respondingClient.Supplies >= 1 {
			if clientId == eventClientId {
				WriteMessage(w, "So many snakes! Regardless, you quickly apply your handy anti-snake cream.\n -1 Supplies")
			} else {
				message := fmt.Sprintf("Thankful it was not yourself, you give your pal %s some anti-snake cream.\n -1 Supplies", eventClientId)
				WriteMessage(w, message)
			}
			respondingClient.Supplies--
			return true
		} else {
			//Not enough supplies
			w.WriteHeader(http.StatusBadRequest)
			WriteMessage(w, "You do not have enough supplies.")
			return false
		}
	case Bandits:
		fmt.Printf("Client %s has been attacked by bandits!\n", eventClientId)
		if respondingClient.Bullets >= 1 {
			if clientId == eventClientId {
				WriteMessage(w, "Bandits! However, you are always prepared and fight them off!\n -1 Bullets")
			} else {
				message := fmt.Sprintf("Bandits! You spring into action to help your pal %s. %s thanks you for your heroic act!\n -1 Bullets", eventClientId, eventClientId)
				WriteMessage(w, message)
			}
			respondingClient.Bullets--
			return true
		} else {
			//Not enough bullets
			w.WriteHeader(http.StatusBadRequest)
			WriteMessage(w, "You do not have enough bullets.")
			return false
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Printf("Event %d not handled\n", event)
		return false
	}
	return false
}

func IgnoreEvent(w http.ResponseWriter, event Event, clientId string) {
	eventClient := clientMap[eventClientId]
	switch event {
	case Dysentery:
		eventClient.IsAlive = false
		fmt.Printf("Client %s died of dysentery.\n", eventClientId)
		message := fmt.Sprintf("Alas, no one could save %s from dysentery.", eventClientId)
		WriteMessage(w, message)
	case Food:
		fmt.Printf("Client %s ignored the food event.\n", eventClientId)
		WriteMessage(w, "Your hands stay in your pockets, as you all watch the animals.\n")
	case BadWater:
		fmt.Printf("Client %s died from drinking bad water!\n", eventClientId)
		eventClient.IsAlive = false
		message := fmt.Sprintf("Nature bested %s with its temptations! Think before you drink! RIP %s.", eventClientId, eventClientId)
		WriteMessage(w, message)
	case BrokenWheel:
		fmt.Printf("The party is left stranded on the trail...\n")
		//Everyone dies
		WriteMessage(w, "Without a wagon, you all spend your days stranded. With bitter smiles, you all spend your limited days together...\n")
		for _, clients := range clientMap {
			clients.IsAlive = false
		}
	case Starvation:
		fmt.Printf("Client %s starved to death!\n", eventClientId)
		message := fmt.Sprintf("The long trial cannot be travelled on an empty stomach. RIP %s.", eventClientId)
		WriteMessage(w, message)
		eventClient.IsAlive = false
	case Town:
		fmt.Printf("Client %s has reached a town!\n", eventClientId)
		//Give them a random amount of items from the town
		foodNum := rand.Intn(3)
		waterNum := rand.Intn(3)
		suppNum := rand.Intn(3)
		bullNum := rand.Intn(3)
		eventClient.Food += foodNum
		eventClient.Water += waterNum
		eventClient.Supplies += suppNum
		eventClient.Bullets += bullNum
		message := fmt.Sprintf("%s enters a town, and they shower you in kindess.\n +%d Food, +%d Water, +%d Supplies, +%d Bullets", eventClientId, foodNum, waterNum, suppNum, bullNum)
		WriteMessage(w, message)
	case SnakeBite:
		fmt.Printf("Client %s has died by a snake bite!\n", eventClientId)
		message := fmt.Sprintf("Watch your step, for nature has fangs! RIP %s", eventClientId)
		WriteMessage(w, message)
		eventClient.IsAlive = false
	case Bandits:
		fmt.Printf("Client %s has been attacked by bandits!\n", eventClientId)
		//Lose a random amount of items from the town
		foodLoss := 0
		waterLoss := 0
		suppLoss := 0
		bullLoss := 0
		if eventClient.Food >= 1 {
			foodLoss = rand.Intn(2)
			eventClient.Food -= foodLoss
		}
		if eventClient.Water >= 1 {
			waterLoss = rand.Intn(2)
			eventClient.Water -= waterLoss
		}
		if eventClient.Supplies >= 1 {
			suppLoss = rand.Intn(2)
			eventClient.Supplies -= suppLoss
		}
		if eventClient.Bullets >= 1 {
			bullLoss = rand.Intn(2)
			eventClient.Bullets -= bullLoss
		}
		//Small chance of 1/50 to lose your life (mean bandits)
		chance := rand.Intn(50)
		if chance == 0 {
			eventClient.IsAlive = false
			message := fmt.Sprintf("The trail can bring the worst in people. And %s unfortunately crossed paths with such people. Revenge lingers in the rest of the group's mind....RIP %s", eventClientId, eventClientId)
			WriteMessage(w, message)
		} else {
			message := fmt.Sprintf("The trail can make others desperate. Bandits attack, and runaway with %s's goods!\n -%d Food, -%d Water, -%d Supplies, -%d Bullets", eventClientId, foodLoss, waterLoss, suppLoss, bullLoss)
			WriteMessage(w, message)
		}
	default:
		fmt.Printf("Event %d not handled\n", event)
	}
}

func CheckForEvent() bool {
	return true
}
