package c2

import (
	"bytes"
	"fmt"
	"io"
	"net"

	"github.com/ngn13/ezcat/server/agent"
	"github.com/ngn13/ezcat/server/builder"
	"github.com/ngn13/ezcat/server/log"
)

type Server struct {
	Listener *net.TCPListener
	Agents   *agent.List
	Build    *builder.Struct
}

func New(build *builder.Struct, list *agent.List) (s *Server) {
	return &Server{
		Build:  build,
		Agents: list,
	}
}

func (s *Server) HandleReq(con *net.TCPConn, buf *bytes.Buffer) ([]byte, error) {
	var (
		packet Packet
		err    error
	)

	if err = packet.Read(buf); err != nil {
		return nil, err
	}

	if !packet.IsRequest() {
		return nil, fmt.Errorf("packet is not a request")
	}

	buf.Reset()

	switch packet.Command() {
	case COMMAND_REGISTER:
		if packet.Size == builder.ID_LEN {
			return nil, fmt.Errorf("invalid ID len")
		}

		if s.Build.IsValid(string(packet.Data)) {
			return nil, fmt.Errorf("invalid ID")
		}

		agent := s.Agents.New()
		packet.Reset()

		packet.Session = agent.Session
		packet.SetFlags(PACKET_TYPE_RES, COMMAND_SUCCESS)
		packet.WorkID = 0
		packet.Size = 0
	}

	if err = packet.Write(buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
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

		if size, err = con.Read(req); err != nil && err != net.ErrClosed && err != io.EOF {
			log.Debg("failed to read TCP connection from %s: %s", con.RemoteAddr().String(), err.Error())
			continue
		}

		if err == net.ErrClosed || err == io.EOF {
			log.Debg("closing TCP connection from: %s", con.RemoteAddr().String())
			con.Close()
			break
		}

		log.Debg("received %d bytes from %s", size, con.RemoteAddr().String())

		if res, err = s.HandleReq(con, bytes.NewBuffer(req)); err != nil {
			log.Debg("failed to handle request from %s: %s", con.RemoteAddr().String(), err.Error())
			con.Close()
			return
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

	if s.Listener, err = net.ListenTCP("tcp4", addr); err != nil {
		return err
	}

	go s.Handle()
	return nil
}
