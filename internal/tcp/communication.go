package tcp

import (
	"net"
	"time"
)

const (
	bufferSize  = 8196
	dialingTime = 15 * time.Second
)

func receive(conn net.Conn) {}

func send(conn net.Conn) {}
