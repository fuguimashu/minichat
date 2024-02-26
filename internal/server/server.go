package server

import (
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"time"

	"minichat/internal/user"
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

	currentUser := user.New(conn.RemoteAddr().String(), conn)

	s.OnlineClient.Store(currentUser.Name, currentUser)

	buf := make([]byte, 4096)

	for {

		n, err := conn.Read(buf)
		if err != nil && err != io.EOF {
			fmt.Println("conn read error:" + err.Error())
			return
		}

		if n == 0 {
			fmt.Println(currentUser.Name + " 离开了")
			s.OnlineClient.Delete(currentUser.Name)
			return
		}

		// '/c' hello,world 公聊
		if string(buf)[0] == '/' && string(buf)[1:2] == "c" {
			s.Message <- fmt.Sprintf("[%s %s]:%s", time.Now().Format("2006-01-02 15:04:05"), currentUser.Name, string(buf[2:n]))
		}

		// '/m username' how are your? 私聊
		if string(buf)[0] == '/' && string(buf)[1:2] == "m" {
			peer := strings.Split(string(buf[:n]), " ")[1]
			peer = strings.Trim(peer, "\r\n")
			fmt.Printf("peer> %q\n", peer)
			if u, ok := s.OnlineClient.Load(peer); ok {
				msg := fmt.Sprintf("[%s %s]:%s", time.Now().Format("2006-01-02 15:04:05"), currentUser.Name, string(buf[4+len(peer):n]))
				u.(*user.User).Conn.Write([]byte(msg))
			} else {
				msg := peer + " 找不到该用户"
				currentUser.Conn.Write([]byte(msg + "\n"))
			}

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
			s.OnlineClient.Delete(currentUser.Name)

			// 新用户名
			currentUser.Name = strings.Split(string(buf[:n]), " ")[1]
			currentUser.Name = strings.Trim(currentUser.Name, "\r\n")

			// 存储新用户
			s.OnlineClient.Store(currentUser.Name, currentUser)
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
