package c2

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/ngn13/ezcat/server/log"
)

const (
	PACKET_MAX_SIZE = 255 + 4 + 3
	PACKET_TYPE_REQ = 0
	PACKET_TYPE_RES = 1
	PACKET_VERSION  = 0
)

const (
	COMMAND_FAILURE  = 0
	COMMAND_SUCCESS  = 1
	COMMAND_REGISTER = 2
	COMMAND_KILL     = 3
	COMMAND_RUN      = 4
	COMMAND_ASK      = 5
	COMMAND_INFO     = 6
	COMMAND_NONE     = 7
	COMMAND_AUTH     = 8
)

type Packet struct {
	Header struct {
		Flags   uint8  // version (3 bits), 1 (type), 4 (command)
		Session uint32 // agent session
		JobID   uint16 // id of the work this packet is assoicated with
		Size    uint8  // data size
	}
	Data []byte
}

func (p *Packet) Reset() {
	p.Header.Session = 0
	p.Header.JobID = 0
	p.Header.Flags = 0
	p.Header.Size = 0

	p.Data = []byte{}
}

func (p *Packet) Version() uint8 {
	return (p.Header.Flags >> 5) & 0b11
}

func (p *Packet) IsRequest() bool {
	return (p.Header.Flags>>4)&1 == PACKET_TYPE_REQ
}

func (p *Packet) Command() uint8 {
	return p.Header.Flags & 0b1111
}

func (p *Packet) SetFlags(typ uint8, cmd uint8) {
	p.Header.Flags = 0
	p.Header.Flags |= (cmd & 0b1111)
	p.Header.Flags |= (typ & 1) << 4
	p.Header.Flags |= (PACKET_VERSION & 0b111) << 5
}

func (p *Packet) Read(b *bytes.Buffer) (err error) {
	if err = binary.Read(b, binary.BigEndian, &p.Header); err != nil {
		log.Debg("failed to read packet header: %s", err.Error())
		return err
	}

	if p.Version() != PACKET_VERSION {
		return fmt.Errorf("version mismatch")
	}

	var size int
	p.Data = make([]byte, p.Header.Size)

	if size, err = b.Read(p.Data); err != nil {
		log.Debg("failed to read data: %s", err.Error())
		return err
	}

	if size != int(p.Header.Size) {
		log.Debg("failed to read all the packet data (%d/%v)", size, p.Header.Size)
		return fmt.Errorf("failed to read all the packet data")
	}

	return nil
}

func (p *Packet) Write(b *bytes.Buffer) (err error) {
	if err = binary.Write(b, binary.BigEndian, &p.Header); err != nil {
		log.Debg("failed to write packet header: %s", err.Error())
		return err
	}

	var size int

	if size, err = b.Write(p.Data[:p.Header.Size]); err != nil {
		log.Debg("failed to write data: %s", err.Error())
		return err
	}

	if size != int(p.Header.Size) {
		log.Debg("failed to write all the packet data (%d/%v)", size, p.Header.Size)
		return fmt.Errorf("failed to write all the packet data")
	}

	return nil
}
