package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"reflect"
	"strings"

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
		host  string
		port  int
		sock  string
		vFlag bool
	)

	flag.StringVar(&host, "host", "localhost", "host")
	flag.IntVar(&port, "port", 6700, "port")
	flag.IntVar(&port, "p", 6700, "port")
	flag.StringVar(&sock, "socket", "", "socket")
	flag.StringVar(&sock, "s", "", "socket")
	flag.BoolVar(&vFlag, "version", false, "version")

	flag.Parse()

	if vFlag {
		fmt.Printf("memdb-cli version: %s\n", version)
		return
	}

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
			err := enc.Encode(
				map[string]interface{}{
					"cmd":   tokens[0],
					"key":   tokens[1],
					"value": tokens[2],
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
