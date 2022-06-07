package main

import (
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"

	"github.com/For-ACGN/Real-ESRGAN-GUI/theme"
)

func main() {
	App := app.New()
	App.Settings().SetTheme(new(theme.Chinese))

	window := App.NewWindow("Real-ESRGAN-GUI")
	window.Resize(fyne.Size{
		Width:  800,
		Height: 600,
	})
	window.CenterOnScreen()

	inputLab := widget.NewLabel("输入图片路径:")
	inputLab.Move(fyne.NewPos(0, 0))
	inputLab.Resize(fyne.Size{Width: 10, Height: 10})

	inputPath := widget.NewEntry()
	inputPath.Move(fyne.NewPos(110, 3))
	inputPath.Resize(fyne.Size{Width: 500, Height: 38})

	inputButton := widget.NewButton("浏览文件", func() {
		win := App.NewWindow("打开文件")
		win.Resize(fyne.Size{Width: 800, Height: 600})
		win.SetFixedSize(true)
		win.CenterOnScreen()
		win.Show()

		fd := dialog.NewFileOpen(func(url fyne.URIReadCloser, err error) {
			if url == nil || err != nil {
				return
			}
			path := strings.Replace(url.URI().String(), "file://", "", 1)
			inputPath.SetText(path)

			win.Close()
		}, win)
		fd.Resize(fyne.Size{Width: 800, Height: 600})

		filter := storage.NewExtensionFileFilter([]string{".png", ".jpg", ".jpeg", ".bmp"})
		fd.SetFilter(filter)
		wd, _ := os.Getwd()
		loc, _ := storage.ListerForURI(storage.NewFileURI(wd))
		fd.SetLocation(loc)

		fd.Show()
	})
	inputButton.Move(fyne.NewPos(630, 3))
	inputButton.Resize(fyne.Size{Width: 80, Height: 38})

	outputLab := widget.NewLabel("输出图片路径:")
	outputLab.Move(inputLab.Position().Add(fyne.NewDelta(0, 50)))
	outputLab.Resize(inputLab.Size())

	outputPath := widget.NewEntry()
	outputPath.Move(inputPath.Position().Add(fyne.NewDelta(0, 50)))
	outputPath.Resize(inputPath.Size())

	outputButton := widget.NewButton("浏览文件", func() {
		win := App.NewWindow("保存文件")
		win.Resize(fyne.Size{Width: 800, Height: 600})
		win.SetFixedSize(true)
		win.CenterOnScreen()
		win.Show()

		fd := dialog.NewFileSave(func(url fyne.URIWriteCloser, err error) {
			if url == nil || err != nil {
				return
			}
			path := strings.Replace(url.URI().String(), "file://", "", 1)
			outputPath.SetText(path)

			win.Close()
		}, win)
		fd.Resize(fyne.Size{Width: 800, Height: 600})

		filter := storage.NewExtensionFileFilter([]string{".png", ".jpg", ".jpeg", ".bmp"})
		fd.SetFilter(filter)
		wd, _ := os.Getwd()
		loc, _ := storage.ListerForURI(storage.NewFileURI(wd))
		fd.SetLocation(loc)

		fd.Show()
	})
	outputButton.Move(inputButton.Position().Add(fyne.NewDelta(0, 50)))
	outputButton.Resize(inputButton.Size())

	cont := container.NewWithoutLayout()
	cont.Add(inputLab)
	cont.Add(inputPath)
	cont.Add(inputButton)
	cont.Add(outputLab)
	cont.Add(outputPath)
	cont.Add(outputButton)
	window.SetContent(cont)

	window.ShowAndRun()
}
