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

type LabelCard struct {
	Label *widget.Label
	Card  *widget.Card
}

type GUI struct {
	db             *DataBaseInMemory
	iteration      int
	window         fyne.Window
	nodeLabelsCard map[int]LabelCard // Map of node ID -> Label for dynamic updates
	Grid           *fyne.Container
}

func NewGUI(db *DataBaseInMemory, window fyne.Window) *GUI {
	gui := &GUI{
		db:             db,
		iteration:      0,
		window:         window,
		nodeLabelsCard: make(map[int]LabelCard),
		Grid:           container.NewGridWithColumns(4), // Adjust columns as needed
	}

	window.SetContent(gui.Grid)
	return gui
}

func (g *GUI) ClearGrid() {
	g.Grid.Objects = []fyne.CanvasObject{}
	g.nodeLabelsCard = make(map[int]LabelCard)
	g.Grid.Refresh()
}

func (g *GUI) UpdateState() {
	for {
		time.Sleep(2 * time.Second) // Replace `StateQueryWaitTime` with an actual value
		g.ClearGrid()
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
			if labelCard, exists := g.nodeLabelsCard[nodeID]; exists {
				labelCard.Label.SetText(state)
			} else {
				// Create a new card for the node if it doesn't exist
				label := widget.NewLabel(state)

				fixedSizeContainer := container.NewScroll(container.NewVBox(label))
				fixedSizeContainer.SetMinSize(fyne.NewSize(250, 250))

				card := widget.NewCard(fmt.Sprintf("Node %d", nodeID), "", fixedSizeContainer)
				g.nodeLabelsCard[nodeID] = LabelCard{
					Label: label,
					Card:  card,
				}
				g.Grid.Add(card)

			}
		}
		for nodeID, labelCard := range g.nodeLabelsCard {
			if _, exists := stateMap[nodeID]; !exists {
				g.Grid.Remove(labelCard.Card)
				delete(g.nodeLabelsCard, nodeID)
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
	state += "Successor Data Replicas Are: " + "\n"
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
