#!/bin/bash

# Deploy only core services (auth, user, gateway)
# Usage: ./deploy-core.sh

echo "Deploying core services only..."
./deploy-all.sh --services="auth,user,user-grpc,gateway"

