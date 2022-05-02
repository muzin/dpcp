package dpcp

import (
	"net"
	"strconv"
)

const (
	FIXED_HEAD_LENGTH uint8 = 32
)

// Message
type Message struct {

	// 版本号
	version uint8

	// 协议类型
	protocol ProtocolType

	// 是否存在头数据
	hasHeader bool

	// 数据类型 data=0/info=1
	dataType DataType

	// 数据代理 是否开启数据代理  enable=true/disable=false
	dataProxy bool

	// 总长度
	length uint32

	// 首部长度
	headLength uint16

	// 会话标识
	sessionNo int64

	// 数据序号
	serialNo int64

	// 头数据长度
	headerLength uint16

	// 数据长度
	dataLength uint16

	// 源端口
	srcPort uint16

	// 目标端口号
	destPort uint16

	// 源地址
	srcAddr []byte

	// 目标地址
	destAddr []byte

	// crc16 校验和
	crc16 uint16

	// 头数据
	headers map[string]string

	// 数据
	data []byte

	// 最大消息长度 默认为：-1 不限长度，需要大于66字节
	maxLength int
}

// NewMessage 创建 Message
func NewMessage() *Message {
	msg := &Message{
		version:   Version,
		protocol:  IPv4,
		headers:   make(map[string]string),
		maxLength: -1,
	}
	return msg
}

// ToByteArray 将 Message 对象 转换为字节数组
func (this *Message) ToByteArray() []byte {
	return MessageToByteArray(this)
}

func (this *Message) HeadersToByteArray() []byte {
	ret := ""
	headers := this.headers
	for k, v := range headers {
		ret += k + ": " + v + "\r\n"
	}
	headerLength := 0
	if len([]rune(ret)) > int(MAX_HEADER_LENGTH)+2 {
		headerLength = int(MAX_HEADER_LENGTH)
	} else {
		retLen := len([]rune(ret))
		if retLen == 0 {
			headerLength = 0
		} else {
			headerLength = len([]rune(ret)) - 2
		}
	}
	if headerLength == 0 {
		return make([]byte, 0)
	} else {
		return []byte(ret[0:headerLength])
	}
}

// SetData 写数据
func (this *Message) SetData(bytes []byte) uint16 {

	availableLength := int(MAX_DATA_LENGTH)
	if this.maxLength > 0 {
		availableLength = this.maxLength - int(this.GetHeadLength()+this.GetHeaderLength())
		if availableLength < 0 {
			availableLength = 0
		}
	}

	if availableLength > 0 {
		if len(bytes) <= availableLength {
			this.data = make([]byte, len(bytes))
		} else {
			this.data = make([]byte, availableLength)
		}
		copy(this.data, bytes)
		this.dataLength = uint16(len(this.data))
		return this.dataLength
	} else {
		return 0
	}
}

// WriteData 写数据
func (this *Message) WriteData(bytes []byte) uint16 {
	ret := this.SetData(bytes)
	return ret
}

// AppendData 向 Data 后追加数据
//
// @param 字节数组
//
// @return 追加的字节长度
func (this *Message) AppendData(bytes []byte) uint16 {
	availableLength := int(MAX_DATA_LENGTH - this.dataLength)
	if this.maxLength > 0 {
		availableLength = this.maxLength - int(this.GetHeadLength()+this.GetHeaderLength()+this.dataLength)
		if availableLength < 0 {
			availableLength = 0
		}
	}

	if availableLength > 0 {
		var newData []byte
		if availableLength > len(bytes) {
			newData = append(this.data, bytes...)
			this.data = newData
			return uint16(len(bytes))
		} else {
			newData = append(this.data, bytes[:availableLength]...)
			this.data = newData
			return uint16(availableLength)
		}
	} else {
		return 0
	}
}

// Next 根据 Message 创建下一个 Message
func (this *Message) Next() *Message {
	message := NewMessage()
	message.version = this.version
	message.protocol = this.protocol
	message.hasHeader = false
	message.dataType = DATA
	message.dataProxy = this.dataProxy
	message.sessionNo = this.sessionNo
	message.srcPort = this.srcPort
	message.srcAddr = this.srcAddr
	message.destPort = this.destPort
	message.destAddr = this.destAddr
	return message
}

func (this *Message) SetVersion(version uint8) {
	this.version = version
}

func (this *Message) GetVersion() uint8 {
	return this.version
}

func (this *Message) SetProtocol(protocol ProtocolType) {
	this.protocol = protocol
}

func (this *Message) GetProtocol() ProtocolType {
	return this.protocol
}

func (this *Message) HasHeader() bool {
	return this.hasHeader
}

func (this *Message) SetDataType(dataType DataType) {
	this.dataType = dataType
}

func (this *Message) GetDataType() DataType {
	return this.dataType
}

func (this *Message) DataProxy(dataProxy bool) {
	this.dataProxy = dataProxy
}

func (this *Message) IsDataProxy() bool {
	return this.dataProxy
}

