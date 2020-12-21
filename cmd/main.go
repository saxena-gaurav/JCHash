package main

import (
	"flag"
	"jchash/hash/service"
	"os"
)

var c service.Config

func main() {
	{
		flag.IntVar(&c.ListenPort, "p", 8080, "Port to listen on")
		flag.IntVar(&c.HashDelay, "d", 5, "Amount to delay hash response in seconds")
	}
	flag.Parse()
	if err := service.NewApplication(&c).Start(); err != nil {
		os.Exit(1)
	}
}
