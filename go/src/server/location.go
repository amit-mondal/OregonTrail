package main

import (
	"fmt"
	"math"
)

// All location constants in meters
const (
	EventDistance float64 = 10
	GameDistance  float64 = 200
)

type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

// haversin(θ) function
func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

func Delta(loc1, loc2 Location) float64 {
	// convert to radians
	// must cast radius as float to multiply later
	var la1, lo1, la2, lo2, r float64
	la1 = loc1.Lat * math.Pi / 180
	lo1 = loc1.Lon * math.Pi / 180
	la2 = loc2.Lat * math.Pi / 180
	lo2 = loc2.Lon * math.Pi / 180

	r = 6378100 // Earth radius in METERS

	// calculate
	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)

	return 2 * r * math.Asin(math.Sqrt(h))
}

// Returns true if location was updated
func UpdateLocation() bool {
	// Always reset the running average location when exiting this function
	defer func() {
		currAvgLocation.Lat = 0
		currAvgLocation.Lon = 0
	}()

	// Divide summed location by number of clients to get average location
	numClients := float64(len(clientMap))

	fmt.Printf("Summed Lat: %d, Lon: %d, numClients: %d\n", currAvgLocation.Lat, currAvgLocation.Lon, numClients)
	currAvgLocation.Lat /= numClients
	currAvgLocation.Lon /= numClients
	fmt.Printf("Avgd Lat: %d, Lon: %d\n", currAvgLocation.Lat, currAvgLocation.Lon)
	// If it's empty then use it
	if (Location{}) == lastEventLocation {
		lastEventLocation = currAvgLocation
		return true
	}
	// Otherwise, check if they've travelled far enough
	dist := Delta(currAvgLocation, lastEventLocation)
	fmt.Printf("Distance from current to last loc: %d\n", dist)
	if dist >= EventDistance {
		lastEventLocation = currAvgLocation
		distanceTravelled += dist
		fmt.Printf("Distance travelled: %d\n", distanceTravelled)
		return true
	}
	// Otherwise don't update the event location and return false
	return false
}
