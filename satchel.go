package main

import (
	"flag"
	"fmt"
	"github.com/SudoQ/satchel/service"
	"runtime"
)

var port string
var limit uint

func init() {
	flag.StringVar(&port, "port", "8080", "Port number of the satchel service")
	flag.StringVar(&port, "p", "8080", "Port number of the satchel service (shorthand)")

	flag.UintVar(&limit, "limit", 60, "Limit of requests per hour")
	flag.UintVar(&limit, "l", 60, "Limit of requests per hour (shorthand)")

	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Usage = func() {
		fmt.Printf("Usage: satchel [options] URL\n")
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	if len(flag.Args()) != 1 {
		flag.Usage()
		return
	}
	url := flag.Arg(0)
	srv := service.New(url, port, limit)

	srv.Run()
}
