package main

import (
	"flag"
	"time"

	"github.com/golang/glog"
)

const (
	RESTART_LENGTH  = 50000 * time.Second
	SHUTDOWN_LENGTH = 1 * time.Second
)

func main() {
	flag.Parse()
	for {
		s := Server{}
		s.Initialize()
		s.Run()
		time.Sleep(RESTART_LENGTH)
		glog.Infof("Shutting down server...")
		if err := s.Shutdown(); err != nil {
			glog.Infof("Server shutdown failed! Sending panic: %s", err)
			panic(err)
		}
		time.Sleep(SHUTDOWN_LENGTH)
		glog.Info("Restarting...")
	}
}
