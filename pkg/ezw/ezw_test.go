package ezw

import (
	"bytes"
	"fmt"
	"slices"
	"testing"

	"github.com/jnafolayan/sip/pkg/signal"
)

func TestGetQuadrants(t *testing.T) {
	s := signal.Signal2D{
		{1, 2, 3, 4},
		{5, 6, 7, 8},
	}

	e := NewEncoder()
	e.Init(s, 3)
	ll, hl, lh, hh := e.getQuadrantsForLevel(1)
	if !slices.Equal(ll, []int{0, 0, 2, 1}) {
		t.Errorf("Expected ll=%v, got=%v\n.", []int{0, 0, 2, 1}, ll)
	}
	if !slices.Equal(hl, []int{2, 0, 2, 1}) {
		t.Errorf("Expected hl=%v, got=%v\n.", []int{2, 0, 2, 1}, hl)
	}
	if !slices.Equal(lh, []int{0, 1, 2, 1}) {
		t.Errorf("Expected lh=%v, got=%v\n.", []int{0, 1, 2, 1}, lh)
	}
	if !slices.Equal(hh, []int{2, 1, 2, 1}) {
		t.Errorf("Expected hh=%v, got=%v\n.", []int{2, 1, 2, 1}, hh)
	}
}

func TestFlattenSource(t *testing.T) {
	s := signal.Signal2D{
		{1, 2, 5, 6},
		{3, 4, 7, 8},
	}

	e := NewEncoder()
	e.Init(s, 3)
	expected := []FlatSignalCoeff{
		{0, 0, 1},
		{0, 1, 2},
		{0, 2, 5},
		{0, 3, 6},
		{1, 0, 3},
		{1, 1, 4},
		{1, 2, 7},
		{1, 3, 8},
	}
	if !slices.Equal(e.dominantList, expected) {
		t.Errorf("Expected dominantList=%v, got=%v.\n", expected, e.dominantList)
	}
}

func TestIsZerotree(t *testing.T) {
	s := signal.Signal2D{
		{127, 69, 24, 73, 13, 5, -8, 5},
		{-37, -18, -18, 8, -6, 7, 15, 4},
		{44, -87, -15, 21, 8, -11, 14, -3},
		{65, 18, 29, -56, 0, -2, 3, 7},
		{34, 38, -18, 17, 3, -9, -2, 1},
		{-27, -41, 11, -5, 0, -1, 0, -3},
		{6, 17, 5, -19, 2, 0, -3, 1},
		{32, 26, -7, 5, -1, -5, 7, 4},
	}
	e := NewEncoder()
	e.Init(s, 3)

	tests := []struct {
		row      int
		col      int
		expected bool
	}{
		{0, 0, false},
		{0, 1, false},
		{1, 0, false},
		{1, 1, true},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d.", i), func(t *testing.T) {
			coeff := FlatSignalCoeff{
				Row:   tt.row,
				Col:   tt.col,
				Value: s[tt.row][tt.col],
			}
			actual := e.checkIsZerotree(coeff)
			if actual != tt.expected {
				t.Errorf("Expected checkIsZerotree(row:%d, col:%d)=%v, got=%v\n", tt.row, tt.col, tt.expected, actual)
			}
		})
	}
}

func TestIsZerotreeDescendant(t *testing.T) {
	s := signal.Signal2D{
		{127, 69, 24, 73, 13, 5, -8, 5},
		{-37, -18, -18, 8, -6, 7, 15, 4},
		{44, -87, -15, 21, 8, -11, 14, -3},
		{65, 18, 29, -56, 0, -2, 3, 7},
		{34, 38, -18, 17, 3, -9, -2, 1},
		{-27, -41, 11, -5, 0, -1, 0, -3},
		{6, 17, 5, -19, 2, 0, -3, 1},
		{32, 26, -7, 5, -1, -5, 7, 4},
	}
	e := NewEncoder()
	e.Init(s, 3)

	tests := []struct {
		row      int
		col      int
		expected bool
	}{
		{0, 0, false},
		{0, 1, false},
		{1, 0, false},
		{1, 1, false},
		{2, 2, true},
		{2, 3, true},
		{4, 4, true},
		{6, 5, true},
		{7, 7, true},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d.", i), func(t *testing.T) {
			coeff := FlatSignalCoeff{
				Row:   tt.row,
				Col:   tt.col,
				Value: s[tt.row][tt.col],
			}
			actual := e.checkIsZerotreeDescendant(coeff)
			if actual != tt.expected {
				t.Errorf("Expected checkIsZerotreeDescendant(row:%d, col:%d)=%v, got=%v\n", tt.row, tt.col, tt.expected, actual)
			}
		})
	}
}

func TestEncoder(t *testing.T) {
	s := signal.Signal2D{
		{127, 69, 24, 73, 13, 5, -8, 5},
		{-37, -18, -18, 8, -6, 7, 15, 4},
		{44, -87, -15, 21, 8, -11, 14, -3},
		{65, 18, 29, -56, 0, -2, 3, 7},
		{34, 38, -18, 17, 3, -9, -2, 1},
		{-27, -41, 11, -5, 0, -1, 0, -3},
		{6, 17, 5, -19, 2, 0, -3, 1},
		{32, 26, -7, 5, -1, -5, 7, 4},
	}
	e := NewEncoder()
	e.Init(s, 3)
	e.Next()

	var out bytes.Buffer
	e.Flush(&out)
	t.Errorf("%v\n", out.String())
	e.Next()
	e.Flush(&out)
	t.Errorf("%v\n", out.String())
	e.Next()
	e.Flush(&out)
	t.Errorf("%v\n", out.String())
}
