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

	// 2.create data channel
	dataChannel, err := peerConnection.CreateDataChannel("fool", nil)
	if err != nil {
		return
	}

	dataChannel.OnOpen(func() {
		fmt.Println("data channel has opened")
		i := -1000
		for range time.NewTicker(time.Second * 5).C {
			i++
			if err := dataChannel.SendText(fmt.Sprintf("offer Hello %d", i)); err != nil {
				fmt.Println(err)
			}
		}
	})

	dataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
		fmt.Printf("Message from DataChannel '%s': '%s'\n", dataChannel.Label(), string(msg.Data))
	})

	// 3.create offer
	offer, err := peerConnection.CreateOffer(nil)
	if err != nil {
		panic(err)
		return
	}
	// 4.set local description
	err = peerConnection.SetLocalDescription(offer)
	if err != nil {
		panic(err)
		return
	}

	// 5.print answer
	fmt.Println("OFFER:")
	println(utils.Encode(offer))

	// 6.input answer
	println("INPUT ANSWER:")
	var answer webrtc.SessionDescription
	answerStr, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		panic(err)
	}
	utils.Decode(answerStr, &answer)
	// 7.set remote description
	if err = peerConnection.SetRemoteDescription(answer); err != nil {
		panic(err)
	}
	select {}
}
