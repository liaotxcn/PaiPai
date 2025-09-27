#!/bin/bash
set -euo pipefail  # 启用严格模式，命令失败时退出，未设置变量时退出，管道命令失败时退出

# 设置为中文输出
LANG=zh_CN.UTF-8
export LANG

# 获取脚本所在目录
script_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
# 项目根目录
project_root="$script_dir/../../"
# 环境变量文件路径
env_file="$project_root/.env.common"

# 加载环境变量文件
load_env_file() {
    if [ -f "$env_file" ]; then
        echo "正在加载环境变量文件: $env_file"
        # 导出环境变量
        while IFS='=' read -r key value || [[ -n "$key" ]]; do
            # 跳过注释行和空行
            if [[ ! "$key" =~ ^# && -n "$key" ]]; then
                # 移除值两边的引号
                # 检查第一个字符是否为引号
                first_char="${value:0:1}"
                last_char="${value: -1}"
                
                # 处理双引号
                if [[ "$first_char" == '"' ]]; then
                    value="${value:1}"
                fi
                if [[ "$last_char" == '"' ]]; then
                    value="${value%?}"
                fi
                
                # 处理单引号
                if [[ "$first_char" == "'" ]]; then
                    value="${value:1}"
                fi
                if [[ "$last_char" == "'" ]]; then
                    value="${value%?}"
                fi
                export "$key"="$value"
            fi
        done < "$env_file"
    else
        echo "警告: 未找到环境变量文件 $env_file"
    fi
}

# 检查必要的命令
sanity_check() {
    if ! command -v docker &> /dev/null; then
        echo "错误: 未安装docker命令，请先安装Docker"
        exit 1
    fi
    
    # jq命令是可选的，用于格式化JSON输出
    jq_installed=true
    if ! command -v jq &> /dev/null; then
        jq_installed=false
        echo "警告: 未安装jq命令，JSON输出将保持原始格式"
    fi
}

# 打印分隔线
print_separator() {
    echo "======================================================================"
}

# 执行主函数
main() {
    # 切换到项目根目录
    if cd "$project_root" 2>/dev/null; then
        echo "切换到项目根目录: $project_root"
    else
        echo "警告: 无法切换到项目根目录，将在当前目录执行"
    fi

    # 欢迎信息
echo -e "\nPaiPai项目服务状态查看工具\n"
print_separator

# 查看Docker版本
echo "Docker 版本信息："
docker --version
if command -v docker-compose &> /dev/null; then
    docker-compose --version
else
    docker compose --version
fi
print_separator

# 查看运行中的容器
echo "运行中的容器："
docker ps -f "name=pai-pai-"
echo -e "\n全部容器（包括已停止）："
docker ps -a -f "name=pai-pai-"
print_separator

# 查看PaiPai相关镜像
echo "PaiPai项目相关镜像："
docker images | grep "pai-pai" || echo "未找到PaiPai相关镜像"
print_separator

# 查看Docker Compose服务状态
echo "Docker Compose服务状态："
if command -v docker-compose &> /dev/null; then
    docker-compose ps 2>/dev/null || echo "无法获取Docker Compose服务状态"
else
    docker compose ps 2>/dev/null || echo "无法获取Docker Compose服务状态"
fi
print_separator

# 查看Docker资源使用情况
echo "Docker资源使用情况："
docker stats --no-stream | head -n 15 || echo "无法获取Docker资源使用情况"
print_separator

# 查看网络情况
echo "PaiPai网络情况："
if docker network ls | grep "pai-pai" &> /dev/null; then
    echo "找到PaiPai网络"
    if [ "$jq_installed" = true ]; then
        echo -e "\nPaiPai网络详情："
        docker network inspect pai-pai --format '{{json .Containers}}' | jq . || echo "无法格式化网络详情"
    else
        echo -e "\nPaiPai网络详情（未格式化）："
        docker network inspect pai-pai --format '{{json .Containers}}' || echo "无法获取网络详情"
    fi
else
    echo "未找到PaiPai网络"
fi
print_separator

# 健康检查状态
echo "容器健康检查状态："
if docker ps -aq -f "name=pai-pai-" &> /dev/null; then
    docker inspect --format='{{.Name}}: {{.State.Health.Status}}' $(docker ps -aq -f "name=pai-pai-") 2>/dev/null || echo "部分容器无健康检查信息"
else
    echo "无PaiPai容器"
fi
print_separator

# 显示端口映射信息
echo "PaiPai服务端口映射："
if docker ps -q -f "name=pai-pai-" &> /dev/null; then
    docker port $(docker ps -q -f "name=pai-pai-") 2>/dev/null || echo "部分容器无端口映射信息"
else
    echo "无运行中的PaiPai服务"
fi
print_separator

# 帮助信息
echo "使用说明："
echo "1. 此脚本用于统一查看PaiPai项目的所有Docker镜像和服务状态"
echo "2. 可以通过修改脚本来自定义查看内容"
echo "3. 常用命令："
echo "   - 停止所有服务: make stop-all"
echo "   - 本地启动所有服务: make local-run"
echo "   - 查看单个容器日志: docker logs <容器名>"
echo "   - 查看MySQL连接: docker exec -it mysql mysql -uroot -p123456"
echo "   - 查看Redis连接: docker exec -it redis redis-cli -a paipai"

# 恢复原目录
cd - > /dev/null 2>&1

# 执行成功提示
echo -e "\n状态检查完成！"
}

# 执行脚本
load_env_file
sanity_check
main

# 明确设置退出状态码
exit 0