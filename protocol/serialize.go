package protocol

import (
	"encoding/binary"
	"strings"
	"unsafe"
)

const (
	HeadLen    uint32 = 4
	MsgTypeLen uint32 = 2
	KeepLen    uint32 = 2
)

// 序列化消息
func serializeMsg(msg map[string]string) []byte {
	buf := make([]byte, 0, 100)
	for k, v := range msg {
		for i, _ := range k {
			buf = append(buf, k[i])
		}
		buf = append(buf, '@', '=')
		for i, _ := range v {
			buf = append(buf, v[i])
		}
		buf = append(buf, '/')
	}
	buf = append(buf, 0)
	return buf
}

//// 序列化消息
//func serializeMsg(msg map[string]string) []byte {
//	var buf bytes.Buffer
//	for k, v := range msg {
//		buf.WriteString(k + "@=" + v + "/")
//	}
//	buf.WriteByte(0)
//	return buf.Bytes()
//}

// 反序列化消息
func unserializeMsg(str *string) map[string]string {
	m := make(map[string]string)
	if (*str)[len(*str)-1:] != string(0) {
		return m
	}
	// 截取最后的空字符和/
	*str = (*str)[:len(*str)-2]
	slis := strings.Split(*str, "/")
	if len(slis) <= 1 {
		return m
	}
	// 分隔符@=
	for _, v := range slis {
		s := strings.Split(v, "@=")
		if len(s) <= 1 {
			continue
		}
		m[s[0]] = s[1]
	}
	return m
}

// 组合协议
func combinationProtocolHead(msg map[string]string) []byte {
	if len(msg) == 0 {
		return make([]byte, 0)
	}
	// 序列化
	body := serializeMsg(msg)
	// fmt.Println(s, len(s))

	// 计算头长度
	head := make([]byte, HeadLen*2+MsgTypeLen+KeepLen)
	// head := make([]byte, 0)

	// 计算总长度

	totalLen := uint32(len(head) + len(body))
	// fmt.Println(totalLen)
	binary.LittleEndian.PutUint32(head, totalLen-HeadLen)
	binary.LittleEndian.PutUint32(head[4:], totalLen-HeadLen)
	binary.LittleEndian.PutUint16(head[8:], 689)
	binary.LittleEndian.PutUint16(head[10:], 0)
	// fmt.Println("head", len(head), head)
	// fmt.Println("body", len(body), body)
	// fmt.Println("sum", len(append(head, body...)), append(head, body...))
	return append(head, body...)
}

func MsgToByte(msg map[string]string) []byte {
	if len(msg) == 0 {
		return make([]byte, 0)
	}
	return combinationProtocolHead(msg)
}

// 解析数据
func ByteToMsg(data []byte) map[string]string {
	str := (*string)(unsafe.Pointer(&data))
	// 反序列化
	return unserializeMsg(str)
}
