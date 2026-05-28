package cipher

import (
	"crypto/sha256"
	"vsatanasov/custom-streaming-algorithm/pkg/lfsr"
)

const (
	//x^8 + x^4 + x^3 + x^2 + 1
	// 100011101
	poly1 = 0b100011101
	//x^8 + x^6 + x^4 + x^3 + x^2 + x + 1
	//101011111
	poly2 = 0b101011111
	// x^8 + x^6 + x^5 + x^4 + 1
	//101110001
	poly3 = 0b101110001
	//x^8 + x^7 + x^6 + x^5 + x^4 + x^2 + 1
	//111110101
	poly4    = 0b111110101
	keyLen   = 32
	sboxSize = 16 // 16×16 = 256 cells
)

type Cipher struct {
	lsfrs    [4]lfsr.LFSR
	sbox     *SBox
	outCount int
}

func (c *Cipher) GetRegisters() [4]lfsr.LFSR {
	return c.lsfrs
}

func (c *Cipher) Encode(message []byte) []byte {
	result := make([]byte, len(message))
	for i := range message {
		result[i] = c.encodeByte(message[i])
	}

	return result
}

func (c *Cipher) encodeByte(b byte) byte {
	keyByte := c.sbox.NextByte()
	c.outCount++
	if c.outCount == sboxSize*sboxSize {
		c.refreshSBox()
		c.outCount = 0
	}
	return b ^ keyByte
}

func (c *Cipher) refreshSBox() {
	for row := range sboxSize {
		for col := range sboxSize {
			c.sbox.XORAt(row, col, c.Тick())
		}
	}
}

func (c *Cipher) Тick() uint8 {
	l1 := c.lsfrs[0].NextBit()
	l2 := c.lsfrs[1].NextBit()
	l3 := c.lsfrs[2].NextBit()
	l4 := c.lsfrs[3].NextBit()
	b := (l1 & l2 & l3) ^ (l2 & l4) ^ l1 ^ l3
	return b ^ majorityBits(l1, l2, l3, l4)
}

func (c *Cipher) warmup() {
	for range 100 {
		c.Тick()
	}
}

func New(iv []byte) *Cipher {
	sha := sha256.Sum256(iv)

	vectors := sha[0:4]
	states := make([]uint, 0, 4)
	for _, v := range vectors {
		states = append(states, uint(int64FromBytes([]byte{v})))
	}
	sboxBytes := makeSBoxSeed(iv)
	cph := &Cipher{
		lsfrs: [4]lfsr.LFSR{
			lfsr.New(states[0], poly1),
			lfsr.New(states[1], poly2),
			lfsr.New(states[2], poly3),
			lfsr.New(states[3], poly4),
		},
		sbox: NewSbox(sboxBytes, sboxSize),
	}
	cph.warmup()

	return cph
}

func makeSBoxSeed(iv []byte) []byte {
	seed := make([]byte, 0, sboxSize*sboxSize)
	h := sha256.Sum256(iv)
	for len(seed) < sboxSize*sboxSize {
		seed = append(seed, h[:]...)
		h = sha256.Sum256(h[:])
	}
	return seed[:sboxSize*sboxSize]
}

func int64FromBytes(bytes []byte) int64 {
	key := int64(0)
	for i := range bytes {
		b := bytes[i]
		key |= int64(b)

		if i != len(bytes)-1 {
			key <<= 8
		}
	}

	return key
}

func majority(b byte) uint8 {
	c := 0
	for i := range 8 {
		c += int((b >> i) & 1)
	}

	if c >= 4 {
		return 1
	}

	return 0
}

func majorityBits(b uint8, b1 uint8, b2 uint8, b3 uint8) uint8 {
	return ((b & b1) ^ (b & b2) ^ (b & b3) ^ (b1 & b2) ^ (b1 & b3) ^ (b2 & b3))
}
