package task

import "time"

// Task 表示一个发送任务
type Task struct {
	TaskID       int64 `json:"taskId"`     // 任务ID
	MessageCount int64 `json:"messageId"`  // 消息ID
	PayloadLen   int64 `json:"payloadLen"` // 消息长度
	MicroSec     int64 `json:"microSec"`   // 微秒时间戳
}

type Message struct {
	TaskID   int64 `json:"taskId"`
	MsgID    int64 `json:"messageId"`
	MicroSec int64 `json:"microSec"` // 微秒时间戳
	Payload  []byte
}

// taskId 从1开始,messageId从1开始
type TaskResult struct {
	TaskID      int64           `json:"taskId"`      // 任务ID
	TotalTime   time.Duration   `json:"totalTime"`   // 总耗时
	Latencies   []time.Duration `json:"latencies"`   // 每条消息的延迟
	P50         time.Duration   `json:"p50"`         // P50延迟
	P90         time.Duration   `json:"p90"`         // P90延迟
	P99         time.Duration   `json:"p99"`         // P99延迟
	CPUUsage    float64         `json:"cpuUsage"`    // CPU使用率
	MemoryUsage uint64          `json:"memoryUsage"` // 内存使用量
}
