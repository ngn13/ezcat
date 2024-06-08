package agent

// Struct and vars based on RFC 1035

// 2.3.4. Size limits
var (
  LABEL_LIMIT int = 63
  NAME_LIMIT  int = 255
  UDP_LIMIT   int = 512
)

type DNS_Packet struct {
  Header      DNS_Header
  Questions   []DNS_QD
  Answers     []DNS_RR
  Authorities []DNS_RR
  Additionals []DNS_RR
}

// 4.1.1. Header section format
type DNS_Header struct {
  ID      uint16 // Transaction ID (16 bits)
  Flags   uint16 // Flags (16 bits)
  QDCount uint16 // Question Count (16 bits)
  ANCount uint16 // Answer Count (16 bits)
  NSCount uint16 // Authority Count (16 bits)
  ARCount uint16 // Additional Info Count (16 bits)
}

// 4.1.2. Question section format
type DNS_QD struct {
  Qname   []string
  Qtype   uint16
  Qclass  uint16
}

// 4.1.3. Resource record format 
type DNS_RR struct {
  Name     []string
  Type     uint16
  Class    uint16
  TTL      uint32
  RDLength uint16
  RData    []byte
}
