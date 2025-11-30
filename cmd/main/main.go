package main

import (
	"crypto/rand"
	"fmt"
	"os"
	"strconv"
	"vsatanasov/custom-streaming-algorithm/pkg/cipher"
)

func main() {
	args := os.Args[1:]
	if len(args) < 3 {
		panic(`3 arguments must be provided: 
		 - key for the cipher 
		 - the lenght in bits of generate PSR
		 - the output file destination`)
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
	iv := rand.Text()
	c := cipher.New([]byte(key), []byte(iv))

	curentSize := 0
	currentByte := byte(0)
	currentByteIndex := 0
	totalBytesWritten := 0
	buffer := make([]byte, 0, 1024)
	for curentSize < size {
		b := c.Тick()
		currentByte = (currentByte << 1) | (b & 1)
		currentByteIndex++
		if currentByteIndex == 8 {
			buffer = append(buffer, currentByte)
			currentByte = byte(0)
			currentByteIndex = 0
		} else if curentSize == size-1 {
			currentByte <<= (7 - currentByteIndex)
			buffer = append(buffer, currentByte)
		}
		if len(buffer) == 1024 || curentSize == size-1 {
			n, err := file.Write(buffer)
			buffer = make([]byte, 0, 1024)
			if err != nil {
				panic("Could not write to file")
			}
			totalBytesWritten += n
		}

		curentSize++
	}

	fmt.Printf("%v total bytes written\n", totalBytesWritten)

}
