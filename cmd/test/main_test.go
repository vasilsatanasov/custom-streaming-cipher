package test

import (
	"fmt"
	"testing"
	"vsatanasov/custom-streaming-algorithm/pkg/cipher"
	"vsatanasov/custom-streaming-algorithm/pkg/lfsr"
)

func TestBits(t *testing.T) {
	x := 0b101
	fmt.Printf("%b\n", x)
	fmt.Printf("%b\n", (x>>1)&1)
	fmt.Printf("%b\n", x&1)

}

func TestLfsr(t *testing.T) {
	//poly = x^8 + x^7 + x^6 + x + 1
	theLfsr := lfsr.New(0b10101011, 0b111000011)
	fmt.Println(theLfsr.ToString())
	for range 100 {
		result := theLfsr.NextBit()
		fmt.Println(theLfsr.ToString())
		fmt.Printf("result: %b\n", result)
	}

}

func TestLfsrWorksAsExpected(t *testing.T) {
	//poly = x^8 + x^7 + x^6 + x + 1
	register := lfsr.New(0b10101011, 0b111000011)
	bit := register.NextBit()
	if bit != 1 {
		t.Error("Unexpected next bit value")
	}

	if register.GetState() != uint8(0b10010110) {
		t.Errorf("Unexpected state value %b", register.GetState())
	}
}

func TestCipherCreate(t *testing.T) {
	iv := []byte("c!ph3r")
	c := cipher.New(iv)
	if c == nil {
		t.Error("Could not create Cipher")
	}

	if len(c.GetRegisters()) != 4 {
		t.Error("Cipher must have 4 registers")
	}
}

func TestCipherEncodeDecode(t *testing.T) {
	iv := []byte("c!ph3r")
	c := cipher.New(iv)
	c1 := cipher.New(iv)

	msg := []byte("Az obicham mach i boza")
	encoded := c.Encode(msg)
	decoded := c1.Encode(encoded)

	if string(encoded) == string(msg) {
		t.Error("Cipher is not encoding")
	}

	if string(msg) != string(decoded) {
		t.Error(fmt.Sprintf("Cipher is not working, Expected %s, got %s", string(msg), string(decoded)))
	}
}

func BenchmarkSBoxNextByte(b *testing.B) {
	sb := cipher.NewSbox(make([]byte, 256), 16)
	b.ResetTimer()
	for range b.N {
		sb.NextByte()
	}
}

func BenchmarkTick(b *testing.B) {
	c := cipher.New([]byte("bench"))
	b.ResetTimer()
	for range b.N {
		c.Тick()
	}
}

func BenchmarkEncodeByte(b *testing.B) {
	c := cipher.New([]byte("bench"))
	b.ResetTimer()
	for range b.N {
		c.Encode([]byte{0xAB})
	}
}

func BenchmarkEncode1KB(b *testing.B) {
	c := cipher.New([]byte("bench"))
	data := make([]byte, 1024)
	b.ResetTimer()
	for range b.N {
		c.Encode(data)
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
