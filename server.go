package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	serverIp   string
	serverPort int
	OnlineMap  map[string]*User
	mapLock    sync.RWMutex
	Message    chan string
}

func NewServer(serverIp string, serverPort int) *Server {
	return &Server{
		serverIp:   serverIp,
		serverPort: serverPort,
		OnlineMap:  make(map[string]*User),
		Message:    make(chan string),
	}
}

func (server *Server) ListenMessage() {
	for {
		msg := <-server.Message
		server.mapLock.Lock()
		for _, user := range server.OnlineMap {
			user.C <- msg
		}
		server.mapLock.Unlock()
	}
}

func (server *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	server.Message <- sendMsg
}

func (server *Server) Handle(conn net.Conn) {
	user := NewUser(conn, server)
	user.Online()

	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("conn读取消息错误, ", err)
				return
			}
			msg := string(buf[:n-1])
			user.DoMessage(msg)
		}
	}()
}

func (server *Server) Start() {
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.serverIp, server.serverPort))
	if err != nil {
		fmt.Println("服务器启动失败, ", err)
		return
	}
	defer func(listen net.Listener) {
		err := listen.Close()
		if err != nil {

		}
	}(listen)

	go server.ListenMessage()
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen accept error", err)
			continue
		}
		go server.Handle(conn)
	}
}
