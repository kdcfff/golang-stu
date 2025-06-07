package main

import (
	"net"
)

//  User 用户结构体
type User struct {
	Name string
	Addr string
	C    chan string
	con  net.Conn
}

// NewUser 创建并初始化一个新的用户对象。
// 该函数接收一个网络连接对象作为参数，用于后续与用户进行通信。
// 返回值是一个指向User结构的指针，表示新创建的用户实例。
func NewUser(con net.Conn) *User {
    // 获取客户端的地址信息，将其用作用户的名称和地址。
    userAddr := con.RemoteAddr().String()
    
    // 初始化一个User结构体实例。
    // 用户名称和地址都设置为客户端的地址信息。
    // 创建一个字符串类型的通道，用于向用户发送消息。
    // 将网络连接对象赋值给用户实例，以便后续进行读写操作。
    user := &User{
        Name: userAddr,
        Addr: userAddr,
        C:    make(chan string),
        con:  con,
    }
    // 启动一个goroutine来监听用户发送的消息。
	go user.ListenMessage()
    // 返回初始化后的用户实例。
    return user
}
//  ListenMessage 监听用户发送的消息。
func (this *User) ListenMessage() { 
	for {
		msg := <-this.C
		this.con.Write([]byte(msg + "\n"))
	}
}

