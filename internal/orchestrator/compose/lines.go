package compose

import (
	"bytes"
	"sync"
)

// maxLineBytes borne une ligne rendue : au-delà, elle est coupée, pour qu'une
// sortie sans saut de ligne ne grossisse pas sans fin en mémoire.
const maxLineBytes = 8 << 10

// lineWriter découpe un flux en lignes et les passe à out. Un retour chariot
// termine aussi une ligne : c'est ainsi que les barres de progression
// réécrivent la leur.
type lineWriter struct {
	stream OutputStream
	out    OutputFunc

	mu  sync.Mutex
	buf []byte
}

func newLineWriter(stream OutputStream, out OutputFunc) *lineWriter {
	return &lineWriter{stream: stream, out: out}
}

func (w *lineWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	n := len(p)
	for len(p) > 0 {
		i := bytes.IndexAny(p, "\r\n")
		if i < 0 {
			w.buf = append(w.buf, p...)
			for len(w.buf) >= maxLineBytes {
				w.emit(w.buf[:maxLineBytes])
				w.buf = append(w.buf[:0], w.buf[maxLineBytes:]...)
			}
			break
		}
		w.buf = append(w.buf, p[:i]...)
		w.emit(w.buf)
		w.buf = w.buf[:0]
		p = p[i+1:]
	}
	return n, nil
}

// flush rend la ligne en cours, s'il en reste une.
func (w *lineWriter) flush() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.emit(w.buf)
	w.buf = w.buf[:0]
}

// emit rend une ligne ; les lignes vides (dont la moitié d'un « \r\n ») sont
// ignorées.
func (w *lineWriter) emit(line []byte) {
	if len(bytes.TrimSpace(line)) == 0 {
		return
	}
	w.out(w.stream, string(line))
}
