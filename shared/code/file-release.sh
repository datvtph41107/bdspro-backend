#!/bin/bash
start=$(date +%s)

# Deploy file-service lên release server (103.172.239.89)
# Dựa trên file-deploy.sh; workspace giống release.sh

REMOTE_USER="root"
REMOTE_HOST=103.172.239.89
WORKSPACE_NAME=golang-microservice
REMOTE_PATH=/root/$WORKSPACE_NAME

cd ../..

service_path=file-service
path_jar=main
payload_build=file

cd $service_path
REMOTE_DIR=${REMOTE_PATH}/${service_path}

echo "--- GO BUILDING ---"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .
echo "--- FINISH BUILDING ---"

echo "Đang đẩy file lên server ${REMOTE_HOST}..."
echo $(pwd)

ssh ${REMOTE_USER}@${REMOTE_HOST} "
mkdir -p ${REMOTE_DIR}
mkdir -p ${REMOTE_DIR}/config
"

scp ${path_jar} ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}

# file-service dùng develop.yml (chưa có production.yml)
scp config/develop.yml ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}/config/develop.yml
scp -r config ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}
scp Dockerfile ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}/Dockerfile
scp -r docs ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}

if [ $? -ne 0 ]; then
  echo "Có lỗi khi đẩy file lên server"
  exit 1
fi

ssh ${REMOTE_USER}@${REMOTE_HOST} "
cd ${REMOTE_PATH}
docker compose build --pull ${payload_build}
docker compose up -d --force-recreate ${payload_build}
"

echo "Deployment và build thành công!"
end=$(date +%s)
elapsed=$((end - start))

echo "Thời gian thực thi: $elapsed giây"
