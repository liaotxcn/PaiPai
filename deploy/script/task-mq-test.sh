#!/bin/bash
reso_addr='registry.cn-shenzhen.aliyuncs.com/paipai/task-mq-dev'
tag='latest'

container_name="pai-pai-task-mq-test"

docker stop ${container_name}

docker rm ${container_name}

docker rmi ${reso_addr}:${tag}

docker pull ${reso_addr}:${tag}


# 如果需要指定配置文件
# docker run -p 10001:8080 --network paipai -v /paipai/config/user-rpc:/user/conf/ --name=${container_name} -d ${reso_addr}:${tag}
# task-mq不对外服务，不需要指定启动端口
docker run --name=${container_name} -d ${reso_addr}:${tag}