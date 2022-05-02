package dpcp

import (
	"fmt"
	"testing"
)

func TestMessage_ToByteArray(t *testing.T) {
	t.Run("", func(t *testing.T) {

		message := NewMessage()
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

		srcAddr := message.GetSrc()
		destAddr := message.GetDest()
		fmt.Printf("%v\t%v\n", srcAddr, destAddr)

		byteArray := message.ToByteArray()

		fmt.Println(byteArray)

		newMessage := ParseMessage(byteArray)

		fmt.Println(newMessage)

		t.Logf("want")
	})
}