func (this *Message) GetLength() uint32 {
	this.length = uint32(this.GetHeadLength() + this.GetHeaderLength() + this.GetDataLength())
	return this.length
}

func (this *Message) GetHeadLength() uint16 {
	this.headLength = uint16(int(FIXED_HEAD_LENGTH) + (this.getIpLength() * 2) + 2)
	return this.headLength
}

func (this *Message) SetSessionNo(sessionNo int64) {
	this.sessionNo = sessionNo
}

func (this *Message) GetSessionNo() int64 {
	return this.sessionNo
}

func (this *Message) SetSerialNo(serialNo int64) {
	this.serialNo = serialNo
}

func (this *Message) GetSerialNo() int64 {
	return this.serialNo
}

func (this *Message) GetHeaderLength() uint16 {
	return this.getHeaderLength()
}

func (this *Message) GetDataLength() uint16 {
	return this.dataLength
}

func (this *Message) SetSrcPort(srcPort uint16) {
	this.srcPort = srcPort
}

func (this *Message) GetSrcPort() uint16 {
	return this.srcPort
}

func (this *Message) SetDestPort(destPort uint16) {
	this.destPort = destPort
}

func (this *Message) GetDestPort() uint16 {
	return this.destPort
}

func (this *Message) SetSrcAddr(srcAddr []byte) {
	ipLength := this.getIpLength()
	addrLen := len(srcAddr)
	sIdx := ipLength - addrLen
	if sIdx < 0 {
		sIdx = 0
	}
	newAddr := make([]byte, ipLength)
	for i := 0; i < addrLen; i++ {
		newAddr[sIdx+i] = srcAddr[i]
	}
	this.srcAddr = newAddr
}

func (this *Message) GetSrcAddr() []byte {
	return this.srcAddr
}

func (this *Message) SetDestAddr(destAddr []byte) {
	ipLength := this.getIpLength()
	addrLen := len(destAddr)
	sIdx := ipLength - addrLen
	if sIdx < 0 {
		sIdx = 0
	}
	newAddr := make([]byte, ipLength)
	for i := 0; i < addrLen; i++ {
		newAddr[sIdx+i] = destAddr[i]
	}
	this.destAddr = newAddr
}

func (this *Message) GetDestAddr() []byte {
	return this.destAddr
}

func (this *Message) GetSrc() string {
	return this.formatAddr(this.srcAddr, this.srcPort)
}

func (this *Message) GetDest() string {
	return this.formatAddr(this.destAddr, this.destPort)
}

func (this *Message) formatAddr(addr []byte, port uint16) string {
	protocol := this.GetProtocol()
	if protocol == IPv4 {
		return net.IPv4(addr[0], addr[1], addr[2], addr[3]).String() + ":" + strconv.Itoa(int(port))
	} else if protocol == IPv6 {
		return "[" + net.IP(addr).String() + "]:" + strconv.Itoa(int(port))
	} else {
		return ""
	}
}

func (this *Message) getIpLength() int {
	ipLen := 0
	if this.protocol == IPv6 {
		ipLen = int((IPv6 + 2) * 2)
	} else {
		ipLen = int(IPv4)
	}
	return ipLen
}

func (this *Message) getHeaderLength() uint16 {
	var ret = 0
	headers := this.headers
	for k, v := range headers {
		ret += len([]byte(k)) + 2 + len([]byte(v)) + 2
	}
	if ret > int(MAX_HEADER_LENGTH)+2 {
		ret = int(MAX_HEADER_LENGTH)
	} else {
		ret -= 2
	}
	return uint16(ret)
}

func (this *Message) GetHeaders() map[string]string {
	return this.headers
}

func (this *Message) AddHeader(key string, value string) {
	this.headers[key] = value
	this.hasHeader = true
	this.headerLength = this.getHeaderLength()
}

func (this *Message) RemoveHeader(key string) {
	delete(this.headers, key)
	length := len(this.headers)
	if length == 0 {
		this.removeAllHeader()
	} else {
		this.headerLength = this.getHeaderLength()
	}
}

func (this *Message) removeAllHeader() {
	this.headers = make(map[string]string)
	this.hasHeader = false
	this.headerLength = 0
}

func (this *Message) ClearHeaders() {
	this.removeAllHeader()
}

func (this *Message) getHeader(key string) string {
	return this.headers[key]
}

func (this *Message) GetData() []byte {
	return this.data
}

func (this *Message) ClearData() {
	this.SetData([]byte{})
}

// ReverseAddress 反转地址
func (this *Message) ReverseAddress() {
	tmpPort := this.srcPort
	this.srcPort = this.destPort
	this.destPort = tmpPort

	tmpAddr := this.srcAddr
	this.srcAddr = this.destAddr
	this.destAddr = tmpAddr
}

// GetCRC16 获取 CRC16 校验和
func (this *Message) GetCRC16() uint16 {
	return this.crc16
}
