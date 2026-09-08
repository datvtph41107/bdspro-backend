#!/usr/bin/env python3
from __future__ import annotations

import re
from pathlib import Path

ROOT = Path(__file__).resolve().parents[3]
CANONICAL = [
    ROOT / 'README.md',
    ROOT / 'CONTRIBUTING.md',
    ROOT / 'BACKEND-ARCHITECTURE.md',
    ROOT / 'BACKEND-ACCEPTANCE.md',
    ROOT / 'documents' / 'DEVELOPER-OPERATING-GUIDE.md',
    ROOT / 'documents' / 'DEVELOPMENT-ARCHITECTURE-STANDARD.md',
    ROOT / 'documents' / 'SERVICE-MAP.md',
    ROOT / 'documents' / 'MIGRATIONS.md',
    ROOT / 'documents' / 'OPERATIONAL-MAPS.md',
    ROOT / 'documents' / 'FINAL-ACCEPTANCE-20260904.md',
    ROOT / 'documents' / 'IMPORT-CHECKLIST.md',
    ROOT / 'documents' / 'CHANGE-SUMMARY-20260904.md',
]
CANONICAL.extend(sorted(ROOT.glob('*-service/README.md')))

link_re = re.compile(r'(?<!!)\[[^\]]+\]\(([^)]+)\)')
errors: list[str] = []
checked = 0
for doc in CANONICAL:
    if not doc.exists():
        errors.append(f'missing canonical doc: {doc.relative_to(ROOT)}')
        continue
    checked += 1
    text = doc.read_text(encoding='utf-8')
    for raw in link_re.findall(text):
        target = raw.strip().split()[0].strip('<>')
        if not target or target.startswith(('#', 'http://', 'https://', 'mailto:')):
            continue
        path_text = target.split('#', 1)[0]
        if not path_text:
            continue
        resolved = (doc.parent / path_text).resolve()
        try:
            resolved.relative_to(ROOT.resolve())
        except ValueError:
            errors.append(f'{doc.relative_to(ROOT)}: link escapes repo: {target}')
            continue
        if not resolved.exists():
            errors.append(f'{doc.relative_to(ROOT)}: broken link: {target}')

if errors:
    raise SystemExit('\n'.join(errors))
print(f'canonical markdown links PASS ({checked} docs checked)')
