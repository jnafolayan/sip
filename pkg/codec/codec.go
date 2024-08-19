package codec

import (
	"fmt"

	"github.com/jnafolayan/sip/pkg/cdf97"
	"github.com/jnafolayan/sip/pkg/haar"
	"github.com/jnafolayan/sip/pkg/wavelet"
)

type CodecOptions struct {
	Wavelet              wavelet.WaveletType
	ThresholdingFactor   int
	ThresholdingStrategy string
	DecompositionLevel   int
}

var DefaultCodecOpts = CodecOptions{
	Wavelet:            wavelet.WaveletHaar,
	ThresholdingFactor: 50,
	DecompositionLevel: 1,
}

func GetWaveletFamily(opts CodecOptions) (wavelet.Wavelet, error) {
	// Get wavelet family
	var w wavelet.Wavelet
	switch opts.Wavelet {
	case wavelet.WaveletHaar:
		w = &haar.HaarWavelet{Level: opts.DecompositionLevel}
	case wavelet.WaveletCDF97:
		w = &cdf97.CDF97Wavelet{Level: opts.DecompositionLevel}
	default:
		return nil, fmt.Errorf("unrecognized wavelet: %s", opts.Wavelet)
	}
	return w, nil
}
