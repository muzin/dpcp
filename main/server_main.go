package main

import (
	"fmt"
	"github.com/muzin/dpcp"
	"github.com/muzin/go_rt/net"
	"github.com/muzin/go_rt/try"
	"time"
)

func main() {

	dpcpServer := dpcp.NewDpcpServer()

	dpcpServer.SetMTU(65535)

	dpcpServer.OnConnect(func(sess *dpcp.Session) {

		fmt.Printf("sess connect %v\n", sess.GetId())

		now := time.Now()

		//byteCnt := 0
		//pkgCnt := 0

		sess.OnData(func(args ...interface{}) {
			bytes := args[0].([]byte)
			//byteCnt += len(bytes)
			//pkgCnt += 1
			//fmt.Printf("on data: %v %v\r", byteCnt, pkgCnt)

			sess.Write(bytes)
		})

		//sess.OnMessage(func(msg dpcp.Message){
		//
		//	//fmt.Printf("msg: %v\n", msg)
		//
		//	//if msg.GetDataType() == dpcp.DATA {
		//	//sess.SendMsg(msg)
		//	//}else {
		//	//
		//	//}
		//
		//})

		sess.OnClose(func(args ...interface{}) {
			since := time.Since(now)
			fmt.Printf("socket close: %v\n", since)
		})

		sess.OnError(func(args ...interface{}) {
			since := time.Since(now)
			throwable := args[0].(try.Throwable)
			fmt.Printf("socket error: %v %v\n", since, throwable)
		})

		sess.OnTimeout(func(args ...interface{}) {
			since := time.Since(now)
			fmt.Printf("socket timeout: %v \n", since)
		})

		sess.OnEnd(func(args ...interface{}) {
			since := time.Since(now)
			fmt.Printf("socket end: %v\n", since)

		})

		//sessMsg := sess.NewMessage()
		//sessMsg.SetSessionNo()
		// dpcp version mtu
		// dpcpVersionInfo := []byte("dpcp " + string(dpcp.Version) + " " + string(dpcpServer.GetMTU() >> 8))
		// fmt.Println(string(dpcpVersionInfo))
		// sess.Write(dpcpVersionInfo)
	})

	dpcpServer.OnListen(func(args ...interface{}) {
		fmt.Printf("listening...\n")
	})

	dpcpServer.OnError(func(args ...interface{}) {
		throwable := args[0].(try.Throwable)

		fmt.Printf("error: %v\n", throwable)
	})

	//fmt.Printf("want\n")

	dpcpServer.Listen(7747, "127.0.0.1")

	net.ExitAfterSocketEnd()

}
