#!/bin/bash
reso_addr='registry.cn-shenzhen.aliyuncs.com/paipai/im-api-dev'
tag='latest'

container_name="pai-pai-im-api-test"

docker stop ${container_name}

docker rm ${container_name}

docker rmi ${reso_addr}:${tag}

docker pull ${reso_addr}:${tag}


# 如果需要指定配置文件
# docker run -p 10001:8080 --network paipai -v /paipai/config/user-rpc:/user/conf/ --name=${container_name} -d ${reso_addr}:${tag}
docker run -p 8882:8882  --name=${container_name} -d ${reso_addr}:${tag}