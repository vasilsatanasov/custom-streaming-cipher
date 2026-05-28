package main

import (
	"bufio"
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
		panic(fmt.Sprintf("Could not create file %s", outPath))
	}

	defer outFile.Close()

	c := cipher.New([]byte(key))

	reader := bufio.NewReaderSize(inFile, 64*1024)
	writer := bufio.NewWriterSize(outFile, 64*1024)
	defer writer.Flush()

	buffer := make([]byte, 64*1024)
	t := time.Now()
	for {
		n, err := reader.Read(buffer)
		if n > 0 {
			encrypted := c.Encode(buffer[:n])
			_, werr := writer.Write(encrypted)
			if werr != nil {
				panic(werr)
			}
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
	}

	t1 := time.Now()
	elapsed := t1.Sub(t)
	fmt.Printf("Elapsed %vs", elapsed.Seconds())

}
