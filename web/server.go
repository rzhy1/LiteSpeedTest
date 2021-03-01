package web

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func ServeFile() {
	http.Handle("/", http.FileServer(http.Dir("web/gui/")))
	http.HandleFunc("/test", updateTest)
	fmt.Println("Start server at http://127.0.0.1:10871")
	http.ListenAndServe(":10871", nil)
}

func updateTest(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		// log.Printf("recv: %s", message)
		links, err := parseLinks(message)
		if err != nil {
			break
		}
		p := ProfileTest{
			Conn:        c,
			MessageType: mt,
			Links:       links,
			Options: ProfileTestOptions{
				Concurrency: 5,
				Timeout:     20 * time.Second,
			},
		}
		go p.testAll(ctx)
		// err = c.WriteMessage(mt, getMsgByte(0, "gotspeed"))
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}