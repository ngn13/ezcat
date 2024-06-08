package agent

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func UtilReadQname(buffer *bytes.Buffer) ([]string, error) {
  var labels []string
	var res string = ""
	c, err := buffer.ReadByte()

	for ; err == nil && int(c) > 0; c, err = buffer.ReadByte() {
    label_len := int(c)

		if label_len > LABEL_LIMIT {
			return labels, fmt.Errorf("illegal label size (%d)", label_len)
		}

    if len(res)+label_len > NAME_LIMIT {
			return labels, fmt.Errorf("illegal name size (%d)", len(res)+label_len)
		}

		part := string(buffer.Next(label_len))
    labels = append(labels, part)

		if res == "" {
			res = part
			continue
		}

		res += "." + part
	}

	return labels, err
}

func UtilWriteQname(buffer *bytes.Buffer, qname []string) error {
  var total int = 0

  for _, p := range qname {
    label_len := len(p)
    if label_len > LABEL_LIMIT {
      return fmt.Errorf("illegal label size")
    }
    
    total += label_len

    if total > NAME_LIMIT {
      return fmt.Errorf("illegal name size")
    }

    buffer.WriteByte(byte(label_len))
    buffer.Write([]byte(p))
  }

  buffer.WriteByte(byte(0))
  return nil
}

func UtilWriteRR(buffer *bytes.Buffer, rr *DNS_RR) error {
  err := UtilWriteQname(buffer, rr.Name)
  if err != nil {
    return err
  }
  
  binary.Write(buffer, binary.BigEndian, rr.Type)
  binary.Write(buffer, binary.BigEndian, rr.Class)
  binary.Write(buffer, binary.BigEndian, rr.TTL)
  binary.Write(buffer, binary.BigEndian, rr.RDLength)
  binary.Write(buffer, binary.BigEndian, rr.RData)
  return nil
}

func UtilPacketToBuffer(packet *DNS_Packet) (*bytes.Buffer, error) {
  buffer := new(bytes.Buffer)
  binary.Write(buffer, binary.BigEndian, packet.Header)
  var i uint16 = 0

  for i = 0; i < packet.Header.QDCount; i++ {
    err := UtilWriteQname(buffer, packet.Questions[i].Qname)
    if err != nil {
      return nil, err
    }

    binary.Write(buffer, binary.BigEndian, packet.Questions[i].Qtype)
    binary.Write(buffer, binary.BigEndian, packet.Questions[i].Qclass)
  }

  for i = 0; i < packet.Header.ANCount; i++ {
    err := UtilWriteRR(buffer, &packet.Answers[i])
    if err != nil {
      return nil, err
    }
  }

  for i = 0; i < packet.Header.NSCount; i++ {
    err := UtilWriteRR(buffer, &packet.Authorities[i])
    if err != nil {
      return nil, err
    }
  }

  for i = 0; i < packet.Header.ARCount; i++ {
    err := UtilWriteRR(buffer, &packet.Additionals[i])
    if err != nil {
      return nil, err
    }
  }

  return buffer, nil
}
