package ezw

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"io"

	"github.com/jnafolayan/sip/pkg/codec"
	"github.com/jnafolayan/sip/pkg/signal"
	"github.com/jnafolayan/sip/pkg/wavelet"
)

type ImageDecoder struct {
	wavelet      wavelet.Wavelet
	destSize     image.Rectangle
	codecOptions codec.CodecOptions
}

func NewImageDecoder(destSize image.Rectangle, codecOpts codec.CodecOptions) *ImageDecoder {
	return &ImageDecoder{
		destSize:     destSize,
		codecOptions: codecOpts,
	}
}

func (id *ImageDecoder) Init(src string) error {
	w, err := codec.GetWaveletFamily(id.codecOptions)
	if err != nil {
		return err
	}
	id.wavelet = w
	return nil
}

func (id *ImageDecoder) DecodeFrame(r io.Reader) ([]signal.Signal2D, error) {
	var marker byte
	var err error

	channels := make([]signal.Signal2D, 3)

	buf := bufio.NewReader(r)

	err = binary.Read(buf, binary.BigEndian, &marker)
	if err != nil {
		return nil, err
	}
	if marker != StartOfImageMarker {
		return nil, errors.New("expected start of image marker")
	}

	var width, height uint16
	err = binary.Read(buf, binary.BigEndian, &width)
	if err != nil {
		return nil, err
	}
	err = binary.Read(buf, binary.BigEndian, &height)
	if err != nil {
		return nil, err
	}

	err = binary.Read(buf, binary.BigEndian, &marker)
	if err != nil {
		return nil, err
	}
	if marker != StartOfChannelMarker {
		return nil, errors.New("expected start of channel marker")
	}
	channelIndex := 0
	for channelIndex <= 2 && marker != EndOfImageMarker {
		var threshold uint8
		err = binary.Read(buf, binary.BigEndian, &threshold)
		if err != nil {
			return nil, err
		}

		upperT := threshold * 2
		midT := threshold + (upperT-threshold)/2

		channel := signal.New(int(width), int(height))
		channels[channelIndex] = channel

		for {
			var symbol uint8
			err = binary.Read(buf, binary.BigEndian, &symbol)
			if err != nil {
				return nil, err
			}

			// Row and col
			var row, col uint16
			err = binary.Read(buf, binary.BigEndian, &row)
			if err != nil {
				return nil, err
			}
			err = binary.Read(buf, binary.BigEndian, &col)
			if err != nil {
				return nil, err
			}

			err = binary.Read(buf, binary.BigEndian, &marker)
			if err != nil {
				return nil, err
			}
			if marker == StartOfChannelMarker || marker == EndOfImageMarker {
				break
			}

			var coeff float64
			switch symbol {
			case SymbolCodes[SymbolZR]:
				coeff = 0
			case SymbolCodes[SymbolIZ]:
				fmt.Println("Isolate zerotree")
				coeff = 0
			case SymbolCodes[SymbolPS]:
				fmt.Println("SymbolPS")
				coeff = +float64(threshold)
			case SymbolCodes[SymbolNG]:
				fmt.Println("SymbolNG")
				coeff = -float64(threshold)
			case SymbolCodes[SymbolLow]:
				coeff = float64(threshold+midT) / 2
			case SymbolCodes[SymbolHigh]:
				coeff = float64(midT+upperT) / 2
			default:
				fmt.Printf("unknown symbol %q", symbol)
			}
			channel[row][col] = coeff

			err = buf.UnreadByte()
			if err != nil {
				return nil, err
			}
		}

		channelIndex++
	}

	if marker != EndOfImageMarker {
		return nil, errors.New("expected end of image marker")
	}

	return channels, nil
}
