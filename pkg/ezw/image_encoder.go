package ezw

import (
	"fmt"
	"io"
	"sync"

	"github.com/jnafolayan/sip/internal/imageutils"
	"github.com/jnafolayan/sip/pkg/codec"
	"github.com/jnafolayan/sip/pkg/signal"
	"github.com/jnafolayan/sip/pkg/wavelet"
)

type ImageEncoder struct {
	encoders     []*Encoder
	wavelet      wavelet.Wavelet
	codecOptions codec.CodecOptions
}

func NewImageEncoder(codecOpts codec.CodecOptions) *ImageEncoder {
	return &ImageEncoder{
		codecOptions: codecOpts,
	}
}

func (ie *ImageEncoder) Init(src string) error {
	w, err := codec.GetWaveletFamily(ie.codecOptions)
	if err != nil {
		return err
	}
	ie.wavelet = w

	channels := getImageChannels(src)
	channels = transformChannels(w, channels, ie.codecOptions)
	encoders := createEncoders(channels, ie.codecOptions)
	ie.encoders = encoders

	return nil
}

func (ie *ImageEncoder) Tick(w io.Writer) error {
	w.Write([]byte{StartOfImageMarker})

	for _, e := range ie.encoders {
		w.Write([]byte{StartOfChannelMarker})

		err := e.Next()
		if err != nil {
			continue
		}

		e.Flush(w)
	}

	w.Write([]byte{EndOfImageMarker})

	return nil
}

func getImageChannels(src string) []signal.Signal2D {
	img, err := imageutils.ReadImage(src)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	width, height := img.Bounds().Dx(), img.Bounds().Dy()

	imageData := imageutils.ConvertImageToImageData(img)
	imageChannels := imageutils.GetImageChannelsFromImageData(imageData, width, height)

	return imageChannels
}

func transformChannels(w wavelet.Wavelet, channels []signal.Signal2D, opts codec.CodecOptions) []signal.Signal2D {
	transformed := make([]signal.Signal2D, len(channels))

	wg := &sync.WaitGroup{}
	wg.Add(len(channels))

	for i, channel := range channels {
		go func() {
			transformed[i] = w.Transform(channel)
			if opts.ThresholdingStrategy == "hard" {
				transformed[i] = w.HardThreshold(transformed[i], opts.ThresholdingFactor)
			} else if opts.ThresholdingStrategy == "soft" {
				transformed[i] = w.SoftThreshold(transformed[i], opts.ThresholdingFactor)
			}
			wg.Done()
		}()
	}
	wg.Wait()

	return transformed
}

// createEncoders creates and initialized EZW encoders for each
// channel in an image.
func createEncoders(channels []signal.Signal2D, opts codec.CodecOptions) []*Encoder {
	encoders := make([]*Encoder, len(channels))
	for i := range channels {
		e := NewEncoder()
		e.Init(channels[i], opts)
		encoders[i] = e
	}

	return encoders
}
