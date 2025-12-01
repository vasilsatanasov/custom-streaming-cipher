package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"time"
	"vsatanasov/custom-streaming-algorithm/pkg/cipher"
)

func main() {
	args := os.Args[1:]
	if len(args) < 3 {
		panic(`3 arguments must be provided: 
		 - key for the cipher 
		 - the input file destination
		 - the output file destination`)
	}
	key := args[0]

	inPath := args[1]
	inFile, err := os.Open(inPath)
	if err != nil {
		panic(fmt.Sprintf("Could not open file %s", inPath))
	}

	defer inFile.Close()

	outPath := args[2]
	outFile, err := os.Create(outPath)
	if err != nil {
		panic(fmt.Sprintf("Could not create file %s", inPath))
	}

	defer outFile.Close()
	iv := rand.Text()
	t := time.Now()
	c := cipher.New([]byte(key), []byte(iv))

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
		encrypted := c.Encode(buffer)
		outFile.Write(encrypted)
	}

	t1 := time.Now()
	elapsed := t1.Sub(t)
	fmt.Printf("Elapsed %vs", elapsed.Seconds())

}
