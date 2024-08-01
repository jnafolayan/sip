package main

import (
	"fmt"
	"image"
	"image/color"
	"syscall/js"

	"github.com/jnafolayan/sip/pkg/codec"
	"github.com/jnafolayan/sip/pkg/wavelet"
)

func compressImage(this js.Value, args []js.Value) interface{} {
	// imageData, width, height, opts
	imageData := args[0]
	width, height := args[1].Int(), args[2].Int()
	opts := args[3]

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	var x, y int
	var r, g, b, a uint8
	for i := 0; i < width*height; i += 4 {
		x = i % width
		y = i / width
		r = uint8(imageData.Index(i + 0).Int())
		g = uint8(imageData.Index(i + 1).Int())
		b = uint8(imageData.Index(i + 2).Int())
		a = uint8(imageData.Index(i + 3).Int())
		img.Set(x, y, color.RGBA{r, g, b, a})
	}

	waveletFamily := opts.Get("waveletFamily").String()
	decompLevel := opts.Get("decompLevel").Int()
	threshold := opts.Get("threshold").Int()

	codecOpts := codec.CodecOptions{
		Wavelet:            wavelet.WaveletType(waveletFamily),
		DecompositionLevel: decompLevel,
		ThresholdingFactor: threshold,
	}

	compressed, result := codec.Encode(img, codecOpts)
	return map[string]interface{}{
		"Compressed": convertImageToDataArray(compressed),
		"Result": map[string]interface{}{
			"PSNR": result.PSNR,
		},
	}
}

func convertImageToDataArray(img image.Image) []interface{} {
	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()
	result := make([]interface{}, w*h*4)

	var r, g, b, a uint32
	var rr, gg, bb, aa uint8
	offset := 0

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, a = img.At(x, y).RGBA()
			rr = uint8(r >> 8)
			gg = uint8(g >> 8)
			bb = uint8(b >> 8)
			aa = uint8(a >> 8)
			// offset = (x + y*w) * 4
			result[offset+0] = rr
			result[offset+1] = gg
			result[offset+2] = bb
			result[offset+3] = aa
			offset += 4
		}
	}

	return result
}

func main() {
	fmt.Println("WASM Go inited")
	js.Global().Set("Sip_CompressImage", js.FuncOf(compressImage))

	<-make(chan bool)
}
