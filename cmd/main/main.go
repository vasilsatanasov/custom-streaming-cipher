package main

import (
	"fmt"
	"os"
	"strconv"
	"vsatanasov/custom-streaming-algorithm/pkg/cipher"
)

func main() {
	args := os.Args[1:]
	if len(args) < 3 {
		panic("3 arguments must be provided: key for the cipher and lenght in bits of generate PSR")
	}
	key := args[0]
	size, err := strconv.Atoi(args[1])
	if err != nil {
		panic(fmt.Sprintf("Invalid value for PSR length: %v", args[1]))
	}

	filePath := args[2]
	file, err := os.Create(filePath)
	if err != nil {
		panic(fmt.Sprintf("Could not create file %s", filePath))
	}

	defer file.Close()

	c := cipher.New([]byte(key))

	curentSize := 0
	currentByte := byte(0)
	currentByteIndex := 0
	totalBytesWritten := 0
	for curentSize < size {
		b := c.Тick()
		currentByte = (currentByte << 1) | (b & 1)
		currentByteIndex++
		if currentByteIndex == 8 {
			//fmt.Printf("%b\n", currentByte)
			n, err := file.Write([]byte{currentByte})
			if err != nil {
				panic("Could not write to file")
			}
			currentByte = byte(0)
			currentByteIndex = 0
			totalBytesWritten += n
		} else if curentSize == size-1 {
			currentByte <<= (7 - currentByteIndex)
			//fmt.Printf("%b\n", currentByte)
			n1, err := file.Write([]byte{currentByte})
			if err != nil {
				panic("Could not write to file")
			}
			totalBytesWritten += n1
		}
		curentSize++
	}

	fmt.Printf("%v total bytes written\n", totalBytesWritten)

}
