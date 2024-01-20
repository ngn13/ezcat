package util

import (
	"encoding/base64"
	"math/rand"
	"net"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)


var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")

func MakeRandom(size int) string {
  ret := make([]rune, 42)
  for i := range ret {
    ret[i] = chars[rand.Intn(len(chars))]
  }

  return string(ret)
}

func RenderErr(c *fiber.Ctx, code int, err error) error{
  log.Errorf("Code: %d Error: %s", code, err)
  switch code {
  case 404:
    return c.Render("error", fiber.Map{
      "msg": "not found",
    })
  case 500:
    return c.Render("error", fiber.Map{
      "msg": "server error",
    })
  
  case 403:
    return c.Render("error", fiber.Map{
      "msg": "forbidden",
    })
  }

  return c.Render("error", fiber.Map{
    "msg": code,
  })
}

func SendEnc(c *fiber.Ctx, plain string) error {
  enc := base64.StdEncoding.EncodeToString([]byte(plain))
  return c.SendString(enc)
}

func Reverse(str string) (result string) { 
  for _, v := range str { 
    result = string(v) + result 
  } 
  
  return
}

func GetAddr(addrs []net.Addr) string {
  for _, a := range addrs {
    switch v := a.(type) {
    case *net.IPNet:
      if strings.Contains(v.IP.String(), ":") {
        continue
      }
      return v.IP.String()
    case *net.IPAddr:
      if strings.Contains(v.IP.String(), ":") {
        continue
      }
      return v.IP.String()
    }
  }

  return ""
}

func GetIP() string {
  ip := "127.0.0.1"
  foundtun := false
  ifs, err := net.Interfaces()
  if err != nil {
    return ip
  }

  for _, i := range ifs {
    if strings.HasPrefix(i.Name, "lo") || 
      strings.HasPrefix(i.Name, "docker") || 
      foundtun{
      continue
    }

    ips, err := i.Addrs()
    if err != nil || len(ips) == 0{
      continue
    }

    nip := GetAddr(ips)
    if nip == "" {
      continue
    }

    ip = nip
    if strings.HasPrefix(i.Name, "tun") {
      foundtun = true
    } 
  }

  return ip
}
