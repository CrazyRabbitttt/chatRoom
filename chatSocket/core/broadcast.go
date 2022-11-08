package core

// 广播的消息接口， 算是cpp中的抽象基类？
type IBroadcast interface {
	MsgType() int32  // 消息的类型
	MsgRange() int32 // 消息的范围
	Sender() int64   // 消息发送者
	Receiver() int64 // 消息接受者
	Message() []byte // 消息序列化之后的东西
}

type Broadcast struct {
	Client  *Client // 源用户
	Msg     IBroadcast
	clients *[]*Client // 用于广播的目标是啥
}
