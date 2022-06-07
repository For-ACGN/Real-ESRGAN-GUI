package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"

	"github.com/For-ACGN/Real-ESRGAN-GUI/theme"
)

func main() {
	wd, _ := os.Getwd()

	App := app.New()
	App.Settings().SetTheme(new(theme.Chinese))

	window := App.NewWindow("Real-ESRGAN-GUI")
	window.Resize(fyne.Size{Width: 728, Height: 470})
	window.SetFixedSize(true)
	window.CenterOnScreen()

	inputLab := widget.NewLabel("输入图片路径:")
	inputLab.Move(fyne.NewPos(0, 0))
	inputLab.Resize(fyne.Size{Width: 10, Height: 10})

	inputPath := widget.NewEntry()
	inputPath.Move(fyne.NewPos(110, 3))
	inputPath.Resize(fyne.Size{Width: 500, Height: 38})
	inputPath.SetText(filepath.Join(wd, "input.png"))

	inputButton := widget.NewButton("浏览图片", func() {
		win := App.NewWindow("打开图片")
		win.Resize(fyne.Size{Width: 800, Height: 600})
		win.SetFixedSize(true)
		win.CenterOnScreen()
		win.Show()

		fd := dialog.NewFileOpen(func(url fyne.URIReadCloser, err error) {
			if url == nil || err != nil {
				return
			}
			path := strings.Replace(url.URI().String(), "file://", "", 1)
			path = filepath.Clean(path)
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
	outputPath.SetText(filepath.Join(wd, "output.png"))

	outputButton := widget.NewButton("浏览图片", func() {
		win := App.NewWindow("保存图片")
		win.Resize(fyne.Size{Width: 800, Height: 600})
		win.SetFixedSize(true)
		win.CenterOnScreen()
		win.Show()

		fd := dialog.NewFileSave(func(url fyne.URIWriteCloser, err error) {
			if url == nil || err != nil {
				return
			}
			path := strings.Replace(url.URI().String(), "file://", "", 1)
			path = filepath.Clean(path)
			outputPath.SetText(path)

			win.Close()
		}, win)
		fd.Resize(fyne.Size{Width: 800, Height: 600})
		fd.SetFileName("output.png")
		filter := storage.NewExtensionFileFilter([]string{".png", ".jpg", ".jpeg", ".bmp"})
		fd.SetFilter(filter)
		wd, _ := os.Getwd()
		loc, _ := storage.ListerForURI(storage.NewFileURI(wd))
		fd.SetLocation(loc)

		fd.Show()
	})
	outputButton.Move(inputButton.Position().Add(fyne.NewDelta(0, 50)))
	outputButton.Resize(inputButton.Size())

	scaleLab := widget.NewLabel("图片放大倍率:")
	scaleLab.Move(fyne.NewPos(0, 105))
	scaleLab.Resize(fyne.Size{Width: 10, Height: 10})

	scaleText := widget.NewEntry()
	scaleText.SetText("4")
	scaleText.Move(fyne.NewPos(110, 105))
	scaleText.Resize(fyne.Size{Width: 33, Height: 38})

	modelLab := widget.NewLabel("模型:")
	modelLab.Move(fyne.NewPos(150, 105))
	modelLab.Resize(fyne.Size{Width: 10, Height: 10})

	modelOptions := []string{"realesrgan-x4plus-anime", "realesrgan-x4plus"}
	modelSelect := widget.NewSelectEntry(modelOptions)
	modelSelect.Move(fyne.NewPos(200, 105))
	modelSelect.Resize(fyne.Size{Width: 235, Height: 38})
	modelSelect.SetText(modelOptions[0])

	formatLab := widget.NewLabel("输出格式:")
	formatLab.Move(fyne.NewPos(446, 105))
	formatLab.Resize(fyne.Size{Width: 10, Height: 10})

	formatOptions := []string{"png", "jpg", "wepb"}
	formatSelect := widget.NewSelectEntry(formatOptions)
	formatSelect.Move(fyne.NewPos(530, 105))
	formatSelect.Resize(fyne.Size{Width: 80, Height: 38})
	formatSelect.SetText(formatOptions[0])

	runText := widget.NewMultiLineEntry()
	runText.Move(fyne.NewPos(10, 155))
	runText.Resize(fyne.Size{Width: 700, Height: 297})
	runText.Validator = nil

	startButton := widget.NewButton("开始", nil)
	startButton.OnTapped = func() {
		r, w, _ := os.Pipe()
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() { _ = w.Close() }()
			cmd := exec.Command("realesrgan-ncnn-vulkan.exe",
				"-i", inputPath.Text, "-o", outputPath.Text,
				"-s", scaleText.Text, "-n", modelSelect.Text,
			)
			cmd.SysProcAttr = &syscall.SysProcAttr{
				HideWindow: true,
			}
			cmd.Stdout = w
			cmd.Stderr = w
			_ = cmd.Run()

		}()
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() { _ = r.Close() }()
			buf := make([]byte, 512)
			for {
				n, err := r.Read(buf)
				if err != nil {
					return
				}
				runText.SetText(runText.Text + string(buf[:n]))
				runText.CursorRow++
				runText.Refresh()
			}
		}()
		go func() {
			wg.Wait()
			runText.SetText(runText.Text + "Finish\n")
			runText.CursorRow += 2
			runText.Refresh()
			startButton.Enable()
		}()
		startButton.Disable()
	}
	startButton.Move(fyne.NewPos(630, 105))
	startButton.Resize(fyne.Size{Width: 80, Height: 38})

	cont := container.NewWithoutLayout()
	cont.Add(inputLab)
	cont.Add(inputPath)
	cont.Add(inputButton)
	cont.Add(outputLab)
	cont.Add(outputPath)
	cont.Add(outputButton)
	cont.Add(scaleLab)
	cont.Add(scaleText)
	cont.Add(modelLab)
	cont.Add(modelSelect)
	cont.Add(formatLab)
	cont.Add(formatSelect)
	cont.Add(runText)
	cont.Add(startButton)
	window.SetContent(cont)

	window.ShowAndRun()
}
