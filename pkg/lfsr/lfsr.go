package lfsr

import (
	"fmt"
)

type LFSR struct {
	state uint8
	taps  []int
}

func New(state uint) LFSR {
	taps := make([]int, 0)
	for idx := range 8 {
		t := (state >> uint8(idx)) & 1
		if t&1 == 1 {
			taps = append(taps, idx)
		}
	}

	return LFSR{
		state: uint8(state),
		taps:  taps,
	}
}

func (l *LFSR) NextBit() uint8 {
	nextBit := l.state & 1
	feedback := uint8(0)
	for i := 0; i < len(l.taps); i++ {
		tapPos := 7 - l.taps[i]
		feedback ^= uint8((l.state >> uint8(tapPos)) & 1)
	}

	l.state = l.state >> 1
	l.state |= uint8(feedback) << 7

	return nextBit
}

func (l *LFSR) ToString() string {
	return fmt.Sprintf("state: %b\n", l.state)
}
