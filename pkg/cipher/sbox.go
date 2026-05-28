package cipher

type SBox struct {
	matrix  [][]byte
	size    int
	nextCol int
	nextRow int
}

func (s *SBox) NextByte() byte {
	value := s.matrix[s.nextRow][s.nextCol]
	nextCol := (s.nextCol + 1) % s.size
	nextRow := s.nextRow
	if nextCol == 0 {
		nextRow = (s.nextRow + 1) % s.size
	}

	s.nextCol = nextCol
	s.nextRow = nextRow
	if s.nextRow == 0 && s.nextCol == 0 {
		s.shiftRows()
		s.shiftColumns()
	}

	return value
}

func NewSbox(bytes []byte, size int) *SBox {
	return &SBox{
		matrix:  makeSBox(bytes, size),
		size:    size,
		nextRow: 0,
		nextCol: 0,
	}
}

func makeSBox(bytes []byte, size int) [][]byte {
	sBox := make([][]byte, 0)
	for i := range size {
		row := make([]byte, 0)
		for j := range size {
			bytesIdx := i*size + j
			if len(bytes)-1 < bytesIdx {
				row = append(row, 0)
			} else {
				row = append(row, bytes[bytesIdx])
			}
		}
		sBox = append(sBox, row)
	}

	return sBox
}

func (s *SBox) XORAt(row, col int, b byte) {
	s.matrix[row][col] ^= b
}

func (s *SBox) shiftColumns() {

	tempCol := make([]byte, s.size)
	for r := range s.size {
		tempCol[r] = s.matrix[r][0]
		s.matrix[r][0] = s.matrix[r][s.size-1]
	}
	for c := 1; c < s.size; c++ {
		for r := range s.size {
			t := tempCol[r]
			tempCol[r] = s.matrix[r][c]
			s.matrix[r][c] = t
		}

	}
}

func (s *SBox) shiftRows() {
	temp := make([]byte, s.size)
	for r := range s.size {
		temp[r] = s.matrix[0][r]
		s.matrix[0][r] = s.matrix[s.size-1][r]
	}
	for r := 1; r < s.size; r++ {
		for c := range s.size {
			t := temp[c]
			temp[c] = s.matrix[r][c]
			s.matrix[r][c] = t
		}

	}
}
