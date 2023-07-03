package main

import (
	"net"
	"strings"
)

type User struct {
	Name   string
	Addr   string
	C      chan string
	Conn   net.Conn
	server *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		Conn:   conn,
		server: server,
	}
	go user.ListenMessage()
	return user
}

func (user *User) Online() {
	user.server.mapLock.Lock()
	user.server.OnlineMap[user.Name] = user
	user.server.mapLock.Unlock()
	user.server.BroadCast(user, "已上线")
}

func (user *User) Offline() {
	user.server.mapLock.Lock()
	delete(user.server.OnlineMap, user.Name)
	user.server.mapLock.Unlock()
	user.server.BroadCast(user, "已下线")
}

func (user *User) DoMessage(msg string) {
	if msg == "who" {
		user.server.mapLock.Lock()
		for _, user := range user.server.OnlineMap {
			user.server.BroadCast(user, "["+user.Addr+"]：在线\n")
		}
		user.server.mapLock.Unlock()
	} else if strings.HasPrefix(msg, "rename|") {
		msg = strings.Replace(msg, "rename|", "", 0)
		user.server.mapLock.Lock()
		delete(user.server.OnlineMap, user.Name)
		user.Name = msg
		user.server.OnlineMap[msg] = user
		user.server.mapLock.Unlock()
	} else {
		user.server.BroadCast(user, msg)
	}
}

func (user *User) ListenMessage() {
	for {
		msg := <-user.C
		user.Conn.Write([]byte(msg + "\n"))
	}
}
