package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	App := app.New()

	window := App.NewWindow("Real-ESRGAN-GUI")
	window.Resize(fyne.Size{
		Width:  600,
		Height: 400,
	})
	window.CenterOnScreen()

	window.ShowAndRun()
}
