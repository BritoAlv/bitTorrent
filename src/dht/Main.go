package main

import (
	"bittorrent/dht/library"
	"fmt"
	"fyne.io/fyne/app"
	"fyne.io/fyne/container"
	"math/rand"
	"strconv"
	"sync"
	"time"
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
	barrier.Add(1)
	go func() {
		gui.UpdateState()
		barrier.Done()
	}()
	w.ShowAndRun()
}

func AddNode(database *library.DataBaseInMemory, barrier *sync.WaitGroup) {
	randomId := library.GenerateRandomBinaryId()
	fmt.Println("Adding Node ", randomId)
	iString := strconv.Itoa(int(randomId))
	var server = library.NewServerInMemory(database, "Server"+iString)
	var client = library.NewClientInMemory(database, "Client"+iString)
	var monitor = library.NewMonitorHand[library.InMemoryContact]("Monitor" + iString)
	node := library.NewBruteChord[library.InMemoryContact](server, client, monitor, randomId)
	database.AddNode(node, server, client)
	barrier.Add(1)
	go func() {
		node.BeginWorking()
		defer barrier.Done()
	}()
}

func RemoveNode(database *library.DataBaseInMemory, barrier *sync.WaitGroup) {
	if len(database.GetNodes()) > 0 {
		for _, node := range database.GetNodes() {
			barrier.Add(1)
			go func() {
				fmt.Println("Removing Node with ID = ", node.GetId())
				database.RemoveNode(node)
				defer barrier.Done()
			}()
			break
		}
	}
}

func ScenarioBasico() (*library.DataBaseInMemory, *sync.WaitGroup) {
	library.SetLogDirectoryPath("Basic Scenario")
	var database = *library.NewDataBaseInMemory()
	var barrier = sync.WaitGroup{}
	fmt.Println("Nodes are being added and removed randomly every once a while")
	barrier.Add(1)
	go func() {
		defer barrier.Done()
		for {
			time.Sleep(1 * time.Second)
			if rand.Float32() <= 0.3 {
				AddNode(&database, &barrier)
			}
			if rand.Float32() <= 0.1 {
				RemoveNode(&database, &barrier)
			}
		}
	}()
	return &database, &barrier
}

func PutScenario() (*library.DataBaseInMemory, *sync.WaitGroup) {
	N := 10
	library.SetLogDirectoryPath("PutScenario")
	var database = *library.NewDataBaseInMemory()
	var barrier = sync.WaitGroup{}
	var toPut = make(map[library.ChordHash][]byte)
	fmt.Printf("Going to Add N = %v  Nodes", N)
	for i := 0; i < N; i++ {
		AddNode(&database, &barrier)
	}
	for i := 0; i < 50; i++ {
		toPut[library.ChordHash(i)] = []byte{byte(i)}
	}
	time.Sleep(3 * time.Second)
	nodes := database.GetNodes()
	println(len(nodes))
	for key, value := range toPut {
		nodes[0].Put(key, value)
	}
	return &database, &barrier
}

func main() {
	database, barrier := ScenarioBasico()
	StartGUI(database, barrier)
}
