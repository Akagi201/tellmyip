package main

import (
	"net"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/jessevdk/go-flags"
)

var opts struct {
	ListenAddr string `long:"listen" default:"0.0.0.0:8327" description:"TCP listen address and port"`
}

var log *logrus.Logger

func init() {
	log = logrus.New()
	log.Level = logrus.InfoLevel
	f := new(logrus.TextFormatter)
	f.TimestampFormat = "2006-01-02 15:04:05"
	f.FullTimestamp = true
	log.Formatter = f
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	addr := conn.RemoteAddr().String()
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		log.Errorf("TCP get remote addr failed, err: %v", err)
		return
	}
	conn.Write([]byte(host))
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		if !strings.Contains(err.Error(), "Usage") {
			log.Fatalf("cli params parse failed, err: %v", err)
		} else {
			return
		}
	}

	l, err := net.Listen("tcp", opts.ListenAddr)
	if err != nil {
		log.Fatalf("TCP listen failed, err: %v", err)
	}
	log.Infof("TCP listening at: %v", opts.ListenAddr)

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Errorf("TCP accept failed, err: %v", err)
			continue
		}
		go handleConn(conn)
	}
}
