package bridge

import "os"

func GetAddr() string {
  addr := os.Getenv("HTTP_ADDR")
  if addr == ""{
    return "0.0.0.0:6001"
  }
  return addr 
}
