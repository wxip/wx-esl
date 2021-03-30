package wxesl

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/fiorix/go-eventsocket/eventsocket"
)

// ErrDisconnect esl 断线异常
var ErrDisconnect = errors.New("esl 断线")

// WxEslConnection connection
type WxEslConnection struct {
	// esl 连接信息
	host     string
	port     int
	password string

	// 真实的连接
	c *eventsocket.Connection
	// 自动连接
	autoConnect bool
	// 事件队列
	eventChan chan *eventsocket.Event
	// 监听的事件列表
	eventList []string
}

// 执行连接
func (a *WxEslConnection) connect() error {
	address := fmt.Sprintf("%s:%d", a.host, a.port)
	connection, err := eventsocket.Dial(address, a.password)
	if nil != err {
		return err
	}
	a.c = connection

	// 启动监控
	a.monitorConnect()
	// 发送监听指令
	a.sendMonitorEventCommand()
	// 设置事件监听
	a.startEventHandle()

	return nil
}

// 监控连接并重连
func (a *WxEslConnection) monitorConnect() {
	go func(instance *WxEslConnection) {
		for {
			// 连接已停止,终止监控任务
			if !instance.autoConnect {
				break
			}

			if nil == instance.c {
				// 初次连接失败,执行连接
				err := instance.connect()
				if nil != err {
					// 首次连接失败,等待1s
					time.Sleep(time.Second)
					continue
				} else {
					break
				}
			}

			breakOut := false

			// 检测连接是否断开
			for instance.autoConnect {

				resChan := make(chan error)
				go func(chan error) {
					_, connectErr := a.c.Send("bgapi status")
					resChan <- connectErr
				}(resChan)

				var connectErr error
				select {
				case <-time.After(time.Duration(3) * time.Second):
					connectErr = errors.New("timeout")
				case connectErr = <-resChan:
				}

				if nil != connectErr {
					// 连接异常
					instance.c = nil
					err := instance.connect()
					if nil != err {
						// 连接失败,等待1s
						time.Sleep(time.Second)
					} else {
						breakOut = true
					}
					break
				} else {
					time.Sleep(time.Second)
				}
			}
			if breakOut {
				break
			}
		}
	}(a)
}

// Disconnect 断开连接
func (a *WxEslConnection) Disconnect() {
	a.autoConnect = false
	if nil != a.c {
		a.c.Close()
	}
	close(a.eventChan)
}

// NewWxEslConnection new
func NewWxEslConnection(host string, port int, password string, eventList []string) *WxEslConnection {
	instance := &WxEslConnection{
		host:        host,
		port:        port,
		password:    password,
		autoConnect: true,
		eventChan:   make(chan *eventsocket.Event),
		eventList:   eventList,
	}
	instance.connect()
	instance.monitorConnect()
	return instance
}

// GetEventChan 获取事件处理队列
func (a *WxEslConnection) GetEventChan() chan *eventsocket.Event {
	return a.eventChan
}

// SendAsync 发送异步指令
func (a *WxEslConnection) SendAsync(command string) error {
	return a.SendSync(fmt.Sprintf("bgapi %s", command))
}

// SendSync 发送同步指令
func (a *WxEslConnection) SendSync(command string) error {
	if nil != a.c {
		_, err := a.c.Send(command)
		return err
	}
	return ErrDisconnect
}

// SendMsg 发送msg
func (a *WxEslConnection) SendMsg(uuid, cmd, name, arg string) error {
	if nil != a.c {
		a.c.SendMsg(eventsocket.MSG{
			"call-command":     cmd,
			"execute-app-name": name,
			"execute-app-arg":  arg,
		}, uuid, "")
	}
	return ErrDisconnect
}

// 发送监听事件
func (a *WxEslConnection) sendMonitorEventCommand() {
	a.SendSync("events json " + strings.Join(a.eventList, " "))
}

// 启用事件处理
func (a *WxEslConnection) startEventHandle() {
	go func(instance *WxEslConnection) {
		for {
			if !instance.autoConnect {
				break
			}
			if nil != a.c {
				event, err := a.c.ReadEvent()
				if nil == err {
					instance.eventChan <- event
				} else {
					if io.EOF == err {
						// esl 连接已断开
						break
					}
				}
			}
		}
	}(a)
}
