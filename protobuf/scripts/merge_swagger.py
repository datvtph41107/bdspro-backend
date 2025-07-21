#!/usr/bin/env python3
"""
Script to merge all swagger JSON files in the docs directory and add Bearer authentication
Usage: python3 merge_swagger.py
"""

import json
import os
import glob
from pathlib import Path
from typing import Dict, Any, List

def add_bearer_auth_to_operation(operation: Dict[str, Any]) -> None:
    """Add Bearer authentication to an operation"""
    # Add security requirement
    operation['security'] = [{'BearerAuth': []}]
    
    # Add Authorization parameter if not already present
    parameters = operation.get('parameters', [])
    has_auth_param = any(
        param.get('name') == 'Authorization' and param.get('in') == 'header'
        for param in parameters
    )
    
    if not has_auth_param:
        auth_param = {
            'name': 'Authorization',
            'in': 'header',
            'description': 'Bearer token for authentication',
            'required': True,
            'type': 'string'
        }
        parameters.append(auth_param)
        operation['parameters'] = parameters

def add_bearer_auth_to_path_item(path_item: Dict[str, Any]) -> None:
    """Add Bearer authentication to all operations in a path item"""
    for method in ['get', 'post', 'put', 'delete', 'patch', 'options', 'head']:
        if method in path_item and path_item[method]:
            add_bearer_auth_to_operation(path_item[method])

def merge_swagger_files() -> None:
    """Merge all swagger files and add Bearer authentication"""
    
    # Initialize merged swagger
    merged_swagger = {
        'swagger': '2.0',
        'info': {
            'title': 'BDS Pro API Documentation',
            'description': 'Complete API documentation for BDS Pro microservices',
            'version': '1.0.0'
        },
        'basePath': '/v2',
        'consumes': ['application/json'],
        'produces': ['application/json'],
        'paths': {},
        'definitions': {},
        'securityDefinitions': {
            'BearerAuth': {
                'type': 'apiKey',
                'name': 'Authorization',
                'in': 'header',
                'description': 'Bearer token for authentication'
            }
        },
        'security': [{'BearerAuth': []}],
        'tags': []
    }
    
    # Find all swagger JSON files
    docs_dir = Path('docs')
    if not docs_dir.exists():
        print("❌ docs directory not found")
        return
    
    swagger_files = list(docs_dir.rglob('*.swagger.json'))
    print(f"📁 Found {len(swagger_files)} swagger files")
    
    # Process each file
    for file_path in swagger_files:
        print(f"🔄 Processing: {file_path}")
        
        try:
            with open(file_path, 'r', encoding='utf-8') as f:
                swagger_data = json.load(f)
            
            # Merge paths and add Bearer auth
            for path, path_item in swagger_data.get('paths', {}).items():
                add_bearer_auth_to_path_item(path_item)
                merged_swagger['paths'][path] = path_item
            
            # Merge definitions
            for name, definition in swagger_data.get('definitions', {}).items():
                merged_swagger['definitions'][name] = definition
            
            # Merge tags (avoid duplicates)
            existing_tags = {tag['name'] for tag in merged_swagger['tags']}
            for tag in swagger_data.get('tags', []):
                if tag['name'] not in existing_tags:
                    merged_swagger['tags'].append(tag)
                    existing_tags.add(tag['name'])
                    
        except Exception as e:
            print(f"⚠️  Error processing {file_path}: {e}")
            continue
    
    # Write merged swagger to file
    output_file = 'merged_swagger.json'
    with open(output_file, 'w', encoding='utf-8') as f:
        json.dump(merged_swagger, f, indent=2, ensure_ascii=False)
    
    # Copy merged swagger to gateway service docs directory
    gateway_docs_dir = Path('../gateway-service/docs')
    try:
        gateway_docs_dir.mkdir(parents=True, exist_ok=True)
        gateway_output_file = gateway_docs_dir / 'merged_swagger.json'
        with open(gateway_output_file, 'w', encoding='utf-8') as f:
            json.dump(merged_swagger, f, indent=2, ensure_ascii=False)
        print(f"📋 Copied merged swagger to {gateway_output_file}")
    except Exception as e:
        print(f"⚠️  Warning: Could not copy to gateway docs: {e}")
    
    print(f"✅ Successfully created {output_file}")
    print(f"📊 Total paths: {len(merged_swagger['paths'])}")
    print(f"📊 Total definitions: {len(merged_swagger['definitions'])}")
    print(f"📊 Total tags: {len(merged_swagger['tags'])}")

if __name__ == '__main__':
    print("🚀 Starting Swagger merge process...")
    merge_swagger_files()
    print("🎉 Merge process completed successfully!") 