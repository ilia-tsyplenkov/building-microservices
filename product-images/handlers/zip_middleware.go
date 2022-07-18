package handlers

import (
	"compress/gzip"
	"net/http"
	"strings"
)

type WrappedResposeWriter struct {
	rw  http.ResponseWriter
	gzw *gzip.Writer
}

func NewWrappedResponseWriter(rw http.ResponseWriter) *WrappedResposeWriter {
	gz := gzip.NewWriter(rw)
	return &WrappedResposeWriter{rw: rw, gzw: gz}
}
func (w *WrappedResposeWriter) Header() http.Header {
	return w.rw.Header()
}

func (w *WrappedResposeWriter) Write(p []byte) (int, error) {
	return w.gzw.Write(p)
}

func (w *WrappedResposeWriter) WriteHeader(code int) {
	w.rw.WriteHeader(code)
}

func (w *WrappedResposeWriter) Flush() error {
	w.gzw.Flush()
	w.gzw.Close()

	return nil
}

type GzipHandler struct {
}

func NewGzipHandler() *GzipHandler {
	return &GzipHandler{}
}

func (g *GzipHandler) GzipMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			// rw.Write([]byte("hello"))
			// return
			grw := NewWrappedResponseWriter(rw)
			grw.Header().Set("Content-Encoding", "gzip")
			defer grw.Flush()
			next.ServeHTTP(grw, r)
			return
		}
		// call next
		next.ServeHTTP(rw, r)

	})
}
