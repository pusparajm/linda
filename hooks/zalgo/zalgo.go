package zalgo

import (
	"bytes"
	"fmt"

	"github.com/kortschak/zalgo"
)

// Zalgo hook
type Zalgo struct{}

// New Zalgo hook
func New() *Zalgo {
	return new(Zalgo)
}

// Zalgo input
func (hook *Zalgo) Fire(input string) string {
	buf := new(bytes.Buffer)
	z := zalgo.NewCorrupter(buf)

	z.Zalgo = func(n int, r rune, z *zalgo.Corrupter) bool {
		z.Up += 0.1
		z.Middle += complex(0.01, 0.01)
		z.Down += complex(real(z.Down)*0.1, 0)
		return false
	}

	z.Up = complex(0, 0.2)
	z.Middle = complex(0, 0.2)
	z.Down = complex(0.001, 0.3)

	fmt.Fprintf(z, input)
	return buf.String()
}
