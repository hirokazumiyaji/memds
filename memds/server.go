package memds

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func Serve(c *Config) error {
	var (
		l   net.Listener
		err error
	)

	buckets, err = NewBuckets(c.BucketNum)
	if err != nil {
		return err
	}

	l, err = listener(c)
	if err != nil {
		return err
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(
		sig,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	var wg sync.WaitGroup
	closed := false
	go func() {
		<-sig
		cancel()
		closed = true
		l.Close()
	}()

	wg.Add(1)
	go gc(ctx, &wg, c)

	for {
		conn, err := l.Accept()
		if err != nil && closed {
			break
		}
		if err != nil {
			Error(err.Error())
			continue
		}
		wg.Add(1)
		go accept(ctx, &wg, conn)
	}

	if !closed {
		l.Close()
	}

	wg.Wait()
	return nil
}

func listener(c *Config) (net.Listener, error) {
	if c.Sock == "" {
		return net.Listen("tcp", fmt.Sprintf(":%d", c.Port))
	}
	return net.Listen("unix", c.Sock)
}

func accept(ctx context.Context, wg *sync.WaitGroup, c net.Conn) {
	closed := false

	go func() {
		<-ctx.Done()
		closed = true
		c.Close()
	}()

	defer func() {
		if !closed {
			c.Close()
		}
		wg.Done()
	}()

	r := bufio.NewReader(c)

	for {
		line, _, err := r.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil && closed {
			break
		}
		if err != nil {
			Error(fmt.Sprintf("%v", err))
			break
		}

		res := Exec(line)

		_, err = c.Write(res)
		if err != nil {
			Error(fmt.Sprintf("%v", err))
			break
		}
	}
}
