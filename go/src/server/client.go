package main

import (
	"math"
)

type Location struct {
	Lat float64 `json:",string"`
	Lon float64 `json:",string"`
}

type Client struct {
	Id       string   `json:"id"`
	Location Location `json:"location"`
	IsAlive  bool     `json:"is_alive"`
	Food     int      `json:"food"`
	Water    int      `json:"water"`
	Bullets  int      `json:"bullets"`
	Supplies int      `json:"supplies"`
}

func Delta(l1 *Location, l2 *Location) float64 {
	return math.Hypot(l2.Lat-l1.Lat, l2.Lon-l1.Lon)
}
