package dpcp

import (
	"github.com/muzin/go_rt/net"
	"github.com/muzin/go_rt/timer"
	"testing"
)

func TestNewSession(t *testing.T) {
	t.Run("", func(t *testing.T) {

		session := NewSession()

		session.OnConnect(func(args ...interface{}) {
			t.Logf("session connect")
		})

		session.OnMessage(func(msg Message) {

			t.Logf("receive msg: %v\n", msg)

		})

		session.Connect(7747, "127.0.0.1")

		timer.SetTimeout(func() {
			message := session.NewMessage()
			message.SetDataType(DATA)
			//message.SetProtocol(IPv6)
			message.DataProxy(true)
			message.SetSessionNo(GenSnowFlakeId())
			message.SetSerialNo(1)
			message.SetSrcPort(5540)
			message.SetSrcAddr([]byte{192, 168, 1, 103})
			message.SetDestPort(5541)
			message.SetDestAddr([]byte{192, 168, 1, 105})
			message.SetData([]byte{192, 168, 1, 105})

			message.AddHeader("title", "title")
			message.AddHeader("key", "value\U0001F970")

			session.SendMsg(*message)
		}, 1000)

		t.Logf("want")

		timer.SetTimeout(func() {
			session.Close()
		}, 300000)

		net.ExitAfterSocketEnd()
	})
}
