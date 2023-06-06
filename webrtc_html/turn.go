package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/pion/turn/v2"
)

// webrt 中继服务器
func main() {
	publicIP := "192.168.1.22"
	port := 5800 // 13902
	user := "foo"
	pwd := "bar"
	realm := "localhost" // 领域，把这个当做数据库名好了

	// 创建一个 UDP 侦听器然后传递给 turn 服务器，turn 服务器本身不分配任何 upd socket，
	// 需要我们自己创建套接字然后通过 turn 传递给对方。
	// 自己创建的 socket 就允许了我们添加日志记录、存储或修改入站/出站流量
	udpListener, err := net.ListenPacket("udp4", publicIP+":"+strconv.Itoa(port))
	if err != nil {
		log.Printf("Failed to create TURN server listener: %v\n", err)
	}

	// 写死个账号密码用于测试，GenerateAuthKey() 用于将密码进行散列，然后再保存到数据
	usersMap := map[string][]byte{}
	usersMap[user] = turn.GenerateAuthKey(user, realm, pwd)

	turnServer, err := turn.NewServer(turn.ServerConfig{
		Realm: realm,
		// Set AuthHandler callback
		// 每当用户试图通过TURN服务器进行身份验证时，都会调用此函数
		// 返回该用户的密钥，如果未找到用户，则返回 false
		AuthHandler: func(username string, realm string, srcAddr net.Addr) ([]byte, bool) {
			log.Printf("Received connect auth, username=%s, realm=%s\n", username, realm)
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

	// Block until user sends SIGINT or SIGTERM
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	if err = turnServer.Close(); err != nil {
		log.Printf("%v\n", err)
	}

}
