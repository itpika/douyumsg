package protocol

import (
	// "bytes"
	"encoding/binary"
	// "io"
	"strings"
)

const (
	HeadLen    uint32 = 4
	MsgTypeLen uint32 = 2
	KeepLen    uint32 = 2
)

// 序列化消息
func serializeMsg(msg map[string]string) (str string) {
	for k, v := range msg {
		str += k + "@=" + v + "/"
	}
	str += string(0)
	return
}

// 反序列化消息
func unserializeMsg(str string) map[string]string {
	m := make(map[string]string)
	if str[len(str)-1:] != string(0) {
		return m
	}
	// 截取最后的空字符和/
	str = str[:len(str)-2]
	slis := strings.Split(str, "/")
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
	s := serializeMsg(msg)
	// fmt.Println(s, len(s))

	// 计算头长度
	head := make([]byte, HeadLen*2+MsgTypeLen+KeepLen)
	// head := make([]byte, 0)

	body := []byte(s)

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
func ByteToMsg(data []byte) (m map[string]string, err error) {
	str := string(data)

	// 反序列化
	m = unserializeMsg(str)
	return
}

// // 解析数据
// func ByteToMsg(data []byte) (m map[string]string, err error) {
// 	if uint32(len(data)) <= HeadLen*2+MsgTypeLen+KeepLen {
// 		return
// 	}
// 	reader := bytes.NewReader(data)
// 	// fmt.Println(data)

// 	sli := make([]byte, 4)
// 	_, e := reader.ReadAt(sli, 0)
// 	if e != nil {
// 		err = e
// 		return
// 	}
// 	_, e = reader.ReadAt(sli, 4)
// 	if e != nil {
// 		err = e
// 		return
// 	}
// 	sli = make([]byte, 2)
// 	_, e = reader.ReadAt(sli, 8)
// 	if e != nil {
// 		err = e
// 		return
// 	}
// 	sli = make([]byte, len(data))
// 	n, e := reader.ReadAt(sli, 12)
// 	if e != nil && e != io.EOF {
// 		err = e
// 		return
// 	}
// 	// fmt.Println(sli[:n])
// 	str := string(sli[:n])

// 	// 反序列化
// 	m = unserializeMsg(str)
// 	return
// }
