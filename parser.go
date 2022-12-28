package dpcp

import (
	bytes2 "bytes"
	"math/rand"
	"time"
)

const (
	MAX_DATA_LENGTH   uint16 = 1<<16 - 1
	MAX_HEAD_LENGTH   uint16 = 1<<16 - 1
	MAX_HEADER_LENGTH uint16 = (1<<16 - 1) / 2
)

var (
	lastTimeStamp int64 = 0
	sn            int64 = 0
)

// ValidateMessageCrc 校验消息的CRC
func ValidateMessage(bytes []byte) bool {

	msg := NewMessage()

	// 协议 + 是否含有头 + 传输类型 + 数据代理 位
	msg.protocol = ProtocolType(bytes[1] >> 4)

	// 校验 协议 是否正确
	if !IsProtocolType(bytes[1] >> 4) {
		return false
	}

	// 首部长度 位
	headHighByte := uint16(bytes[2])
	headLowByte := uint16(bytes[3])
	msg.headLength = headHighByte<<8 | headLowByte

	ipLength := msg.getIpLength()

	// 总长度 位
	totalHigh1Byte := uint32(bytes[4])
	totalHigh2Byte := uint32(bytes[5])
	totalLow1Byte := uint32(bytes[6])
	totalLow2Byte := uint32(bytes[7])
	msg.length = totalHigh1Byte<<24 | totalHigh2Byte<<16 | totalLow1Byte<<8 | totalLow2Byte

	// Header 长度 位
	headerHighByte := uint16(bytes[24])
	headerLowByte := uint16(bytes[25])
	msg.headerLength = headerHighByte<<8 | headerLowByte

	// Data 长度 位
	dataHighByte := uint16(bytes[26])
	dataLowByte := uint16(bytes[27])
	msg.dataLength = dataHighByte<<8 | dataLowByte

	headLen := uint16(FIXED_HEAD_LENGTH) + uint16(ipLength)*2 + 2
	// 校验 head长度 是否正确
	if msg.headLength != headLen+msg.headerLength {
		return false
	}

	if uint32(msg.headLength+msg.dataLength) != msg.length {
		return false
	}

	// CRC16 校验位
	crc16BitIdx := headLen - 2
	crc16HighByte := uint16(bytes[crc16BitIdx])
	crc16LowByte := uint16(bytes[crc16BitIdx+1])
	msg.crc16 = crc16HighByte<<8 | crc16LowByte

	headBytes := make([]byte, msg.headLength)
	copy(headBytes, bytes)
	headBytes[headLen-2] = 0
	headBytes[headLen-1] = 0
	mbcrc16 := UsMBCRC16(headBytes, int(msg.headLength))
	// 校验 CRC16 head 校验字节是否正确
	if crc16HighByte != uint16(mbcrc16)>>8 {
		return false
	}

	return true
}

// ParseMessageLength 解析消息长度
func ParseMessageLength(bytes []byte) uint32 {
	// 总长度 位
	totalHigh1Byte := uint32(bytes[4])
	totalHigh2Byte := uint32(bytes[5])
	totalLow1Byte := uint32(bytes[6])
	totalLow2Byte := uint32(bytes[7])

	return totalHigh1Byte<<24 | totalHigh2Byte<<16 | totalLow1Byte<<8 | totalLow2Byte
}

func ParseMessageHeadLength(bytes []byte) uint16 {
	// 首部长度 位
	headHighByte := uint16(bytes[2])
	headLowByte := uint16(bytes[3])
	headLength := headHighByte<<8 | headLowByte
	return headLength
}

