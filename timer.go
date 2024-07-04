// Timer struct and methods
// for games written in Go using Ebitengine
// call UpdateAllTimers once every frame
// Author Paul Brace
// July 2024

package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"time"
)

var timers [] *Timer

func UpdateAllTimers(){
	for _, t := range timers {
		t.Update()
	} 
}

// Timer struct for use in games ticks updated every frame
type Timer struct {
	currentTicks 	int
	targetTicks  	int
	active 			bool
	repeat			bool
}

// Creates a new timer if repeating is false then one off if true the triggers every Duration
func NewTimer(seconds float64, repeating bool) *Timer {
	d := time.Duration(seconds) * time.Second
	timer := Timer{
		currentTicks: 	0,
		targetTicks:  	int(d.Milliseconds()) * ebiten.TPS() / 1000,
		active: 		true,
		repeat:			repeating,
	}
	timers = append(timers, &timer)
	return &timer
}

// Updates a single timer (call every frame)
func (t *Timer) Update() {
	if t.active{
		if t.currentTicks < t.targetTicks {
			t.currentTicks++
		}
	}
}

// Checks if the timer has expired and if it has reset it
func (t *Timer) IsReady() bool {
	hasExpired := t.active && t.currentTicks >= t.targetTicks
	if hasExpired {
		t.Reset()
	}
	return hasExpired
}

// Called to reset for repeating timers or deactivate for one off
func (t *Timer) Reset() {
	t.currentTicks = 0
	t.active = t.repeat
}

// function to allow the timer interval to be changed
func (t *Timer) ChangeTime(seconds float64, repeating bool) {
	d := time.Duration(seconds) * time.Second
	t.currentTicks = 0
	t.targetTicks =	int(d.Milliseconds()) * ebiten.TPS() / 1000
	t.active = true
	t.repeat = repeating
}


