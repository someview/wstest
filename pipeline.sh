#!/bin/bash
# run_pipeline.sh

# 启动Go服务端（后台运行）
go run server.go &
SERVER_PID=$!
echo "Server started with PID: $SERVER_PID"

# 等待服务端初始化
sleep 2  # 根据实际启动时间调整

# 运行Go客户端（前台运行）
go run client.go

# 客户端执行完成后绘制图表
echo "Generating server charts..."
python plot_server.py  # 服务端数据可视化

echo "Generating client charts..."
python plot_client.py  # 客户端数据可视化

# 清理服务端进程
kill $SERVER_PID
echo "Server process terminated"
