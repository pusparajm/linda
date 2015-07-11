package swastika

import (
	"errors"
)

const (
	MinLen = 3
)

func New(input string) (string, error) {
	length := len(input)
	if length < MinLen {
		return "", errors.New("Input string length should be more than 2")
	}

	size := length*2 - 1
	mid := (size - 1) / 2

	//reverse := reverseString(input)
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

	return matrix.ToString(), nil
}

func reverseString(s string) string {
	n := len(s)
	runes := make([]rune, n)
	for _, rune := range s {
		n--
		runes[n] = rune
	}
	return string(runes[n:])
}
