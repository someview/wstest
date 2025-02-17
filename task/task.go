package task

import (
	"encoding/json"
	"github.com/shirou/gopsutil/v3/process"
	"os"
	"time"
)

type TaskConfig struct {
	MessageCount int64 `json:"messageId"`
	MessageSize  int64 `json:"MessageSize"`
	Time         int64 `json:"Time"`
}

// GetConfig 从指定路径读取配置文件并返回 TaskConfig
func LoadConfig(path string) *TaskConfig {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	var config TaskConfig
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		panic(err)
	}
	return &config
}

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

type TaskResult struct {
	TaskID       int64         `json:"taskId"`       // 任务ID
	TotalTime    time.Duration `json:"totalTime"`    // 总耗时
	Latencies    []int64       `json:"latencies"`    // 每条消息的延迟，单位微秒
	P50          time.Duration `json:"p50"`          // P50延迟
	P90          time.Duration `json:"p90"`          // P90延迟
	P99          time.Duration `json:"p99"`          // P99延迟
	ErrCount     uint64        `json:"errCount"`     // 错误量
	SysCPUUsage  float64       `json:"SysCpuUsage"`  // 系统态CPU使用率
	UserCPUUsage float64       `json:"UserCpuUsage"` // 用户态CPU使用率
	MemoryUsage  uint64        `json:"memoryUsage"`  // 内存使用量
}

var (
	p, _ = process.NewProcess(int32(os.Getpid()))
)

func SumCpuAndMemUsage(result *TaskResult) func() {
	// 初始化资源监控
	startTime := time.Now()
	startCPU, _ := p.Times()
	startMem, _ := p.MemoryInfo()

	// 返回闭包函数用于defer调用
	return func() {
		// 获取结束状态
		endCPU, _ := p.Times()
		endMem, _ := p.MemoryInfo()
		elapsed := time.Since(startTime).Seconds()
		// 计算CPU
		if startCPU != nil && endCPU != nil && elapsed > 0 {
			result.SysCPUUsage = (endCPU.User - startCPU.User) / elapsed * 100
			result.UserCPUUsage = (endCPU.System - startCPU.System) / elapsed * 100
		}
		// 计算内存
		if startMem != nil && endMem != nil {
			result.MemoryUsage = (endMem.RSS - startMem.RSS) / 1024 / 1024
		}
	}
}
