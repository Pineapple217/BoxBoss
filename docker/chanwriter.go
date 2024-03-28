package docker

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
)

type ChanWriter struct {
	ch chan<- string
}

func NewChanWriter(ch chan<- string) *ChanWriter {
	return &ChanWriter{ch: ch}
}

func (w *ChanWriter) Write(p []byte) (int, error) {
	// n := len(p)
	j, err := json.Marshal(p)
	if err != nil {
		slog.Warn("chanwrite marshal error", "err", err)
		return 0, err
	}
	w.ch <- string(j)
	slog.Info("???", "len_j", len(j), "len_p", len(p))
	return len(p) - 1, nil
	// ik was write error va chan writer aan het debuggen
}

func NewFixLinebreakMiddleware(next io.Writer) io.Writer {
	return fixLinebreakMiddleware{next}
}

type fixLinebreakMiddleware struct {
	next io.Writer
}

func (c fixLinebreakMiddleware) Write(p []byte) (n int, err error) {
	// Replace each newline character with \r\n before writing
	replaced := bytes.ReplaceAll(p, []byte("\n"), []byte("\r\n"))
	return c.next.Write(replaced)
}
