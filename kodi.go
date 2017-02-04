package kodi

import (
	"io"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"strings"

	"golang.org/x/net/websocket"
)

type Server struct {
	rpc *rpc.Server
}

type readWriter struct {
	io.Reader
	io.Writer
}

func (r readWriter) Close() error {
	if c, ok := r.Reader.(io.Closer); ok {
		c.Close()
	}
	if c, ok := r.Writer.(io.Closer); ok {
		c.Close()
	}
	return nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.runConn(repsonseRequest{Reader: r.Body, Writer: w})
	case http.MethodGet:
		r.ParseForm()
		s.runConn(repsonseRequest{Reader: strings.NewReader(r.Form.Get(request)), Writer: w})
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) Listen(l net.Listener) error {
	for {
		c, err := l.Accept()
		if err != nil {
			return err
		}
		go s.runConn(c)
	}
}

func (s *Server) Websocket(c *websocket.Conn) {
	s.runConn(c)
}

func (s *Server) Connect(c io.ReadWriteCloser) {
	s.runConn(c)
}

func (s *Server) runConn(rwc io.ReadWriteCloser) {
	s.rpc.ServeCodec(jsonrpc.NewServerCodec(rwc))
}
