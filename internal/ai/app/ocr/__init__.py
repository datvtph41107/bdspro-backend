"""
OCR Module for National ID Card Processing
"""

# Lazy imports to avoid errors when dependencies are missing
try:
    from .national_card_ocr import NationalCardOCR
except (ImportError, NameError):
    NationalCardOCR = None

try:
    from .easyocr_service import EasyOCRService
except (ImportError, NameError):
    EasyOCRService = None

__all__ = ['NationalCardOCR', 'EasyOCRService']

