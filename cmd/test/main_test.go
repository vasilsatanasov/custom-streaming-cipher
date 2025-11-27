package test

import (
	"fmt"
	"testing"
	"vsatanasov/custom-streamon-algorithm/pkg/cypher"
	"vsatanasov/custom-streamon-algorithm/pkg/lfsr"
)

func TestEncoded(t *testing.T) {

	x := 0b10110
	y := 0 | (x >> 1)
	fmt.Printf("%b\n", y)
	// x^8 + x^4 + x^3 + x^2 + 1
	theLfsr := lfsr.New(0b00011101)
	fmt.Println(theLfsr.ToString())
	for range 100 {

		result := theLfsr.NextBit()
		fmt.Println(theLfsr.ToString())
		fmt.Printf("result: %b\n", result)
	}

}

func TestCypherCreate(t *testing.T) {
	pass := []byte("abcd")
	c := cypher.New(pass)
	if c == nil {
		t.Error("Could not crate cypher")
	}

	if len(c.GetRegisters()) != 4 {
		t.Error("Cypher must have 4 registers")
	}

	if c.GetKey() == 0 {
		t.Error("Cypher must have nonzero key")
	}

}

func TestCypherEncodeDecode(t *testing.T) {
	pass := []byte("abcdefgh")
	c := cypher.New(pass)
	c1 := cypher.New(pass)

	msg := []byte("Az obicham mach i boza")
	fmt.Printf("%v\n", len(msg)*8)
	encoded := c.Encode(msg)
	decoded := c1.Encode(encoded)

	if string(encoded) == string(msg) {
		t.Error("Cypher is not encoding")
	}

	if string(msg) != string(decoded) {
		t.Error("Cypher is not working")
	}

	fmt.Println(c.GetKeySequence())
	fmt.Printf("%v\n", len(c.GetKeySequence()))

}

func TestSimpleEncoding(t *testing.T) {
	message := "test123"
	bytes := []byte(message)
	result := encode(bytes)
	result = encode(result)
	fmt.Println(message)
	fmt.Println(string(result))
}

func encode(bytes []byte) []byte {
	result := make([]byte, 0)
	for _, v := range bytes {
		b := byte(0)
		for i := 7; i >= 0; i-- {
			b |= ((v >> i) & 1) ^ 1
			if i != 0 {
				b = b << 1
			}

		}
		result = append(result, b)
	}

	return result
}
