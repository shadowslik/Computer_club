package middleware

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/time/rate"
	"net/http"
)

func TestStatusRecorder_WriteHeader(t *testing.T) {
	rec := httptest.NewRecorder()
	sr := &statusRecorder{ResponseWriter: rec, status: 200}

	sr.WriteHeader(404)
	assert.Equal(t, 404, sr.status)
	assert.Equal(t, 404, rec.Code)
}

func TestTrackingWriter_Write_CountsBytes(t *testing.T) {
	rec := httptest.NewRecorder()
	tw := &trackingWriter{ResponseWriter: rec}

	n, err := tw.Write([]byte("hello"))
	require.NoError(t, err)
	assert.Equal(t, 5, n)
	assert.Equal(t, int64(5), tw.Written())
}

func TestTrackingWriter_MultipleWrites(t *testing.T) {
	rec := httptest.NewRecorder()
	tw := &trackingWriter{ResponseWriter: rec}

	tw.Write([]byte("abc"))
	tw.Write([]byte("de"))
	assert.Equal(t, int64(5), tw.Written())
}

func TestRateLimitedWriter_Write(t *testing.T) {
	rec := httptest.NewRecorder()
	lim := rate.NewLimiter(rate.Inf, 0) // unlimited
	rlw := &rateLimitedWriter{
		ResponseWriter: rec,
		limiter:        lim,
		ctx:            context.Background(),
	}

	n, err := rlw.Write([]byte("test data"))
	require.NoError(t, err)
	assert.Equal(t, 9, n)
	assert.Equal(t, int64(9), rlw.Written())
}

func TestRateLimitedReader_ReadAndClose(t *testing.T) {
	lim := rate.NewLimiter(rate.Inf, 0)
	reader := io.NopCloser(strings.NewReader("hello world"))
	rlr := &rateLimitedReader{
		reader:  reader,
		limiter: lim,
		ctx:     context.Background(),
	}

	buf := make([]byte, 5)
	n, err := rlr.Read(buf)
	require.NoError(t, err)
	assert.Equal(t, 5, n)
	assert.Equal(t, "hello", string(buf[:n]))

	assert.NoError(t, rlr.Close())
}

func TestStatusRecorder_WritesThrough(t *testing.T) {
	rec := httptest.NewRecorder()
	sr := &statusRecorder{ResponseWriter: rec, status: http.StatusOK}
	sr.Header().Set("X-Test", "yes")
	assert.Equal(t, "yes", rec.Header().Get("X-Test"))
}
