package main

import (
	"bittorrent/dht/library"
	"fmt"
	"fyne.io/fyne/app"
	"fyne.io/fyne/container"
	"strconv"
	"sync"
)

func StartGUI(database *library.DataBaseInMemory, barrier *sync.WaitGroup) {
	a := app.New()
	fmt.Println("App Started")              // Create a new application
	w := a.NewWindow("Chord Network State") // Create a new window
	gui := library.NewGUI(database, w)      // Create the GUI
	scrollContainer := container.Scroll{
		Content: gui.Grid,
	}
	// Set the grid layout as content
	w.SetContent(&scrollContainer)
	// Run state updates in a separate goroutine
	go func() {
		barrier.Add(1)
		gui.UpdateState()
		barrier.Done()
	}()

	w.ShowAndRun()
}

func main() {
	library.SetLogDirectoryPath("Main")
	var database = *library.NewDataBaseInMemory()
	var barrier = sync.WaitGroup{}
	N := 35
	for i := 0; i < N; i++ {
		iString := strconv.Itoa(i)
		var server = library.NewServerInMemory(&database, "Server"+iString)
		var client = library.NewClientInMemory(&database, "Client"+iString)
		node := library.NewBruteChord[library.InMemoryContact](server, client, library.NewMonitorHand[library.InMemoryContact]("Monitor"+iString))
		database.AddNode(node, server, client)
		barrier.Add(1)
		go func() {
			node.BeginWorking()
			defer barrier.Done()
		}()
	}
	fmt.Println("All nodes are started")
	StartGUI(&database, &barrier)
}
