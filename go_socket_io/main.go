package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/googollee/go-socket.io"
)

var allowOriginFunc = func(r *http.Request) bool {
	return true
}

func main() {
	server := socketio.NewServer(nil)

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		log.Println("connected:", s.ID())
		return nil
	})

	server.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
		log.Println("notice:", msg)
		s.Emit("reply", "have "+msg)
	})

	server.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
		fmt.Printf("chat msg\n")
		s.SetContext(msg)
		return "recv " + msg
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		log.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		log.Println("closed", reason)
	})

	go func() {
		if err := server.Serve(); err != nil {
			log.Fatalf("socketio listen error: %s\n", err)
		}
	}()
	defer server.Close()

	//socket
	http.Handle("/demo/", server)

	// 获取程序当前运行的目录
	executable, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	assetPath := filepath.Join(executable, "asset")
	fmt.Printf("runtime:%v\n", assetPath)

	http.Handle("/html", http.StripPrefix("/html", http.FileServer(http.Dir(assetPath))))

	log.Println("Serving at localhost:8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))

}
