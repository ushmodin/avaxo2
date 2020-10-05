package multiwriter

import (
	"io"
	"log"
	"sync"
)

// MultiWriter - Writer decorator. Able write to many writers
// If one of writers return error, it will be drop
type MultiWriter struct {
	writers []io.Writer
	lock    sync.Mutex
}

// NewMultiWriter - create new MultiWriter
func NewMultiWriter() *MultiWriter {
	return &MultiWriter{
		writers: []io.Writer{},
	}
}

func (mw *MultiWriter) Write(p []byte) (int, error) {
	log.Printf("Write %d bytes to %d writers\n", len(p), len(mw.writers))
	var newWriters []io.Writer
	mw.lock.Lock()
	var n int
	for _, w := range mw.writers {
		wn, err := w.Write(p)
		if err == nil {
			newWriters = append(newWriters, w)
			n = wn
		} else {
			log.Println("Writer droped")
		}
	}
	mw.writers = newWriters
	mw.lock.Unlock()
	log.Println("Writer ok")
	return n, nil
}

// AddWriter - add writer to multiwriter
func (mw *MultiWriter) AddWriter(w io.Writer) {
	mw.lock.Lock()
	mw.writers = append(mw.writers, w)
	mw.lock.Unlock()
}
