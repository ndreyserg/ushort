package handlers

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type compressWriter struct {
	http.ResponseWriter
	gw   *gzip.Writer
	code int
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		ResponseWriter: w,
	}
}

func (cw *compressWriter) Write(b []byte) (int, error) {

	ct := cw.Header().Get("Content-Type")
	if ct == "application/json" || ct == "text/html" {
		if cw.gw == nil {
			cw.Header().Add("Content-Encoding", "gzip")
			cw.Header().Del("Content-Length")
			cw.gw = gzip.NewWriter(cw.ResponseWriter)
			if cw.code != 0 {
				cw.ResponseWriter.WriteHeader(cw.code)
				cw.code = 0
			}
		}
		return cw.gw.Write(b)
	}
	if cw.code != 0 {
		cw.ResponseWriter.WriteHeader(cw.code)
		cw.code = 0
	}
	return cw.ResponseWriter.Write(b)
}

func (cw *compressWriter) WriteHeader(code int) {
	if cw.code == 0 {
		cw.code = code
	}
}

func (cw *compressWriter) Close() error {
	if cw.code != 0 {
		cw.ResponseWriter.WriteHeader(cw.code)
		cw.code = 0
	}

	if cw.gw != nil {
		return cw.gw.Close()
	}
	return nil
}

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

func gzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w
		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			cw := newCompressWriter(w)
			ow = cw
			defer cw.Close()
		}

		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")

		if sendsGzip {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			r.Body = cr
			defer cr.Close()
		}

		next.ServeHTTP(ow, r)
	})
}
