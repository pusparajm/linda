package swastika

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
	for i := 0; i < len(input); i++ {
		j := from + i

		if j < m.Size {
			m.Table[line][j] = rune(input[i])
		}
	}

	return m
}

func (m *RuneMatrix) FillVertical(input string, from, column int) *RuneMatrix {
	for i := 0; i < len(input); i++ {
		j := from + i

		if j < m.Size {
			m.Table[j][column] = rune(input[i])
		}
	}

	return m
}

func (m *RuneMatrix) ToString() string {
	result := ""

	for i := 0; i < m.Size; i++ {
		line := ""
		for j := 0; j < m.Size; j++ {
			line += string(m.Table[i][j]) + " "
		}

		result += line + "\n"
	}

	return result
}
