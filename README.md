package main

import (
	"fmt"
	"github.com/dalgurak007/douyumsg"
	"os"
	"runtime"
)

func main() {
	// 根据房间号码获取一个房间
	room := douyumsg.NewRoom("122402")
	// 与服务器建立连接
	go func() {
		err := room.Run("openapi-danmu.douyu.com:8601")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()
	// 获取弹幕消息，传入chan缓冲区大小，返回一个chan
	go func() {
		msg := room.ReceiveBarrage(100)
		for {
			m := <-msg
			fmt.Println("level:", m["level"], m["nn"], ":", m["txt"])
		}
	}()
	// 获取用户进入房间通知
	go func() {
		msg := room.JoinRoom(0)
		for {
			m := <-msg
			fmt.Println("用户：", "level:", m["level"], m["nn"], "进入直播间")
		}
	}()
	// 获取所有消息，同样返回一个chan，需要自己对消息进行过滤处理，格式参考斗鱼弹幕服务器第三方接入协议v1.6.2.pdf
	//for {
	//	msg := <-room.ReceiveAll(100)
	//	fmt.Println(msg)
	//}

	runtime.Goexit()
}
