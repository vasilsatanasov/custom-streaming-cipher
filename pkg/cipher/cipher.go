package cipher

import (
	"crypto/sha256"
	"vsatanasov/custom-streaming-algorithm/pkg/lfsr"
)

const (
	poly1 = 0b10000011
	poly2 = 0b10001001
	poly3 = 0b10001111
	poly4 = 0b10010001
)

type Cipher struct {
	lsfrs      [4]lfsr.LFSR
	key        int64
	nextKeyPos int
}

func (c *Cipher) GetRegisters() [4]lfsr.LFSR {
	return c.lsfrs
}

func (c *Cipher) GetKey() int64 {
	return c.key
}

func (c *Cipher) Encode(message []byte) []byte {
	result := make([]byte, 0)
	for i := range message {
		r := c.encodeByte(message[i])
		result = append(result, r)
	}

	return result
}

func (c *Cipher) encodeByte(b byte) byte {
	b1 := byte(0)

	for i := 7; i >= 0; i-- {
		bit := c.Тick()
		b1 |= (((b >> i) & 1) ^ bit)
		if i > 0 {
			b1 = b1 << 1
		}
	}

	return b1
}

func (c *Cipher) Тick() uint8 {
	l1 := c.lsfrs[0].NextBit()
	l2 := c.lsfrs[1].NextBit()
	l3 := c.lsfrs[2].NextBit()
	l4 := c.lsfrs[3].NextBit()

	b := ((l1 & l2) ^ (l1 & l3) ^ (l1 & l4) ^ (l2 & l3) ^ (l2 & l4) ^ (l3 & l4)) & 1
	k := uint8((c.key >> int64(c.nextKeyPos)) & 1)
	c.nextKeyPos = (c.nextKeyPos + 1) % 64
	return b ^ k
}

func (c *Cipher) warmup() {
	for range 100 {
		c.Тick()
	}
}

func New(key []byte, iv []byte) *Cipher {
	if len(key) < 4 || len(key) > 8 {
		panic("Key must be 4 - 8 bytes")
	}
	k := Int64FromBytes(key)
	sha := sha256.Sum256(iv)

	vectors := sha[0:4]
	states := make([]uint, 0, 4)
	for _, v := range vectors {
		states = append(states, uint(Int64FromBytes([]byte{v})))
	}
	cph := &Cipher{
		lsfrs: [4]lfsr.LFSR{
			lfsr.New(states[0], poly1),
			lfsr.New(states[1], poly2),
			lfsr.New(states[2], poly3),
			lfsr.New(states[3], poly4),
		},
		key:        k,
		nextKeyPos: 0,
	}
	cph.warmup()

	return cph
}

func Int64FromBytes(bytes []byte) int64 {
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
