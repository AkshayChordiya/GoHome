package main

import (
	"flag"
	"fmt"
	"time"
	"github.com/andlabs/ui"
)

func main() {
	hourPtr := flag.Int("h", 7, "Set the hour(s)")
	minPtr := flag.Int("m", 30, "Set the minute(s)")
	flag.Parse()

	hour := *hourPtr
	min := *minPtr

	endTime := getWorkEndTime(hour, min)

	fmt.Printf("Will remind you after %d hour and %d minute", hour, min)
	hour = hour * 60 + min

	fmt.Println()

	ticker := time.NewTicker(time.Second * 5)
	go func() {
		for t := range ticker.C {
			fmt.Println("Remaining Time: ", endTime.Sub(t))
		}
	}()
	<-time.After(time.Duration(hour) * time.Minute)
	ticker.Stop()

	showTimeCompleteWindow()
}

func showTimeCompleteWindow() {
	err := ui.Main(func() {
		home := ui.NewLabel("Time to close your work, buddy!")

		// Layout
		box := ui.NewVerticalBox()
		box.Append(home, false)

		// Window
		window := ui.NewWindow("Go Home!", 200, 100, false)
		window.SetChild(box)
		window.OnClosing(func(*ui.Window) bool {
			ui.Quit()
			return true
		})
		window.Show()
	})
	if err != nil {
		panic(err)
	}
}

func getWorkEndTime(hour, min int) (time.Time) {
	t := time.Now()
	t = t.Add(time.Hour * time.Duration(hour))
	t = t.Add(time.Minute * time.Duration(min))
	return t
}