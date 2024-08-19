package ezw

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"io"
	"sync"

	"github.com/jnafolayan/sip/pkg/codec"
	"github.com/jnafolayan/sip/pkg/signal"
	"github.com/jnafolayan/sip/pkg/wavelet"
)

type ImageDecoder struct {
	wavelet      wavelet.Wavelet
	destSize     image.Rectangle
	channels     []signal.Signal2D
	codecOptions codec.CodecOptions
}

func NewImageDecoder(destSize image.Rectangle, codecOpts codec.CodecOptions) *ImageDecoder {
	channels := make([]signal.Signal2D, 3)
	return &ImageDecoder{
		destSize:     destSize,
		channels:     channels,
		codecOptions: codecOpts,
	}
}

func (id *ImageDecoder) SetDestSize(size image.Rectangle) {
	id.destSize = size
}

func (id *ImageDecoder) Init() error {
	w, err := codec.GetWaveletFamily(id.codecOptions)
	if err != nil {
		return err
	}
	id.wavelet = w
	return nil
}

func (id *ImageDecoder) ReconstructChannels() []signal.Signal2D {
	w := id.wavelet
	channels := id.channels
	width, height := id.destSize.Dx(), id.destSize.Dy()
	fmt.Println(width, height)

	wg := sync.WaitGroup{}
	wg.Add(len(channels))

	recon := make([]signal.Signal2D, len(channels))
	for i, channel := range channels {
		go func() {
			r := w.InverseTransform(channel)
			r = r.Slice(0, 0, width, height)
			recon[i] = r
			wg.Done()
		}()
	}

	wg.Wait()

	return recon
}

func (id *ImageDecoder) DecodeFrame(r io.Reader) error {
	var marker byte
	var err error

	channels := id.channels

	buf := bufio.NewReader(r)

	err = binary.Read(buf, binary.BigEndian, &marker)
	if err != nil {
		return err
	}
	if marker != StartOfImageMarker {
		return errors.New("expected start of image marker")
	}

	var width, height uint16
	err = binary.Read(buf, binary.BigEndian, &width)
	if err != nil {
		return err
	}
	err = binary.Read(buf, binary.BigEndian, &height)
	if err != nil {
		return err
	}

	err = binary.Read(buf, binary.BigEndian, &marker)
	if err != nil {
		return err
	}
	if marker != StartOfChannelMarker {
		return errors.New("expected start of channel marker")
	}
	channelIndex := 0
	for channelIndex <= 2 && marker != EndOfImageMarker {
		var threshold uint8
		err = binary.Read(buf, binary.BigEndian, &threshold)
		if err != nil {
			return err
		}

		T := float64(threshold)
		fmt.Println(T)
		upperT := T * 2
		midT := T + (upperT-T)/2

		if channels[channelIndex] == nil {
			channels[channelIndex] = signal.New(int(width), int(height))
		}

		channel := channels[channelIndex]

		for {
			var symbol uint8
			err = binary.Read(buf, binary.BigEndian, &symbol)
			if err != nil {
				return err
			}

			// Row and col
			var row, col uint16
			err = binary.Read(buf, binary.BigEndian, &row)
			if err != nil {
				return err
			}
			err = binary.Read(buf, binary.BigEndian, &col)
			if err != nil {
				return err
			}

			err = binary.Read(buf, binary.BigEndian, &marker)
			if err != nil {
				return err
			}
			if marker == StartOfChannelMarker || marker == EndOfImageMarker {
				break
			}

			coeff := channel[row][col]
			switch symbol {
			case SymbolCodes[SymbolZR]:
				continue
			case SymbolCodes[SymbolIZ]:
				fmt.Println("Isolate zerotree")
				continue
			case SymbolCodes[SymbolPS]:
				// fmt.Println("SymbolPS")
				coeff += T
			case SymbolCodes[SymbolNG]:
				// fmt.Println("SymbolNG")
				coeff -= T
			case SymbolCodes[SymbolLow]:
				continue
			case SymbolCodes[SymbolHigh]:
				if coeff > 0 {
					coeff += T / 2
				} else {
					coeff -= T / 2
				}
			default:
				fmt.Printf("unknown symbol %q", symbol)
			}
			channel[row][col] = coeff
			_ = midT

			err = buf.UnreadByte()
			if err != nil {
				return err
			}
		}

		channelIndex++
	}

	if marker != EndOfImageMarker {
		return errors.New("expected end of image marker")
	}

	return nil
}
