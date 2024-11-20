package c2

const PACKET_MAX_SIZE = 2 + 255

type Packet struct {
	ID     uint32 // unique ID
	Flags  uint8  // version (3 bits), 1 (type), 4 (command)
	WorkID uint8  // id of the work this packet is assoicated with
	Size   uint8  // data size
	Data   []byte
}
