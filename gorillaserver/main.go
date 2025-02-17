package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
	"wstest/task"
)

var upgrader = websocket.Upgrader{}

func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error during connection upgrade:", err)
		return
	}
	defer conn.Close()
	go func() {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error during message reading:", err)
				break
			}
		}
	}()
	runTask(conn)
}

func runTask(conn *websocket.Conn) {
	conf := task.LoadConfig("./config.json")
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
}

func main() {
	http.HandleFunc("/ws", handler)
	log.Println("Server started at ws://localhost:8080/ws")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
