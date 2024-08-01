package main

import (
	"fmt"
	"syscall/js"
	"time"

	"github.com/jnafolayan/sip/pkg/codec"
	"github.com/jnafolayan/sip/pkg/wavelet"
)

func jsCompressImage(this js.Value, args []js.Value) interface{} {
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

	start := time.Now()
	compressed, result := codec.EncodeImageData(imageData, width, height, codecOpts)
	fmt.Printf("took %fs\n", time.Since(start).Seconds())

	safeCompressed := js.Global().Get("Uint8Array").New(len(compressed))
	js.CopyBytesToJS(safeCompressed, compressed)

	return map[string]interface{}{
		"Compressed": safeCompressed,
		"Result": map[string]interface{}{
			"PSNR":  result.PSNR,
			"Ratio": 0.0,
		},
	}
}

func main() {
	fmt.Println("WASM Go initialized")
	js.Global().Set("Sip_CompressImage", js.FuncOf(jsCompressImage))

	<-make(chan bool)
}
