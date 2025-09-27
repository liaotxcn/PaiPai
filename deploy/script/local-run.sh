#!/bin/bash

# 设置为中文输出
LANG=zh_CN.UTF-8

# 默认使用国内镜像源配置文件
USE_LOCAL_MIRROR=true

# 解析命令行参数
while [[ $# -gt 0 ]]; do
    case $1 in
        --local-mirror|-l)
            USE_LOCAL_MIRROR=true
            shift
            ;;
        *)
            echo "未知参数: $1"
            echo "使用方法: $0 [--local-mirror|-l]"
            echo "  --local-mirror, -l: 使用国内镜像源配置文件"
            exit 1
            ;;
    esac
done

# 显示当前配置
if [ "$USE_LOCAL_MIRROR" = true ]; then
    echo "使用国内镜像源配置文件: docker-compose.local.yaml"
else
    echo "使用标准配置文件: docker-compose.yaml"
fi

# 停止并删除所有相关容器
stop_and_remove_containers() {
    echo "停止并删除所有PaiPai服务容器..."
    containers=("pai-pai-user-rpc-test" "pai-pai-user-api-test" 
                "pai-pai-social-rpc-test" "pai-pai-social-api-test" 
                "pai-pai-im-rpc-test" "pai-pai-im-api-test" 
                "pai-pai-im-ws-test" "pai-pai-task-mq-test")
    
    for container in "${containers[@]}"; do
        if [ "$(docker ps -aq -f name=$container)" ]; then
            echo "- 停止容器: $container"
            docker stop $container > /dev/null 2>&1
            echo "- 删除容器: $container"
            docker rm $container > /dev/null 2>&1
        fi
    done
}

# 构建所有服务镜像
build_all_services() {
    echo "\n构建所有服务镜像..."
    services=("user-rpc-dev" "user-api-dev" 
              "social-rpc-dev" "social-api-dev" 
              "im-rpc-dev" "im-api-dev" 
              "im-ws-dev" "task-mq-dev")
    
    # 保存当前目录
    current_dir=$(pwd)
    # 切换到项目根目录
    cd "$current_dir/../../"
    
    for service in "${services[@]}"; do
        echo "- 构建服务: $service"
        make $service
        if [ $? -ne 0 ]; then
            echo "构建$service失败，请检查错误"
            exit 1
        fi
    done
    
    # 切回原目录
    cd "$current_dir"
}

# 使用国内Docker镜像加速器
setup_docker_mirror() {
    echo "设置Docker国内镜像加速器..."
    # 检查是否有daemon.json文件
    if [ -f "/etc/docker/daemon.json" ]; then
        echo "检测到已存在的daemon.json文件，将保留原有配置并添加镜像加速器"
        # 备份现有文件
        cp /etc/docker/daemon.json /etc/docker/daemon.json.bak
        # 添加镜像加速器配置
        if grep -q "registry-mirrors" /etc/docker/daemon.json; then
            # 如果已经有registry-mirrors配置，替换为新的加速器
            sed -i 's/"registry-mirrors": \[\(.*\)\]/"registry-mirrors": ["https:\/\/docker.1ms.run"]/' /etc/docker/daemon.json
        else
            # 如果没有registry-mirrors配置，添加完整配置
            sed -i 's/}$/,
"registry-mirrors": ["https:\/\/docker.1ms.run"]
}/' /etc/docker/daemon.json
        fi
    else
        echo "创建新的daemon.json文件并配置镜像加速器"
        echo '{"registry-mirrors": ["https://docker.1ms.run"]}' > /etc/docker/daemon.json
    fi
    
    # 重启Docker服务使配置生效
    echo "重启Docker服务使配置生效..."
    if command -v systemctl &> /dev/null; then
        sudo systemctl daemon-reload
        sudo systemctl restart docker
    elif command -v service &> /dev/null; then
        sudo service docker restart
    else
        echo "警告：无法自动重启Docker服务，请手动重启Docker以应用镜像加速器配置"
    fi
    
    # 等待Docker重启
    sleep 5
}

# 预拉取基础设施镜像，带重试机制
pull_infrastructure_images() {
    echo "预拉取基础设施服务镜像..."
    images=("bitnami/etcd:3.5.20" "redis:alpine3.18" "mysql:8.0" "mongo:4.0" 
           "wurstmeister/zookeeper" "wurstmeister/kafka" "apache/apisix-dashboard:3.0.1-alpine" 
           "apache/apisix:latest" "ccr.ccs.tencentyun.com/hyy-yu/sail:latest" "jaegertracing/all-in-one:latest" 
           "elasticsearch:7.17.4" "kibana:7.17.4")
    
    max_retries=3
    retry_interval=5
    
    for image in "${images[@]}"; do
        retry=0
        success=false
        
        while [ $retry -lt $max_retries ] && [ $success = false ]; do
            echo "- 拉取镜像: $image (尝试 $((retry+1))/$max_retries)"
            docker pull $image
            if [ $? -eq 0 ]; then
                echo "- 镜像 $image 拉取成功"
                success=true
            else
                retry=$((retry+1))
                if [ $retry -lt $max_retries ]; then
                    echo "- 镜像拉取失败，$retry 秒后重试..."
                    sleep $retry_interval
                else
                    echo "- 镜像 $image 拉取失败，已达到最大重试次数"
                fi
            fi
        done
    done
}

