package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

// Iterator that reads next byte from file
type FileByteIter struct {
	reader *bufio.Reader
}

func NewFileByteIter(path string) (*FileByteIter, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return &FileByteIter{reader: bufio.NewReader(f)}, nil
}

func (it *FileByteIter) Next() (byte, error) {
	b, err := it.reader.ReadByte()
	return b, err
}

// Brent cycle detection adapted for file iterator
func BrentCycleDetection(next func() (byte, error)) (lambda, mu int, err error) {
	// Read initial byte
	x0, err := next()
	if err != nil {
		return 0, 0, fmt.Errorf("file too small or error reading: %w", err)
	}

	power := 1
	lam := 1

	tortoise := x0
	hare, err := next()
	if err != nil {
		return 0, 0, fmt.Errorf("not enough data: %w", err)
	}

	// Phase 1: Find λ
	for tortoise != hare {
		if power == lam {
			tortoise = hare
			power *= 2
			lam = 0
		}
		hare, err = next()
		if err != nil {
			return 0, 0, fmt.Errorf("period not found (data ended): %w", err)
		}
		lam++
	}

	// Phase 2: Find μ
	mu = 0
	tortoise = x0
	hare = x0

	// advance hare by λ steps
	for i := 0; i < lam; i++ {
		hare, err = next()
		if err != nil {
			return 0, 0, fmt.Errorf("not enough data to verify cycle start: %w", err)
		}
	}

	for tortoise != hare {
		tortoise, err = next()
		if err != nil {
			return 0, 0, fmt.Errorf("unexpected EOF during mu scan: %w", err)
		}
		hare, err = next()
		if err != nil {
			return 0, 0, fmt.Errorf("unexpected EOF during mu scan: %w", err)
		}
		mu++
	}

	return lam, mu, nil
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: period <binary_file>")
	}

	iter, err := NewFileByteIter(os.Args[1])
	if err != nil {
		log.Fatal("Error:", err)
	}

	period, start, err := BrentCycleDetection(iter.Next)
	if err != nil {
		log.Fatal("Error:", err)
	}

	fmt.Printf("Detected period: %d bytes\n", period)
	fmt.Printf("Cycle start at offset: %d\n", start)
}
