#!/bin/bash

echo "Updating routes files..."

# List of route files to update
route_files=(
    "infra/routes/follow.go"
    "infra/routes/friend.go"
    "infra/routes/group.go"
    "infra/routes/profession.go"
    "infra/routes/profile.go"
)

for file in "${route_files[@]}"; do
    if [ -f "$file" ]; then
        echo "Updating $file..."
        
        # Update imports
        sed -i '' 's|services "user/internal/services/user"|"user/internal/usecases"|g' "$file"
        
        # Update struct fields and constructor parameters
        sed -i '' 's/\*services\./\*usecases\./g' "$file"
        sed -i '' 's/services\./usecases\./g' "$file"
        sed -i '' 's/Service/Usecase/g' "$file"
        
        # Update field names
        sed -i '' 's/service:/usecase:/g' "$file"
        sed -i '' 's/\.service\./\.usecase\./g' "$file"
    fi
done

echo "Routes update completed!"
