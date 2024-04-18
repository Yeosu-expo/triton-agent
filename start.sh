#!/bin/bash
# IP 주소를 추출하고 이를 환경변수로 도커 컨테이너에 전달하는 스크립트

# 호스트의 IP 주소를 추출
HOST_IP=$(hostname -I | awk '{print $1}')

# 추출된 IP 주소와 함께 도커 컨테이너 실행
echo "Running Docker container with HOST_IP: $HOST_IP"
#docker run -it --rm --gpus all -e HOST_IP=%IPADDR% --name triton-agent --network triton -p 7000:7000 -p 6934:6934 -v $PWD/../models:/server/models triton-agent 
