package library

import (
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/widget"
	"sort"
	"strconv"
	"time"
)

type GUI struct {
	db         *DataBaseInMemory
	iteration  int
	window     fyne.Window
	nodeLabels map[int]*widget.Label // Map of node ID -> Label for dynamic updates
	Grid       *fyne.Container
}

func NewGUI(db *DataBaseInMemory, window fyne.Window) *GUI {
	gui := &GUI{
		db:         db,
		iteration:  0,
		window:     window,
		nodeLabels: make(map[int]*widget.Label),
		Grid:       container.NewGridWithColumns(5), // Adjust columns as needed
	}

	window.SetContent(gui.Grid)
	return gui
}

func (g *GUI) UpdateState() {
	for {
		time.Sleep(2 * time.Second) // Replace `StateQueryWaitTime` with an actual value

		stateMap := g.PrepareState()

		// Get the keys and sort them
		keys := make([]int, 0, len(stateMap))
		for k := range stateMap {
			keys = append(keys, k)
		}
		sort.Ints(keys)

		// Iterate over the sorted keys
		for _, nodeID := range keys {
			state := stateMap[nodeID]
			if label, exists := g.nodeLabels[nodeID]; exists {
				label.SetText(state)
			} else {
				// Create a new card for the node if it doesn't exist
				label := widget.NewLabel(state)
				card := widget.NewCard(fmt.Sprintf("Node %d", nodeID), "", label)
				g.nodeLabels[nodeID] = label
				g.Grid.Add(card)
			}
		}
		g.window.Content().Refresh()
	}
}

func (g *GUI) PrepareState() map[int]string {
	g.iteration++
	nodes := g.db.GetNodes()
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].GetId() < nodes[j].GetId()
	})

	stateMap := make(map[int]string)
	for _, node := range nodes {
		stateMap[int(node.GetId())] = g.ShowNodeState(node)
	}

	return stateMap
}

func (g *GUI) ShowNodeState(node *BruteChord[InMemoryContact]) string {
	state := "Node: " + strconv.Itoa(int(node.GetId())) + "\n"
	state += "Successor: " + strconv.Itoa(int(node.GetSuccessor().getNodeId())) + "\n"
	state += "Successor Data Replica: " + "\n"
	for key, value := range node.GetSuccessorReplicatedData() {
		state += strconv.Itoa(int(key)) + " -> " + fmt.Sprintf("%v", value) + "\n"
	}
	state += "SuccessorSuccessor: " + strconv.Itoa(int(node.GetSuccessorSuccessor().getNodeId())) + "\n"
	state += "SuccessorSuccessor Data Replica:" + "\n"
	for key, value := range node.GetSuccessorSuccessorReplicatedData() {
		state += strconv.Itoa(int(key)) + " -> " + fmt.Sprintf("%v", value) + "\n"
	}
	state += "Predecessor: " + strconv.Itoa(int(node.GetPredecessor().getNodeId())) + "\n"
	state += "Data stored:\n"
	for key, value := range node.GetAllOwnData() {
		state += strconv.Itoa(int(key)) + " -> " + fmt.Sprintf("%v", value) + "\n"
	}
	return state
}
