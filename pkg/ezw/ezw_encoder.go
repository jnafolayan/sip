package ezw

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"math"

	"github.com/jnafolayan/sip/pkg/signal"
)

type FlatSignalCoeff struct {
	Row, Col int
	Value    signal.SignalCoeff
}

type SignificantCoeff struct {
	FlatSignalCoeff
	Symbol SymbolType
}

type Encoder struct {
	signal          signal.Signal2D
	dominantList    []FlatSignalCoeff
	subordinateList []SignificantCoeff
	threshold       int
	level           int
	output          *bytes.Buffer
}

func NewEncoder() *Encoder {
	return &Encoder{}
}

// Init prepares the encoder for subsequent EZW coding passes
func (e *Encoder) Init(coeffs signal.Signal2D, decompLevel int) error {
	width, height := coeffs.Size()
	e.signal = coeffs
	e.level = decompLevel
	e.dominantList = e.flattenSource()
	e.subordinateList = make([]SignificantCoeff, 0, width*height)
	e.output = new(bytes.Buffer)
	e.threshold = int(math.Pow(2, math.Floor(math.Log2(findMaxCoeff(coeffs)))))

	return nil
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

func (e *Encoder) write(coeff SignificantCoeff) {
	e.output.Write(e.encodeCoefficient(coeff))
}

func (e *Encoder) encodeCoefficient(coeff SignificantCoeff) []byte {
	buf := new(bytes.Buffer)

	// Encode symbol
	symbolBits := SymbolCodes[coeff.Symbol]
	buf.WriteByte(symbolBits)

	// Encode row and col indices
	rowBits := uint16(coeff.Row)
	colBits := uint16(coeff.Col)
	binary.Write(buf, binary.BigEndian, rowBits)
	binary.Write(buf, binary.BigEndian, colBits)

	// Encode the value
	valueBits := float64(coeff.Value)
	binary.Write(buf, binary.BigEndian, valueBits)

	return buf.Bytes()
}

// Flush writes the current output and clears the output buffer.
func (e *Encoder) Flush(w io.Writer) {
	w.Write(e.output.Bytes())
	e.output.Reset()
}

var ErrStopped = errors.New("stopped encoding")

func (e *Encoder) Next() error {
	if e.threshold <= 0 {
		return ErrStopped
	}

	e.SignificancePass()
	e.RefinementPass()
	e.threshold /= 2

	return nil
}

func (e *Encoder) SignificancePass() {
	markedForDeletion := make([]int, 0, 2)
	T := float64(e.threshold)

	for coeffIndex, coeff := range e.dominantList {
		sCoeff := SignificantCoeff{
			FlatSignalCoeff: coeff,
			Symbol:          SymbolNone,
		}
		if math.Abs(coeff.Value) >= T {
			if coeff.Value >= 0 {
				sCoeff.Symbol = SymbolPS
			} else {
				sCoeff.Symbol = SymbolNG
			}
			e.write(sCoeff)
			e.subordinateList = append(e.subordinateList, sCoeff)
			markedForDeletion = append(markedForDeletion, coeffIndex)
		} else {
			if e.checkIsZerotreeDescendant(coeff) {
				// Don't code - it is "predictably insignificant"
				continue
			} else if e.checkIsZerotree(coeff) {
				sCoeff.Symbol = SymbolZR
				e.write(sCoeff)
			} else {
				// Coeff is an isolated zerotree
				sCoeff.Symbol = SymbolIZ
			}
		}
	}

	// Delete coeffs that were added to the subordinate list
	for _, idx := range markedForDeletion {
		e.dominantList = append(e.dominantList[:idx], e.dominantList[idx+1:]...)
	}
}

func (e *Encoder) RefinementPass() {
	var abs signal.SignalCoeff
	T := float64(e.threshold)
	upperT := T * 2
	midT := T + (upperT-T)/2
	for _, coeff := range e.subordinateList {
		abs = math.Abs(coeff.Value)
		if abs >= T && abs < midT {
			coeff.Symbol = SymbolLow
			e.write(coeff)
		} else if abs >= midT && abs <= upperT {
			coeff.Symbol = SymbolHigh
			e.write(coeff)
		}
	}
}

func (e *Encoder) checkIsZerotree(coeff FlatSignalCoeff) bool {
	w, h := e.signal.Size()
	row, col := coeff.Row, coeff.Col

	if e.signal[row][col] >= float64(e.threshold) {
		return false
	}

	if row == 0 && col == 0 {
		// FIXME: should handle this edge case better
		return e.signal[0][0] == 0
	}

	for {
		row *= 2
		col *= 2
		if row >= h || col >= w {
			break
		}

		for y := row; y < row+2; y++ {
			for x := col; x < col+2; x++ {
				if math.Abs(e.signal[y][x]) >= float64(e.threshold) {
					return false
				}
			}
		}
	}

	return true
}

func (e *Encoder) checkIsZerotreeDescendant(coeff FlatSignalCoeff) bool {
	width, height := e.signal.Size()
	row, col := coeff.Row, coeff.Col

	llWidth := (width / (1 << e.level))
	llHeight := (height / (1 << e.level))

	// Return false if coeff is in LL (is a root)
	if row < llHeight && col < llWidth {
		return false
	}

	for {
		row /= 2
		col /= 2
		if math.Abs(e.signal[row][col]) < float64(e.threshold) {
			return true
		}

		if row < llHeight && col < llWidth {
			// Break when we reach the LL subband. We can't search further.
			break
		}
	}

	return false
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

// flattenQuadrant flattens the 2D slice into a 1D slice
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
