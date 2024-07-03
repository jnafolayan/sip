package wavelet

import "github.com/jnafolayan/sip/signal"

type Wavelet interface {
	Transform(signal signal.Signal2D) signal.Signal2D
	InverseTransform(signal signal.Signal2D) signal.Signal2D
}
