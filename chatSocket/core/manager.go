package core

type ClientManager struct {
	ClientMapper map[int64]*Client // uid ==> Client
	Broadcast    chan *Broadcast
	Register     chan *Client
	Unregister   chan *Client
}

var Manager = ClientManager{
	ClientMapper: make(map[int64]*Client),
	Broadcast:    make(chan *Broadcast),
	Register:     make(chan *Client),
	Unregister:   make(chan *Client),
}

// 通过 uid 获得Client
func (m *ClientManager) GetClient(userId int64) *Client {
	client := m.ClientMapper[userId]
	return client // return local object
}

func (m *ClientManager) Do() {
	for {
		select {
		case conn := <-Manager.Register: // 注册连接的行为
			m.register(conn)
		case conn := <-Manager.Unregister: // 断开连接
			m.unregister(conn)
			//case broadcast := <-Manager.Broadcast: // 广播消息
			// TODO; 广播的操作
		}
	}
}

// 传入 *Client 进行注册
func (m *ClientManager) register(conn *Client) {
	// TODO：Log
	Manager.ClientMapper[conn.Id] = conn
	// TODO: 发送消息给到Client， 提示说连接服务成功！
}

func (m *ClientManager) unregister(conn *Client) {
	// TODO: log, 连接断开的log
	if _, ok := Manager.ClientMapper[conn.Id]; ok {
		// 从map中删除key对应的entry， 传入的是key
		delete(Manager.ClientMapper, conn.Id)
	}
}
