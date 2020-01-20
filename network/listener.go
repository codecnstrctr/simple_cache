package network

import (
	"net"
	"net/rpc/jsonrpc"
	"strconv"
)

const (
	defaultProto = "tcp"
)

// Listener represents a tcp network listener
type Listener struct {
	proto string
	addr  string
	done  chan struct{}
}

// NewListener creates new listener instance
func NewListener(port int) *Listener {
	l := new(Listener)
	l.proto = defaultProto
	l.addr = ":" + strconv.Itoa(port)

	return l
}

// Stop stops l
func (l *Listener) Stop() {
	l.done <- struct{}{}
	close(l.done)
}

// Listen is waiting for new connections
// to be established and handles them
func (l *Listener) Listen() error {
	tcpaddr, err := net.ResolveTCPAddr(l.proto, l.addr)
	if err != nil {
		return err
	}

	ln, err := net.ListenTCP(l.proto, tcpaddr)
	if err != nil {
		return err
	}
	defer ln.Close()

listen_loop:
	for {
		select {
		case <-l.done:
			break listen_loop
		default:
			conn, err := ln.AcceptTCP()
			if err != nil {
				continue listen_loop
			}

			go jsonrpc.ServeConn(conn)
		}
	}

	return nil
}
