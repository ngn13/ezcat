package c2

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"

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

func infoCallback(j *agent.Job) {
	var (
		pid int64
		res string = ""
		err error
	)

	if j.Response != nil {
		res = string(j.Response)
	}

	if !j.Success {
		log.Debg("info job failed for %s: %s", j.Agent.Session, res)
		j.Agent.DelJob(j.ID)
		return
	}

	sections := strings.Split(res, "@")

	if len(sections) != 4 {
		log.Debg("info callback failed, invalid section count")
		j.Agent.DelJob(j.ID)
		return
	}

	if pid, err = strconv.ParseInt(sections[3], 10, 32); err != nil {
		log.Debg("info callback failed, cannot parse PID (%s)", err.Error())
		j.Agent.DelJob(j.ID)
		return
	}

	j.Agent.Hostname = sections[0]
	j.Agent.Username = sections[1]
	j.Agent.OS = sections[2]
	j.Agent.PID = int32(pid)

	j.Agent.DelJob(j.ID)
}

func (s *Server) HandleReq(con *net.TCPConn, buf *bytes.Buffer) ([]byte, error) {
	var (
		ag  *agent.Agent
		job *agent.Job

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

	if ag = s.Agents.Find(packet.Header.Session); ag == nil && packet.Command() != COMMAND_REGISTER {
		log.Debg("unknown session (%v), asking for registration", packet.Header.Session)

		packet.Reset()
		packet.SetFlags(PACKET_TYPE_RES, COMMAND_AUTH)

		goto send_packet
	}

	switch packet.Command() {
	case COMMAND_REGISTER:
		if packet.Header.Size != builder.ID_LEN {
			return nil, fmt.Errorf("invalid ID len (%v != %d)", packet.Header.Size, builder.ID_LEN)
		}

		if !s.Build.IsValid(string(packet.Data)) {
			return nil, fmt.Errorf("invalid ID")
		}

		ag = s.Agents.New()
		ag.AddJob(COMMAND_INFO, nil, 0, infoCallback)

		packet.Reset()

		packet.SetFlags(PACKET_TYPE_RES, COMMAND_SUCCESS)
		packet.Header.Session = ag.Session

	case COMMAND_ASK:
		packet.Reset()

		if job = ag.NextJob(); job == nil {
			log.Debg("no available job for session (%v)", packet.Header.Session)

			packet.Header.Session = ag.Session
			packet.SetFlags(PACKET_TYPE_RES, COMMAND_NONE)

			break
		}

		packet.Header.Session = ag.Session
		packet.SetFlags(PACKET_TYPE_RES, job.Command)
		packet.Header.JobID = job.ID

		packet.Header.Size = job.ArgumentSize
		packet.Data = job.Argument

	case COMMAND_SUCCESS:
		if job = ag.GetJob(packet.Header.JobID); job == nil {
			return nil, fmt.Errorf("invalid job ID (%v)", packet.Header.JobID)
		}

		job.Success = true
		job.Waiting = false

		job.Response = packet.Data
		job.ResponseSize = packet.Header.Size

		if job.Callback != nil {
			job.Callback(job)
		}

		return nil, nil

	case COMMAND_FAILURE:
		if job = ag.GetJob(packet.Header.JobID); job == nil {
			return nil, fmt.Errorf("invalid job ID (%v)", packet.Header.JobID)
		}

		job.Success = false
		job.Waiting = false

		job.Response = packet.Data
		job.ResponseSize = packet.Header.Size

		if job.Callback != nil {
			job.Callback(job)
		}

		return nil, nil

	default:
		return nil, fmt.Errorf("invalid commad: %v", packet.Command())
	}

	if ag != nil {
		ag.LastCon = time.Now()
		ag.IP = con.RemoteAddr().String()
	}

send_packet:
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
			continue
		}

		if nil == res {
			log.Debg("no response to send")
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

	if s.Listener, err = net.ListenTCP("tcp4", addr); err != nil {
		return err
	}

	go s.Handle()
	return nil
}
