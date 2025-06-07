package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int
	//  在线用户列表
	OnlineMap map[string]*User
	mapLock   sync.Mutex
	//  广播消息
	msg chan string
}

// 创建一个服务器
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		mapLock:   sync.Mutex{},
		//添加一个广播频道
		msg:       make(chan string),
	}
	return server
}

//  处理用户连接
func (this *Server) handle(conn net.Conn) {
	//fmt.Println("链接建立成功")
	user := NewUser(conn)
	//用户上线加入map
	this.mapLock.Lock()
	//添加到map中
	this.OnlineMap[user.Name] = user
	//解开锁
	this.mapLock.Unlock()

	//向所有用户广播
	this.BroadCast(user, "已上线")
	//	阻塞
	select {}

}
//  广播
func (this *Server) BroadCast(user *User, msg string) { 
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	this.msg <- sendMsg
}

//监听msg 
func (this *Server) ListenMsg() { 
	for {
		msg := <-this.msg
		//  遍历所有用户
		this.mapLock.Lock()
		for _, user := range this.OnlineMap {
			user.C <- msg
		}
		this.mapLock.Unlock()
	}
}
// 启动服务器
func (this *Server) Start() {
	//socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	defer listener.Close()

	//启动监听
	go this.ListenMsg()
	for {
		//accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept err:", err)
			return
		}
		//handle
		go this.handle(conn)
		//close
	}

}
