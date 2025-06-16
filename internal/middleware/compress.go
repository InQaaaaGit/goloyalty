package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// Compress middleware для сжатия ответов
func Compress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем, поддерживает ли клиент сжатие
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		// Создаем gzip writer
		gz, err := gzip.NewWriterLevel(w, gzip.BestCompression)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		defer gz.Close()

		// Устанавливаем заголовки
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Vary", "Accept-Encoding")

		// Создаем response writer с сжатием
		gw := &gzipResponseWriter{
			ResponseWriter: w,
			Writer:         gz,
		}

		next.ServeHTTP(gw, r)
	})
}

type gzipResponseWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (g *gzipResponseWriter) Write(data []byte) (int, error) {
	return g.Writer.Write(data)
}
