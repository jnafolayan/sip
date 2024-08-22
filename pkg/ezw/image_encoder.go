package ezw

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"io"
	"strings"
	"sync"

	"github.com/jnafolayan/sip/internal/imageutils"
	"github.com/jnafolayan/sip/pkg/codec"
	"github.com/jnafolayan/sip/pkg/signal"
	"github.com/jnafolayan/sip/pkg/wavelet"
)

type ImageEncoder struct {
	channels []signal.Signal2D
	encoders []*Encoder
	wavelet  wavelet.Wavelet

	srcWidth     int
	srcHeight    int
	frameWidth   int
	frameHeight  int
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
	ie.channels = channels

	transformed := transformChannels(w, channels, ie.codecOptions)

	encoders := createEncoders(transformed, ie.codecOptions)
	ie.encoders = encoders

	ie.srcWidth, ie.srcHeight = channels[0].Size()
	ie.frameWidth, ie.frameHeight = transformed[0].Size()

	return nil
}

func (ie *ImageEncoder) SrcSize() image.Rectangle {
	return image.Rect(0, 0, ie.srcWidth, ie.srcHeight)
}

func (ie *ImageEncoder) Tick(w io.Writer) error {
	binary.Write(w, binary.BigEndian, StartOfImageMarker)
	binary.Write(w, binary.BigEndian, uint16(ie.frameWidth))
	binary.Write(w, binary.BigEndian, uint16(ie.frameHeight))

	for _, e := range ie.encoders {
		e.SetEncodeMode(EncodeBinary)
		binary.Write(w, binary.BigEndian, StartOfChannelMarker)
		binary.Write(w, binary.BigEndian, uint8(e.threshold))

		err := e.Next()
		if err != nil {
			continue
		}

		e.Flush(w)
	}

	binary.Write(w, binary.BigEndian, EndOfImageMarker)

	return nil
}

func (ie *ImageEncoder) TickJSON(w io.Writer) error {
	frames := make([]string, len(ie.encoders))
	buf := new(bytes.Buffer)
	for i, e := range ie.encoders {
		frames[i] = "null"
		e.SetEncodeMode(EncodeJSON)

		err := e.Next()
		if err != nil {
			continue
		}

		e.Flush(buf)

		frames[i] = buf.String()
		buf.Reset()
	}

	w.Write([]byte(fmt.Sprintf("[%s]", strings.Join(frames, ", "))))
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
