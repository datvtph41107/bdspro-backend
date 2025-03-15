#!/bin/bash

# Kiểm tra xem có tham số đầu vào không
if [ -z "$1" ]; then
  echo "Usage: $0 <project_name>"
  exit 1
fi

# Chọn thư mục dựa vào tham số đầu tiên
service_path=/common
path_jar=''
service_name=$1
case "$1" in
  "gateway")
    service_path=gateway-service
    path_jar=main
    ;;
  "file")
    service_path=file-service
    path_jar=main
    ;;
  "map")
    service_path=map-service
    path_jar=main
    ;;
  "chat")
    service_path=chat-service
    path_jar=main
    ;;
  "payment")
    service_path=payment-service
    path_jar=main
    ;;
  *)
    echo "Tham số không hợp lệ."
    exit 1
    ;;
esac

cd $service_path

echo "--- GO BUILDING ---"
GOOS=linux GOARCH=amd64 go build -o main .
echo "--- FINISH BUILDING ---"

# Cấu hình thông tin server
REMOTE_USER="root"
REMOTE_HOST="14.225.210.29"
REMOTE_PATH=/root/spring-boot-microservice
REMOTE_DIR=${REMOTE_PATH}/${service_path}

# Đẩy các file trong thư mục hiện tại lên server
echo "Đang đẩy file lên server..."
echo $(pwd)

if [ "$1" = "compose" ]; then
  echo "đẩy docker-compose.yml lên server..."
  scp docker-compose.yml ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_PATH}
  exit 0
fi

ssh ${REMOTE_USER}@${REMOTE_HOST} "
mkdir -p ${REMOTE_DIR}
"

scp ${path_jar} ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}
# scp ../common/target/common-1.0.jar ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}
scp Dockerfile ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}/Dockerfile
# rsync -avz --delete ./Dockerfile ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}${service_path}
if [ $? -ne 0 ]; then
  echo "Có lỗi khi đẩy file lên server"
  exit 1
fi

ssh ${REMOTE_USER}@${REMOTE_HOST} "
cd ${REMOTE_DIR}
docker-compose up -d --build ${service_name}
"

# # SSH vào server và chạy lệnh build (giả sử bạn có file build.sh trong thư mục remote)
# echo "Đang chạy build trên server..."
# ssh ${REMOTE_USER}@${REMOTE_HOST} "cd ${REMOTE_DIR} && ./build.sh"
# if [ $? -ne 0 ]; then
#   echo "Build thất bại trên server"
#   exit 1
# fi

echo "Deployment và build thành công!"
