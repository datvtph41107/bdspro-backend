#!/bin/bash

# ============================================================================
# Deploy script for HTTP Server
# ============================================================================
# This script:
# 1. Syncs code to server using rsync
# 2. Builds Docker image on server
# 3. Stops old container
# 4. Runs new container
# 5. Performs health check
# ============================================================================

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# ============================================================================
# Configuration
# ============================================================================

# Server configuration
SERVER_HOST="14.225.210.29"
SERVER_USER="root"
SERVER_PATH=/root/backend-microservice/ai-service
SSH_KEY="${SSH_KEY:-}"  # Optional: path to SSH key

# Service configuration
SERVICE_NAME="http-server"
CONTAINER_NAME="http-server"
DOCKERFILE="Dockerfile.http"
PORT="8218"

# ============================================================================
# Functions
# ============================================================================

print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

check_requirements() {
    print_info "Checking requirements..."
    
    # Check if rsync is installed
    if ! command -v rsync &> /dev/null; then
        print_error "rsync is not installed. Please install it first."
        exit 1
    fi
    
    # Check if SSH key exists (if provided)
    # if [ -n "$SSH_KEY" ] && [ ! -f "$SSH_KEY" ]; then
    #     print_error "SSH key not found: $SSH_KEY"
    #     exit 1
    # fi
    
    print_info "Requirements check passed ✓"
}

build_rsync_command() {
    local RSYNC_CMD="rsync -avz --delete"
    
    # Add SSH key if provided
    # if [ -n "$SSH_KEY" ]; then
    #     RSYNC_CMD="$RSYNC_CMD -e \"ssh -i $SSH_KEY\""
    # fi
    
    # Exclude patterns (without quotes for proper rsync handling)
    RSYNC_CMD="$RSYNC_CMD --exclude=.git"
    RSYNC_CMD="$RSYNC_CMD --exclude=__pycache__"
    RSYNC_CMD="$RSYNC_CMD --exclude=*.pyc"
    RSYNC_CMD="$RSYNC_CMD --exclude=venv"
    RSYNC_CMD="$RSYNC_CMD --exclude=venv/"
    RSYNC_CMD="$RSYNC_CMD --exclude=.venv"
    RSYNC_CMD="$RSYNC_CMD --exclude=.venv/"
    RSYNC_CMD="$RSYNC_CMD --exclude=env"
    RSYNC_CMD="$RSYNC_CMD --exclude=env/"
    RSYNC_CMD="$RSYNC_CMD --exclude=context"
    RSYNC_CMD="$RSYNC_CMD --exclude=.env"
    RSYNC_CMD="$RSYNC_CMD --exclude=zbar"
    RSYNC_CMD="$RSYNC_CMD --exclude=uploads/*"
    RSYNC_CMD="$RSYNC_CMD --exclude=data/chroma_db/*"
    RSYNC_CMD="$RSYNC_CMD --exclude=logs/*"
    RSYNC_CMD="$RSYNC_CMD --exclude=*.ipynb"
    RSYNC_CMD="$RSYNC_CMD --exclude=.idea"
    RSYNC_CMD="$RSYNC_CMD --exclude=files"
    RSYNC_CMD="$RSYNC_CMD --exclude=.vscode"
    RSYNC_CMD="$RSYNC_CMD --exclude=node_modules"
    RSYNC_CMD="$RSYNC_CMD --exclude=*.pyc"
    RSYNC_CMD="$RSYNC_CMD --exclude=.DS_Store"
    RSYNC_CMD="$RSYNC_CMD --exclude=*.swp"
    RSYNC_CMD="$RSYNC_CMD --exclude=*.swo"
    
    echo "$RSYNC_CMD"
}

build_ssh_command() {
    local SSH_CMD="ssh"
    
    if [ -n "$SSH_KEY" ]; then
        SSH_CMD="$SSH_CMD -i $SSH_KEY"
    fi
    
    SSH_CMD="$SSH_CMD $SERVER_USER@$SERVER_HOST"
    
    echo "$SSH_CMD"
}

