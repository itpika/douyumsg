package main

import (
	"client"
	"log"
	"os"
	"runtime"
)

func errHandler(err error) {
	if err != nil {
		log.Println("[ERROR]:", err)
		os.Exit(1)
	}
}

func main() {
	room := client.NewRoom("122402")
	go room.Run()
	go func() {
		msg := room.ReceiveBarrage(0)
		for {
			m := <-msg
			log.Println("level:", m["level"], m["nn"], ":", m["txt"])
		}
	}()
	go func() {
		msg := room.JoinRoom(0)
		for {
			m := <-msg
			log.Println("用户：", "level:", m["level"], m["nn"], "进入直播间")
		}
	}()
	runtime.Goexit()
}
