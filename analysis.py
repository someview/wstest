import pandas as pd
import matplotlib.pyplot as plt
import seaborn as sns
from datetime import datetime

# 配置可视化样式
plt.style.use('ggplot')
sns.set_palette("husl")

def analyze(csv_path):
    # 读取数据
    df = pd.read_csv(csv_path, parse_dates=['timestamp'])

    # 预处理时间序列
    df['time_window'] = df['timestamp'].dt.floor('1S')  # 按秒聚合

    # 创建画布
    fig, axes = plt.subplots(3, 1, figsize=(14, 12))

    # 绘制延迟百分位数趋势
    latency_fields = ['P50', 'P90', 'P99']
    df.groupby('time_window')[latency_fields].mean().plot(
        ax=axes[0],
        title='Latency Percentiles Trend (Microseconds)',
        ylabel='μs'
    )

    # 绘制CPU使用率
    cpu_fields = ['UserCPUUsage', 'SysCPUUsage']
    df.groupby('time_window')[cpu_fields].mean().plot(
        ax=axes[1],
        kind='area',
        stacked=True,
        title='CPU Usage Breakdown',
        ylabel='% Usage'
    )

    # 绘制内存使用
    df.groupby('time_window')['MemoryUsage'].max().plot(
        ax=axes[2],
        kind='bar',
        title='Memory Usage Peaks',
        ylabel='MB'
    )

    plt.tight_layout()
    plt.savefig('performance_report.png')
    plt.show()

if __name__ == '__main__':
    analyze('metrics.csv')
