package main

import (
	"flag"
	"fmt"

	"github.com/hirokazumiyaji/memds/memds"
)

var (
	version string
)

func main() {
	var (
		port       int
		sock       string
		bucketNum  int
		configPath string
		config     *memds.Config
		vFlag      bool
		err        error
	)

	flag.IntVar(&port, "port", 6700, "listen port")
	flag.IntVar(&port, "p", 6700, "listen port")
	flag.StringVar(&sock, "sock", "", "socket")
	flag.StringVar(&sock, "s", "", "socket")
	flag.IntVar(&bucketNum, "bucket_num", 10, "bucket num")
	flag.IntVar(&bucketNum, "bn", 10, "bucket num")
	flag.StringVar(&configPath, "config", "", "config path")
	flag.StringVar(&configPath, "c", "", "config path")
	flag.BoolVar(&vFlag, "version", false, "version")
	flag.BoolVar(&vFlag, "v", false, "version")

	flag.Parse()

	if vFlag {
		fmt.Printf("memds version: %s\n", version)
		return
	}

	if configPath == "" {
		config = new(memds.Config)
		config.Port = port
		config.Sock = sock
		config.BucketNum = bucketNum
	} else {
		config, err = memds.LoadConfig(configPath)
		if err != nil {
			memds.Error(err.Error())
			return
		}
	}

	memds.Info("start memds")
	if err := memds.Serve(config); err != nil {
		memds.Error(err.Error())
	}
}
