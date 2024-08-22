package ezw

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
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

func (id *ImageDecoder) DecodeJSONFrame(r io.Reader) error {
	frames := make([]JSONFrame, 3)
	err := json.NewDecoder(r).Decode(&frames)
	if err != nil {
		return fmt.Errorf("error decoding json frame: %w", err)
	}

	for idx, jsonFrame := range frames {
		fmt.Printf("Channel %d: %d coeffs, %d T\n", idx, len(jsonFrame.Coefficients), jsonFrame.Threshold)
		for _, coeff := range jsonFrame.Coefficients {
			if id.channels[idx] == nil {
				id.channels[idx] = signal.New(jsonFrame.FrameWidth, jsonFrame.FrameHeight)
			}
			channel := id.channels[idx]

			symbol, row, col, value := coeff.Symbol, coeff.Row, coeff.Col, coeff.Value
			T := float64(jsonFrame.Threshold)
			upperT := T * 2
			midT := T + (upperT-T)/2
			_ = midT

			v := channel[row][col]
			switch SymbolType(symbol) {
			case SymbolZR:
				continue
			case SymbolIZ:
				fmt.Println("Isolate zerotree")
				continue
			case SymbolPS:
				// fmt.Println("SymbolPS")
				v = T
			case SymbolNG:
				// fmt.Println("SymbolNG")
				v = -T
			case SymbolLow:
				continue
			case SymbolHigh:
				if v < 0 {
					v -= T / 2
				} else {
					v += T / 2
				}
			default:
				fmt.Printf("unknown symbol: %q\n", symbol)
			}

			v = float64(value)
			channel[row][col] = v
		}
	}

	return nil
}

func (id *ImageDecoder) DecodeBinaryFrame(r io.Reader) error {
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
