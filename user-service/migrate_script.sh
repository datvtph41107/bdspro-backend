#!/bin/bash

# Script to migrate services to usecases and repo to postgres

echo "Starting migration..."

# Create usecases directory if not exists
mkdir -p internal/usecases

# Create postgres directory if not exists  
mkdir -p infra/postgres

# List of files to migrate
services_files=(
    "internal/services/user/certification.go"
    "internal/services/user/friend.go" 
    "internal/services/user/group.go"
    "internal/services/user/profession.go"
    "internal/services/user/profile.go"
    "internal/services/user/report.go"
)

repo_files=(
    "infra/repo/admin.go"
    "infra/repo/certification.go"
    "infra/repo/contact.go"
    "infra/repo/friend.go"
    "infra/repo/group.go"
    "infra/repo/profession.go"
    "infra/repo/profile.go"
    "infra/repo/rate.go"
    "infra/repo/report.go"
)

echo "Migrating services to usecases..."
for file in "${services_files[@]}"; do
    if [ -f "$file" ]; then
        # Extract filename without path
        filename=$(basename "$file")
        # Convert to usecase filename
        usecase_name=$(echo "$filename" | sed 's/\.go$/_usecase.go/')
        
        echo "Processing $file -> internal/usecases/$usecase_name"
        
        # Copy file and modify package name and imports
        cp "$file" "internal/usecases/$usecase_name"
        
        # Replace package name
        sed -i '' 's/package user_service/package usecases/' "internal/usecases/$usecase_name"
        
        # Replace imports
        sed -i '' 's|user_repo "user/infra/repo"|"user/infra/postgres"|g' "internal/usecases/$usecase_name"
        sed -i '' 's/\*user_repo\./\*postgres\./g' "internal/usecases/$usecase_name"
        sed -i '' 's/user_repo\./postgres\./g' "internal/usecases/$usecase_name"
        
        # Replace service struct names and constructor names
        struct_name=$(grep -o 'type [A-Z][a-zA-Z]*Service struct' "$file" | sed 's/type //' | sed 's/ struct//')
        if [ ! -z "$struct_name" ]; then
            usecase_struct=$(echo "$struct_name" | sed 's/Service$/Usecase/')
            sed -i '' "s/$struct_name/$usecase_struct/g" "internal/usecases/$usecase_name"
            
            # Replace constructor
            constructor_name=$(echo "$struct_name" | sed 's/^/New/')
            new_constructor=$(echo "$usecase_struct" | sed 's/^/New/')
            sed -i '' "s/$constructor_name/$new_constructor/g" "internal/usecases/$usecase_name"
        fi
    fi
done

echo "Migrating repos to postgres..."
for file in "${repo_files[@]}"; do
    if [ -f "$file" ]; then
        # Extract filename without path
        filename=$(basename "$file")
        # Convert to postgres filename
        postgres_name=$(echo "$filename" | sed 's/\.go$/_postgres.go/')
        
        echo "Processing $file -> infra/postgres/$postgres_name"
        
        # Copy file and modify package name
        cp "$file" "infra/postgres/$postgres_name"
        
        # Replace package name
        sed -i '' 's/package repo/package postgres/' "infra/postgres/$postgres_name"
        
        # Replace struct names and constructor names
        struct_name=$(grep -o 'type [A-Z][a-zA-Z]*Repo struct' "$file" | sed 's/type //' | sed 's/ struct//')
        if [ ! -z "$struct_name" ]; then
            postgres_struct=$(echo "$struct_name" | sed 's/Repo$/Postgres/')
            sed -i '' "s/$struct_name/$postgres_struct/g" "infra/postgres/$postgres_name"
            
            # Replace constructor
            constructor_name=$(echo "$struct_name" | sed 's/^/New/')
            new_constructor=$(echo "$postgres_struct" | sed 's/^/New/')
            sed -i '' "s/$constructor_name/$new_constructor/g" "infra/postgres/$postgres_name"
        fi
    fi
done

echo "Migration completed!"
