package user

import (
	"net"
)

type User struct {
	Name string
	Conn net.Conn
	Ch   chan string
}

func (u *User) ListenMessage() {
	for {
		msg := <-u.Ch
		u.Conn.Write([]byte(msg))
	}
}

func New(name string, conn net.Conn) *User {
	u := &User{name, conn, make(chan string)}
	go u.ListenMessage()
	return u
}
