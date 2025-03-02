package Manager

import (
	"bittorrent/dht/library/BruteChord/Core"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"sort"
	"strconv"
	"time"
)

type labelCard struct {
	Label *widget.Label
	Card  *widget.Card
}

type GUI struct {
	manager        IManagerRPC
	window         fyne.Window
	nodeLabelsCard map[Core.ChordHash]labelCard // Map of node ID -> Label for dynamic updates
	grid           *fyne.Container
	paused         bool
	pausedButton   *widget.Button
}

func NewGUI(manager IManagerRPC) *GUI {
	a := app.NewWithID("ChordNetworkState") // Create a new application
	window := a.NewWindow("Chord Network State")
	gui := &GUI{
		manager:        manager,
		window:         window,
		nodeLabelsCard: make(map[Core.ChordHash]labelCard),
		grid:           container.NewGridWithColumns(4),
		paused:         false,
	}

	gui.pausedButton = widget.NewButton("Pause", func() {
		gui.paused = !gui.paused
		if gui.paused {
			gui.pausedButton.SetText("Resume")
		} else {
			gui.pausedButton.SetText("Pause")
		}
	})

	content := container.NewVBox(gui.pausedButton, gui.grid)
	scrollContainer := container.NewScroll(content)
	window.SetContent(scrollContainer)
	window.Resize(fyne.NewSize(800, 600))

	return gui
}

func (g *GUI) updateState() {
	for {
		time.Sleep(1 * time.Second)
		if g.paused {
			continue
		}
		stateMap := g.prepareState()

		// Get the keys and sort them
		keys := make([]Core.ChordHash, 0, len(stateMap))
		for k := range stateMap {
			keys = append(keys, k)
		}
		Core.Sort(keys)

		// Iterate over the sorted keys
		for _, nodeID := range keys {
			state := stateMap[nodeID]
			if labelC, exists := g.nodeLabelsCard[nodeID]; exists {
				labelC.Label.SetText(state)
			} else {
				// Create a new card for the node if it doesn't exist
				label := widget.NewLabel(state)

				fixedSizeContainer := container.NewScroll(container.NewVBox(label))
				fixedSizeContainer.SetMinSize(fyne.NewSize(250, 250))

				card := widget.NewCard(fmt.Sprintf("Node %d", nodeID), "", fixedSizeContainer)
				g.nodeLabelsCard[nodeID] = labelCard{
					Label: label,
					Card:  card,
				}
				g.grid.Add(card)

			}
		}
		for nodeID, labelCard := range g.nodeLabelsCard {
			if _, exists := stateMap[nodeID]; !exists {
				g.grid.Remove(labelCard.Card)
				delete(g.nodeLabelsCard, nodeID)
			}
		}
		objects := g.grid.Objects
		g.grid.Objects = nil
		sort.Slice(objects, func(i, j int) bool {
			one, _ := strconv.Atoi(objects[i].(*widget.Card).Title[5:])
			two, _ := strconv.Atoi(objects[j].(*widget.Card).Title[5:])
			return one < two
		})
		for _, object := range objects {
			g.grid.Add(object)
		}
		g.window.Content().Refresh()
	}
}

func (g *GUI) prepareState() map[Core.ChordHash]string {
	result := make(map[Core.ChordHash]string)
	for _, key := range g.manager.GetActiveNodesIds() {
		result[key] = g.manager.GetNodeStateRPC(key) // this is an RPC.
	}
	return result
}

func (g *GUI) Start() {
	// Create a new window // Run state updates in a separate goroutine
	go g.updateState()
	g.window.ShowAndRun()
}
