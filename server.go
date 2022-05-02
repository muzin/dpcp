package dpcp

import (
	"github.com/muzin/go_rt/collection/hash_map"
	rt_net "github.com/muzin/go_rt/net"
)

const (
	DEFAULT_MTU int = 4096
)

type DpcpServer struct {
	rt_net.Server

	sessMap *hash_map.HashMap

	// 最大传输单元
	mtu int
}

func NewDpcpServer() *DpcpServer {
	server := &DpcpServer{
		Server:  rt_net.NewTCPServer(),
		sessMap: hash_map.NewHashMap(),
		mtu:     DEFAULT_MTU,
	}
	return server
}

// Listen 监听端口
//
// Listen(port, [host]])
func (this *DpcpServer) Listen(args ...interface{}) {
	this.Server.Listen(args...)
}

func (this *DpcpServer) OnListen(listener func(args ...interface{})) {
	this.Server.OnListen(listener)
}

// OnConnect
func (this *DpcpServer) OnConnect(listener func(session *Session)) {
	this.Server.OnConnect(func(args ...interface{}) {
		socket := args[0].(rt_net.Socket)

		socket.SetBufferSize(this.GetMTU())

		sess := &Session{
			Socket: socket,
			id:     GenSnowFlakeId(),
			from:   socket.LocalAddr(),
			dest:   socket.RemoteAddr(),
			mtu:    this.GetMTU(),
		}
		sess.msgChannel = sess.newMsgChannel()

		// 当 session 结束后 从 session 字典中 去掉
		sess.OnEnd(func(args ...interface{}) {
			this.sessMap.Remove(sess.GetId())
		})

		var sess_iptr interface{} = sess
		this.sessMap.Put(sess.GetId(), sess_iptr)

		if listener != nil {
			listener(sess)
		}

		//go socket.ConnectHandle()
	})
}

// GetMTU 获取最大传输单元
func (this *DpcpServer) GetMTU() int {
	return this.mtu
}

// SetMTU 设置最大传输单元
func (this *DpcpServer) SetMTU(mtu int) {
	this.mtu = mtu
}
