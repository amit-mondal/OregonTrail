package main

type State int

const (
	WaitForGameStart State = 0
	WaitForCheckIn   State = 1
	WaitForReceive   State = 2
	WaitForDecision  State = 3
)
