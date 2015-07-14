package artist

import (
	"unicode/utf8"
)

func NewSwastika(input string) string {
	length := utf8.RuneCountInString(input)
	size := length*2 - 1
	mid := (size - 1) / 2
	matrix := NewRuneMatrix(size)

	// From top to bottom, left to right
	matrix.FillHorizontal(input, mid, 0)
	matrix.FillHorizontal(input, 0, mid)
	matrix.FillHorizontal(input, mid, mid)
	matrix.FillHorizontal(input, 0, size-1)

	// From left to right, from to bottom
	matrix.FillVertical(input, 0, 0)
	matrix.FillVertical(input, 0, mid)
	matrix.FillVertical(input, mid, mid)
	matrix.FillVertical(input, mid, size-1)

	return matrix.ToString()
}
