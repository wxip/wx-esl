package wxesl

import (
	"fmt"
	"testing"
	"time"

	"github.com/cloudflare/cfssl/log"
	"github.com/fiorix/go-eventsocket/eventsocket"
)

func GetWxEsl() *WxEsl {
	return NewWxEsl("10.10.10.74", 8021, "ClueCon", []string{CHANNEL_HANGUP_COMPLETE, CHANNEL_ANSWER})
}

func TestAll(t *testing.T) {
	instance := GetWxEsl()
	eventChan := instance.GetEventChan()

	go func(chan *eventsocket.Event) {
		for {
			data, ok := <-eventChan
			if !ok {
				log.Info("event handle stop")
				break
			}
			fmt.Println(data.String())

			if data.Header["Event-Name"] == CHANNEL_ANSWER {
				uuid := data.Header["Unique-Id"]
				log.Info("uuid is ", uuid)
				time.Sleep(time.Second)
				instance.SendMsg(uuid.(string), "hangup", "", "")
			}

		}
	}(eventChan)

	instance.SendAsync("originate user/1000 &echo")

	instance.SendExecuteMsg("uuid", "name", "arg")

	time.Sleep(time.Second * 10)
	instance.Disconnect()
	time.Sleep(time.Second * 100)
}
