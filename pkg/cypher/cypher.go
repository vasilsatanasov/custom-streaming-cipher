package cypher

import (
	"strconv"
	"vsatanasov/custom-streamon-algorithm/pkg/lfsr"
)

const (
	poly1 = 0b10000011
	poly2 = 0b10001001
	poly3 = 0b10001111
	poly4 = 0b10010001
)

type Cypher struct {
	lsfrs       [4]lfsr.LFSR
	key         int64
	nextKeyPos  int
	keySequence string
}

func (c *Cypher) GetRegisters() [4]lfsr.LFSR {
	return c.lsfrs
}

func (c *Cypher) GetKey() int64 {
	return c.key
}

func (c *Cypher) GetKeySequence() string {
	return c.keySequence
}

func (c *Cypher) Encode(message []byte) []byte {
	result := make([]byte, 0)
	for i := range message {
		r := c.encodeByte(message[i])
		result = append(result, r)
	}

	return result
}

func (c *Cypher) encodeByte(b byte) byte {
	b1 := byte(0)

	for i := 7; i >= 0; i-- {
		bit := c.tick()
		c.keySequence += strconv.Itoa(int(bit & 1))
		b1 |= (((b >> i) & 1) ^ bit)
		if i > 0 {
			b1 = b1 << 1
		}
	}

	return b1
}

func (c *Cypher) tick() uint8 {
	b := c.lsfrs[0].NextBit() ^ c.lsfrs[1].NextBit() ^ c.lsfrs[2].NextBit() ^ c.lsfrs[3].NextBit()
	k := uint8((c.key >> int64(c.nextKeyPos)) & 1)
	c.nextKeyPos = (c.nextKeyPos + 1) % 64
	return b ^ k
}

func New(key []byte) *Cypher {
	k := keyFromBytes(key)
	return &Cypher{
		lsfrs: [4]lfsr.LFSR{
			lfsr.New(poly1),
			lfsr.New(poly2),
			lfsr.New(poly3),
			lfsr.New(poly4),
		},
		key:         k,
		nextKeyPos:  0,
		keySequence: "",
	}
}

func keyFromBytes(bytes []byte) int64 {
	if len(bytes) < 4 {
		panic("Key must be at least 4 bytes")
	}

	key := int64(0)
	for i := 3; i >= 0; i-- {
		b := bytes[i]
		for j := 7; j >= 0; j-- {
			key |= int64((b >> j) & 1)
			key = key << 1
		}
	}

	return key
}
