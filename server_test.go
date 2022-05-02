package dpcp

import (
	"github.com/muzin/go_rt/net"
	"github.com/muzin/go_rt/try"
	"testing"
	"time"
)

func TestDpcpServer_Listen(t *testing.T) {
	t.Run("1", func(t *testing.T) {

		dpcpServer := NewDpcpServer()

		dpcpServer.OnConnect(func(sess *Session) {

			t.Logf("sess connect %v\n", sess.GetId())

			now := time.Now()

			sess.OnMessage(func(msg Message) {

				t.Logf("msg: %v\n", msg)

				sess.SendMsg(msg)

			})

			sess.OnClose(func(args ...interface{}) {
				since := time.Since(now)
				t.Logf("socket close: %v\n", since)
			})

			sess.OnError(func(args ...interface{}) {
				since := time.Since(now)
				throwable := args[0].(try.Throwable)
				t.Logf("socket error: %v %v\n", since, throwable)
			})

			sess.OnTimeout(func(args ...interface{}) {
				since := time.Since(now)
				t.Logf("socket timeout: %v \n", since)
			})

			sess.OnEnd(func(args ...interface{}) {
				since := time.Since(now)
				t.Logf("socket end: %v\n", since)

				//socket.Destroy()
				//socket = nil
			})

		})

		dpcpServer.OnListen(func(args ...interface{}) {
			t.Logf("listening...")
		})

		dpcpServer.Listen(7747, "127.0.0.1")

		t.Logf("want")

		net.ExitAfterSocketEnd()

	})
}
