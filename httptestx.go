package httptestx

import (
	"net/http"
	"net/http/httptest"
)

var server *HTTPTestServerExt

func init() {
	server = NewServer()
}

func Server() *HTTPTestServerExt {
	return server
}

func Serve() *HTTPTestServerExt {
	return server.Serve()
}

func URL() string {
	return server.URL()
}

func Close() {
	server.Close()
}

type HTTPTestServerExt struct {
	s      *httptest.Server
	h      http.Header
	status int
	b      []byte
}

func NewServer() *HTTPTestServerExt {
	h := http.HandlerFunc(nil)
	return &HTTPTestServerExt{
		s: httptest.NewServer(h),
		b: []byte(""),
		h: make(http.Header),
	}
}

func (s *HTTPTestServerExt) Serve() *HTTPTestServerExt {
	s.s.Close()
	s.s = httptest.NewServer(s.BuildHandler())
	return s
}

func (s *HTTPTestServerExt) BuildHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		for k, vals := range s.h {
			for _, v := range vals {
				w.Header().Add(k, v)
			}
		}

		w.WriteHeader(s.status)

		_, err := w.Write(s.b)
		if err != nil {
			panic(err)
		}
	}
}

func (s *HTTPTestServerExt) Status(code int) *HTTPTestServerExt {
	s.s.Close()
	s.status = code
	s.s = httptest.NewServer(s.BuildHandler())
	return s
}

func (s *HTTPTestServerExt) AddHeader(key, value string) *HTTPTestServerExt {
	s.s.Close()
	s.h.Add(key, value)
	s.s = httptest.NewServer(s.BuildHandler())
	return s
}

func (s *HTTPTestServerExt) URL() string {
	return s.s.URL
}

func (s *HTTPTestServerExt) Close() {
	s.s.Close()
}
