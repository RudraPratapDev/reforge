package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/RudraPratapDev/reforge/resp"
)

type server struct {
	lnAddr string
	ln     net.Listener
	quitCh chan struct{}
}

func newServer(lnAddr string) *server {
	return &server{
		lnAddr: lnAddr,
		quitCh: make(chan struct{}),
	}
}

func (s *server) start() error {
	ln, err := net.Listen("tcp", s.lnAddr)
	if err != nil {
		return err
	}
	defer ln.Close()

	s.ln = ln
	go s.acceptLoop()
	<-s.quitCh

	return nil
}

func (s *server) Shutdown() {
	close(s.quitCh)
	s.ln.Close()
}

func (s *server) acceptLoop() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}
			log.Printf("accept error: %v", err)
			continue
		}
		log.Printf("client connected: %s", conn.RemoteAddr())
		go func() {
			if err := s.handleConn(conn); err != nil {
				log.Printf("connection error (%s): %v", conn.RemoteAddr(), err)
			}
		}()

	}
}

func (s *server) handleConn(conn net.Conn) error {
	defer conn.Close()
	decoder := resp.NewDecoder(conn)
	for {

		value, err := decoder.Decode()
		if err == io.EOF {
			log.Printf("client disconnected: %s", conn.RemoteAddr())
			return nil
		}
		if err != nil {
			return err
		}
		fmt.Println(value)

	}
}

func main() {
	srv := newServer("localhost:3000")
	log.Fatal(srv.start())

}
