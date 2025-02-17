package main

import (
	"bytes"
	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
	"log"
	"log/slog"
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
	runTaskLoop(conn)
}

func runTaskLoop(conn *websocket.Conn) {
	// todo 结构化task
	conf := task.LoadConfig("./config.json")
	ticker := time.NewTicker(time.Second)
	timer := time.NewTimer(time.Second * time.Duration(conf.Time))
	defer ticker.Stop()
	taskId, payload := int64(0), bytes.Repeat([]byte{1}, int(conf.MessageSize))
	reporter := task.NewReporter("metrics.csv")
	go reporter.Run()
	for {
		select {
		case <-timer.C:
			slog.Info("end", slog.Time("任务结束时间", time.Now()))
		case <-ticker.C:
			taskId++
			summaryRes := runTask(taskId, conf.MessageCount, conn, payload)
			reporter.Report(summaryRes)
		}
	}
}

func runTask(taskId int64, runtimes int64, conn *websocket.Conn, payload []byte) (summaryRes *task.TaskResult) {
	summaryRes = &task.TaskResult{TaskID: taskId}
	defer task.SumCpuAndMemUsage(summaryRes)()
	msgId := int64(0)
	for i := int64(0); i < runtimes; i++ {
		start := time.Now().UnixMicro()
		msgId++
		msg := task.Message{TaskID: taskId, MsgID: msgId, MicroSec: time.Now().UnixMicro(), Payload: payload}
		msgBuf, err := jsoniter.Marshal(msg)
		if err != nil {
			summaryRes.ErrCount++
		}
		err = conn.WriteMessage(websocket.TextMessage, msgBuf)
		if err != nil {
			log.Println("Error during message writing:", err)
			summaryRes.ErrCount++
			return
		}
		end := time.Now().UnixMicro()
		summaryRes.Latencies = append(summaryRes.Latencies, end-start)
	}
	return
}

func main() {
	http.HandleFunc("/ws", handler)
	log.Println("Server started at ws://localhost:8080/ws")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
