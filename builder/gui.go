package builder

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

func HandleGui(ctx context.Context, srv *http.Server, wg *sync.WaitGroup, cwd string) error {
	app := app.New()
	window := app.NewWindow("Markdown Resume Builder")

	sourceFile := widget.NewEntry()
	sourceFile.Validator = validation.NewRegexp(`[a-zA-Z]:[\\\/](?:[a-zA-Z0-9-_]+[\\\/])*([a-zA-Z0-9-_]+\.md)`, "not a valid Windows Markdown path")

	openFile := widget.NewButton("Browse Markdown Files", func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, window)
				return
			}
			if reader == nil {
				return
			}
			chosenURI := reader.URI().String()[7:]
			sourceFile.SetText(chosenURI)
		}, window)
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".md"}))
		fd.Show()
	})

	mode := widget.NewRadioGroup([]string{"HTML", "PDF"}, func(string) {})
	mode.SetSelected("HTML")
	mode.Horizontal = true
	mode.Required = true

	destinationFolder := widget.NewEntry()
	destinationFolder.Validator = validation.NewRegexp(`[a-zA-Z]:[\\\/](?:[a-zA-Z0-9-_]+[\\\/])*([a-zA-Z0-9-_])`, "not a valid Windows Folder path")

	openFolder := widget.NewButton("Browse Folders", func() {
		fd := dialog.NewFolderOpen(func(reader fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, window)
				return
			}
			if reader == nil {
				return
			}
			chosenURI := reader.String()[7:]
			destinationFolder.SetText(chosenURI)
		}, window)
		fd.Show()
	})

	fontAwesomeKit := widget.NewEntry()
	fontAwesomeKit.SetPlaceHolder("https://kit.fontawesome.com/placeholder.js")

	photoFile := widget.NewEntry()

	openPhoto := widget.NewButton("Browse Photos", func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, window)
				return
			}
			if reader == nil {
				return
			}
			chosenURI := reader.URI().String()[7:]
			photoFile.SetText(chosenURI)
		}, window)
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".jpg", ".jpeg"}))
		fd.Show()
	})

	separator := widget.NewSeparator()

	statusText := widget.NewLabel("")
	statusText.Alignment = fyne.TextAlignCenter
	statusText.TextStyle = fyne.TextStyle{Bold: true}

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Resume Markdown File", Widget: sourceFile},
			{Widget: openFile, HintText: "Browse your computer to find the Markdown file containing your resume"},
			{Widget: separator},
			{Text: "Output format", Widget: mode, HintText: "The output format of your resume"},
			{Widget: separator},
			{Text: "Output folder", Widget: destinationFolder},
			{Widget: openFolder, HintText: "Specify the folder where you want to save your resume."},
			{Widget: separator},
			{Text: "Font Awesome Kit (Optional)", Widget: fontAwesomeKit, HintText: "Your Font Awesome Kit URL, in case you want to use icons"},
			{Widget: separator},
			{Text: "Resume Photo (Optional)", Widget: photoFile},
			{Widget: openPhoto, HintText: "Browse your computer to find the your photo"},
			{Widget: separator},
		},
		OnCancel: func() {
			fmt.Println("Cancelled")
		},
	}

	form.OnSubmit = func() {
		statusText.SetText("Building resume...")
		form.Disable()
		err := buildResume(ctx, srv, wg, cwd, strings.TrimSpace(sourceFile.Text), mode.Selected, strings.TrimSpace(destinationFolder.Text), strings.TrimSpace(fontAwesomeKit.Text), strings.TrimSpace(photoFile.Text))
		form.Enable()
		if err != nil {
			statusText.SetText(fmt.Sprintf("Error building resume: %v", err))
		} else {
			statusText.SetText(fmt.Sprintf("%s resume built successfully", mode.Selected))
		}
	}

	grid := container.New(layout.NewGridLayout(1), form, statusText)
	window.SetContent(grid)
	window.Resize(fyne.NewSize(0, 0))
	window.ShowAndRun()
	return nil
}
