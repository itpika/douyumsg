# go-douyuchatmsg for golang
___
### 斗鱼弹幕客户端

功能：
- 获取斗鱼服务器弹幕消息，礼物消息，动态消息等
#### 安装

```
go get -u github.com/dalgurak007/douyumsg
```
#### 导入

```
import "github.com/dalgurak007/douyumsg"
```
#### 快速开始

```

func main() {
  // 根据房间号码获取一个房间
  room := douyumsg.NewRoom("122402")
  // 与服务器建立连接
  go room.Run()
  // 获取弹幕消息，传入chan缓冲区大小，返回一个chan
  go func() {
    msg := room.ReceiveBarrage(0)
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
  // room.ReceiveAll(0)
  runtime.Goexit()
}

```

##### 斗鱼服务消息格式参考

* https://github.com/dalgurak007/go-douyuchatmsg/blob/master/%E6%96%97%E9%B1%BC%E5%BC%B9%E5%B9%95%E6%9C%8D%E5%8A%A1%E5%99%A8%E7%AC%AC%E4%B8%89%E6%96%B9%E6%8E%A5%E5%85%A5%E5%8D%8F%E8%AE%AEv1.6.2.pdf


___
作者： dalgurak007
欢迎大家加入开发。