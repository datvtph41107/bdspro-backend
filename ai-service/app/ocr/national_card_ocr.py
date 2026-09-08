"""
Module để đọc MRZ và OCR text từ ảnh căn cước công dân
"""

from __future__ import annotations  # Enable postponed evaluation of annotations

import io
import re
import time
from typing import Dict, Optional, Tuple, TYPE_CHECKING, Any
from urllib.request import urlopen
from urllib.parse import urlparse

# Define Image type stub to avoid NameError during class definition
if TYPE_CHECKING:
    from PIL import Image
else:
    # Create a dummy class to avoid NameError
    class Image:
        @staticmethod
        def open(*args, **kwargs):
            raise ImportError("PIL/Pillow is not installed")

try:
    import pytesseract
    from PIL import Image as PILImage
    Image = PILImage  # Use real Image
    TESSERACT_AVAILABLE = True
except ImportError:
    TESSERACT_AVAILABLE = False
    # Image is already defined as dummy class above
    print("⚠️  Warning: pytesseract or PIL not installed. OCR will not work.")
    print("   Install with: pip install pytesseract pillow")

try:
    from readmrz import MrzDetector, MrzReader
    READMRZ_AVAILABLE = True
except ImportError:
    READMRZ_AVAILABLE = False
    print("⚠️  Warning: readmrz not installed. MRZ reading will use fallback method.")
    print("   Install with: pip install readmrz")


