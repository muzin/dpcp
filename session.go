package dpcp

import (
	rt_net "github.com/muzin/go_rt/net"
	"github.com/muzin/go_rt/try"
	"net"
	"strconv"
)

type Session struct {
	rt_net.Socket

	// 会话标识
	id int64

	// 目标地址
	dest net.Addr

	// 来源地址
	from net.Addr

	// 会话序列
	serialNo int64

	// 最大传输单元
	mtu int

	// 数据缓冲区
	_dataBuffer []byte

	// 消息通道
	msgChannel chan Message
}

func NewSession() *Session {
	sess := &Session{
		Socket: rt_net.NewTCPSocket(),
		id:     GenSnowFlakeId(),
		mtu:    -1,
	}
	sess.msgChannel = sess.newMsgChannel()
	return sess
}

func (this *Session) newMsgChannel() chan Message {
	msgChan := make(chan Message, 100)
	return msgChan
}

func (this *Session) msgChannelProcess() {

}

// Connect
//
// Connect(port, [host, [options]])
func (this *Session) Connect(args ...interface{}) {
	this.Socket.Connect(args...)
}

// SendMsg 发送消息
//
func (this *Session) SendMsg(msg Message) {
	this.serialNo++
	msg.SetSerialNo(this.serialNo)
	bytes := msg.ToByteArray()
	this.Socket.Write(bytes)
}

// OnMessage 消息回调
//
func (this *Session) OnMessage(listener func(msg Message)) {

	this.Socket.OnData(func(args ...interface{}) {

		bytes := args[0].([]byte)

		this.ProcessMessage(bytes, listener)

		//// 如果 存在上一次遗留的数据，拼接上次的数据进行解析
		//if this._dataBuffer != nil {
		//	bytes = append(this._dataBuffer, bytes...)
		//	this._dataBuffer = nil
		//}
		//
		//// 消息校验不通过关闭socket
		//validateMessage := ValidateMessage(bytes)
		//if !validateMessage {
		//	this.sendMsgCrcFail()
		//	this.Socket.Close()
		//}
		//
		//// 解析 数据长度
		//msgLength := ParseMessageLength(bytes)
		//
		//bytesLen := len(bytes)
		//if int(msgLength) > bytesLen {
		//	dataBuffer = bytes
		//	return
		//}else{
		//	msg := ParseMessage(bytes)
		//	msgLength = msg.length
		//	// 如果 获取消息后byte有剩余
		//	if bytesLen - int(msgLength) > 0{
		//
		//	}
		//
		//	if msg != nil {
		//		listener(*msg)
		//	}else{
		//		this.Socket.Close()
		//	}
		//}
	})
}

func (this *Session) ProcessMessage(bytes []byte, listener func(msg Message)) int {

	// 如果 存在上一次遗留的数据，拼接上次的数据进行解析
	if this._dataBuffer != nil {
		bytes = append(this._dataBuffer, bytes...)
		this._dataBuffer = nil
	}

	ret := 0
	bytesLen := len(bytes)

	if bytesLen < int(FIXED_HEAD_LENGTH) {
		this._dataBuffer = bytes
		return ret
	}

	headLength := ParseMessageHeadLength(bytes)
	if bytesLen < int(headLength) {
		this._dataBuffer = bytes
		return ret
	}

	// 消息校验不通过关闭socket
	validateMessage := ValidateMessage(bytes)
	if !validateMessage {
		this.sendMsgCrcFail()
		this.Socket.Close()
	}

	// 解析 数据长度
	msgLength := ParseMessageLength(bytes)

	if int(msgLength) > bytesLen {
		this._dataBuffer = bytes
		return 0
	} else {
		msg := ParseMessage(bytes)
		msgLength = msg.length
		ret += int(msgLength)
		if msg != nil {

			try.Try(func() {
				listener(*msg)
			}, try.CatchUncaughtException(func(throwable try.Throwable) {
				this.Emit("error", throwable)
			}))

			// 如果 获取消息后byte有剩余
			if bytesLen-int(msgLength) > 0 {
				cnt := this.ProcessMessage(bytes[msgLength:], listener)
				ret += cnt
			}

		} else {
			this.Socket.Close()
		}
	}

	return ret
}

// Write 写数据
//
// Write(data []byte[, len int[, index int]])
func (this *Session) Write(args ...interface{}) int {
	return this.Socket.Write(args...)
}

func (this *Session) NewMessage() *Message {
	msg := NewMessage()
	msg.SetSessionNo(this.id)
	this.serialNo++
	msg.SetSerialNo(this.serialNo)
	msg.maxLength = this.mtu
	return msg
}

func (this *Session) GetId() int64 {
	return this.id
}

func (this *Session) GetMTU() int {
	return this.mtu
}

func (this *Session) SetMTU(mtu int) {
	this.mtu = mtu
	this.Socket.SetBufferSize(mtu)
}

// Close 关闭会话
//
func (this *Session) Close() {
	this.Socket.Close()
}

func (this *Session) sendMsgCrcFail() {
	msg := this.NewMessage()
	msg.SetDataType(INFO)
	msg.AddHeader("Info_Type", "Error")
	msg.AddHeader("Error", "Message CRC verification failed")
	this.SendMsg(*msg)
}

func (this *Session) sendMsgLengthExceedsMtu() {
	msg := this.NewMessage()
	msg.SetDataType(INFO)
	msg.AddHeader("Info_Type", "Error")
	msg.AddHeader("Error", "Message length exceeds MTU "+strconv.Itoa(this.mtu))
	this.SendMsg(*msg)
}
