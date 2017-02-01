package main

import (
	"flag"
	"fmt"
	"time"
	// For Native UI
	"github.com/andlabs/ui"
	// For Database
	"github.com/peterbourgon/diskv"
)

// The database instance
var db *diskv.Diskv

const GO_HOME_KEY = "go-time"

// Starting point of the app
func main() {
	// Read flags
	hour, min, message := readFlags()
	endTime := getWorkEndTime(hour, min)

	fmt.Println("Going Home at", endTime.Format(time.Kitchen), "after", endTime.Sub(time.Now()))

	// Create new ticker
	ticker := time.NewTicker(time.Minute * 1)
	go remainingTimeTicker(*ticker, endTime)

	// Let's keep the app running till working hour is complete
	<-time.After(time.Duration(hour * 60 + min) * time.Minute)

	// Completion hours are complete, let's stop the timer
	ticker.Stop()

	// Show dialog to GoHome !!!!
	showTimeCompleteWindow(message)
}

// Read the flags from command line
// h flag for hours
// m flag for minutes
// message flag for message to display on completion
func readFlags() (hour int, min int, message string) {
	hourPtr := flag.Int("h", 7, "Set the hour(s)")
	minPtr := flag.Int("m", 30, "Set the minute(s)")
	messagePtr := flag.String("message", "Time to close your work, buddy!", "Set the message(s)")
	flag.Parse()

	hour = *hourPtr
	min = *minPtr
	message = *messagePtr
	return
}

// Get the exact time when the working hours will end.
// Initially it checks for saved instance of end time
// in the database if found returns it else builds fresh end time.
// It returns the exact time (Ex. HH:MM => 05:30)
func getWorkEndTime(hour, min int) (time.Time) {
	t := time.Now()
	t = t.Add(time.Hour * time.Duration(hour))
	t = t.Add(time.Minute * time.Duration(min))


	db = diskv.New(diskv.Options{
		BasePath:     "time-store",
	})
	value, err := db.Read(GO_HOME_KEY)
	if err != nil {
		db.Write(GO_HOME_KEY, []byte(t.Format(time.UnixDate)))
		return t
	} else {
		current_time, _ := time.Parse(time.UnixDate, string(value))
		return current_time
	}
}

// Prints the remaining time periodically using Ticker
// endTime is the exact time
func remainingTimeTicker(ticker time.Ticker, endTime time.Time) {
	for t := range ticker.C {
		fmt.Println("Going Home at", endTime.Format(time.Kitchen), "after", endTime.Sub(t))
	}
}

// Show dialog to notify user about completion of his/her working hours
// and can go home to live :-P
// It shows the message provided in the parameter
func showTimeCompleteWindow(message string) {
	err := ui.Main(func() {
		home := ui.NewLabel(message)

		// Layout
		box := ui.NewVerticalBox()
		box.Append(home, false)

		// Window
		window := ui.NewWindow("Go Home!", 200, 100, false)
		window.SetChild(box)
		window.OnClosing(func(*ui.Window) bool {
			db.Erase(GO_HOME_KEY)
			ui.Quit()
			return true
		})
		window.Show()
	})
	if err != nil {
		panic(err)
	}
}