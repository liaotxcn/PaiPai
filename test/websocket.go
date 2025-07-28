package main

import (
	"github.com/gorilla/websocket" // websocket库
	"log"
	"net/http"
)

// WebSocket应用实例

var upgrader = websocket.Upgrader{}

func serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil) // http-->websocket 返回websocket连接对象
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer conn.Close()
	for {
		mt, message, err := conn.ReadMessage() // 处理websocket连接对象消息
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = conn.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

//func main() {
//	http.HandleFunc("/ws", serveWs)
//	fmt.Println("websocket已启动")
//	log.Fatal(http.ListenAndServe(":8080", nil))
//}
