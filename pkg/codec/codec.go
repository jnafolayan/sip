package codec

import "github.com/jnafolayan/sip/pkg/wavelet"

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
