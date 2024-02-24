package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:2577")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	//设置信号处理函数，接收到 SIGINT 或 SIGTERM 信号时关闭连接
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		fmt.Println("Bye~")
		conn.Close()
		os.Exit(0)
	}()

	go io.Copy(conn, os.Stdin)
	io.Copy(os.Stdout, conn)
}
