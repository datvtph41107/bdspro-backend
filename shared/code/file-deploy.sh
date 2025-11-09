#!/bin/bash
start=$(date +%s)

# Kiểm tra xem có tham số đầu vào không
if [ -z "$1" ]; then
  echo "Usage: $0 <project_name>"
  exit 1
fi

# Cấu hình thông tin server
REMOTE_USER="root"
REMOTE_HOST=14.225.210.29
WORKSPACE_NAME=backend-microservice
REMOTE_PATH=/root/$WORKSPACE_NAME

cd ../..

# Chọn thư mục dựa vào tham số đầu tiên
service_path=/common
path_jar=''
service_name=$1
payload_build=$1
case "$1" in
  "file")
    service_path=file-service
    path_jar=main
    ;;
  *)
    echo "Tham số không hợp lệ."
    exit 1
    ;;
esac

cd $service_path
REMOTE_DIR=${REMOTE_PATH}/${service_path}

echo "--- GO BUILDING ---"
GOOS=linux GOARCH=amd64 go build -o main .
echo "--- FINISH BUILDING ---"


# Đẩy các file trong thư mục hiện tại lên server
echo "Đang đẩy file lên server..."
echo $(pwd)

ssh ${REMOTE_USER}@${REMOTE_HOST} "
mkdir -p ${REMOTE_DIR}
mkdir -p ${REMOTE_DIR}/config
"

scp ${path_jar} ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}

scp config/develop.yml ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}/config/develop.yml
scp ../shared/code/Dockerfile.develop ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}/Dockerfile

if [ "$1" = "file" ]; then
  scp -r docs ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}
fi
if [ $? -ne 0 ]; then
  echo "Có lỗi khi đẩy file lên server"
  exit 1
fi



ssh ${REMOTE_USER}@${REMOTE_HOST} "
cd ${REMOTE_DIR}
docker compose build --pull --no-cache ${payload_build}
docker compose up -d --force-recreate ${payload_build}
"

echo "Deployment và build thành công!"
end=$(date +%s)
elapsed=$((end - start))

echo "Thời gian thực thi: $elapsed giây"