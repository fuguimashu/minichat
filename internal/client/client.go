package client

import "net"

type Client struct {
	Name string
	Conn net.Conn
}

func New(name string, conn net.Conn) *Client {
	return &Client{name, conn}
}
