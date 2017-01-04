package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/hirokazumiyaji/memds/log"
	"github.com/hirokazumiyaji/memds/memds"
	"github.com/ugorji/go/codec"
)

var (
	version string
	mh      codec.MsgpackHandle
)

func init() {
	mh.MapType = reflect.TypeOf(map[string]interface{}(nil))
}

func main() {
	var (
		host     string
		port     int
		sock     string
		logLevel string
		vFlag    bool
	)

	flag.StringVar(&host, "host", "localhost", "host")
	flag.IntVar(&port, "port", 6700, "port")
	flag.IntVar(&port, "p", 6700, "port")
	flag.StringVar(&sock, "socket", "", "socket")
	flag.StringVar(&sock, "s", "", "socket")
	flag.StringVar(&logLevel, "log", "info", "log level")
	flag.StringVar(&logLevel, "l", "info", "log level")
	flag.BoolVar(&vFlag, "version", false, "version")

	flag.Parse()

	if vFlag {
		fmt.Printf("memdb-cli version: %s\n", version)
		return
	}

	log.SetLevel(logLevel)

	var (
		conn net.Conn
		err  error
	)
	if sock == "" {
		conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	} else {
		conn, err = net.Dial("unix", sock)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	defer conn.Close()

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("memds> ")
		l, _, err := reader.ReadLine()
		if err != nil {
			fmt.Println(err)
			continue
		}
		cmd := string(l)
		if cmd == "QUIT" || cmd == "quit" || cmd == "EXIT" || cmd == "exit" {
			break
		}

		switch {
		case strings.HasPrefix(cmd, "get") || strings.HasPrefix(cmd, "GET"):
			tokens := strings.Split(cmd, " ")
			if len(tokens) < 2 {
				fmt.Println("wrong number of arguments for 'get' command")
				continue
			}
			if tokens[0] != "get" && tokens[0] != "GET" {
				fmt.Printf("Unknown command '%s'\n", tokens[0])
				continue
			}
			var b []byte
			enc := codec.NewEncoderBytes(&b, &mh)
			err := enc.Encode(
				map[string]interface{}{
					"cmd": tokens[0],
					"key": tokens[1],
				},
			)
			if err != nil {
				fmt.Println(err)
				continue
			}

			b = append(b, '\n')
			_, err = conn.Write(b)
			if err != nil {
				fmt.Println(err)
				continue
			}

			r, _, err := bufio.NewReader(conn).ReadLine()
			if err != nil {
				fmt.Println(err)
			}

			res := make(map[string]interface{})
			dec := codec.NewDecoderBytes(r, &mh)
			err = dec.Decode(&res)
			if err != nil {
				fmt.Println(err)
				continue
			}

			s, ok := res["status"]
			if !ok {
				fmt.Println("response format error")
				continue
			}
			sb, ok := s.(bool)
			if !ok {
				fmt.Println("response format error")
				continue
			}

			if sb {
				switch v := res["value"].(type) {
				case string:
					fmt.Printf("%v\n", v)
				case []uint8:
					fmt.Printf("%v\n", memds.Uint8ArrayToString(v))
				default:
					fmt.Printf("%v\n", v)
				}
			} else {
				switch v := res["msg"].(type) {
				case string:
					fmt.Printf("%v\n", v)
				case []uint8:
					fmt.Printf("%v\n", memds.Uint8ArrayToString(v))
				default:
					fmt.Printf("%v\n", v)
				}
			}
		case strings.HasPrefix(cmd, "set") || strings.HasPrefix(cmd, "SET"):
			tokens := strings.Split(cmd, " ")
			if len(tokens) < 3 {
				fmt.Println("wrong number of arguments for 'set' command")
				continue
			}
			if tokens[0] != "set" && tokens[0] != "SET" {
				fmt.Printf("Unknown command '%s'\n", tokens[0])
				continue
			}
			var b []byte
			enc := codec.NewEncoderBytes(&b, &mh)

			sv := map[string]interface{}{
				"cmd":   tokens[0],
				"key":   tokens[1],
				"value": tokens[2],
			}

			if len(tokens) >= 4 {
				v, err := strconv.ParseInt(tokens[3], 10, 64)
				if err != nil {
					fmt.Printf("command '%s' arguments format error: %v\n", tokens[0], err)
					continue
				}
				sv["expire"] = v
			}

			err := enc.Encode(sv)
			if err != nil {
				fmt.Printf("send data encode error: %v\n", err)
				continue
			}
			b = append(b, '\n')
			_, err = conn.Write(b)
			if err != nil {
				fmt.Printf("send error: %v\n", err)
				continue
			}

			r, _, err := bufio.NewReader(conn).ReadLine()
			if err != nil {
				fmt.Printf("response error: %v\n", err)
				continue
			}

			res := make(map[string]interface{})
			dec := codec.NewDecoderBytes(r, &mh)
			err = dec.Decode(&res)
			if err != nil {
				fmt.Printf("response decode error: %v\n", err)
				continue
			}

			_, ok := res["status"]
			if !ok {
				fmt.Println("response format error")
				continue
			}

			switch v := res["msg"].(type) {
			case string:
				fmt.Printf("%v\n", v)
			case []uint8:
				fmt.Printf("%v\n", memds.Uint8ArrayToString(v))
			default:
				fmt.Printf("%v\n", v)
			}
		case strings.HasPrefix(cmd, "del") || strings.HasPrefix(cmd, "DEL"):
			tokens := strings.Split(cmd, " ")
			if len(tokens) < 2 {
				fmt.Println("wrong number of arguments for 'del' command")
				continue
			}
			if tokens[0] != "del" && tokens[0] != "DEL" {
				fmt.Printf("Unknown command '%s'\n", tokens[0])
				continue
			}
			var b []byte
			enc := codec.NewEncoderBytes(&b, &mh)
			err := enc.Encode(
				map[string]interface{}{
					"cmd": tokens[0],
					"key": tokens[1],
				},
			)
			if err != nil {
				fmt.Printf("send data encode error: %v\n", err)
				continue
			}

			b = append(b, '\n')
			_, err = conn.Write(b)
			if err != nil {
				fmt.Printf("send error: %v\n", err)
				continue
			}

			r, _, err := bufio.NewReader(conn).ReadLine()
			if err != nil {
				fmt.Printf("response error: %v", err)
				continue
			}

			res := make(map[string]interface{})
			dec := codec.NewDecoderBytes(r, &mh)
			err = dec.Decode(&res)
			if err != nil {
				fmt.Printf("response decode error: %v\n", err)
				continue
			}

			s, ok := res["status"]
			if !ok {
				fmt.Println("response format error")
				continue
			}
			_, ok = s.(bool)
			if !ok {
				fmt.Println("response format error")
				continue
			}

			switch v := res["msg"].(type) {
			case string:
				fmt.Printf("%v\n", v)
			case []uint8:
				fmt.Printf("%v\n", memds.Uint8ArrayToString(v))
			default:
				fmt.Printf("%v\n", v)
			}
		default:
			tokens := strings.Split(cmd, " ")
			fmt.Printf("Unknown command '%s'\n", tokens[0])
		}
	}
}
