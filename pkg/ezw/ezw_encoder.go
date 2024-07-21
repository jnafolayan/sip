package ezw

import (
	"math"

	"github.com/jnafolayan/sip/pkg/signal"
)

type SignificantCoeff struct {
	FlatSignalCoeff
	Symbol SymbolType
}

type FlatSignalCoeff struct {
	Row, Col int
	Value    signal.SignalCoeff
}

type SymbolType int

const (
	SymbolNone SymbolType = iota
	SymbolPS              // Positive siginificant
	SymbolNG              // Negative significant
	SymbolZR              // Zerotree root
	SymbolIZ              // Isolated zero
)

type Encoder struct {
	signal          signal.Signal2D
	dominantList    []FlatSignalCoeff
	subordinateList []SignificantCoeff
	threshold       int
	level           int
}

func findMaxCoeff(s signal.Signal2D) signal.SignalCoeff {
	var maxCoeff signal.SignalCoeff = 0

	for row := range s {
		for col := range s[row] {
			if math.Abs(s[row][col]) > maxCoeff {
				maxCoeff = math.Abs(s[row][col])
			}
		}
	}

	return maxCoeff
}

// Init prepares the encoder for subsequent EZW coding passes
func (e *Encoder) Init(s signal.Signal2D, level int) {
	w, h := s.Size()
	e.signal = s
	e.level = level
	e.dominantList = e.flattenSource()
	e.subordinateList = make([]SignificantCoeff, w*h)
	e.threshold = int(math.Pow(2, math.Floor(math.Log2(findMaxCoeff(s)))))
}

func (e *Encoder) SignificancePass() {
	l := e.level
	for l >= 1 {
		l--
	}
}

// flattenSource generates a list containing coefficients by scanning the source
// 2D signal using Morton scan.
func (e *Encoder) flattenSource() []FlatSignalCoeff {
	w, h := e.signal.Size()
	l := e.level
	result := make([]FlatSignalCoeff, 0, w*h)

	for l >= 1 {
		ll, hl, lh, hh := e.getQuadrantsForLevel(l)
		if l == e.level {
			// Append the average coeffs once
			result = append(result, e.flattenQuadrant(ll)...)
		}
		result = append(result, e.flattenQuadrant(hl)...)
		result = append(result, e.flattenQuadrant(lh)...)
		result = append(result, e.flattenQuadrant(hh)...)
		l--
	}

	return result
}

func (e *Encoder) flattenQuadrant(q []int) []FlatSignalCoeff {
	endRow := q[1] + q[3]
	endCol := q[0] + q[2]
	result := make([]FlatSignalCoeff, 0, q[2]*q[3])
	for row := q[1]; row < endRow; row++ {
		for col := q[0]; col < endCol; col++ {
			result = append(result, FlatSignalCoeff{
				Row:   row,
				Col:   col,
				Value: e.signal[row][col],
			})
		}
	}
	return result
}

// getQuadrantsForLevel returns the subbands for a decomposition level using Morton
// scan.
// quadrant = [x, y, width, height]
func (e *Encoder) getQuadrantsForLevel(level int) ([]int, []int, []int, []int) {
	w, h := e.signal.Size()
	halfWidth := w / (1 << level)
	halfHeight := h / (1 << level)

	ll := []int{0, 0, halfWidth, halfHeight}
	hl := []int{halfWidth, 0, halfWidth, halfHeight}
	lh := []int{0, halfHeight, halfWidth, halfHeight}
	hh := []int{halfWidth, halfHeight, halfWidth, halfHeight}

	return ll, hl, lh, hh
}
