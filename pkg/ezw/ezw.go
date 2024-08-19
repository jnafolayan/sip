package ezw

type SymbolType int

const (
	SymbolNone SymbolType = iota
	SymbolPS              // Positive siginificant
	SymbolNG              // Negative significant
	SymbolZR              // Zerotree root
	SymbolIZ              // Isolated zero

	// Used in the refinement pass
	SymbolLow
	SymbolHigh
)

var SymbolCodes = map[SymbolType]uint8{
	SymbolZR:   0b0,
	SymbolPS:   0b10,
	SymbolLow:  0b110,
	SymbolNG:   0b1110,
	SymbolHigh: 0b11110,
	SymbolIZ:   0b11111,
}

const (
	StartOfImageMarker   = 0b10000000
	StartOfChannelMarker = 0b01000000
	EndOfImageMarker     = 0b11000000
)
