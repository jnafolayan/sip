package main

import (
	"fmt"
	"image"
	"image/color"
	"syscall/js"
	"time"

	"github.com/jnafolayan/sip/pkg/codec"
	"github.com/jnafolayan/sip/pkg/wavelet"
)

func getImageChannels(this js.Value, args []js.Value) interface{} {
	imageData := args[0]
	width, height := args[1].Int(), args[2].Int()
	img := convertImageDataToImage(imageData, width, height)

	return img
}

func compressImage(this js.Value, args []js.Value) interface{} {
	// imageData, width, height, opts
	jsImageData := args[0]
	width, height := args[1].Int(), args[2].Int()
	opts := args[3]

	imageData := make([]uint8, width*height*4)
	js.CopyBytesToGo(imageData, jsImageData)

	waveletFamily := opts.Get("waveletFamily").String()
	decompLevel := opts.Get("decompLevel").Int()
	threshold := opts.Get("threshold").Int()

	codecOpts := codec.CodecOptions{
		Wavelet:            wavelet.WaveletType(waveletFamily),
		DecompositionLevel: decompLevel,
		ThresholdingFactor: threshold,
	}
	fmt.Println(codecOpts)

	start := time.Now()
	compressed, result := codec.EncodeImageData(imageData, width, height, codecOpts)
	fmt.Printf("took %fs\n", time.Since(start).Seconds())

	safeCompressed := make([]interface{}, width*height*4)
	for i := range compressed {
		safeCompressed[i] = compressed[i]
	}

	return map[string]interface{}{
		"Compressed": safeCompressed,
		"Result": map[string]interface{}{
			"PSNR": result.PSNR,
		},
	}
}

func convertImageDataToImage(imageData js.Value, width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	var x, y int
	var r, g, b, a uint8
	for i := 0; i < imageData.Length(); i += 4 {
		x = (i / 4) % width
		y = (i / 4) / width
		r = uint8(imageData.Index(i + 0).Int())
		g = uint8(imageData.Index(i + 1).Int())
		b = uint8(imageData.Index(i + 2).Int())
		a = uint8(imageData.Index(i + 3).Int())
		img.Set(x, y, color.RGBA{r, g, b, a})
	}
	return img
}

func main() {
	fmt.Println("WASM Go inited")
	js.Global().Set("Sip_CompressImage", js.FuncOf(compressImage))
	js.Global().Set("Sip_GetImageChannels", js.FuncOf(getImageChannels))

	<-make(chan bool)
}
