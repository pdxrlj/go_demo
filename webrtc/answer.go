package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/pion/webrtc/v3"

	"webrtc/utils"
)

func main() {
	// 1.create peer connection
	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		panic(err)
		return
	}
	defer func() {
		if err := peerConnection.Close(); err != nil {
			panic(err)
		}
	}()

	// 2. on data channel
	peerConnection.OnDataChannel(func(dataChannel *webrtc.DataChannel) {
		dataChannel.OnOpen(func() {
			println("data channel has opened")
			i := -1000
			for range time.NewTicker(time.Second * 5).C {
				i++
				if err := dataChannel.SendText(fmt.Sprintf("answer Hello %d", i)); err != nil {
					fmt.Println(err)
				}
			}
		})

		dataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
			fmt.Printf("Message from DataChannel '%s': '%s'\n", dataChannel.Label(), string(msg.Data))
		})
	})

	// 3.input offer
	println("请输入OFFER:")
	offerStr, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		panic(err)
	}
	var offer webrtc.SessionDescription
	utils.Decode(offerStr, &offer)

	//4. set remote description
	err = peerConnection.SetRemoteDescription(offer)
	if err != nil {
		panic(err)
	}
	//5. create answer
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		panic(err)
	}
	//6. set local description
	err = peerConnection.SetLocalDescription(answer)
	if err != nil {
		panic(err)
	}

	//7. gather complete
	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)
	<-gatherComplete
	//8. print answer
	fmt.Println("ANSWER:")
	println(utils.Encode(peerConnection.LocalDescription()))
	select {}
}
