package ezw

import (
	"slices"
	"testing"

	"github.com/jnafolayan/sip/pkg/signal"
)

func TestGetQuadrants(t *testing.T) {
	s := signal.Signal2D{
		{1, 2, 3, 4},
		{5, 6, 7, 8},
	}

	e := &Encoder{}
	e.Init(s, 1)
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

	e := &Encoder{}
	e.Init(s, 1)
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
