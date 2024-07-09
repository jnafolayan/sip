package wavelet

import "github.com/jnafolayan/sip/pkg/signal"

type Wavelet interface {
	GetDecompositionLevel() int
	Transform(signal signal.Signal2D) signal.Signal2D
	InverseTransform(signal signal.Signal2D) signal.Signal2D
}
