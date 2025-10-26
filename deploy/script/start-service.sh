#!/bin/bash

# 设置为中文输出
LANG=zh_CN.UTF-8

# 设置必要的环境变量 - 在容器网络内使用服务名
export HOST_IP=127.0.0.1  # 仅用于外部访问
export ETCD_ADDR=etcd:2379  # 在docker网络内使用容器名
export REDIS_ADDR=redis:6379
export MYSQL_ADDR=mysql:3306

echo "开始启动PaiPai服务容器..."

# 检查pai-pai网络是否存在，如果不存在则创建
check_network() {
    echo "检查pai-pai网络..."
    if ! docker network inspect pai-pai > /dev/null 2>&1; then
        echo "创建pai-pai网络..."
        docker network create pai-pai
        if [ $? -ne 0 ]; then
            echo "错误：创建pai-pai网络失败"
            exit 1
        fi
    else
        echo "pai-pai网络已存在"
    fi
}

# 启动容器函数
start_container() {
    local container_name=$1
    local port_mapping=$2
    local image_name=$3
    local service_type=$4
    
    # 先检查容器是否已存在，如果存在则停止并删除
    if [ "$(docker ps -aq -f name=$container_name)" ]; then
        echo "停止并删除已存在的容器: $container_name"
        docker stop $container_name > /dev/null 2>&1
        docker rm $container_name > /dev/null 2>&1
    fi
    
    echo "启动容器: $container_name"
    
    # 所有服务使用相同的环境变量
    docker run -p $port_mapping \
              --network pai-pai \
              --name="$container_name" \
              -e ETCD_ADDR="etcd:2379" \
              -e REDIS_ADDR="redis:6379" \
              -e MYSQL_ADDR="mysql:3306" \
              -d $image_name
    
    if [ $? -eq 0 ]; then
        echo "✅ $container_name 启动成功"
        # 等待服务启动
        sleep 10
    else
        echo "❌ $container_name 启动失败"
        return 1
    fi
}

# 启动服务函数
start_services() {
    # 确保核心服务在pai-pai网络中
    echo "连接核心服务到pai-pai网络..."
    docker network connect pai-pai etcd 2>/dev/null || echo "etcd已在网络中"
    docker network connect pai-pai redis 2>/dev/null || echo "redis已在网络中"
    docker network connect pai-pai mysql 2>/dev/null || echo "mysql已在网络中"
    
    # 按照依赖顺序启动
    echo "启动用户RPC服务..."
    start_container "pai-pai-user-rpc-test" "10000:10000" "paipai-user-rpc-test" "rpc"
    
    echo "等待用户RPC服务注册..."
    sleep 20
    
    # 检查用户RPC是否成功注册
    echo "检查etcd中的服务注册..."
    docker exec etcd etcdctl --endpoints=http://localhost:2379 get --prefix "user.rpc"
    
    echo "启动用户API服务..."
    start_container "pai-pai-user-api-test" "10001:8080" "paipai-user-api-test" "api"
    sleep 5
    
    echo "启动社交RPC服务..."
    start_container "pai-pai-social-rpc-test" "10002:10002" "paipai-social-rpc-test" "rpc"
    sleep 20
    
    echo "检查社交RPC服务注册..."
    docker exec etcd etcdctl --endpoints=http://localhost:2379 get --prefix "social.rpc"
    
    echo "启动社交API服务..."
    start_container "pai-pai-social-api-test" "10003:8080" "paipai-social-api-test" "api"
    sleep 5
    
    echo "启动即时通讯RPC服务..."
    start_container "pai-pai-im-rpc-test" "10004:10004" "paipai-im-rpc-test" "rpc"
    sleep 20
    
    echo "启动即时通讯API服务..."
    start_container "pai-pai-im-api-test" "10005:8080" "paipai-im-api-test" "api"
    sleep 5
    
    echo "启动任务队列服务..."
    start_container "pai-pai-task-mq-test" "10006:10006" "paipai-task-mq-test" "rpc"
}

# 检查并创建网络
check_network

# 启动所有服务
start_services

echo "所有服务启动完成！"
echo "当前运行的服务容器："
docker ps --filter "name=pai-pai-"