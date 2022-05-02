package dpcp

type DataType uint8

// enum DataType
// 定义传输消息的类型
//
const (
	DATA DataType = 0

	INFO DataType = 1
)

func (this DataType) String() string {
	switch this {
	case DATA:
		return "DATA"
	case INFO:
		return "INFO"
	default:
		return ""
	}
}
