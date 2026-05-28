package lfsr

import (
	"fmt"
)

type LFSR struct {
	state uint8
	taps  uint8
}

func New(state uint, polynomial uint) LFSR {
	taps := tapsUint8FormPolynomial(polynomial)
	return LFSR{
		state: uint8(state),
		taps:  taps,
	}
}

func (l *LFSR) NextBit() uint8 {
	nextBit := l.state & 1
	l.state = l.state >> 1
	if nextBit == 1 {
		l.state ^= l.taps
	}

	return nextBit
}

func (l *LFSR) GetState() uint8 {
	return l.state
}

func (l *LFSR) ToString() string {
	return fmt.Sprintf("state: %b\ntaps: %b", l.state, l.taps)
}

func tapsFormPolynomial(poly uint) []int {
	taps := make([]int, 0)
	for idx := range 8 {
		t := (uint8(poly) >> uint8(idx)) & 1
		if t&1 == 1 {
			taps = append(taps, idx)
		}
	}

	return taps
}

func tapsUint8FormPolynomial(poly uint) uint8 {
	var taps uint8 = 0
	for i := range 8 {
		shift := 7 - i

		taps = taps | uint8(poly>>shift)&1
		if i != 7 {
			taps <<= 1
		}
	}

	return taps
}
