package core

import (
	"fmt"
	"github.com/gorilla/websocket"
	chatlog "logger"
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

func (c *Client) Read() {
	// 最后这个连接删除的时候就调用析构呗
	defer func() {
		Manager.Unregister <- c
		_ = c.Conn.Close()
		close(c.SendPipe)
	}()

	// ping pong 检测是不是连接还是存在的
	_ = c.Conn.SetReadDeadline(time.Now().Add(pongWait)) // 读的 deadline ==> 60s
	c.Conn.SetPongHandler(func(string) error {
		_ = c.Conn.SetReadDeadline(time.Now().Add(pongWait)) // 收到 pong 继续维持生命周期
		return nil
	})

	for {
		m_type, data, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				chatlog.Lg().Errorf("消息读取:网络错误 %v \n", err)
			} else {
				chatlog.Lg().Errorf("pong超时：%v ", err)
			}
			break
		}
		if m_type != websocket.BinaryMessage {
			chatlog.Lg().Errorln("消息类型读取出错")
			break
		}
		// 目前已经是正确的读取到数据了
		messagebase := &MessageBase{}

		err = FromProtoc[*MessageBase](data, messagebase)

		//err = messagebase.FromProtoc(data)
		if err != nil {
			chatlog.Lg().Errorln("消息的序列化出错")
			break
		}
		chatlog.Lg().Errorf("-----处理消息类型: %v", messagebase.Impl.MsgType)

		Router.ExecHandler(&Context{
			Client: c,
			Msg:    messagebase,
		})

		if err != nil {
			chatlog.Lg().Errorf("消息路由处理出错: %v", err)
			break
		}
	}
}

func (c *Client) Write() {

}
