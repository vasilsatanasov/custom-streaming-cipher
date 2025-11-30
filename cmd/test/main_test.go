package test

import (
	"encoding/binary"
	"fmt"
	"testing"
	"vsatanasov/custom-streaming-algorithm/pkg/cipher"
	"vsatanasov/custom-streaming-algorithm/pkg/lfsr"
)

func TestKeyFromBytes(t *testing.T) {
	key := []byte("test")
	result := cipher.Int64FromBytes(key)
	check := int64(binary.BigEndian.Uint32(key))
	if result != check {
		t.Errorf("%b and %b are not equal", result, check)
	}
}

func TestLfsr(t *testing.T) {
	//poly = x^7 + x^4 + x^3 + x^2 + 1
	theLfsr := lfsr.New(0b10101011, 0b00011101)
	fmt.Println(theLfsr.ToString())
	for range 100 {
		result := theLfsr.NextBit()
		fmt.Println(theLfsr.ToString())
		fmt.Printf("result: %b\n", result)
	}

}

func TestLfsrWorksAsExpected(t *testing.T) {
	//poly = x^7 + x + 1
	register := lfsr.New(0b10101011, 0b10000011)
	bit := register.NextBit()
	if bit != 1 {
		t.Error("Unexpected next bit value")
	}

	if register.GetState() != uint8(0b11010101) {
		t.Errorf("Unexpected state value %b", register.GetState())
	}
}

func TestCipherCreate(t *testing.T) {
	pass := []byte("abcd")
	iv := []byte("c!ph3r")
	c := cipher.New(pass, iv)
	if c == nil {
		t.Error("Could not crate Cipher")
	}

	if len(c.GetRegisters()) != 4 {
		t.Error("Cipher must have 4 registers")
	}

	if c.GetKey() == 0 {
		t.Error("Cipher must have nonzero key")
	}

}

func TestCipherEncodeDecode(t *testing.T) {
	pass := []byte("abcdefgh")
	iv := []byte("c!ph3r")
	c := cipher.New(pass, iv)
	c1 := cipher.New(pass, iv)

	msg := []byte("Az obicham mach i boza")
	encoded := c.Encode(msg)
	decoded := c1.Encode(encoded)

	if string(encoded) == string(msg) {
		t.Error("Cipher is not encoding")
	}

	if string(msg) != string(decoded) {
		t.Error("Cipher is not working")
	}
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
