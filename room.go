package douyumsg

import (
	"log"
	"net"
	"protocol"
	"time"
)

const (
	addr    = "openbarrage.douyutv.com:8601" // 服务器地址
	heartbe = "30s"                          // 心跳时间
)

// room对象可以与服务器建立tcp连接，并与之通信
type Room struct {
	RoomId                                  string
	conn                                    net.Conn
	login                                   chan bool
	barrageSwitch, allMsgSwitch, joinSwitch bool
	barrage, allMsg, join                   chan map[string]string
	bool
	logout chan bool
}

// 返回一个room指针
func NewRoom(roomId string) *Room {
	return &Room{RoomId: roomId, login: make(chan bool), logout: make(chan bool)}
}

// 运行这个room
func (r *Room) Run() error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()
	r.conn = conn

	conn.Write(protocol.MsgToByte(map[string]string{
		"type":   "loginreq",
		"roomid": r.RoomId,
	}))
	go receiveMsg(r)
	go r.keepConnection()
	<-r.logout
	return nil
}

// 接收服务器返回的消息
func receiveMsg(r *Room) {
	b := make([]byte, 4039)
	for {
		n, err := r.conn.Read(b)
		if err != nil {
			log.Println("[ERROR]:", err)
			return
		}
		data, err := protocol.ByteToMsg(b[:n])
		if err != nil {
			log.Println("[ERROR]:", err)
			return
		}
		switch data["type"] {
		case "loginres":
			r.login <- true
			// 加入组
			r.conn.Write(protocol.MsgToByte(map[string]string{
				"type": "joingroup",
				"rid":  r.RoomId,
				"gid":  "-9999",
			}))
		case "chatmsg":
			// 弹幕发送
			if r.barrageSwitch {
				r.barrage <- data
			}
		case "uenter":
			// 进入放假
			if r.joinSwitch {
				r.join <- data
			}

		default:
			continue
		}
		if r.allMsgSwitch {
			r.allMsg <- data
		}
	}
}

// 用户进入直播间
func (r *Room) JoinRoom(chanSize int) <-chan map[string]string {
	r.joinSwitch = true
	r.join = make(chan map[string]string, chanSize)
	return r.join
}

// 接收弹幕消息
func (r *Room) ReceiveBarrage(chanSize int) <-chan map[string]string {
	r.barrageSwitch = true
	r.barrage = make(chan map[string]string, chanSize)
	return r.barrage
}

// 接收所有消息
func (r *Room) ReceiveAll(chanSize int) <-chan map[string]string {
	r.allMsgSwitch = true
	r.allMsg = make(chan map[string]string, chanSize)
	return r.allMsg
}

// 客户端与服务器保持连接
func (r *Room) keepConnection() {
	<-r.login
	close(r.login)
	for {
		r.conn.Write(protocol.MsgToByte(map[string]string{
			"type": "mrkl",
		}))
		t, _ := time.ParseDuration(heartbe)
		time.Sleep(t)
	}
}
