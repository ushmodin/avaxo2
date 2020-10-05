package util

import (
	"errors"
	"io"
	"net/http"
)

// FlushWriter decorator for http.ResponseWriter add immediate flush after write
type FlushWriter struct {
	io.Writer
	w io.Writer
	f http.Flusher
}

// NewFlushWriter create new FlushWriter
func NewFlushWriter(w http.ResponseWriter) (FlushWriter, error) {
	f, ok := w.(http.Flusher)
	if !ok {
		return FlushWriter{}, errors.New("Not implements http.Flusher")
	}
	return FlushWriter{
		f: f,
		w: w,
	}, nil
}

func (fw FlushWriter) Write(p []byte) (n int, err error) {
	n, err = fw.w.Write(p)
	if fw.f != nil {
		fw.f.Flush()
	}
	return
}
