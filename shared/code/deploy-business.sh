#!/bin/bash

# Deploy only business services (excluding core services)
# Usage: ./deploy-business.sh [--parallel]

echo "Deploying business services only..."
./deploy-all.sh --services="bdspro,crm,notification,social,chat,chathttp,file,map,payment,membership,relay,appointment,org,transaction,tqd,marketing,feedback" "$@"

