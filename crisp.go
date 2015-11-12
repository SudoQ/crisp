package main

import (
	"flag"
	"fmt"
	"github.com/SudoQ/crisp/service"
	"github.com/SudoQ/crisp/logging"
	"runtime"
)

var port string
var limit uint

func init() {
	flag.StringVar(&port, "port", "8080", "Port number of the crisp service")
	flag.StringVar(&port, "p", "8080", "Port number of the crisp service (shorthand)")

	flag.UintVar(&limit, "limit", 60, "Limit of requests per hour")
	flag.UintVar(&limit, "l", 60, "Limit of requests per hour (shorthand)")

	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	flag.Parse()

	if len(flag.Args()) != 1 {
		logging.Info(
			fmt.Sprintf("Usage: \n\t%s\n\t%s",
				"crisp [-p=port] [-l=limit] <url>",
				"crisp -f=<file>"))
		return
	}
	url := flag.Arg(0)
	srv := service.New(url, port, limit)

	srv.Run()
}
