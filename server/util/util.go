package util

import (
	"encoding/base64"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/ngn13/ezcat/server/log"
)


var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")

func MakeRandom(size int) string {
  ret := make([]rune, size)
  for i := range ret {
    ret[i] = chars[rand.Intn(len(chars))]
  }

  return string(ret)
}

func Error(c *fiber.Ctx, err string) error {
  return c.JSON(fiber.Map{
    "error": err,
  })
}

func ErrorCode(c *fiber.Ctx, code int) error {
  switch code {
  case 400:
    return c.Status(code).JSON(fiber.Map{
      "error": "Bad request",
    })

  case 401:
    return c.Status(code).JSON(fiber.Map{
      "error": "You are not logged in",
    })

  default:
    return c.Status(code).JSON(fiber.Map{
      "error": "Internal error",
    })
  }
}

func ErrorInternal(c *fiber.Ctx, err string) error {
  log.Err("Error: %s", err)
  return ErrorCode(c, 500)
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

func GetToken(c *fiber.Ctx) string {
  return c.Get("Authorization")
}

func CORS(c *fiber.Ctx) error {
  c.Set("Access-Control-Allow-Origin", "*")
  c.Set("Access-Control-Allow-Headers", 
  "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
  c.Set("Access-Control-Allow-Methods", "OPTIONS, PUT, DELETE, GET")

  if c.Method() == "OPTIONS" {
    return c.SendString("")
  }

  return c.Next()
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
  enip := os.Getenv("SHELLIP")
  if enip != "" {
    return enip
  }

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

func Access(p string) bool {
  _, err := os.Stat(p)
  if err != nil {
    return false
  }
  return true
}

func Rot13Byte(c byte) byte {
  if (c >= 'A' && c <= 'Z'){
    return 'A' + ((c - 'A' + 13) % 26)
  }else if (c >= 'a' && c <= 'z'){
    return 'a' + ((c - 'a' + 13) % 26)
  }
  return c
}

func Rot13(s string) string {
  var result []byte
  bytes := []byte(s)

  for _, b := range bytes {
    result = append(result, Rot13Byte(b))
  }

  return string(result)
}

func CopyFile(dst string, src string) error {
  srcf, err := os.Open(src)
  if err != nil {
    return err
  }
  defer srcf.Close()

  dstf, err := os.Create(dst)
  if err != nil {
    return err
  }
  defer dstf.Close()

  _, err = io.Copy(dstf, srcf)
  if err != nil {
    return err
  }

  err = dstf.Sync()
  if err != nil {
    return err
  }

  return nil
}

func CopySimple(dst string, src string) error {
  st, err := os.Stat(src)
  if err != nil {
    return err
  } 

  if !st.IsDir(){
    return CopyFile(dst, src)
  }

  err = os.Mkdir(dst, os.ModePerm)
  if err != nil {
    return err
  }

  entries, err := os.ReadDir(src)
  if err != nil {
    return err
  }

  for _, e := range entries {
    src_new := path.Join(src, e.Name())
    dst_new := path.Join(dst, e.Name())

    err = CopySimple(dst_new, src_new)
    if err != nil {
      return err
    }
  }

  return nil
} 

func BoolToStr(r bool) string{
  if r {
    return "true"
  }
  return "false"
}

func RunBuild(dir string, oscode string) (string, error) {
  cmd := exec.Command("/bin/bash", "build.sh", oscode)
  cmd.Dir = dir

  output, err := cmd.CombinedOutput()

  if err != nil {
    return string(output), err
  }

  return string(output), nil
}

func IsValidPort(port int) bool{
  return port > 0 && port <= 65535
}

func ParseAddr(addr string) (string, int, error) {
  cols := strings.Split(addr, ":")
  if len(cols) != 2{
    return "", -1, fmt.Errorf("invalid address format")
  }

  port, err := strconv.Atoi(cols[1])
  if err != nil || !IsValidPort(port) {
    return "", -1, fmt.Errorf("invalid port number")
  }

  return cols[0], port, nil
}

func ContainsString(s string, l []string) bool{
  for _, e := range l {
    if s == e {
      return true
    }
  }
  return false
}

func CleanTemp() {
  entries, err := os.ReadDir("/tmp")
  if err != nil {
    return
  }

  for _, e := range entries {
    name := e.Name()
    if !strings.HasPrefix(name, "ezcat_") {
      continue
    }
    os.RemoveAll(path.Join("/tmp", name))
  }
}
