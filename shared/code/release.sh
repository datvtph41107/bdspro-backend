#!/bin/bash
start=$(date +%s)

# Kiểm tra xem có tham số đầu vào không
if [ -z "$1" ]; then
  echo "Usage: $0 <project_name>"
  exit 1
fi

# Cấu hình thông tin server
REMOTE_USER="root"
# REMOTE_HOST="14.225.210.29"
REMOTE_HOST="103.172.239.89"
WORKSPACE_NAME=golang-microservice
REMOTE_PATH=/root/$WORKSPACE_NAME

# if [ "$1" = "file" ]; then
#   REMOTE_HOST=14.225.210.29
# fi


if [ "$1" = "compose" ]; then
  echo "đẩy docker compose file lên server..." $REMOTE_PATH
  ssh ${REMOTE_USER}@${REMOTE_HOST} "
mkdir -p ${REMOTE_PATH}
"
  # mkdir -p $WORKSPACE_NAME
  scp docker-compose.yml ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_PATH}
  exit 0
fi

if [ "$1" = "tqd-tool" ]; then
  REMOTE_DIR=$REMOTE_PATH/tqd-service
  scp Dockerfile.tqd-tool ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}/Dockerfile.tqd-tool
  ssh ${REMOTE_USER}@${REMOTE_HOST} "if [ "$1" = "tqd" ]; then
  scp ../shared/code/Dockerfile.tqd.production ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}/Dockerfile
  scp -r template ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}
elif
  mkdir -p ${REMOTE_DIR}
  cd ${REMOTE_DIR}
  docker build -f Dockerfile.tqd-tool -t tippecanoe-builder .
"
  exit 0
fi

cd ../..

if [ "$1" = "nginx" ]; then
ssh ${REMOTE_USER}@${REMOTE_HOST} "
mkdir -p ${REMOTE_PATH}/nginx
"
cd nginx
 scp -r conf ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_PATH}/nginx/conf
 
 ssh ${REMOTE_USER}@${REMOTE_HOST} "
 mkdir -p ${REMOTE_PATH}
cd ${REMOTE_PATH}
docker stop nginx
docker rm nginx
docker compose up -d --build nginx
"
 exit 0
fi

if [ "$1" = "certbot" ]; then
 exit 0
fi
# if [ "$1" = "pmtiles" ]; then
#   echo "đẩy pmtiles file lên server..." $REMOTE_PATH
#   ssh ${REMOTE_USER}@${REMOTE_HOST} "
# mkdir -p ${REMOTE_PATH}/pmtiles
# "
#   cd tqd-service
#   # scp output.pmtiles ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_PATH}/pmtiles/output.pmtiles
#   # scp config.toml ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_PATH}/pmtiles/config.toml
#   # scp Dockerfile.pmtiles ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_PATH}/pmtiles/Dockerfile

#   ssh ${REMOTE_USER}@${REMOTE_HOST} "
# cd ${REMOTE_PATH}/pmtiles
# docker compose build --no-cache tileserver
# docker rm -f tileserver
# docker compose up -d tileserver
# "

#   exit 0
# fi
# if [ "$1" = "assistant" ]; then
#   service_path=assistant-service
#   REMOTE_DIR=${REMOTE_PATH}/${service_path}
  
#   echo "--- PUSHING ALL FILES TO SERVER ---"
#   ssh ${REMOTE_USER}@${REMOTE_HOST} "
#     mkdir -p ${REMOTE_DIR}
#     rm -rf ${REMOTE_DIR}
#   "
  
#   # Push toàn bộ file lên server (trừ các file không cần thiết)
#   rsync -avz --exclude='.git' --exclude='tmp' --exclude='files' --exclude='main' --exclude='assistant-service' \
#     ${service_path}/ ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}/
  
#   # Push shared folder (cần cho go.mod dependencies)
#   rsync -avz --exclude='.git' --exclude='node_modules' \
#     ./shared ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_PATH}
  
#   # Push Dockerfile (sử dụng Dockerfile.develop hoặc tạo Dockerfile build từ source)
#   scp ${service_path}/Dockerfile ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}/Dockerfile
  
#   echo "--- BUILDING IN DOCKER ---"
#   # Build từ context lớn hơn (từ REMOTE_PATH) để có thể truy cập shared folder
#   # Dockerfile cần sửa để copy shared từ ./shared thay vì ../shared
#   ssh ${REMOTE_USER}@${REMOTE_HOST} "
#     cd ${REMOTE_PATH}
#     docker compose -f docker-compose.yml build --no-cache assistant
#     docker compose -f docker-compose.yml up -d assistant
#   "
  
#   echo "Deployment và build thành công!"
#   end=$(date +%s)
#   elapsed=$((end - start))
#   echo "Thời gian thực thi: $elapsed giây"
#   exit 0
# fi

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
    local_path=chat-v1-service
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
  "assistant")
    service_path=assistant-service
    path_jar=main
    payload_build=assistant
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
  "pmtiles")
    service_path=tqd-service
    path_jar=main
    payload_build=pmtiles
    ;;
  "feedback")
    service_path=feedback-service
    path_jar=main
    payload_build=feedback
    ;;
  "transaction")
    service_path=transaction-service
    path_jar=main
    payload_build=transaction
    ;;
  "hub")
    service_path=hub-service
    path_jar=main
    payload_build=hub
    ;;
  *)
    echo "Tham số không hợp lệ."
    exit 1
    ;;
esac

if [ -n "$local_path" ]; then
  cd "$local_path"
else
  cd "$service_path"
fi

REMOTE_DIR=${REMOTE_PATH}/${service_path}

echo "--- GO BUILDING ---"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .
echo "--- FINISH BUILDING ---"

# Đẩy các file trong thư mục hiện tại lên server
echo "Đang đẩy file lên server..."
echo $(pwd)

if [ "$1" = "compose" ]; then
  echo "đẩy docker compose file lên server..."
  scp docker-compose.yml ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_PATH}
  exit 0  
fi

ssh ${REMOTE_USER}@${REMOTE_HOST} "
mkdir -p ${REMOTE_DIR}
mkdir -p ${REMOTE_DIR}/config
"

scp ${path_jar} ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}
# scp ../common/target/common-1.0.jar ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}

scp config/production.yml ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}/config/production.yml
if [ "$1" = "tqd" ]; then
  scp ../shared/code/Dockerfile.tqd.production ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}/Dockerfile.tqd
  scp -r template ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}
else
  scp ../shared/code/Dockerfile.production ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}/Dockerfile
fi

# rsync -avz --delete ./Dockerfile ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}${service_path}

if [ "$1" = "notification" ]; then
  scp config/bdspro-5e599-firebase-adminsdk-fbsvc-5aff064187.json ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}/bdspro-5e599-firebase-adminsdk-fbsvc-5aff064187.json
fi
if [ "$1" = "chat" ]; then
  scp -r configs ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}
fi
CONFIG_COPY_SERVICES=("relay" "org" "payment" "appointment" "assistant")
for service in "${CONFIG_COPY_SERVICES[@]}"; do
  if [ "$1" = "$service" ]; then
    scp -r config ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}
    break
  fi
done

scp -r config ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}

if [ "$1" = "swagger-chat" ]; then
  scp -r proto ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}
fi
if [ "$1" = "gateway" ]; then
  # scp -r docs ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}
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
docker compose build --no-cache ${payload_build}
docker compose up -d ${payload_build}
"
# docker compose build --no-cache ${payload_build}

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