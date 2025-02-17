# 修正后的 pipeline.sh
#!/bin/bash

# 进入服务端目录并启动
(cd server && go run server.go) &
SERVER_PID=$!
echo "Server started with PID: $SERVER_PID"

# 等待服务端初始化
sleep 2

# 进入客户端目录并运行
(cd client && go run client.go)

# 等待服务端退出
sleep 2

# 生成分析报告
echo "生成服务端报告..."
python analysis.py -i ./server/server_metrics.csv -o ./server/report.png

echo "生成客户端报告..."
python analysis.py -i ./client/client_metrics.csv -o ./client/report.png


