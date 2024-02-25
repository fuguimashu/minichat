package server

import (
	"fmt"
	"io"
	"minichat/internal/user"
	"net"
	"strings"
	"sync"
	"time"
)

type Server struct {
	Ip           string
	Port         int
	OnlineClient sync.Map
	Message      chan string
}

func New(ip string, port int) *Server {
	s := &Server{
		Ip:           ip,
		Port:         port,
		OnlineClient: sync.Map{},
		Message:      make(chan string),
	}

	go s.ListenChat()

	return s
}

func (s *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("net listen error:" + err.Error())
		return
	}
	defer listener.Close()

	fmt.Println("TCP Server Start Successful.")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept:" + err.Error())
			continue
		}

		go s.handler(conn)
	}
}

func (s *Server) handler(conn net.Conn) {
	defer conn.Close()
	// conn.Write([]byte(fmt.Sprintf("Welcome %s\n", conn.RemoteAddr().String())))

	user := user.New(conn.RemoteAddr().String(), conn)

	s.OnlineClient.Store(user.Name, user)

	buf := make([]byte, 4096)

	for {

		n, err := conn.Read(buf)
		if err != nil && err != io.EOF {
			fmt.Println("conn read error:" + err.Error())
			return
		}

		if n == 0 {
			fmt.Println(user.Name + " 离开了")
			s.OnlineClient.Delete(user.Name)
			return
		}

		// '/c' hello,world 公聊
		if string(buf)[0] == '/' && string(buf)[1:2] == "c" {
			s.Message <- fmt.Sprintf("[%s %s]%s", time.Now().Format("2006-01-02 15:04:05"), user.Name, string(buf)[2:n])
		}

		// '/who' 显示在线用户
		if string(buf)[0] == '/' && string(buf)[1:4] == "who" {
			s.OnlineClient.Range(func(key, value any) bool {
				fmt.Println(key, value)
				conn.Write([]byte(key.(string) + "\n"))
				return true
			})
		}

		// '/rename 张三' 修改用户名
		if string(buf)[0] == '/' && string(buf)[1:7] == "rename" {

			// '/rename ' 用户名不为空
			if len(strings.Split(string(buf[:n]), " ")) <= 1 {
				continue
			}

			// 删除旧的key
			s.OnlineClient.Delete(user.Name)

			// 新用户名
			user.Name = strings.Split(string(buf[:n]), " ")[1]
			user.Name = strings.Trim(user.Name, "\r\n")

			// 存储新用户
			s.OnlineClient.Store(user.Name, user)
			// conn.Write([]byte("修改用户成功!\n"))
		}
	}
}

func (s *Server) ListenChat() {

	for {
		message := <-s.Message
		s.OnlineClient.Range(func(key, value any) bool {
			u, ok := s.OnlineClient.Load(key)
			if ok {
				u.(*user.User).Ch <- message
			}
			return ok
		})
	}
}