# 运行基础设施服务
start_infrastructure() {
    echo "\n启动基础设施服务..."
    
    # 检查Docker Hub连接
    echo "检查Docker Hub连接..."
    docker pull hello-world > /dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo "警告：Docker Hub连接超时，尝试配置国内镜像加速器"
        echo "是否配置国内Docker镜像加速器？(y/n, 默认y)"
        read -t 10 use_mirror
        use_mirror=${use_mirror:-y}
        
        if [ "$use_mirror" = "y" ] || [ "$use_mirror" = "Y" ]; then
            setup_docker_mirror
        fi
        
        echo "尝试预拉取基础设施镜像，这可能需要一些时间..."
        pull_infrastructure_images
    fi
    
    # 选择配置文件
    if [ "$USE_LOCAL_MIRROR" = true ]; then
        COMPOSE_FILE="../../docker-compose.local.yaml"
        echo "使用国内镜像源配置文件启动基础设施服务"
    else
        COMPOSE_FILE="../../docker-compose.yaml"
    fi
    
    # 检查是否有docker-compose命令，如果没有则使用docker compose
    if command -v docker-compose &> /dev/null; then
        docker-compose -f $COMPOSE_FILE up -d etcd redis mysql mongo zookeeper kafka apisix apisix-dashboard
    else
        docker compose -f $COMPOSE_FILE up -d etcd redis mysql mongo zookeeper kafka apisix apisix-dashboard
    fi
    
    if [ $? -ne 0 ]; then
        echo "启动基础设施服务失败，请检查错误"
        echo "可能的解决方案："
        echo "1. 检查网络连接是否正常"
        echo "2. 手动配置Docker镜像加速器"
        echo "3. 尝试手动拉取所需镜像: docker pull 镜像名称"
        echo "4. 检查docker-compose.yaml文件中的镜像配置"
        exit 1
    fi
    
    # 等待基础设施服务启动
    echo "等待基础设施服务初始化..."
    sleep 10
}

# 运行本地构建的服务容器
run_local_containers() {
    echo "\n运行本地构建的服务容器..."
    
    # User RPC
    echo "- 启动user-rpc服务"
    docker run -p 10000:10000 \
        --health-cmd="nc -z localhost 10000 || exit 1" \
        --health-interval=30s \
        --health-timeout=5s \
        --health-retries=3 \
        --name="pai-pai-user-rpc-test" -d paipai-user-rpc-test
    
    # User API
    echo "- 启动user-api服务"
    docker run -p 10001:8080 \
        --health-cmd="curl -sf http://localhost:8080/api/v1/user/ping || exit 1" \
        --health-interval=30s \
        --health-timeout=5s \
        --health-retries=3 \
        --name="pai-pai-user-api-test" -d paipai-user-api-test
    
    # Social RPC
    echo "- 启动social-rpc服务"
    docker run -p 10002:8080 \
        --health-cmd="nc -z localhost 8080 || exit 1" \
        --health-interval=30s \
        --health-timeout=5s \
        --health-retries=3 \
        --name="pai-pai-social-rpc-test" -d paipai-social-rpc-test
    
    # Social API
    echo "- 启动social-api服务"
    docker run -p 10003:8080 \
        --health-cmd="curl -sf http://localhost:8080/api/v1/social/ping || exit 1" \
        --health-interval=30s \
        --health-timeout=5s \
        --health-retries=3 \
        --name="pai-pai-social-api-test" -d paipai-social-api-test
    
    # IM RPC
    echo "- 启动im-rpc服务"
    docker run -p 10004:8080 \
        --health-cmd="nc -z localhost 8080 || exit 1" \
        --health-interval=30s \
        --health-timeout=5s \
        --health-retries=3 \
        --name="pai-pai-im-rpc-test" -d paipai-im-rpc-test
    
    # IM API
    echo "- 启动im-api服务"
    docker run -p 10005:8080 \
        --health-cmd="curl -sf http://localhost:8080/api/v1/im/ping || exit 1" \
        --health-interval=30s \
        --health-timeout=5s \
        --health-retries=3 \
        --name="pai-pai-im-api-test" -d paipai-im-api-test
    
    # IM WS
    echo "- 启动im-ws服务"
    docker run -p 10006:8080 \
        --health-cmd="nc -z localhost 8080 || exit 1" \
        --health-interval=30s \
        --health-timeout=5s \
        --health-retries=3 \
        --name="pai-pai-im-ws-test" -d paipai-im-ws-test
    
    # Task MQ
    echo "- 启动task-mq服务"
    docker run -p 10007:8080 \
        --health-cmd="nc -z localhost 8080 || exit 1" \
        --health-interval=30s \
        --health-timeout=5s \
        --health-retries=3 \
        --name="pai-pai-task-mq-test" -d paipai-task-mq-test
}

# 显示运行状态
show_status() {
    echo "\n所有服务启动完成！"
    echo "\n当前运行的容器："
    docker ps --filter "name=pai-pai-"
    
    echo "\n使用指南："
    echo "1. 访问API服务：http://localhost:10001 (user-api), http://localhost:10003 (social-api), http://localhost:10005 (im-api)"
    echo "2. 查看APISIX仪表盘：http://localhost:9000"
    echo "3. 停止所有服务：make stop-all"
    echo "4. 重新构建并运行：cd deploy/script && ./local-run.sh"
}

# 主函数
main() {
    stop_and_remove_containers
    build_all_services
    start_infrastructure
    run_local_containers
    show_status
}

# 执行主函数
main