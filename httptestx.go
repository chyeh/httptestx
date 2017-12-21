package httptestx

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
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

func (s *HTTPTestServerExt) Body(body io.Reader) *HTTPTestServerExt {
	b, err := ioutil.ReadAll(body)
	if err != nil {
		panic(err)
	}
	s.b = b
	s.s = httptest.NewServer(s.BuildHandler())
	return s
}

func (s *HTTPTestServerExt) BodyString(str string) *HTTPTestServerExt {
	s.b = []byte(str)
	s.s = httptest.NewServer(s.BuildHandler())
	return s
}

func (s *HTTPTestServerExt) JSON(data interface{}) *HTTPTestServerExt {
	s.h.Set("Content-Type", "application/json")
	b, err := readAndDecode(data, "json")
	if err != nil {
		panic(err)
	}
	s.b = b
	s.s = httptest.NewServer(s.BuildHandler())
	return s
}

func (s *HTTPTestServerExt) XML(data interface{}) *HTTPTestServerExt {
	s.h.Set("Content-Type", "application/xml")
	b, err := readAndDecode(data, "xml")
	if err != nil {
		panic(err)
	}
	s.b = b
	s.s = httptest.NewServer(s.BuildHandler())
	return s
}

func readAndDecode(data interface{}, kind string) ([]byte, error) {
	buf := &bytes.Buffer{}

	switch data.(type) {
	case string:
		buf.WriteString(data.(string))
	case []byte:
		buf.Write(data.([]byte))
	default:
		var err error
		if kind == "xml" {
			err = xml.NewEncoder(buf).Encode(data)
		} else {
			err = json.NewEncoder(buf).Encode(data)
		}
		if err != nil {
			return nil, err
		}
	}

	return ioutil.ReadAll(buf)
}

func (s *HTTPTestServerExt) URL() string {
	return s.s.URL
}

func (s *HTTPTestServerExt) Close() {
	s.s.Close()
}
