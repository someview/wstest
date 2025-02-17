package task

import (
	"encoding/csv"
	"os"
	"sort"
	"strconv"
	"sync"
)

type Reporter struct {
	receiver chan *TaskResult
	writer   *csv.Writer
	file     *os.File
	mutex    sync.Mutex
}

func NewReporter(outputPath string) *Reporter {
	f, err := os.Create(outputPath)
	if err != nil {
		panic(err)
	}

	writer := csv.NewWriter(f)

	// 写入CSV头
	header := []string{
		"TaskID", "TotalTime(ms)", "P50(μs)", "P90(μs)", "P99(μs)",
		"Errors", "UserCPU(%)", "SysCPU(%)", "Memory(MB)",
	}
	writer.Write(header)
	writer.Flush()

	return &Reporter{
		receiver: make(chan *TaskResult, 1000),
		writer:   writer,
		file:     f,
	}
}

func (r *Reporter) Report(result *TaskResult) {
	r.receiver <- result
}

func (r *Reporter) Run() {
	defer r.file.Close()

	for res := range r.receiver {
		// 计算百分位数
		sort.Slice(res.Latencies, func(i, j int) bool {
			return res.Latencies[i] < res.Latencies[j]
		})

		p50 := percentile(res.Latencies, 50)
		p90 := percentile(res.Latencies, 90)
		p99 := percentile(res.Latencies, 99)

		// 准备CSV行数据
		record := []string{
			strconv.FormatInt(res.TaskID, 10),
			strconv.FormatInt(res.TotalTime.Milliseconds(), 10),
			strconv.FormatInt(p50, 10),
			strconv.FormatInt(p90, 10),
			strconv.FormatInt(p99, 10),
			strconv.FormatUint(res.ErrCount, 10),
			strconv.FormatFloat(res.UserCPUUsage, 'f', 2, 64),
			strconv.FormatFloat(res.SysCPUUsage, 'f', 2, 64),
			strconv.FormatUint(res.MemoryUsage, 10),
		}

		// 线程安全写入
		r.mutex.Lock()
		r.writer.Write(record)
		r.writer.Flush()
		r.mutex.Unlock()
	}
}

func (r *Reporter) Stop() {
	close(r.receiver)
}

// 百分位数计算函数
func percentile(data []int64, p int) int64 {
	if len(data) == 0 {
		return 0
	}
	index := (len(data) - 1) * p / 100
	return data[index]
}
