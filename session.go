package main

import (
	"fmt"
	"net"
	"strings"
)

type Session struct {
	conn      net.Conn
	id        int64
	uID       byte
	ipAddress string
}

func (s Session) WriteLong(i int64) {
	s.WriteByte_(byte(i >> 56))
	s.WriteByte_(byte(i >> 48))
	s.WriteByte_(byte(i >> 40))
	s.WriteByte_(byte(i >> 32))
	s.WriteByte_(byte(i >> 24))
	s.WriteByte_(byte(i >> 16))
	s.WriteByte_(byte(i >> 8))
	s.WriteByte_(byte(i & 0xFF))
}

func (s Session) WriteByte_(b byte) {
	_, err := s.conn.Write([]byte{b})
	if err != nil {
		fmt.Println("Error writing byte(s) to client session.")
		fmt.Println(err.Error())
	}
}

func (s Session) String() string {
	return s.ipAddress
}

func NewSession(conn net.Conn) Session {
	s := Session{conn, -1, 0xFF, strings.Split(conn.RemoteAddr().String(), ":")[0]}
	fmt.Printf("Accepting connection from: %s\n", s)
	return s
}
