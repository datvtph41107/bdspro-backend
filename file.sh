#!/bin/bash

# Cấu hình thông tin server
REMOTE_USER="root"
REMOTE_HOST=14.225.210.29
WORKSPACE_NAME=backend-microservice
REMOTE_PATH=/root/$WORKSPACE_NAME
service_path=file-service
REMOTE_DIR=${REMOTE_PATH}/${service_path}

echo "--- GO BUILDING ---"
GOOS=linux GOARCH=amd64 go build -o main .
echo "--- FINISH BUILDING ---"

cd file-service

ssh ${REMOTE_USER}@${REMOTE_HOST} "
mkdir -p ${REMOTE_DIR}
"
scp -r docs ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}
scp Dockerfile ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}/Dockerfile
scp ffmpeg.Dockerfile ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}/ffmpeg.Dockerfile
scp main ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}/main

rebuild_image=""
if [ "$1" = "rebuild" ]; then
  rebuild_image="docker build -t file-ffmpeg -f ffmpeg.Dockerfile ."
fi
ssh ${REMOTE_USER}@${REMOTE_HOST} "
cd ${REMOTE_DIR}
${rebuild_image}

docker build --no-cache -t file-service . -f Dockerfile
docker stop go-file
docker rm go-file

docker run --name go-file \
  -p 8002:8002 \
  -v $(pwd)/files:/files \
  file-service
"

echo "Deployment và build thành công!"
