package c2

import (
	"bytes"
	"net"

	"github.com/ngn13/ezcat/server/log"
)

type Server struct {
	Listener *net.TCPListener
}

func (s *Server) HandleReq(con *net.TCPConn, buf *bytes.Buffer) ([]byte, error) {
	return []byte{}, nil
}

func (s *Server) HandleCon(con *net.TCPConn) {
	for {
		var (
			req  []byte
			res  []byte
			size int
			err  error
		)

		req = make([]byte, PACKET_MAX_SIZE)

		if size, err = con.Read(req); err != nil && err != net.ErrClosed {
			log.Debg("failed to read TCP connection from %s: %s", con.RemoteAddr().String(), err.Error())
			continue
		}

		if err == net.ErrClosed {
			log.Debg("closing TCP connection from: %s", con.RemoteAddr().String())
			break
		}

		log.Debg("received %d bytes from %s", size, con.RemoteAddr().String())

		if res, err = s.HandleReq(con, bytes.NewBuffer(req)); err != nil {
			log.Debg("failed to handle request from %s: %s", con.RemoteAddr().String(), err.Error())
			continue
		}

		if _, err = con.Write(res); err != nil {
			log.Debg("failed to send the response to %s: %s", con.RemoteAddr().String(), err.Error())
			continue
		}
	}
}

func (s *Server) Handle() error {
	defer s.Listener.Close()

	for {
		var (
			con *net.TCPConn
			err error
		)

		if con, err = s.Listener.AcceptTCP(); err != nil {
			log.Debg("failed to accept the TCP connection: %s", err.Error())
			continue
		}

		log.Debg("accepted connection from %s", con.RemoteAddr().String())
		go s.HandleCon(con)
	}
}

func (s *Server) Listen(_addr string) (err error) {
	var addr *net.TCPAddr

	if addr, err = net.ResolveTCPAddr("tcp", _addr); err != nil {
		return err
	}

	if s.Listener, err = net.ListenTCP("tcp", addr); err != nil {
		return err
	}

	go s.Handle()
	return nil
}
