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
		window:         window,
		nodeLabelsCard: make(map[int]LabelCard),
		Grid:           container.NewGridWithColumns(4), // Adjust columns as needed
	}

	window.SetContent(gui.Grid)
	window.Resize(fyne.NewSize(800, 600))
	return gui
}

func (g *GUI) UpdateState() {
	for {
		time.Sleep(2 * time.Second)
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
		objects := g.Grid.Objects
		g.Grid.Objects = nil
		sort.Slice(objects, func(i, j int) bool {
			one, _ := strconv.Atoi(objects[i].(*widget.Card).Title[5:])
			two, _ := strconv.Atoi(objects[j].(*widget.Card).Title[5:])

			return one < two
		})
		for _, object := range objects {
			g.Grid.Add(object)
		}

		g.window.Content().Refresh()
	}
}

func (g *GUI) PrepareState() map[int]string {
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
	return node.State()
}