sync_code() {
    print_info "Syncing code to server..."
    
    local RSYNC_CMD=$(build_rsync_command)
    local SOURCE="./"
    local DEST="$SERVER_USER@$SERVER_HOST:$SERVER_PATH/"
    
    if [ -n "$SSH_KEY" ]; then
        eval "$RSYNC_CMD $SOURCE $DEST"
    else
        eval "$RSYNC_CMD $SOURCE $DEST"
    fi
    
    # Verify code was synced by checking a key file
    local SSH_CMD=$(build_ssh_command)
    if $SSH_CMD "test -f $SERVER_PATH/http_server.py"; then
        print_info "Code synced successfully ✓"
        # Show file modification time to verify
        local FILE_TIME=$($SSH_CMD "stat -c %y $SERVER_PATH/http_server.py 2>/dev/null || stat -f %Sm $SERVER_PATH/http_server.py 2>/dev/null")
        print_info "http_server.py last modified: $FILE_TIME"
    else
        print_error "Code sync failed! http_server.py not found on server."
        exit 1
    fi
}

create_directories_on_server() {
    print_info "Creating necessary directories on server..."
    
    local SSH_CMD=$(build_ssh_command)
    $SSH_CMD "rm -rf $SERVER_PATH"
    
    # Create directories if they don't exist (don't delete existing data)
    $SSH_CMD "mkdir -p $SERVER_PATH/uploads"
    $SSH_CMD "mkdir -p $SERVER_PATH/data"
    $SSH_CMD "mkdir -p $SERVER_PATH/logs"
    
    print_info "Directories created ✓"
}

stop_old_container() {
    print_info "Stopping and removing old container..."
    
    local SSH_CMD=$(build_ssh_command)
    
    # Stop and remove old container if exists
    $SSH_CMD "cd $SERVER_PATH && \
        (docker-compose stop $SERVICE_NAME 2>/dev/null || true) && \
        (docker-compose rm -f $SERVICE_NAME 2>/dev/null || true) && \
        (docker stop $CONTAINER_NAME 2>/dev/null || true) && \
        (docker rm -f $CONTAINER_NAME 2>/dev/null || true)"
    
    print_info "Old container stopped and removed ✓"
}

build_docker_image() {
    print_info "Building Docker image on server (no cache)..."
    
    local SSH_CMD=$(build_ssh_command)
    local BUILD_TIMESTAMP=$(date +%s)
    
    # Prune build cache and remove old images
    print_info "Cleaning up Docker cache and old images..."
    $SSH_CMD "cd $SERVER_PATH && \
        docker builder prune -f && \
        (docker rmi $SERVICE_NAME:latest 2>/dev/null || true) && \
        (docker images | grep $SERVICE_NAME | awk '{print \$3}' | xargs -r docker rmi -f 2>/dev/null || true)"
    
    # Build with no cache and build arg to ensure fresh build
    print_info "Building fresh Docker image..."
    $SSH_CMD "cd $SERVER_PATH && \
        docker build \
            --no-cache \
            --pull \
            --build-arg BUILD_TIMESTAMP=$BUILD_TIMESTAMP \
            -f $DOCKERFILE \
            -t $SERVICE_NAME:latest \
            ."
    
    # Verify image was created
    if $SSH_CMD "docker images | grep -q $SERVICE_NAME"; then
        print_info "Docker image built successfully ✓"
        local IMAGE_ID=$($SSH_CMD "docker images $SERVICE_NAME:latest --format '{{.ID}}'")
        print_info "Image ID: $IMAGE_ID"
    else
        print_error "Docker image build failed!"
        exit 1
    fi
}

start_container() {
    print_info "Starting container..."
    
    local SSH_CMD=$(build_ssh_command)
    
    # Use docker-compose if available, otherwise use docker run
    $SSH_CMD "cd $SERVER_PATH && \
        if [ -f docker-compose.yml ]; then
            docker-compose build --no-cache $SERVICE_NAME && \
            docker-compose up -d --force-recreate --no-deps $SERVICE_NAME
        else
            docker run -d \
                --name $CONTAINER_NAME \
                --restart unless-stopped \
                -p $PORT:$PORT \
                -v $SERVER_PATH/uploads:/app/uploads \
                -v $SERVER_PATH/data:/app/data \
                -v $SERVER_PATH/logs:/app/logs \
                -e HTTP_PORT=$PORT \
                $SERVICE_NAME:latest
        fi"
    
    # Verify container is running
    sleep 2
    if $SSH_CMD "docker ps | grep -q $CONTAINER_NAME"; then
        print_info "Container started successfully ✓"
        # Show container info
        local CONTAINER_ID=$($SSH_CMD "docker ps | grep $CONTAINER_NAME | awk '{print \$1}'")
        print_info "Container ID: $CONTAINER_ID"
    else
        print_error "Container failed to start!"
        show_logs
        exit 1
    fi
}

