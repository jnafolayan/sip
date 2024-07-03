package signal

import "testing"

func TestSignalClone(t *testing.T) {
	s := createDummySignal()
	sClone := s.Clone()
	if s[0][0] != sClone[0][0] || s[0][1] != sClone[0][1] ||
		s[1][0] != sClone[1][0] || s[1][1] != sClone[1][1] {
		t.Fatal("Cloned signal did not equal original")
	}
}

func TestSignalString(t *testing.T) {
	s := createDummySignal()
	str := s.String(s.Bounds())
	expected := `   0.0 |   1.0 
   2.0 |   3.0 `
	if str != expected {
		t.Fatalf("String was incorrect, got:\n%s\nwant:\n%s.", str, expected)
	}
}

func TestSignalBounds(t *testing.T) {
	s := createDummySignal()
	b := s.Bounds()
	size := 2
	if b.Min.X != 0 || b.Max.X != size || b.Min.Y != 0 || b.Max.Y != size {
		t.Errorf("Bounds was incorrect, got: %dx%d, want: %dx%d.", b.Dx(), b.Dy(), size, size)
	}
}

func createDummySignal() Signal2D {
	return Signal2D{
		{0, 1},
		{2, 3},
	}
}