// ParseMessage 解析消息
func ParseMessage(bytes []byte) *Message {

	if len(bytes) < int(FIXED_HEAD_LENGTH) {
		return nil
	}

	msg := NewMessage()

	// 版本 位
	msg.version = bytes[0]

	// 协议 + 是否含有头 + 传输类型 + 数据代理 位
	msg.protocol = ProtocolType(bytes[1] >> 4)

	hasHeaderBit := bytes[1] >> 3 & 0x01
	if hasHeaderBit == 1 {
		msg.hasHeader = true
	}

	dataTypeBit := bytes[1] >> 2 & 0x01
	msg.dataType = DataType(dataTypeBit)

	dataProxyBit := bytes[1] >> 1 & 0x01
	if dataProxyBit == 1 {
		msg.dataProxy = true
	}

	// 获取IP地址长度
	ipLength := msg.getIpLength()

	// 首部长度 位
	headHighByte := uint16(bytes[2])
	headLowByte := uint16(bytes[3])
	msg.headLength = headHighByte<<8 | headLowByte

	// 总长度 位
	totalHigh1Byte := uint32(bytes[4])
	totalHigh2Byte := uint32(bytes[5])
	totalLow1Byte := uint32(bytes[6])
	totalLow2Byte := uint32(bytes[7])
	msg.length = totalHigh1Byte<<24 | totalHigh2Byte<<16 | totalLow1Byte<<8 | totalLow2Byte

	// 会话标识 位
	for i := 1; i <= 8; i++ {
		msg.sessionNo |= int64(bytes[7+i]) << (64 - i*8)
	}

	// 序列 位
	for i := 1; i <= 8; i++ {
		msg.serialNo |= int64(bytes[15+i]) << (64 - i*8)
	}

	// Header 长度 位
	headerHighByte := uint16(bytes[24])
	headerLowByte := uint16(bytes[25])
	msg.headerLength = headerHighByte<<8 | headerLowByte

	// Data 长度 位
	dataHighByte := uint16(bytes[26])
	dataLowByte := uint16(bytes[27])
	msg.dataLength = dataHighByte<<8 | dataLowByte

	// Src Port 位
	srcPortHighByte := uint16(bytes[28])
	srcPortLowByte := uint16(bytes[29])
	msg.srcPort = srcPortHighByte<<8 | srcPortLowByte

	// Dest Port 位
	destPortHighByte := uint16(bytes[30])
	destPortLowByte := uint16(bytes[31])
	msg.destPort = destPortHighByte<<8 | destPortLowByte

	// Src IP Address 位
	msg.srcAddr = make([]byte, ipLength)
	for i := 0; i < ipLength; i++ {
		msg.srcAddr[i] = bytes[32+i]
	}

	// Dest IP Address 位
	msg.destAddr = make([]byte, ipLength)
	destIpAddrBitIdx := 32 + ipLength
	for i := 0; i < ipLength; i++ {
		msg.destAddr[i] = bytes[destIpAddrBitIdx+i]
	}

	// CRC16 校验位
	crc16BitIdx := destIpAddrBitIdx + ipLength
	crc16HighByte := uint16(bytes[crc16BitIdx])
	crc16LowByte := uint16(bytes[crc16BitIdx+1])
	msg.crc16 = crc16HighByte<<8 | crc16LowByte

	// headerBytes
	headerBitIdx := crc16BitIdx + 2
	if msg.headerLength == 0 {
		msg.headers = make(map[string]string)
	} else {
		msg.headers = bytesToMap(bytes[headerBitIdx : headerBitIdx+int(msg.headerLength)])
	}

	// dataBytes
	dataBitIdx := headerBitIdx + int(msg.headerLength)
	msg.data = make([]byte, msg.dataLength)
	copy(msg.data, bytes[dataBitIdx:dataBitIdx+int(msg.dataLength)])

	return msg
}

// bytesToMap 字节 转 map
func bytesToMap(bytes []byte) map[string]string {
	maps := make(map[string]string)
	sep := []byte("\r\n")
	kvSep := []byte(": ")
	byteSplits := bytes2.Split(bytes, sep)
	for i := 0; i < len(byteSplits); i++ {
		bytesItem := byteSplits[i]
		kvBytes := bytes2.Split(bytesItem, kvSep)
		if len(kvBytes) == 1 {
			maps[string(kvBytes[0])] = ""
		} else if len(kvBytes) == 2 {
			maps[string(kvBytes[0])] = string(kvBytes[1])
		}
	}
	return maps
}

