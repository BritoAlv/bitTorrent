package library

import (
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/widget"
	"sort"
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
	return gui
}

func (g *GUI) UpdateState() {
	for {
		time.Sleep(2 * time.Second) // Replace with actual interval
		stateMap := g.PrepareState()

		// Sort node IDs
		keys := make([]int, 0, len(stateMap))
		for k := range stateMap {
			keys = append(keys, k)
		}
		sort.Ints(keys)

		// Store existing scroll positions
		scrollPositions := make(map[int]int)
		for nodeID, labelCard := range g.nodeLabelsCard {
			if scroll, ok := labelCard.Card.Content.(*container.Scroll); ok {
				scrollPositions[nodeID] = scroll.Offset.Y
			}
		}

		// Clear the grid but keep references
		g.Grid.Objects = nil

		// Re-add nodes in sorted order
		for _, nodeID := range keys {
			state := stateMap[nodeID]
			var labelCard LabelCard
			if existingCard, exists := g.nodeLabelsCard[nodeID]; exists {
				existingCard.Label.SetText(state)
				labelCard = existingCard
			} else {
				// Create a new UI component if it doesn't exist
				label := widget.NewLabel(state)
				scrollContainer := container.NewScroll(container.NewVBox(label))
				scrollContainer.SetMinSize(fyne.NewSize(250, 250))

				card := widget.NewCard(fmt.Sprintf("Node %d", nodeID), "", scrollContainer)
				labelCard = LabelCard{Label: label, Card: card}
				g.nodeLabelsCard[nodeID] = labelCard
			}
			// Restore scroll position if it existed
			if scroll, ok := labelCard.Card.Content.(*container.Scroll); ok {
				if pos, found := scrollPositions[nodeID]; found {
					scroll.Offset.Y = pos
				}
			}
			g.Grid.Add(labelCard.Card)
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