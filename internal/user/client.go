package user

import "net"

type User struct {
	Name string
	Conn net.Conn
}

func New(name string, conn net.Conn) *User {
	return &User{name, conn}
}
