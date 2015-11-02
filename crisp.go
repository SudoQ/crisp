package main

import (
	"flag"
	"fmt"
	"github.com/SudoQ/crisp/service"
	"log"
	"runtime"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var configFilename string
var port string
var limit uint

func init() {
	flag.StringVar(&configFilename, "file", "config.json", "Path to configuration file")
	flag.StringVar(&configFilename, "f", "config.json", "Path to configuration file (shorthand)")

	flag.StringVar(&port, "port", "8080", "Port number of the crisp service")
	flag.StringVar(&port, "p", "8080", "Port number of the crisp service (shorthand)")

	flag.UintVar(&limit, "limit", 60, "Limit of requests per hour")
	flag.UintVar(&limit, "l", 60, "Limit of requests per hour (shorthand)")

	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	flag.Parse()

	useConfig := false
	flag.Visit(func(flg *flag.Flag) {
		if flg.Name == "f" {
			useConfig = true
		}
	})

	var srv *service.Service
	if useConfig {
		var err error
		srv, err = service.NewFromFile(configFilename)
		if err != nil {
			log.Fatal(err)
			return
		}
	} else {
		var url string
		if len(flag.Args()) != 1 {
			log.Println(
				fmt.Sprintf("Usage: \n\t%s\n\t%s",
					"crisp [-p=port] [-l=limit] <url>",
					"crisp -f=<file>"))
			return
		}
		url = flag.Arg(0)
		srv = service.New(url, port, limit)
	}

	log.Println(srv.Info())
	srv.Run()
}
