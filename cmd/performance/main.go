package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"golang.org/x/crypto/chacha20"
)

func main() {

	if len(os.Args) != 3 {
		panic("provide in and out path")
	}
	inFilePath := os.Args[1]
	inFile, err := os.Open(inFilePath)
	if err != nil {
		panic("Could not read in file")
	}

	defer inFile.Close()

	outFilePath := os.Args[2]
	outFile, err := os.Create(outFilePath)
	if err != nil {
		panic("Could not create out file")
	}

	defer outFile.Close()

	key := []byte("12345678901234567890123456789012")
	nonce := []byte("123456789012")
	t := time.Now()
	cipher, _ := chacha20.NewUnauthenticatedCipher(key, nonce)
	buffer := make([]byte, 1024)
	for {
		_, err := inFile.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}
		ciphertext := make([]byte, len(buffer))
		cipher.XORKeyStream(ciphertext, buffer)
		outFile.Write(ciphertext)
	}

	t1 := time.Now()
	elapsed := t1.Sub(t)
	fmt.Printf("Elapsed %vs", elapsed.Seconds())
}
