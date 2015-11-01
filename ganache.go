package main

import (
	"flag"
	"github.com/SudoQ/ganache/service"
	"log"
	"runtime"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var configFilename string
var url string
var port string
var limit uint

func init() {
	flag.StringVar(&configFilename, "f", "config.json", "Path to configuration file")
	flag.StringVar(&url, "url", "http://whatthecommit.com/index.txt", "URL to cache")
	flag.StringVar(&port, "p", "8080", "Port number of the ganache service")
	flag.UintVar(&limit, "l", 60, "Limit of requests per hour")
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
		srv = service.New(url, port, limit)
	}
	log.Println(srv.Info())
	srv.Run()
}
