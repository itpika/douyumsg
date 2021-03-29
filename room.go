package douyumsg

import (
	"encoding/binary"
	"github.com/itpika/douyumsg/lib/common"
	"github.com/itpika/douyumsg/lib/logger"
	"github.com/itpika/douyumsg/protocol"
	"io"
	"net"
	"strconv"
	"sync"
	"time"
)

type msgChannel struct {
	channel chan map[string]string
	open    bool
}

// room对象可以与服务器建立tcp连接，并与之通信
type Room struct {
	RoomId string
	conn   net.Conn
	// 心跳间隔: seconds
	heart                                                            int64
	barrageSwitch, allMsgSwitch, userEnterSwitch, giftSwitch         bool
	barrageChanSize, allMsgChanSize, userEnterChanSize, giftChanSize int64
	barrage, allMsg, userEnter, gift                                 chan map[string]string
	exit                                                             bool
	wg                                                               sync.WaitGroup
}

/*
	构建一个room，返回这个room指针
*/
func NewRoom(roomId string) *Room {
	return &Room{RoomId: roomId}
}

/*
	设置心态消息频率(秒)
*/
func (r *Room) SetHeart(heartSecond int64) {
	r.heart = heartSecond
}

/*
	运行这个room
*/
func (r *Room) Run(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		logger.Err(err)
		return err
	}
	r.conn = conn
	// 登录弹幕服务器
	if _, err := conn.Write(protocol.MsgToByte(map[string]string{
		"type":   "loginreq",
		"roomid": r.RoomId,
		"time":   strconv.Itoa(int(time.Now().Unix())),
	})); err != nil {
		logger.Err(err)
		return err
	}
	// 心跳
	go r.keepConnection()
	// 接收消息
	go r.receiveMsg()
	return nil
}

/*
	停止这个room
*/
func (r *Room) Stop() {
	r.exit = true
	r.wg.Wait()
	r.conn.Close()
	logger.Infof("The leave room %s successful\n", r.RoomId)
}

/*
	room客户端与服务器保持连接
*/
func (r *Room) keepConnection() error {
	time.Sleep(time.Second * 3)
	r.wg.Add(1)

	for {
		if r.exit {
			// 登出
			r.conn.Write(protocol.MsgToByte(map[string]string{
				"type": "logout",
			}))
			logger.Info("keepConnection exit")
			r.wg.Done()
			break
		}
		// 发送心跳消息
		if _, err := r.conn.Write(protocol.MsgToByte(map[string]string{
			"type": "mrkl",
		})); err != nil {
			return err
		}
		if r.heart <= 0 {
			r.heart = common.Heartbe
		}
		time.Sleep(time.Second * time.Duration(r.heart))
	}
	return nil
}

/*
	接收服务器返回的消息
*/
func (r *Room) receiveMsg() {
	r.wg.Add(1)
	h := make([]byte, protocol.HeadLen*2+protocol.MsgTypeLen+protocol.KeepLen)
	for {
		if r.exit {
			if r.userEnterSwitch {
				close(r.userEnter)
			}
			if r.allMsgSwitch {
				close(r.allMsg)
			}
			if r.barrageSwitch {
				close(r.barrage)
			}
			if r.giftSwitch {
				close(r.gift)
			}
			r.wg.Done()
			logger.Info("receiveMsg exit")
			break
		}
		// 读取协议头
		n, err := r.conn.Read(h)
		if err != nil {
			if err == io.EOF {
				continue
			}
			logger.Err(err)
			continue
		}
		// 读取body
		b := make([]byte, int(binary.LittleEndian.Uint32(h[0:4]))-int(protocol.HeadLen+protocol.MsgTypeLen+protocol.KeepLen))
		n, err = r.conn.Read(b)
		if err != nil {
			logger.Err(err)
			continue
		}

		// log.Println("data", len(b[:n]), b[:n])
		// return
		data := protocol.ByteToMsg(b[:n])

		switch data["type"] {

		case "chatmsg":
			// 弹幕发送
			if r.barrageSwitch {
				r.barrage <- data
			}
		case "uenter":
			// 进入房间
			if r.userEnterSwitch {
				r.userEnter <- data
			}
		case "dgb":
			// 赠送礼物
			if r.giftSwitch {
				r.gift <- data
			}
		case "loginres":
			// 登陆弹幕服务器后，加入房间组
			r.conn.Write(protocol.MsgToByte(map[string]string{
				"type": "joingroup",
				"rid":  r.RoomId,
				"time": strconv.Itoa(int(time.Now().Unix())),
			}))
		default:
			continue
		}
		if r.allMsgSwitch {
			r.allMsg <- data
		}

		for i, _ := range h {
			h[i] = 0
		}
	}
}

/*
	设置用户赠送礼物消息channel大小
*/
func (r *Room) SetgiftChanSize(chanSize int64) {
	r.giftChanSize = chanSize
}

/*
	设置用户进入直播间消息channel大小
*/
func (r *Room) SetUserEnterChanSize(chanSize int64) {
	r.userEnterChanSize = chanSize
}

/*
	设置弹幕消息channel大小
*/
func (r *Room) SetBarrageChanSize(chanSize int64) {
	r.barrageChanSize = chanSize
}

/*
	设置全部消息channel大小
*/
func (r *Room) SetAllMsgChanSize(chanSize int64) {
	r.allMsgChanSize = chanSize
}

/*
	赠送礼物
*/
func (r *Room) Gify() <-chan map[string]string {
	r.giftSwitch = true
	if r.giftChanSize <= 0 {
		r.giftChanSize = common.GiftChanSize
	}
	r.gift = make(chan map[string]string, r.giftChanSize)
	return r.gift
}

/*
	用户进入直播间
*/
func (r *Room) UserEnter() <-chan map[string]string {
	r.userEnterSwitch = true
	if r.userEnterChanSize <= 0 {
		r.userEnterChanSize = common.UserEnterChanSize
	}
	r.userEnter = make(chan map[string]string, r.userEnterChanSize)
	return r.userEnter
}

/*
	接收弹幕消息
*/
func (r *Room) ReceiveBarrage() <-chan map[string]string {
	r.barrageSwitch = true
	if r.barrageChanSize <= 0 {
		r.barrageChanSize = common.BarrageChanSize
	}
	r.barrage = make(chan map[string]string, r.barrageChanSize)
	return r.barrage
}

/*
	接收所有消息
*/
func (r *Room) ReceiveAll() <-chan map[string]string {
	r.allMsgSwitch = true
	if r.allMsgChanSize <= 0 {
		r.allMsgChanSize = common.AllMsgChanSize
	}
	r.allMsg = make(chan map[string]string, r.allMsgChanSize)
	return r.allMsg
}