class NationalCardOCR:
    """Class để xử lý OCR cho căn cước công dân"""
    
    def __init__(self):
        """Khởi tạo OCR engine"""
        if not TESSERACT_AVAILABLE:
            raise ImportError("pytesseract and PIL are required for OCR")
        
        # Cấu hình Tesseract
        # Có thể set path nếu Tesseract không ở PATH
        # pytesseract.pytesseract.tesseract_cmd = r'/usr/local/bin/tesseract'
        
        # Khởi tạo MRZ detector và reader nếu có
        self.mrz_detector = None
        self.mrz_reader = None
        if READMRZ_AVAILABLE:
            try:
                self.mrz_detector = MrzDetector()
                self.mrz_reader = MrzReader()
            except Exception as e:
                print(f"⚠️  Warning: Failed to initialize readmrz: {e}")
                self.mrz_detector = None
                self.mrz_reader = None
    
    def download_image(self, image_url: str) -> "Image.Image":
        """Download ảnh từ URL và convert sang PIL Image"""
        try:
            with urlopen(image_url) as response:
                image_data = response.read()
            image = Image.open(io.BytesIO(image_data))
            return image
        except Exception as e:
            raise ValueError(f"Failed to download image from {image_url}: {str(e)}")
    
    def extract_mrz(self, back_image: "Image.Image") -> Optional[str]:
        """
        Trích xuất MRZ từ ảnh mặt sau
        
        MRZ thường nằm ở phần dưới cùng của ảnh (20-30% dưới)
        """
        try:
            # Nếu có readmrz, dùng nó
            if self.mrz_detector and self.mrz_reader:
                try:
                    # Crop phần dưới của ảnh (MRZ thường ở đó)
                    width, height = back_image.size
                    crop_y = int(height * 0.7)  # Bắt đầu từ 70% chiều cao
                    mrz_region = back_image.crop((0, crop_y, width, height))
                    
                    # Detect và đọc MRZ
                    cropped = self.mrz_detector.crop_area(mrz_region)
                    if cropped is not None:
                        result = self.mrz_reader.process(cropped)
                        if result and 'mrz' in result:
                            return result['mrz']
                except Exception as e:
                    print(f"⚠️  readmrz failed, using fallback: {e}")
            
            # Fallback: Dùng Tesseract với config đặc biệt cho MRZ
            width, height = back_image.size
            crop_y = int(height * 0.7)
            mrz_region = back_image.crop((0, crop_y, width, height))
            
            # OCR với config cho MRZ
            # PSM 6: Assume a single uniform block of text
            # Whitelist: Chỉ cho phép A-Z, 0-9, <
            mrz_text = pytesseract.image_to_string(
                mrz_region,
                config='--psm 6 -c tessedit_char_whitelist=ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789<',
                lang='eng'
            )
            
            # Clean và validate MRZ
            mrz_lines = []
            for line in mrz_text.split('\n'):
                line = line.strip()
                # MRZ thường có ít nhất 25 ký tự và nhiều ký tự <
                if len(line) >= 25 and line.count('<') > 3:
                    mrz_lines.append(line)
            
            if len(mrz_lines) >= 2:
                return '\n'.join(mrz_lines[:2])  # MRZ thường có 2 dòng
            
            return None
            
        except Exception as e:
            print(f"⚠️  Error extracting MRZ: {e}")
            return None
    
    def parse_mrz(self, mrz: str) -> Dict[str, Optional[str]]:
        """
        Parse MRZ string thành thông tin
        
        MRZ format cho căn cước Việt Nam thường có 2 dòng:
        - Dòng 1: ID number, date of birth, etc.
        - Dòng 2: Expiry date, etc.
        """
        info = {}
        
        if not mrz:
            return info
        
        lines = mrz.strip().split('\n')
        if len(lines) < 2:
            return info
        
        line1 = lines[0]
        line2 = lines[1]
        
        # Extract ID number (thường ở đầu dòng 1, 9-12 ký tự)
        id_match = re.search(r'[A-Z0-9]{9,12}', line1)
        if id_match:
            info['id_number'] = id_match.group()
        
        # Extract date of birth (format: YYMMDD)
        dob_match = re.search(r'\d{6}', line1)
        if dob_match:
            dob = dob_match.group()
            if len(dob) == 6:
                year = "20" + dob[0:2]
                month = dob[2:4]
                day = dob[4:6]
                info['date_of_birth'] = f"{day}/{month}/{year}"
        
        # Extract expiry date từ line 2
        expiry_match = re.search(r'\d{6}', line2)
        if expiry_match:
            expiry = expiry_match.group()
            if len(expiry) == 6:
                year = "20" + expiry[0:2]
                month = expiry[2:4]
                day = expiry[4:6]
                info['expiry_date'] = f"{day}/{month}/{year}"
        
        return info
    
    def ocr_image(self, image: "Image.Image", lang: str = 'vie+eng') -> str:
        """OCR ảnh và trả về text"""
        try:
            text = pytesseract.image_to_string(image, lang=lang)
            return text
        except Exception as e:
            raise ValueError(f"OCR failed: {str(e)}")
    
    def parse_front_side(self, text: str) -> Dict[str, Optional[str]]:
        """Parse thông tin từ text đã OCR của mặt trước"""
        info = {}
        lines = text.split('\n')
        
        for line in lines:
            line = line.strip()
            if not line:
                continue
            
            # Tìm số CMND/CCCD
            id_match = re.search(r'\d{9,12}', line)
            if id_match and 'id_number' not in info:
                info['id_number'] = id_match.group()
            
            # Tìm họ tên
            if ('họ' in line.lower() and 'tên' in line.lower()) or 'họ và tên' in line.lower():
                parts = line.split(':')
                if len(parts) > 1:
                    name = parts[1].strip()
                    if name:
                        info['full_name'] = name
            
            # Tìm ngày sinh
            if 'ngày sinh' in line.lower() or 'sinh ngày' in line.lower():
                date_match = re.search(r'(\d{1,2})[/-](\d{1,2})[/-](\d{4})', line)
                if date_match:
                    day, month, year = date_match.groups()
                    info['date_of_birth'] = f"{day}/{month}/{year}"
            
            # Tìm giới tính
            if 'giới tính' in line.lower():
                if 'nam' in line.lower():
                    info['gender'] = 'Nam'
                elif 'nữ' in line.lower():
                    info['gender'] = 'Nữ'
            
            # Tìm quốc tịch
            if 'quốc tịch' in line.lower():
                parts = line.split(':')
                if len(parts) > 1:
                    info['nationality'] = parts[1].strip()
            
            # Tìm quê quán
            if 'quê quán' in line.lower():
                parts = line.split(':')
                if len(parts) > 1:
                    info['place_of_origin'] = parts[1].strip()
            
            # Tìm nơi thường trú
            if 'nơi thường trú' in line.lower() or 'thường trú' in line.lower():
                parts = line.split(':')
                if len(parts) > 1:
                    info['place_of_residence'] = parts[1].strip()
            
            # Tìm dân tộc
            if 'dân tộc' in line.lower():
                parts = line.split(':')
                if len(parts) > 1:
                    info['ethnic'] = parts[1].strip()
            
            # Tìm tôn giáo
            if 'tôn giáo' in line.lower():
                parts = line.split(':')
                if len(parts) > 1:
                    info['religion'] = parts[1].strip()
        
        return info
    
    def parse_back_side(self, text: str) -> Dict[str, Optional[str]]:
        """Parse thông tin từ text đã OCR của mặt sau"""
        info = {}
        lines = text.split('\n')
        
        for line in lines:
            line = line.strip()
            if not line:
                continue
            
            # Tìm ngày cấp
            if 'ngày cấp' in line.lower() or 'cấp ngày' in line.lower():
                date_match = re.search(r'(\d{1,2})[/-](\d{1,2})[/-](\d{4})', line)
                if date_match:
                    day, month, year = date_match.groups()
                    info['issue_date'] = f"{day}/{month}/{year}"
            
            # Tìm nơi cấp
            if 'nơi cấp' in line.lower():
                parts = line.split(':')
                if len(parts) > 1:
                    info['issue_place'] = parts[1].strip()
            
            # Tìm đặc điểm nhận dạng
            if 'đặc điểm' in line.lower() or 'nhận dạng' in line.lower():
                parts = line.split(':')
                if len(parts) > 1:
                    info['features'] = parts[1].strip()
            
            # Tìm địa chỉ
            if 'địa chỉ' in line.lower():
                parts = line.split(':')
                if len(parts) > 1:
                    info['address'] = parts[1].strip()
        
        return info
    
    def extract_national_card_info(
        self, 
        front_image_url: str, 
        back_image_url: str
    ) -> Tuple[Dict[str, Optional[str]], Optional[str], float]:
        """
        Trích xuất thông tin từ ảnh căn cước
        
        Returns:
            Tuple[info_dict, mrz_string, processing_time_ms]
        """
        start_time = time.time()
        
        try:
            # Download ảnh
            front_image = self.download_image(front_image_url)
            back_image = self.download_image(back_image_url)
            
            # Extract MRZ từ ảnh mặt sau
            mrz_text = self.extract_mrz(back_image)
            mrz_parsed = self.parse_mrz(mrz_text) if mrz_text else {}
            
            # OCR cả 2 ảnh
            front_text = self.ocr_image(front_image, lang='vie+eng')
            back_text = self.ocr_image(back_image, lang='vie+eng')
            
            # Parse thông tin
            info = {}
            
            # Parse từ MRZ trước (chính xác hơn)
            if mrz_parsed:
                info.update(mrz_parsed)
            
            # Parse từ mặt trước
            front_info = self.parse_front_side(front_text)
            for key, value in front_info.items():
                if value and key not in info:  # Không override MRZ data
                    info[key] = value
            
            # Parse từ mặt sau
            back_info = self.parse_back_side(back_text)
            info.update(back_info)
            
            processing_time = (time.time() - start_time) * 1000
            
            return info, mrz_text, processing_time
            
        except Exception as e:
            raise ValueError(f"Failed to extract national card info: {str(e)}")

