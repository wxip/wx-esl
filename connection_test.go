package wxesl

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func GetConnection() *WxEslConnection {
	return NewWxEslConnection("10.10.10.74", 8021, "ClueCon", []string{CHANNEL_HANGUP_COMPLETE})
}

func TestNewWxEslConnection(t *testing.T) {
	NewWxEslConnection("10.10.10.74", 8021, "ClueCon", []string{})
	time.Sleep(time.Second * 10)
}

func TestGetEventChan(t *testing.T) {
	instance := GetConnection()
	eventChan := instance.GetEventChan()

	go func() {
		for {
			data, ok := <-eventChan
			if !ok {
				break
			}
			fmt.Println(data.String())
		}
	}()
	time.Sleep(time.Second * 10)
}

func TestSendAsync(t *testing.T) {
	instance := GetConnection()
	err := instance.SendAsync("originate user/1000 &echo")
	if nil != err {
		log.Fatal(err)
		t.FailNow()
	}
}
