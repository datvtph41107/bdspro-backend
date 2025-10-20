#!/bin/bash
start=$(date +%s)

# Kiểm tra xem có tham số đầu vào không
if [ -z "$1" ]; then
  echo "Usage: $0 <project_name>"
  exit 1
fi

# Cấu hình thông tin server
REMOTE_USER="root"
REMOTE_HOST="14.225.205.222"
WORKSPACE_NAME=backend-microservice
REMOTE_PATH=/root/$WORKSPACE_NAME

# if [ "$1" = "file" ]; then
#   REMOTE_HOST=14.225.210.29
# fi


if [ "$1" = "compose" ]; then
  REMOTE_HOST=14.225.205.222
  echo "đẩy docker-compose.yml lên server..." $REMOTE_PATH
  ssh ${REMOTE_USER}@${REMOTE_HOST} "
mkdir -p ${REMOTE_PATH}
"
  # mkdir -p $WORKSPACE_NAME
  scp docker-compose.yml ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_PATH}
  exit 0
fi

cd ../..

# Chọn thư mục dựa vào tham số đầu tiên
service_path=/common
path_jar=''
service_name=$1
payload_build=$1
case "$1" in
  "gateway")
    service_path=gateway-service
    path_jar=main
    ;;
  "user")
    service_path=user-service
    path_jar=main
    payload_build=user
    ;;
  "auth")
    service_path=auth-service
    path_jar=main
    payload_build=auth
    ;;
  "user-grpc")
    service_path=user-service
    path_jar=main
    payload_build=user-grpc
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
  "notification")
    service_path=notification-service
    path_jar=main
    ;;
  "crm")
    service_path=crm-service
    path_jar=main
    ;;
  "membership")
    service_path=membership-service
    path_jar=main
    ;;
  "bdspro")
    service_path=bdspro-service
    path_jar=main
    ;;
  "bdspro-grpc")
    service_path=bdspro-service
    path_jar=main
    payload_build=bdspro-grpc
    ;;
  "chathttp")
    service_path=chat-service
    path_jar=main
    payload_build=http-server
    ;;
  "relay")
    service_path=relay-service
    path_jar=main
    payload_build=relay
    ;;
  "swagger-chat")
    service_path=chat-service
    path_jar=main
    payload_build=swagger-chat
    ;;
  "social")
    service_path=social-service
    path_jar=main
    payload_build=social
    ;;
  "appointment")
    service_path=appointment-service
    path_jar=main
    payload_build=appointment
    ;;
  "org")
    service_path=organization-service
    path_jar=main
    payload_build=organization
    ;;
  "tqd")
    service_path=tqd-service
    path_jar=main
    payload_build=tqd
    ;;
  "feedback")
    service_path=feedback-service
    path_jar=main
    payload_build=feedback
    ;;
  "marketing")
    service_path=marketing-service
    path_jar=main
    payload_build=marketing
    ;;
  "transaction")
    service_path=transaction-service
    path_jar=main
    payload_build=transaction
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

if [ "$1" = "compose" ]; then
  echo "đẩy docker-compose.yml lên server..."
  scp docker-compose.yml ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_PATH}
  exit 0  
fi

ssh ${REMOTE_USER}@${REMOTE_HOST} "
mkdir -p ${REMOTE_DIR}
mkdir -p ${REMOTE_DIR}/config
"

scp ${path_jar} ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}
# scp ../common/target/common-1.0.jar ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}

scp config/develop.yml ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}/config/develop.yml
scp ../shared/code/Dockerfile.develop ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}/Dockerfile

# rsync -avz --delete ./Dockerfile ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}${service_path}

if [ "$1" = "notification" ]; then
  scp config/bdspro-5e599-firebase-adminsdk-fbsvc-5aff064187.json ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}/bdspro-5e599-firebase-adminsdk-fbsvc-5aff064187.json
fi
if [ "$1" = "chat" ]; then
  scp -r configs ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}
fi
if [ "$1" = "relay" ]; then
  scp -r config ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}
fi
if [ "$1" = "org" ]; then
  scp -r config ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}
fi
if [ "$1" = "payment" ]; then
  scp -r config ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}
fi
if [ "$1" = "appointment" ]; then
  scp -r config ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}
fi

if [ "$1" = "swagger-chat" ]; then
  scp -r proto ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}
fi
if [ "$1" = "gateway" ]; then
  scp -r docs ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}
  scp Dockerfile ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}/Dockerfile
fi
if [ "$1" = "file" ]; then
  scp -r docs ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}
fi
if [ $? -ne 0 ]; then
  echo "Có lỗi khi đẩy file lên server"
  exit 1
fi



ssh ${REMOTE_USER}@${REMOTE_HOST} "
cd ${REMOTE_DIR}
docker-compose build --no-cache ${payload_build}
docker-compose up -d ${payload_build}
"
# docker-compose build --no-cache ${payload_build}

# # SSH vào server và chạy lệnh build (giả sử bạn có file build.sh trong thư mục remote)
# echo "Đang chạy build trên server..."
# ssh ${REMOTE_USER}@${REMOTE_HOST} "cd ${REMOTE_DIR} && ./build.sh"
# if [ $? -ne 0 ]; then
#   echo "Build thất bại trên server"
#   exit 1
# fi

echo "Deployment và build thành công!"
end=$(date +%s)
elapsed=$((end - start))

echo "Thời gian thực thi: $elapsed giây"