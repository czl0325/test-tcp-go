package main

import "net"

type User struct {
	Name string
	Addr string
	C    chan string
	Conn net.Conn
}

func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		Conn: conn,
	}
	go user.ListenMessage()
	return user
}

func (user *User) ListenMessage() {
	for {
		msg := <-user.C
		user.Conn.Write([]byte(msg + "\n"))
	}
}
