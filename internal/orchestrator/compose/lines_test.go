package compose

import (
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type collected struct {
	mu    sync.Mutex
	lines []string
}

func (c *collected) out(stream OutputStream, line string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	prefix := "out:"
	if stream == Stderr {
		prefix = "err:"
	}
	c.lines = append(c.lines, prefix+line)
}

func TestLineWriter_SplitsLines(t *testing.T) {
	var c collected
	w := newLineWriter(Stderr, c.out)

	// Les morceaux arrivent coupés n'importe où, comme les lectures d'un tube.
	for _, chunk := range []string{" Container a Crea", "ting \n Container a Created\r\n", "\n", "Pull: 10%\rPull: 100%\r", "fin sans saut"} {
		n, err := w.Write([]byte(chunk))
		assert.NoError(t, err)
		assert.Equal(t, len(chunk), n, "Write rend toujours tout ce qu'il reçoit")
	}
	assert.Equal(t, []string{"err: Container a Creating ", "err: Container a Created", "err:Pull: 10%", "err:Pull: 100%"}, c.lines,
		"un retour chariot termine une ligne ; les lignes vides sont ignorées")

	w.flush()
	assert.Equal(t, "err:fin sans saut", c.lines[len(c.lines)-1], "flush rend la dernière ligne incomplète")
	w.flush()
	assert.Len(t, c.lines, 5, "un second flush n'ajoute rien")
}

func TestLineWriter_CutsVeryLongLines(t *testing.T) {
	var c collected
	w := newLineWriter(Stdout, c.out)

	_, _ = w.Write([]byte(strings.Repeat("x", 2*maxLineBytes+10)))
	assert.Len(t, c.lines, 2, "une sortie sans saut de ligne est rendue par morceaux bornés")
	for _, l := range c.lines {
		assert.Len(t, l, len("out:")+maxLineBytes)
	}
	w.flush()
	assert.Equal(t, "out:"+strings.Repeat("x", 10), c.lines[2])
}
