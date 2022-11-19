package main

import (
	"io/ioutil"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

type config struct {
	EditWidget *widget.Entry
	PreviewWidget *widget.RichText
	CurrentFile fyne.URI
	SaveMenuItem *fyne.MenuItem
}

var cfg config

func main() {
	// create the app
	a := app.New()
	a.Settings().SetTheme(&myTheme{})
	// create the window
	window := a.NewWindow("Markdown Editor")

	// get UI
	edit, preview := cfg.makeUI()
	cfg.createMenuItems(window)
	// set the window content
	window.SetContent(container.NewHSplit(edit, preview))
	// show the window and run the app
	window.Resize(fyne.NewSize(800, 600))
	window.CenterOnScreen()
	window.ShowAndRun()
}

func (app *config) makeUI() (*widget.Entry, *widget.RichText) {
	// create the widgets
	edit := widget.NewMultiLineEntry()
	preview := widget.NewRichTextFromMarkdown("")

	app.EditWidget = edit
	app.PreviewWidget = preview

	edit.OnChanged = preview.ParseMarkdown 
	// return the widgets 
	return edit, preview
}

func (app *config) createMenuItems(window fyne.Window) {
	// create the menu items
	open := fyne.NewMenuItem("Open", app.openFile(window))
	save := fyne.NewMenuItem("Save",  app.saveFile(window))
	app.SaveMenuItem = save
	app.SaveMenuItem.Disabled = true
	saveAs := fyne.NewMenuItem("Save As", app.saveAsFile(window))
	quit := fyne.NewMenuItem("Quit", func() {
		cfg.openFile(window)
	})

	cfg.SaveMenuItem = save

	// create the menu
	fileMenu := fyne.NewMenu("File", open, save, saveAs, quit)
	// set the menu
	window.SetMainMenu(fyne.NewMainMenu(fileMenu))
}

var filter = storage.NewExtensionFileFilter([]string{".md", ".MD"})

// saveAs function
func (app *config) saveAsFile(window fyne.Window) func(){
	return func ()  {
		saveDialog := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err != nil {
				dialog.ShowError(err, window)
				return
			}
			if writer == nil {
				return
			}

			if !strings.HasSuffix(strings.ToLower(writer.URI().String()), ".md") {
				dialog.ShowInformation(
					"Invalid file extension",
					"Please use a .md file extension",
					window,
				)
			}

			writer.Write([]byte(app.EditWidget.Text))
			app.CurrentFile = writer.URI()
			defer writer.Close()

			window.SetTitle("Markdown Editor - " + app.CurrentFile.Name())

			app.SaveMenuItem.Disabled = false
		}, window)
		saveDialog.SetFileName("Untitled.md")
		saveDialog.SetFilter(filter)
		saveDialog.Show()
	}
}

// open file function
func (app *config) openFile(window fyne.Window) func() {
	return func() {
		openDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, window)
				return
			}
			if reader == nil {
				return
			}
			defer reader.Close()

			data, err := ioutil.ReadAll(reader)
			if err != nil {
				dialog.ShowError(err, window)
				return
			}

			app.EditWidget.SetText(string(data))

			app.CurrentFile = reader.URI()
			window.SetTitle("Markdown Editor - " + app.CurrentFile.Name())
			app.SaveMenuItem.Disabled = false
		}, window)
		openDialog.SetFilter(filter)
		openDialog.Show()
	}
}

// save file func
func (app *config) saveFile(window fyne.Window) func() {
	return func() {
		if app.CurrentFile != nil {
			writer, err := storage.Writer(app.CurrentFile)
			if err != nil {
				dialog.ShowError(err, window)
				return
			}
			writer.Write([]byte(app.EditWidget.Text))
			defer writer.Close()
		}
	}
}