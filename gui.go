package main

import (
	"fmt"
	"path/filepath"

	"github.com/andlabs/ui"
)

// SimpleGui - using if main program run without arguments
func SimpleGui() {
	err := ui.Main(func() {
		packMsg := fmt.Sprintf("Pack mode:\nselect any file in target dir")
		unpackMsg := fmt.Sprintf("Unpack mode:\nselect file of import")
		btnPack := ui.NewButton(packMsg)
		btnUnpack := ui.NewButton(unpackMsg)
		box := ui.NewHorizontalBox()
		box.Append(btnPack, true)
		box.Append(btnUnpack, true)
		window := ui.NewWindow("ap operation helper", 300, 100, false)
		window.SetMargined(true)
		window.SetChild(box)
		btnPack.OnClicked(func(*ui.Button) {
			wSelectFile := ui.NewWindow("select dir", 300, 100, false)
			inputName = ui.OpenFile(wSelectFile)

			if inputName != "" {
				inputName = filepath.Dir(inputName) + separator
				checkOutput()
				ui.Quit()
			}
		})
		btnUnpack.OnClicked(func(*ui.Button) {
			unpack = true
			wSelectFile := ui.NewWindow("select file", 300, 100, false)
			inputName = ui.OpenFile(wSelectFile)

			if inputName != "" {
				checkOutput()
				ui.Quit()
			}
		})
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
