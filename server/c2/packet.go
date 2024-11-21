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
	COMMAND_SUCCESS  = 0
	COMMAND_REGISTER = 1
)

type Packet struct {
	Flags   uint8  // version (3 bits), 1 (type), 4 (command)
	Session uint32 // agent session
	WorkID  uint16 // id of the work this packet is assoicated with
	Size    uint8  // data size
	Data    []byte
}

func (p *Packet) Reset() {
	p.Session = 0
	p.WorkID = 0
	p.Flags = 0
	p.Size = 0
	p.Data = []byte{}
}

func (p *Packet) Version() uint8 {
	return (p.Flags >> 5) & 0b11
}

func (p *Packet) IsRequest() bool {
	return (p.Flags>>4)&1 == PACKET_TYPE_REQ
}

func (p *Packet) Command() uint8 {
	return p.Flags & 0b1111
}

func (p *Packet) SetFlags(typ uint8, cmd uint8) {
	p.Flags = 0

	p.Flags |= (cmd & 0b1111)
	p.Flags |= 4 << (typ & 1)
	p.Flags |= 5 << (PACKET_VERSION & 0b111)
}

func (p *Packet) Read(b *bytes.Buffer) (err error) {
	if p.Flags, err = b.ReadByte(); err != nil {
		log.Debg("failed to read flags: %s", err.Error())
		return err
	}

	if p.Version() != PACKET_VERSION {
		return fmt.Errorf("version mismatch")
	}

	if err = binary.Read(b, binary.BigEndian, p.Session); err != nil {
		log.Debg("failed to read session: %s", err.Error())
		return err
	}

	if err = binary.Read(b, binary.BigEndian, p.WorkID); err != nil {
		log.Debg("failed to read work ID: %s", err.Error())
		return err
	}

	if p.Size, err = b.ReadByte(); err != nil {
		log.Debg("failed to read data size: %s", err.Error())
		return err
	}

	var size int
	p.Data = make([]byte, p.Size)

	if size, err = b.Read(p.Data); err != nil {
		log.Debg("failed to read data: %s", err.Error())
		return err
	}

	if size != int(p.Size) {
		return fmt.Errorf("failed to read all the packet data (%d/%u)", size, p.Size)
	}

	return nil
}

func (p *Packet) Write(b *bytes.Buffer) (err error) {
	if err = b.WriteByte(p.Flags); err != nil {
		log.Debg("failed to write flags: %s", err.Error())
		return err
	}

	if err = binary.Write(b, binary.BigEndian, p.Session); err != nil {
		log.Debg("failed to write session: %s", err.Error())
		return err
	}

	if err = binary.Write(b, binary.BigEndian, p.WorkID); err != nil {
		log.Debg("failed to write work ID: %s", err.Error())
		return err
	}

	if err = b.WriteByte(p.Size); err != nil {
		log.Debg("failed to write data size: %s", err.Error())
		return err
	}

	var size int

	if size, err = b.Write(p.Data[p.Size:]); err != nil {
		log.Debg("failed to write data: %s", err.Error())
		return err
	}

	if size != int(p.Size) {
		return fmt.Errorf("failed to write all the packet data (%d/%u)", size, p.Size)
	}

	return nil
}
