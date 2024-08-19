package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/jnafolayan/sip/pkg/codec"
	"github.com/jnafolayan/sip/pkg/wavelet"
)

var baseDir = "/Users/jnafolayan/workspace/projects/sip/resources"

func configForm(win fyne.Window, codec *codec.CodecOptions) *fyne.Container {
	wavelet := widget.NewSelect([]string{"haar", "cdf97"}, func(w string) {
		codec.Wavelet = wavelet.WaveletType(w)
	})
	wavelet.SetSelected(string(codec.Wavelet))

	levelLabelData := binding.NewInt()
	levelLabel := binding.NewSprintf("Level of decomposition (%d)", levelLabelData)
	level := widget.NewSlider(0, 10)
	level.Step = 1
	level.OnChanged = func(f float64) {
		codec.DecompositionLevel = int(f)
		levelLabelData.Set(int(f))
	}
	level.SetValue(float64(codec.DecompositionLevel))

	thresholdLabelData := binding.NewInt()
	thresholdLabel := binding.NewSprintf("Level of decomposition (%d)", thresholdLabelData)
	threshold := widget.NewSlider(0, 255)
	threshold.Step = 1
	threshold.OnChanged = func(f float64) {
		codec.ThresholdingFactor = int(f)
		thresholdLabelData.Set(int(f))
	}
	threshold.SetValue(float64(codec.ThresholdingFactor))

	thresholdStrategy := widget.NewSelect([]string{"soft", "hard"}, func(s string) {
		codec.ThresholdingStrategy = s
	})
	thresholdStrategy.SetSelected(codec.ThresholdingStrategy)

	return container.New(
		layout.NewFormLayout(),

		widget.NewLabel("Wavelet"),
		wavelet,

		widget.NewLabelWithData(levelLabel),
		level,

		widget.NewLabelWithData(thresholdLabel),
		threshold,

		widget.NewLabel("Thresholding strategy"),
		thresholdStrategy,
	)
}

func homeView(win fyne.Window) *fyne.Container {
	codecOpts := codec.CodecOptions{
		Wavelet:              "haar",
		ThresholdingFactor:   10,
		ThresholdingStrategy: "hard",
		DecompositionLevel:   1,
	}

	fileSelectDialog := dialog.NewFileOpen(func(f fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, win)
			return
		}

		if f == nil {
			return
		}

		win.SetContent(streamView(win, f.URI(), codecOpts))
	}, win)

	fileSelectDialog.Resize(fyne.NewSize(800, 640))
	fileSelectDialog.SetFilter(storage.NewMimeTypeFileFilter([]string{"image/*"}))
	l, _ := storage.ListerForURI(storage.NewFileURI(baseDir))
	fileSelectDialog.SetLocation(l)

	uploadButton := widget.NewButtonWithIcon("Select image to stream", theme.UploadIcon(), func() {
		fileSelectDialog.Show()
	})

	content := container.NewVBox(
		layout.NewSpacer(),
		container.NewHBox(
			layout.NewSpacer(),
			configForm(win, &codecOpts),
			layout.NewSpacer(),
		),
		container.NewHBox(
			layout.NewSpacer(),
			uploadButton,
			layout.NewSpacer(),
		),
		layout.NewSpacer(),
	)
	return content
}
