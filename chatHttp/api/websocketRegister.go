package api

import (
	"chatSocket/core"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"goHttp/pkg/utils"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool { // 不检查源站点啥的
		return true
	},
}

func WebSocketRegister(c *gin.Context) {
	// 获得 userId
	userId := utils.GetId(c)
	if userId == -1 {
		return
	}
	// websocket 协议：（就是将http升级为能够主动发信息的协议）
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	// Client 对应的构造需要 [userId, Conn]
	client := core.NewClient(userId, conn) // websocket conn handler
	core.Manager.Register <- client        // 放入到全局的 client 队列中， 那么 manager 中的管理操作就会被唤醒处理 register

}















