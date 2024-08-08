# dameng
## 模块简介
### dm
dameng-database-driver
达梦官方数据库驱动，根据官方提供压缩包解压所得，未进行任何修改
https://eco.dameng.com/download/
### dm_exporter
dameng-database-exporter
达梦数据库 prometheus-metrics-exporter
## 构建方式
仅在linux/arm64,linux/amd64上做兼容测试，详细架构适配请参考达梦驱动
### 构建二进制包
make build
### 构建docker镜像
make docker-build