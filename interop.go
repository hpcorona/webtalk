package main;

import (
	"os"
	"fmt"
	"net"
)

type ConnHandler func(*net.TCPConn)

func Start(address string, ch ConnHandler) (l *net.TCPListener, err os.Error) {
	var tcpAddr *net.TCPAddr
	tcpAddr, err = net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return
	}
	
	l, err = net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return
	}
	
	go accept(l, ch)
	
	return
}

func accept(listener *net.TCPListener, ch ConnHandler) {
	var c *net.TCPConn
	var err os.Error
	for {
		c, err = listener.AcceptTCP()
		if (err != nil) {
			fmt.Printf("[error] Accepting TCP connection: " + err.String())
			continue
		}
		
		go ch(c)
	}
}


