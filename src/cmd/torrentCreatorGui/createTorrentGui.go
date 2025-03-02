package main

import (
	"bittorrent/torrent"
	"encoding/json"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"os"
)

type TorrentInput struct {
	FilePath        string `json:"filePath"`
	TorrentName     string `json:"torrentName"`
	TrackerLocation string `json:"trackerLocation"`
}

const inputFile = "inputs.json"

func loadInputs() ([]TorrentInput, error) {
	var inputs []TorrentInput
	data, err := os.ReadFile(inputFile)
	if err != nil {
		if os.IsNotExist(err) {
			return inputs, nil
		}
		return nil, err
	}
	err = json.Unmarshal(data, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func saveInput(input TorrentInput) error {
	inputs, err := loadInputs()
	if err != nil {
		return err
	}
	inputs = append(inputs, input)
	data, err := json.Marshal(inputs)
	if err != nil {
		return err
	}
	err = os.WriteFile(inputFile, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	myApp := app.NewWithID("Torrent Creator")
	myWindow := myApp.NewWindow("Torrent File Creator")
	myWindow.Resize(fyne.NewSize(600, 400)) // Larger window for better UI

	filePathEntry := widget.NewEntry()
	filePathEntry.SetPlaceHolder("Select a file...")

	torrentNameEntry := widget.NewEntry()
	torrentNameEntry.SetPlaceHolder("Enter the torrent name")

	trackerLocationEntry := widget.NewEntry()
	trackerLocationEntry.SetPlaceHolder("Enter the tracker location")

	resultLabel := widget.NewLabel("")

	// Improved File Selection Button
	selectFileButton := widget.NewButton("Browse...", func() {
		fileDialog := dialog.NewFileOpen(func(uri fyne.URIReadCloser, err error) {
			if uri != nil {
				filePathEntry.SetText(uri.URI().Path())
			}
		}, myWindow)
		fileDialog.Show()
	})

	// Load previous inputs
	inputs, err := loadInputs()
	if err != nil {
		resultLabel.SetText("Error loading inputs: " + err.Error())
	}

	// List of previously saved inputs
	inputList := widget.NewList(
		func() int {
			return len(inputs)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(inputs[i].FilePath + " - " + inputs[i].TorrentName + " - " + inputs[i].TrackerLocation)
		},
	)

	// Handle list selection
	inputList.OnSelected = func(id widget.ListItemID) {
		selectedInput := inputs[id]
		filePathEntry.SetText(selectedInput.FilePath)
		torrentNameEntry.SetText(selectedInput.TorrentName)
		trackerLocationEntry.SetText(selectedInput.TrackerLocation)
		inputList.Unselect(id)
	}

	// Submit Button
	submitButton := widget.NewButton("Create Torrent", func() {
		filePath := filePathEntry.Text
		torrentName := torrentNameEntry.Text
		trackerLocation := trackerLocationEntry.Text

		err := torrent.CreateTorrentFile(filePath, torrentName, trackerLocation)
		if err != nil {
			resultLabel.SetText("Error: " + err.Error())
		} else {
			resultLabel.SetText("Torrent file created successfully")
			newInput := TorrentInput{FilePath: filePath, TorrentName: torrentName, TrackerLocation: trackerLocation}
			saveErr := saveInput(newInput)
			if saveErr != nil {
				resultLabel.SetText("Error saving input: " + saveErr.Error())
			}
		}
	})

	// UI Layout
	content := container.NewVBox(
		widget.NewLabel("Select a file:"),
		container.NewBorder(nil, nil, nil, selectFileButton, filePathEntry), // Entry with button in the same row
		widget.NewLabel("Torrent Name:"),
		torrentNameEntry,
		widget.NewLabel("Tracker Location:"),
		trackerLocationEntry,
		submitButton,
		resultLabel,
		inputList,
	)
	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}
