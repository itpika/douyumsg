@@ -1,56 +1,44 @@
# go-douyuchatmsg for golang
___
### 斗鱼弹幕客户端
package main

功能：
- 获取斗鱼服务器弹幕消息，礼物消息，动态消息等
#### 安装

```
go get -u github.com/itpika/douyumsg
```
#### 导入

```
import "github.com/itpika/douyumsg"
```
#### 快速开始

```
package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/itpika/douyumsg/lib/logger"

	"github.com/itpika/douyumsg"
)

func main() {
	// 根据房间号码获取一个房间
	room := douyumsg.NewRoom("276200")
	// 与服务器建立连接
	err := room.Run("openapi-danmu.douyu.com:8601")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	room.SetHeart(7)              // 设置心跳时间间隔,，默认30s
	room.SetBarrageChanSize(200)  // 设置弹幕消息channel大小，默认100
	room.SetUserEnterChanSize(20) // 设置弹幕消息channel大小，默认50
	room.SetAllMsgChanSize(400)   // 设置弹幕消息channel大小，默认300
	room.SetgiftChanSize(10)      // 设置礼物消息channel大小，默认30
	// 获取弹幕消息
	go func() {
		msg := room.ReceiveBarrage()
		for {
			m := <-msg
			if m == nil {
				println("弹幕消息队列关闭")
				break
			}
			logger.Infof("等级:[%s] %s : %s\n", m["level"], m["nn"], m["txt"])
		}
	}()
	// 获取礼物消息
	go func() {
		msg := room.Gify()
		for {
			m := <-msg
			if m == nil {
				println("礼物消息队列关闭")
				break
			}
			logger.Infof("等级:[%s] %s 赠送了%s个礼物， %s连击\n", m["level"], m["nn"], m["gfcnt"], m["hits"])

		}
	}()
	// 获取用户进入房间消息
	go func() {
		msg := room.UserEnter()
		for {
			m := <-msg
			if m == nil {
				println("进入房间队列关闭")
				break
			}
			logger.Infof("等级:[%s] %s : 进入直播间\n", m["level"], m["nn"])
		}
	}()
	// 获取所有消息，同样返回一个chan，需要自己对消息进行过滤处理，格式参考斗鱼弹幕服务器第三方接入协议v1.6.2.pdf
	//for {
	//	msg := <-room.ReceiveAll(100)
	//	fmt.Println(msg)
	//}
	// room.Stop() 终止消息接收
	runtime.Goexit()
}


```

##### 斗鱼服务消息格式参考

* https://github.com/itpika/douyumsg/blob/master/%E6%96%97%E9%B1%BC%E5%BC%B9%E5%B9%95%E6%9C%8D%E5%8A%A1%E5%99%A8%E7%AC%AC%E4%B8%89%E6%96%B9%E6%8E%A5%E5%85%A5%E5%8D%8F%E8%AE%AEv1.6.2.pdf

* 官方文档 https://open.douyu.com/source/api/63

___
作者： itpika
欢迎大家加入开发。
