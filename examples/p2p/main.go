package main

import (
	"fmt"
	"github.com/chet-hub/FunRTC"
	"github.com/pion/webrtc/v3"
)

func main() {

	var apply,e = FunRTC.ToConnect("Lily", func(channel *webrtc.DataChannel) {
		channel.OnMessage(func(msg webrtc.DataChannelMessage) {
			fmt.Print("Tom get message: " + string(msg.Data) + "\n")
		})
		channel.OnOpen(func() {
			channel.SendText("Hi,Lily. I am Tom")
		})
	},nil,"","")
	if e !=nil{
		panic(e)
	}

	approve,e := FunRTC.Accept("Tom",apply, func(channel *webrtc.DataChannel) {
		channel.OnMessage(func(msg webrtc.DataChannelMessage) {
			fmt.Print("Lily get message: " + string(msg.Data) + "\n")
		})
		channel.OnOpen(func() {
			channel.SendText("Hi,Tom. I am Lily")
		})
	},nil,"","")
	if e !=nil{
		panic(e)
	}

	e = FunRTC.DoConnect("Lily",approve)
	if e !=nil{
		panic(e)
	}


	select {}

}
