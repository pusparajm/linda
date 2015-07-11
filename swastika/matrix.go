package swastika

import (
	"unicode/utf8"
)

type RuneMatrix struct {
	Table [][]rune
	Size  int
}

func NewRuneMatrix(size int) *RuneMatrix {
	m := new(RuneMatrix)
	m.Size = size
	m.Table = make([][]rune, size)

	for i := 0; i < size; i++ {
		m.Table[i] = make([]rune, size)
	}

	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			m.Table[i][j] = ' '
		}
	}

	return m
}

func (m *RuneMatrix) FillHorizontal(input string, from, line int) *RuneMatrix {
	b := []byte(input)
	i := 0
	for len(b) > 0 {
		j := from + i

		if j < m.Size {
			char, size := utf8.DecodeRune(b)
			m.Table[line][j] = char
			b = b[size:]
			i++
		} else {
			break
		}
	}

	return m
}

func (m *RuneMatrix) FillVertical(input string, from, column int) *RuneMatrix {
	b := []byte(input)
	i := 0
	for len(b) > 0 {
		j := from + i

		if j < m.Size {
			char, size := utf8.DecodeRune(b)
			m.Table[j][column] = char
			b = b[size:]
			i++
		} else {
			break
		}
	}

	return m
}

func (m *RuneMatrix) runeToString(r rune) string {
	buf := make([]byte, 3)
	utf8.EncodeRune(buf, r)
	return string(buf)
}

func (m *RuneMatrix) ToString() string {
	result := ""

	for i := 0; i < m.Size; i++ {
		line := ""
		for j := 0; j < m.Size; j++ {
			line += m.runeToString(m.Table[i][j]) + " "
		}

		result += line + "\n"
	}

	return result
}
