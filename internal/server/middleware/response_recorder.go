package middleware

import (
	"context"
	"io"
	"net/http"

	"golang.org/x/time/rate"
)

type responseWriterWithSize interface {
	http.ResponseWriter
	Written() int64
}

type rateLimitedWriter struct {
	http.ResponseWriter
	limiter *rate.Limiter
	ctx     context.Context
	written int64
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (sr *statusRecorder) WriteHeader(code int) {
	sr.status = code
	sr.ResponseWriter.WriteHeader(code)
}

func (w *rateLimitedWriter) Write(p []byte) (int, error) {
	total := 0
	chunkSize := 32 * 1024
	for len(p) > 0 {
		chunk := p
		if len(chunk) > chunkSize {
			chunk = p[:chunkSize]
		}
		if err := w.limiter.WaitN(w.ctx, len(chunk)); err != nil {
			return total, err
		}
		n, err := w.ResponseWriter.Write(chunk)
		total += n
		w.written += int64(n)
		if err != nil {
			return total, err
		}
		p = p[n:]
	}
	return total, nil
}

func (w *rateLimitedWriter) Written() int64 { return w.written }

type rateLimitedReader struct {
	reader  io.ReadCloser
	limiter *rate.Limiter
	ctx     context.Context
}

func (r *rateLimitedReader) Read(p []byte) (int, error) {
	n, err := r.reader.Read(p)
	if n > 0 {
		if waitErr := r.limiter.WaitN(r.ctx, n); waitErr != nil {
			return n, waitErr
		}
	}
	return n, err
}

func (r *rateLimitedReader) Close() error { return r.reader.Close() }

type trackingWriter struct {
	http.ResponseWriter
	written int64
}

func (tw *trackingWriter) Write(b []byte) (int, error) {
	n, err := tw.ResponseWriter.Write(b)
	tw.written += int64(n)
	return n, err
}

func (tw *trackingWriter) Written() int64 { return tw.written }
