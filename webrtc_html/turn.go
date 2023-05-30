package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/pion/turn/v2"
)

var turnServer *turn.Server

func startTurn(publicIP string, listenPort int, realm string, username string, password string) {
	udpListener, err := net.ListenPacket("udp4", "0.0.0.0:"+strconv.Itoa(listenPort))
	if err != nil {
		log.Panicf("Failed to create TURN server listener: %s", err)
	}

	usersMap := map[string][]byte{}
	usersMap[username] = turn.GenerateAuthKey(username, realm, password)

	turnServer, err = turn.NewServer(turn.ServerConfig{
		Realm: realm,
		// Set AuthHandler callback
		// This is called everytime a user tries to authenticate with the TURN server
		// Return the key for that user, or false when no user is found
		AuthHandler: func(username string, realm string, srcAddr net.Addr) ([]byte, bool) {
			fmt.Printf("Received connect auth, username=%s, realm=%s\n", username, realm)
			// framework will check auth key
			if key, ok := usersMap[username]; ok {
				return key, true
			}
			return nil, false
		},

		// PacketConnConfigs is a list of UDP Listeners and the configuration around them
		PacketConnConfigs: []turn.PacketConnConfig{
			{
				PacketConn: udpListener,
				RelayAddressGenerator: &turn.RelayAddressGeneratorStatic{
					RelayAddress: net.ParseIP(publicIP), // Claim that we are listening on IP passed by user (This should be your Public IP)
					Address:      "0.0.0.0",             // But actually be listening on every interface
				},
			},
		},
	})

	fmt.Printf("turn server public ip=%s, listen port=%d, realm=%s, username=%s, password=%s\n",
		publicIP, listenPort, realm, username, password)

	if err != nil {
		log.Panic(err)
	}
}

func main() {
	var listenAddress string
	listenAddress = "0.0.0.0:8084"
	http.Handle("/", http.FileServer(http.Dir(".")))

	fmt.Printf("HTTP server listen on http://%s\n", listenAddress)

	startTurn("127.0.0.1", 13902, "localhost", "foo", "bar")

	panic(http.ListenAndServe(listenAddress, nil))
}