func MessageToByteArray(message *Message) []byte {

	this := message

	headerBytes := HeadersToByteArray(*this)
	dataBytes := this.data
	dataLength := this.dataLength

	headLength := this.GetHeadLength()
	totalLength := this.GetLength()
	headerLength := this.GetHeaderLength()
	ipLength := this.getIpLength()

	bytes := make([]byte, headLength)

	// 版本 位
	bytes[0] = this.version

	// 协议 + 是否含有头 + 传输类型 + 数据代理 位
	var hasHeaderBytes byte = 0
	var dataProxyBytes byte = 0
	if this.hasHeader == true {
		hasHeaderBytes = 1
	} else {
		hasHeaderBytes = 0
	}
	if this.dataProxy == true {
		dataProxyBytes = 1
	} else {
		dataProxyBytes = 0
	}
	bytes[1] = uint8(this.protocol)<<4 | hasHeaderBytes<<3 | uint8(this.dataType)<<2 | dataProxyBytes<<1

	// 首部长度 位
	headTotalLen := (headLength + headerLength)
	bytes[2] = byte(headTotalLen >> 8)
	bytes[3] = byte(headTotalLen & 0xFF)

	// 总长度 位
	bytes[4] = byte(totalLength >> 24 & 0xFF)
	bytes[5] = byte(totalLength >> 16 & 0xFF)
	bytes[6] = byte(totalLength >> 8 & 0xFF)
	bytes[7] = byte(totalLength & 0xFF)

	// 会话标识 位
	for i := 1; i <= 8; i++ {
		bytes[7+i] = byte(this.sessionNo >> (64 - i*8) & 0xFF)
	}

	// 序列 位
	for i := 1; i <= 8; i++ {
		bytes[15+i] = byte(this.serialNo >> (64 - i*8) & 0xFF)
	}

	// Header 长度 位
	bytes[24] = byte(headerLength >> 8)
	bytes[25] = byte(headerLength & 0xFF)

	// Data 长度 位
	bytes[26] = byte(dataLength >> 8)
	bytes[27] = byte(dataLength & 0xFF)

	// Src Port 位
	bytes[28] = byte(this.srcPort >> 8)
	bytes[29] = byte(this.srcPort & 0xFF)

	// Dest Port 位
	bytes[30] = byte(this.destPort >> 8)
	bytes[31] = byte(this.destPort & 0xFF)

	// Src IP Address 位
	for i := 0; i < ipLength; i++ {
		bytes[32+i] = this.srcAddr[i]
	}

	// Dest IP Address 位
	destIpAddrBitIdx := 32 + ipLength
	for i := 0; i < ipLength; i++ {
		bytes[destIpAddrBitIdx+i] = this.destAddr[i]
	}

	// Header 位
	// 由 append 拼接

	// Data 位
	// 由 append 拼接

	bytes = append(append(bytes, headerBytes...), dataBytes...)

	// CRC16 校验位
	mbcrc16 := UsMBCRC16(bytes, int(totalLength))
	headmbcrc16 := UsMBCRC16(bytes, int(headTotalLen))
	this.crc16 = (uint16(headmbcrc16) & 0xff00) | (uint16(mbcrc16) & 0xff)
	crc16BitIdx := destIpAddrBitIdx + ipLength
	bytes[crc16BitIdx] = byte(this.crc16 >> 8 & 0xFF)
	bytes[crc16BitIdx+1] = byte(this.crc16 & 0xFF)

	this = nil

	return bytes
}

// HeadersToByteArray 将 message 中 header 部分转换成 []byte
func HeadersToByteArray(message Message) []byte {
	ret := ""
	headers := message.headers
	for k, v := range headers {
		ret += k + ": " + v + "\r\n"
	}
	headerLength := 0
	if len([]byte(ret)) > int(MAX_HEADER_LENGTH)+2 {
		headerLength = int(MAX_HEADER_LENGTH)
	} else {
		retLen := len([]byte(ret))
		if retLen == 0 {
			headerLength = 0
		} else {
			headerLength = len([]byte(ret)) - 2
		}
	}
	if headerLength == 0 {
		return make([]byte, 0)
	} else {
		return []byte(ret[0:headerLength])
	}
}

// GenSnowFlake 雪花算法生成ID
func GenSnowFlakeId() int64 {
	intn := rand.Intn(1024)
	return GenSnowFlakeIdWithMachineId(intn)
}

func GenSnowFlakeIdWithMachineId(machineId int) int64 {
	// 如果想让时间戳范围更长，也可以减去一个日期
	curTimeStamp := time.Now().UnixNano() / 1e6

	if curTimeStamp == lastTimeStamp {
		// 2的12次方 -1 = 4095，每毫秒可产生4095个ID
		if sn > 4095 {
			time.Sleep(time.Millisecond)
			curTimeStamp = time.Now().UnixNano() / 1e6
			sn = 0
		}
	} else {
		sn = 0
	}
	sn++
	lastTimeStamp = curTimeStamp
	// 应为时间戳后面有22位，所以向左移动22位
	curTimeStamp = curTimeStamp << 22
	machineId = machineId << 12
	// 通过与运算把各个部位连接在一起
	return curTimeStamp | int64(machineId) | sn
}
