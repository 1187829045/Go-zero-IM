#!/bin/bash
reso_addr='registry.cn-hangzhou.aliyuncs.com/llb-easy-chat/im-api-dev'
tag='latest'

container_name="llb-chat-im-api-test"

docker stop ${container_name}

docker rm ${container_name}

docker rmi ${reso_addr}:${tag}

docker pull ${reso_addr}:${tag}


# 如果需要指定配置文件的
# docker run -p 10001:8080 --network imooc_easy-im -v /easy-im/config/user-rpc:/user/conf/ --name=${container_name} -d ${reso_addr}:${tag}
docker run -p 8882:8882  --name=${container_name} -d ${reso_addr}:${tag}