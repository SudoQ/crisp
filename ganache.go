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

var configFile string

func init() {
	flag.StringVar(&configFile, "f", "config.json", "Path to configuration file")
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	flag.Parse()

	srv, _ := service.NewFromFile(configFile)
	log.Println(srv.Info())
	config, err := srv.JSON()
	log.Println(string(config))
	check(err)
	err = srv.SaveConfig("config.json")
	check(err)
	err = srv.LoadConfig("config.json")
	check(err)
	srv.Run()
}
