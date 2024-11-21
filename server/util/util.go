package util

import (
	"encoding/base64"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")

func MakeRandom(size int) string {
	ret := make([]rune, size)
	for i := range ret {
		ret[i] = chars[rand.Intn(len(chars))]
	}

	return string(ret)
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
			foundtun {
			continue
		}

		ips, err := i.Addrs()
		if err != nil || len(ips) == 0 {
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
	if c >= 'A' && c <= 'Z' {
		return 'A' + ((c - 'A' + 13) % 26)
	} else if c >= 'a' && c <= 'z' {
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

	if !st.IsDir() {
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

func BoolToStr(r bool) string {
	if r {
		return "true"
	}
	return "false"
}

func RunBuild(dir string, oscode string, options map[string]string) (string, error) {
	var (
		output []byte
		err    error
	)

	cmd := exec.Command("/bin/bash", "build.sh", oscode)
	cmd.Dir = dir

	if nil != options {
		for k, v := range options {
			opt := fmt.Sprintf("%s=%s", k, v)
			cmd.Args = append(cmd.Args, opt)
		}
	}

	if output, err = cmd.CombinedOutput(); err != nil {
		return string(output), err
	}

	return string(output), nil
}

func IsValidPort(port int) bool {
	return port > 0 && port <= math.MaxUint16
}

func ParseAddr(addr string) (string, uint16, error) {
	cols := strings.Split(addr, ":")

	if len(cols) != 2 {
		return "", 0, fmt.Errorf("invalid address format")
	}

	port, err := strconv.Atoi(cols[1])

	if err != nil || !IsValidPort(port) {
		return "", 0, fmt.Errorf("invalid port number")
	}

	return cols[0], uint16(port), nil
}

func ContainsString(s string, l []string) bool {
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

func ToUint16(s string) (uint16, error) {
	id, err := strconv.ParseUint(s, 10, 32)
	return uint16(id), err
}

func Rand16() uint16 {
	return uint16(rand.Intn(math.MaxUint16))
}
