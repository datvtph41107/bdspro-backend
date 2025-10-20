"""
Generate Training Data - Create Q&A pairs from knowledge base
"""

import sys
import json
from pathlib import Path
from typing import List, Dict
import re

# Add parent directory to path
sys.path.insert(0, str(Path(__file__).parent.parent))

def extract_sections(content: str) -> List[Dict]:
    """
    Extract sections from markdown content
    
    Returns:
        List of {title, content, level} dicts
    """
    sections = []
    lines = content.split('\n')
    
    current_section = None
    current_content = []
    
    for line in lines:
        # Check if it's a header
        if line.startswith('#'):
            # Save previous section
            if current_section:
                sections.append({
                    'title': current_section['title'],
                    'level': current_section['level'],
                    'content': '\n'.join(current_content).strip()
                })
            
            # Start new section
            level = len(line) - len(line.lstrip('#'))
            title = line.lstrip('#').strip()
            current_section = {'title': title, 'level': level}
            current_content = []
        else:
            if current_section:
                current_content.append(line)
    
    # Save last section
    if current_section and current_content:
        sections.append({
            'title': current_section['title'],
            'level': current_section['level'],
            'content': '\n'.join(current_content).strip()
        })
    
    return sections

def generate_qa_pairs(sections: List[Dict], source_file: str) -> List[Dict]:
    """
    Generate Q&A pairs from sections
    
    Returns:
        List of {question, answer, source, context} dicts
    """
    qa_pairs = []
    
    for section in sections:
        title = section['title']
        content = section['content']
        
        if not content or len(content) < 50:
            continue
        
        # Generate questions based on title patterns
        questions = []
        
        # Pattern 1: "What is X?"
        if any(keyword in title.lower() for keyword in ['architecture', 'overview', 'structure', 'design']):
            questions.append(f"{title} là gì?")
            questions.append(f"Giải thích về {title}")
        
        # Pattern 2: "How to X?"
        if any(keyword in title.lower() for keyword in ['setup', 'install', 'create', 'implement', 'build']):
            questions.append(f"Làm sao để {title.lower()}?")
            questions.append(f"Cách {title.lower()} như thế nào?")
        
        # Pattern 3: "What are X?"
        if any(keyword in title.lower() for keyword in ['components', 'features', 'services', 'modules']):
            questions.append(f"Có những {title.lower()} nào?")
            questions.append(f"Liệt kê các {title.lower()}")
        
        # Pattern 4: "How does X work?"
        if any(keyword in title.lower() for keyword in ['flow', 'process', 'workflow', 'mechanism']):
            questions.append(f"{title} hoạt động như thế nào?")
            questions.append(f"Giải thích cách hoạt động của {title.lower()}")
        
        # Pattern 5: Generic questions
        questions.append(f"Thông tin về {title}")
        questions.append(f"Cho tôi biết về {title.lower()}")
        
        # Create Q&A pairs
        for question in questions[:3]:  # Limit to 3 questions per section
            qa_pairs.append({
                'question': question,
                'answer': content,
                'source': source_file,
                'section': title,
                'metadata': {
                    'level': section['level'],
                    'type': 'knowledge_base'
                }
            })
    
    return qa_pairs