show_logs() {
    print_info "Showing container logs..."
    
    local SSH_CMD=$(build_ssh_command)
    
    $SSH_CMD "cd $SERVER_PATH && \
        if [ -f docker-compose.yml ]; then
            docker-compose logs --tail=50 $SERVICE_NAME
        else
            docker logs --tail=50 $CONTAINER_NAME
        fi"
}

show_status() {
    print_info "Container status:"
    
    local SSH_CMD=$(build_ssh_command)
    
    $SSH_CMD "cd $SERVER_PATH && \
        if [ -f docker-compose.yml ]; then
            docker-compose ps $SERVICE_NAME
        else
            docker ps -a | grep $CONTAINER_NAME
        fi"
}

# ============================================================================
# Main deployment flow
# ============================================================================

verify_code_in_container() {
    print_info "Verifying code in container..."
    
    local SSH_CMD=$(build_ssh_command)
    
    # Get file hash from server
    local SERVER_HASH=$($SSH_CMD "md5sum $SERVER_PATH/http_server.py 2>/dev/null | awk '{print \$1}' || md5 $SERVER_PATH/http_server.py 2>/dev/null | awk '{print \$4}'")
    
    # Get file hash from container
    local CONTAINER_HASH=$($SSH_CMD "docker exec $CONTAINER_NAME md5sum /app/http_server.py 2>/dev/null | awk '{print \$1}' || docker exec $CONTAINER_NAME md5 /app/http_server.py 2>/dev/null | awk '{print \$4}'")
    
    if [ -n "$SERVER_HASH" ] && [ -n "$CONTAINER_HASH" ]; then
        if [ "$SERVER_HASH" = "$CONTAINER_HASH" ]; then
            print_info "Code verification passed ✓ (Hash: ${SERVER_HASH:0:8}...)"
        else
            print_warn "Code hash mismatch!"
            print_warn "Server:  ${SERVER_HASH:0:8}..."
            print_warn "Container: ${CONTAINER_HASH:0:8}..."
            print_warn "Container may be using old code. Consider rebuilding."
        fi
    else
        print_warn "Could not verify code hash (md5sum/md5 not available)"
    fi
}

main() {
    print_info "Starting deployment..."
    print_info "Server: $SERVER_USER@$SERVER_HOST"
    print_info "Path: $SERVER_PATH"
    print_info "Service: $SERVICE_NAME"
    echo ""
    
    # Check requirements
    check_requirements
    
    # Create directories
    create_directories_on_server
    
    # Sync code first (ensure latest code is on server)
    sync_code
    
    # Stop old container
    stop_old_container
    
    # Build Docker image (no cache to ensure fresh build)
    build_docker_image
    
    # Start container
    start_container
    
    # Wait a bit for container to start
    print_info "Waiting for container to start..."
    sleep 5
    
    # Verify container is running
    local SSH_CMD=$(build_ssh_command)
    if $SSH_CMD "docker ps | grep -q $CONTAINER_NAME"; then
        print_info "Container is running ✓"
        
        # Verify code in container matches server
        verify_code_in_container
    else
        print_error "Container failed to start!"
        show_logs
        exit 1
    fi
}

# ============================================================================
# Script execution
# ============================================================================

# Parse command line arguments
case "${1:-deploy}" in
    deploy)
        main
        ;;
    logs)
        show_logs
        ;;
    status)
        show_status
        ;;
    sync)
        sync_code
        ;;
    build)
        build_docker_image
        ;;
    restart)
        stop_old_container
        start_container
        ;;
    *)
        echo "Usage: $0 [deploy|logs|status|sync|build|restart]"
        echo ""
        echo "Commands:"
        echo "  deploy   - Full deployment (default)"
        echo "  logs     - Show container logs"
        echo "  status   - Show container status"
        echo "  sync     - Sync code only"
        echo "  build    - Build Docker image only"
        echo "  restart  - Restart container"
        echo ""
        echo "Environment variables:"
        echo "  SERVER_HOST  - Server hostname or IP (default: your-server.com)"
        echo "  SERVER_USER  - SSH user (default: root)"
        echo "  SERVER_PATH  - Deployment path (default: /opt/ai-service)"
        echo "  SSH_KEY      - Path to SSH private key (optional)"
        echo ""
        echo "Example:"
        echo "  SERVER_HOST=192.168.1.100 SERVER_USER=deploy ./deploy.sh"
        exit 1
        ;;
esac

