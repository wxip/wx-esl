# WX ESL
## 增加特性
- 自动重连
- 简化操作
## 功能列表
- 获取包
```
go get -u github.com/DreamChaserT/wx-esl
```
- 连接服务器
```
instance := NewWxEsl("127.0.0.1", 8021, "ClueCon", []string{CHANNEL_HANGUP_COMPLETE, CHANNEL_ANSWER})
```
- 获取事件消息
```
eventChan := instance.GetEventChan()
go func(chan *eventsocket.Event) {
	for {
		data, ok := <-eventChan
		if !ok {
			log.Info("event handle stop")
			break
		}
		fmt.Println(data.String())
		uuid := data.Header["Unique-Id"]
		log.Info("uuid is ", uuid)
	}
}(eventChan)
```
- 发送异步指令
```
instance.SendAsync("originate user/1000 &echo")
```
- 发送Msg
```
instance.SendMsg("uuid", "cmd", "name", "arg")

# 挂断电话
instance.SendMsg(uuid, "hangup", "", "")
# 播音
instance.SendMsg("uuid","execute","playback","/tmp/test.wav")
```
- 发送Execute Msg
```
instance.SendExecuteMsg("uuid","name","arg")

# 播音
instance.SendExecuteMsg("uuid","playback","/tmp/test.wav")
```
- 断开连接
```
instance.Disconnect()
```