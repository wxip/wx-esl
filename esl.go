package wxesl

import (
	"github.com/fiorix/go-eventsocket/eventsocket"
)

// WxEsl 主要功能类,用于与fs交互
type WxEsl struct {
	c *WxEslConnection
}

// NewWxEsl 创建一个新的连接实例
func NewWxEsl(host string, port int, password string, eventList []string) *WxEsl {
	return &WxEsl{
		c: NewWxEslConnection(host, port, password, eventList),
	}
}

// Disconnect 断开esl连接
func (a *WxEsl) Disconnect() {
	if nil != a.c {
		a.c.Disconnect()
	}
}

// GetEventChan 获取任务队列,阻塞队列,请异步处理
func (a *WxEsl) GetEventChan() chan *eventsocket.Event {
	if nil != a.c {
		return a.c.GetEventChan()
	}
	return nil
}

// SendAsync 发送异步消息
func (a *WxEsl) SendAsync(cmd string) error {
	if nil != a.c {
		return a.c.SendAsync(cmd)
	}
	return ErrDisconnect
}

// SendMsg 发送msg
func (a *WxEsl) SendMsg(uuid, cmd, name, arg string) error {
	if nil != a.c {
		return a.c.SendMsg(uuid, cmd, name, arg)
	}
	return ErrDisconnect
}

// SendMsg 发送 execute msg
func (a *WxEsl) SendExecuteMsg(uuid, name, arg string) error {
	if nil != a.c {
		return a.c.SendMsg(uuid, "execute", name, arg)
	}
	return ErrDisconnect
}
