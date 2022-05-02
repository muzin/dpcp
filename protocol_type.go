package dpcp

type ProtocolType uint8

// enum ProtocolType
const (
	IPv4 ProtocolType = 4

	IPv6 ProtocolType = 6
)

// 是否是 协议类型
func IsProtocolType(val uint8) bool {
	valType := ProtocolType(val)
	if valType == IPv4 || valType == IPv6 {
		return true
	} else {
		return false
	}
}

func (this ProtocolType) String() string {
	switch this {
	case IPv4:
		return "IPv4"
	case IPv6:
		return "IPv6"
	default:
		return ""
	}
}
