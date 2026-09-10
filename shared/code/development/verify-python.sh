#!/usr/bin/env bash
set -euo pipefail
repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)"
cd "$repo_root"

python3 - <<'PY'
from pathlib import Path
import sys

roots = [Path('ai-service')]
errors = []
count = 0
for root in roots:
    if not root.exists():
        continue
    for path in root.rglob('*.py'):
        if any(part in {'.venv', 'venv', '__pycache__'} for part in path.parts):
            continue
        count += 1
        try:
            compile(path.read_text(encoding='utf-8'), str(path), 'exec')
        except Exception as exc:
            errors.append(f'{path}: {exc}')
if errors:
    print('\n'.join(errors), file=sys.stderr)
    raise SystemExit(1)
print(f'ai-service Python syntax PASS ({count} files)')
PY
