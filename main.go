package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"./cache"
	"./network"
)

var (
	port                int
	removeExpiredPeriod int
	dryRun              bool
)

func main() {
	loadFlags()
	fmt.Printf("start listening on port %d\n", port)
	c := cache.NewCache(removeExpiredPeriod)
	defer c.Stop()

	cmd := network.NewCommand(c)
	l := network.NewListener(port, cmd)

	l.Listen()
	defer l.Stop()

	if dryRun {
		fmt.Println("successful dry run")
		return
	}

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	<-sigterm
	fmt.Println("graceful shutdown")
}

func loadFlags() {
	flag.IntVar(&port, "p", 2886, "port to listen")
	flag.IntVar(&removeExpiredPeriod, "r", 60, "period to remove expired keys in seconds")
	flag.BoolVar(&dryRun, "t", false, "dry run")
	flag.Parse()
}