def generate_code_qa_pairs() -> List[Dict]:
    """
    Generate Q&A pairs for common coding tasks
    """
    code_qa = [
        {
            'question': 'Làm sao tạo API mới trong BDSPro?',
            'answer': '''Để tạo API mới trong BDSPro, làm theo các bước sau:

1. **Định nghĩa trong protobuf** (`protobuf/schema/<service>/`):
```proto
service YourService {
  rpc YourMethod(YourRequest) returns (YourResponse);
}
```

2. **Implement handler** (`infra/handler/`):
```go
// @bind: internal/interface/handler.IYourHandler
func (h *YourHandler) YourMethod(ctx context.Context, req *pb.YourRequest) (*pb.YourResponse, error) {
    // Implementation
}
```

3. **Viết usecase** (`internal/usecase/`):
```go
func (u *YourUsecase) Process(ctx context.Context, dto *YourDTO) (*Result, error) {
    // Business logic
}
```

4. **Tạo repository** (`internal/interface/repo/` và `infra/postgre/`):
```go
type IYourRepo interface {
    Create(ctx context.Context, entity *YourEntity) error
}
```

5. **Đăng ký trong gateway** nếu là service mới.

Chi tiết xem trong documentation về Clean Architecture.''',
            'source': 'generated',
            'section': 'API Development',
            'metadata': {'type': 'code_pattern'}
        },
        {
            'question': 'Clean Architecture trong BDSPro được tổ chức như thế nào?',
            'answer': '''BDSPro tuân theo Clean Architecture với các layer rõ ràng:

**1. config/** - Cấu hình (YAML, ENV)

**2. infra/** - Infrastructure layer:
- `client/`: External service clients
- `db/`: Database connections
- `handler/`: HTTP/gRPC handlers
- `mapper/`: DTO ↔ Entity conversion
- `postgre/`: Repository implementations
- `redis/`: Cache implementations
- `oauth/`: OAuth providers

**3. internal/** - Business logic:
- `domain/`: Entities
- `dto/`: Data Transfer Objects
- `interface/`: Interface definitions (repo, service, provider)
- `usecase/`: Business logic
- `utils/`: Utilities
- `validator/`: Input validation

**4. cmd/** - Entry points

**5. wire/** - Dependency injection

Quy tắc: infra không import usecase, chỉ implement interfaces từ internal/interface.''',
            'source': 'generated',
            'section': 'Clean Architecture',
            'metadata': {'type': 'architecture'}
        },
        {
            'question': 'Làm sao lấy thông tin user từ context?',
            'answer': '''Để lấy thông tin user từ context trong BDSPro:

**Import utils:**
```go
import _utils "common/utils"
```

**Lấy ProfileID:**
```go
profileID := _utils.GetProfileIdWithContext(ctx)
```

**Lấy OrganizationID:**
```go
organizationID := _utils.GetOrganizationIdFromContext(ctx)
```

**Lưu ý:**
- Chỉ gọi trong layer usecase hoặc handler
- Không tự parse JWT hay context
- Utils đã xử lý edge cases và type casting

Thông tin này được set từ middleware authentication ở gateway.''',
            'source': 'generated',
            'section': 'Context Management',
            'metadata': {'type': 'code_pattern'}
        },
        {
            'question': 'Cách sử dụng BaseEntity và AuditBase?',
            'answer': '''**BaseEntity** - Base cho tất cả entities:
```go
type Product struct {
    _models.BaseEntity
    Name  string
    Price float64
}
```

Cung cấp:
- `ID` (uint64)
- `CreatedAt`, `UpdatedAt` (timestamps)
- `DeletedAt` (soft delete)
- `AuditBase` (audit fields)

**AuditBase** - Tracking user actions:
```go
type AuditBase struct {
    CreatedBy *uint64 `json:"createdBy"`
    UpdatedBy *uint64 `json:"updatedBy"`
}
```

**GORM Callbacks tự động:**
- `BeforeCreate`: Set CreatedBy và UpdatedBy
- `BeforeUpdate`: Cập nhật UpdatedBy

Lấy user ID từ context key `profileId` hoặc `ProfileIDKey`.

**Lưu ý:** Phải set profileId vào context trước khi gọi DB.''',
            'source': 'generated',
            'section': 'Entity Management',
            'metadata': {'type': 'code_pattern'}
        },
        {
            'question': 'Cách implement repository với CrudRepo?',
            'answer': '''Trong infra layer, **LUÔN dùng `_provider.CrudRepo[T]`**:

```go
// infra/postgre/product_postgre.go
type ProductRepo struct {
    _provider.CrudRepo[_models.Product]
}

func NewProductRepo(db *_db.TransactionRepo) _repo.IProductRepo {
    repo := &ProductRepo{}
    repo.Init(repo, db)  // QUAN TRỌNG!
    return repo
}
```

**Override hooks khi cần:**
```go
func (r *ProductRepo) BeforeSave(c context.Context, id *uint64, entity *_models.Product) error {
    // Custom validation
    return nil
}

func (r *ProductRepo) AfterSave(c context.Context, id *uint64, entity *_models.Product) error {
    // Cache invalidation, audit, etc.
    return nil
}
```

**CrudRepo cung cấp:**
- Create, Update, Delete
- GetByID, GetDetail, GetAll
- GetList (với phân trang)
- GetListByIDs
- Transaction support tự động

Không implement CRUD từ đầu!''',
            'source': 'generated',
            'section': 'Repository Pattern',
            'metadata': {'type': 'code_pattern'}
        }
    ]
    
    return code_qa

def main():
    """Generate training data from knowledge base"""
    
    print("🤖 Training Data Generator")
    print("="*60)
    
    # Output paths
    train_dir = Path(__file__).parent.parent / "data" / "train"
    train_dir.mkdir(parents=True, exist_ok=True)
    
    output_file = train_dir / "qa_dataset.json"
    
    # Knowledge base directory
    kb_dir = Path(__file__).parent.parent / "knowledge_base"
    
    all_qa_pairs = []
    
    # Process markdown files
    if kb_dir.exists():
        print(f"\n📂 Processing knowledge base: {kb_dir}")
        
        for md_file in sorted(kb_dir.glob("*.md")):
            print(f"\n📄 {md_file.name}")
            
            # Read content
            with open(md_file, 'r', encoding='utf-8') as f:
                content = f.read()
            
            # Extract sections
            sections = extract_sections(content)
            print(f"   ✂️  Found {len(sections)} sections")
            
            # Generate Q&A pairs
            qa_pairs = generate_qa_pairs(sections, md_file.name)
            all_qa_pairs.extend(qa_pairs)
            print(f"   ✅ Generated {len(qa_pairs)} Q&A pairs")
    
    # Add code-specific Q&A
    print(f"\n📝 Adding code-specific Q&A pairs...")
    code_qa = generate_code_qa_pairs()
    all_qa_pairs.extend(code_qa)
    print(f"   ✅ Added {len(code_qa)} code Q&A pairs")
    
    # Save to JSON
    print(f"\n💾 Saving to: {output_file}")
    with open(output_file, 'w', encoding='utf-8') as f:
        json.dump(all_qa_pairs, f, ensure_ascii=False, indent=2)
    
    # Summary
    print(f"\n{'='*60}")
    print(f"✅ DONE!")
    print(f"   Total Q&A pairs: {len(all_qa_pairs)}")
    print(f"   Output: {output_file}")
    print(f"{'='*60}\n")
    
    # Show sample
    print("📋 Sample Q&A pair:")
    print(json.dumps(all_qa_pairs[0], ensure_ascii=False, indent=2))

if __name__ == "__main__":
    try:
        main()
    except Exception as e:
        print(f"\n❌ Error: {e}")
        import traceback
        traceback.print_exc()
        sys.exit(1)

