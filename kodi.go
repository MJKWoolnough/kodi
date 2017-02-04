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

type Client interface {
}

type wrapper struct {
	client Client
}

func New(c Client) (*Server, error) {
	s := rpc.NewServer()
	err := s.Register(wrapper{c})
	if err != nil {
		return nil, err
	}
	return &Server{rpc: s}, nil
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
		s.Connect(repsonseRequest{Reader: r.Body, Writer: w})
		r.Body.Close()
	case http.MethodGet:
		r.ParseForm()
		s.Connect(repsonseRequest{Reader: strings.NewReader(r.Form.Get(request)), Writer: w})
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		//s.rpc.ServeHTTP(w, r)
	}
}

func (s *Server) Listen(l net.Listener) error {
	for {
		c, err := l.Accept()
		if err != nil {
			return err
		}
		// check for websocket?
		go s.Connect(c)
	}
}

func (s *Server) Websocket(c *websocket.Conn) {
	s.Connect(c)
}

func (s *Server) Connect(rwc io.ReadWriteCloser) {
	s.rpc.ServeCodec(jsonrpc.NewServerCodec(rwc))
}
