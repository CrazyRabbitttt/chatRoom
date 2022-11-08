package core

import (
	"fmt"
	"github.com/gorilla/websocket"
	"time"
)

const (
	writeWait = 10 * time.Second
	pongWait  = 60 * time.Second // ping pong 用于维护心跳
	pingWait  = (pongWait * 9) / 10
)

func Hello() {
	fmt.Println("Hello from package core")
}

type Client struct {
	Id       int64
	Conn     *websocket.Conn
	SendPipe chan []byte
}

func NewClient(id int64, conn *websocket.Conn) *Client {
	return &Client{Id: id, Conn: conn, SendPipe: make(chan []byte, 512)}
}
