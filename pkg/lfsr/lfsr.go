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

	nextBit := uint8((l.state >> uint8(l.taps[0])) & 1)
	for i := l.taps[1]; i < len(l.taps); i++ {
		nextBit ^= uint8((l.state >> uint8(l.taps[i])) & 1)
	}

	l.state = l.state >> 1
	l.state |= uint8(nextBit) << 7

	return uint8(nextBit)
}

func (l *LFSR) ToString() string {
	return fmt.Sprintf("state: %b\n", l.state)
}
