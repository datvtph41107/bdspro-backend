"""
Service để đọc MRZ (Machine Readable Zone) từ ảnh CCCD mặt sau
Hỗ trợ input từ URL hoặc file upload
"""

import io
import logging
import os
import re
import tempfile
from typing import Optional, Dict, Any, Union, Tuple
from urllib.request import urlopen
from urllib.parse import urlparse

try:
    import pytesseract
    from PIL import Image
    TESSERACT_AVAILABLE = True
except ImportError:
    TESSERACT_AVAILABLE = False

try:
    from passporteye import read_mrz
    PASSPORTEYE_AVAILABLE = True
except ImportError:
    PASSPORTEYE_AVAILABLE = False

try:
    import cv2
    import numpy as np
    OPENCV_AVAILABLE = True
except ImportError:
    OPENCV_AVAILABLE = False

logger = logging.getLogger(__name__)


class MRZScanService:
    """Service để đọc MRZ từ ảnh CCCD mặt sau"""
    
    def __init__(self):
        """Khởi tạo MRZ scan service"""
        self.use_passporteye = PASSPORTEYE_AVAILABLE
        self.use_opencv = OPENCV_AVAILABLE
        
        if not TESSERACT_AVAILABLE:
            logger.warning("pytesseract not available. MRZ reading may not work.")
        
        if self.use_passporteye:
            logger.info("Using passporteye for MRZ detection")
        else:
            logger.info("Using pytesseract fallback for MRZ reading")
    
    def _convert_to_rgb(self, image: Image.Image) -> Image.Image:
        """
        Convert image sang RGB mode (hỗ trợ RGBA và các mode khác)
        
        Args:
            image: PIL Image
            
        Returns:
            PIL Image ở mode RGB
        """
        if image.mode == 'RGB' or image.mode == 'L':
            return image
        
        if image.mode == 'RGBA':
            # Convert RGBA sang RGB với background trắng
            rgb_image = Image.new('RGB', image.size, (255, 255, 255))
            rgb_image.paste(image, mask=image.split()[3])
            return rgb_image
        
        # Convert các mode khác sang RGB
        return image.convert('RGB')
    
    def download_image_from_url(self, image_url: str) -> Tuple[Image.Image, Optional[str]]:
        """
        Download ảnh từ URL và convert sang PIL Image
        
        Args:
            image_url: URL của ảnh
            
        Returns:
            Tuple (PIL Image object, content_type)
            
        Raises:
            ValueError: Nếu không thể download hoặc decode ảnh
        """
        try:
            with urlopen(image_url) as response:
                content_type = response.headers.get('Content-Type', '').lower()
                image_data = response.read()
                image = Image.open(io.BytesIO(image_data))
                return self._convert_to_rgb(image), content_type
        except Exception as e:
            raise ValueError(f"Failed to download image from {image_url}: {str(e)}")
    
    def preprocess_image_for_mrz(self, image: Image.Image) -> Image.Image:
        """
        Preprocess ảnh để tăng độ chính xác khi đọc MRZ
        
        Args:
            image: PIL Image
            
        Returns:
            Preprocessed PIL Image
        """
        if not self.use_opencv:
            return image
        
        try:
            # Convert PIL to OpenCV format
            img_array = np.array(image)
            
            # Convert RGB to BGR (OpenCV format)
            if len(img_array.shape) == 3:
                img_bgr = cv2.cvtColor(img_array, cv2.COLOR_RGB2BGR)
            else:
                img_bgr = img_array
            
            # Crop phần dưới 35% của ảnh (MRZ thường ở đó)
            # Giảm crop ratio để không bỏ sót dòng MRZ
            height, width = img_bgr.shape[:2]
            crop_y = int(height * 0.55)  # Từ 0.7 xuống 0.65 để lấy thêm phần trên
            mrz_region = img_bgr[crop_y:height, 0:width]
            
            # Convert to grayscale
            gray = cv2.cvtColor(mrz_region, cv2.COLOR_BGR2GRAY) if len(mrz_region.shape) == 3 else mrz_region
            
            # Apply threshold để làm rõ text
            _, thresh = cv2.threshold(gray, 0, 255, cv2.THRESH_BINARY + cv2.THRESH_OTSU)
            
            # Convert back to PIL Image
            pil_image = Image.fromarray(thresh)
            
            return pil_image
            
        except Exception as e:
            logger.warning(f"Image preprocessing failed, using original: {e}")
            return image
    
    def read_mrz_with_passporteye(self, image_path: str) -> Optional[Dict[str, Any]]:
        """
        Đọc MRZ sử dụng passporteye (tốt hơn)
        
        Args:
            image_path: Đường dẫn đến file ảnh
            
        Returns:
            Dict chứa MRZ data hoặc None
        """
        try:
            mrz = read_mrz(image_path)
            
            if mrz is None:
                return None
            
            # Extract MRZ data
            mrz_text = mrz.text if hasattr(mrz, 'text') else str(mrz)
            
            # Parse MRZ
            parsed = self._parse_mrz(mrz_text)
            
            return {
                "raw": mrz_text,
                "lines": mrz_text.split('\n')[:2] if mrz_text else [],
                "parsed": parsed,
                "method": "passporteye"
            }
            
        except Exception as e:
            logger.warning(f"passporteye failed: {e}")
            return None
    
    def read_mrz_with_tesseract(self, image: Image.Image) -> Optional[Dict[str, Any]]:
        """
        Đọc MRZ sử dụng Tesseract (fallback)
        
        Args:
            image: PIL Image
        
        Returns:
            Dict chứa MRZ data hoặc None
        """
        try:
            # Crop phần dưới của ảnh trước (MRZ thường ở đó)
            # width, height = image.size
            # crop_y = int(height * 0.65)  # Lấy 40% dưới cùng để đảm bảo không bỏ sót
            # mrz_region = image.crop((0, crop_y, width, height))
            
            # Preprocess ảnh đã crop
            # processed_image = self.preprocess_image_for_mrz(mrz_region)
            
            # Thử nhiều PSM mode để tăng độ chính xác
            # psm_modes = [6, 11, 12]  # 6: single block, 11: sparse text, 12: sparse text with OSD
            all_mrz_lines = []
            
            try:
                # OCR toàn bộ ảnh để lấy tất cả text
                mrz_text = pytesseract.image_to_string(
                    image,
                    lang='eng',
                    config='--psm 6'
                )
                # mrz_text = mrz_text.strip()
                
                print(f"Raw OCR text: {mrz_text}")
                
                # Parse và collect dòng MRZ
                lines = mrz_text.split('\n')
                for line in lines:
                    line = line.strip()
                    # Lấy tất cả dòng có chứa ký tự đặc biệt của MRZ (<, số, chữ in hoa)
                    # MRZ thường có format: chứa <, số, và chữ in hoa
                    if len(line) >= 10 and ('<' in line or re.search(r'[A-Z0-9]{9,}', line)):
                        if line not in all_mrz_lines:
                            all_mrz_lines.append("".join(line.split(' ')))
            except Exception as e:
                logger.debug(f"Tesseract failed: {e}")
            
            logger.debug(f"Found {len(all_mrz_lines)} MRZ lines: {all_mrz_lines}")
            
            # Lấy tất cả dòng MRZ hợp lệ (thường 2-3 dòng)
            if len(all_mrz_lines) >= 2:
                # Lấy tất cả dòng (có thể 2-3 dòng)
                mrz_lines = all_mrz_lines if len(all_mrz_lines) <= 3 else all_mrz_lines[-3:]
                mrz_data = '\n'.join(mrz_lines)
                parsed = self._parse_mrz(mrz_data)
                
                return {
                    "raw": mrz_data,
                    "lines": mrz_lines,
                    "parsed": parsed,
                    "method": "tesseract"
                }
            
            # Nếu chỉ có 1 dòng, vẫn trả về để debug
            if len(all_mrz_lines) == 1:
                logger.warning(f"Only found 1 MRZ line: {all_mrz_lines[0]}")
            
            return None
            
        except Exception as e:
            logger.error(f"Tesseract MRZ reading failed: {e}")
            return None
    
    def _get_image_format_from_content_type(self, content_type: str) -> Optional[str]:
        """Parse format ảnh từ Content-Type header"""
        if not content_type:
                        return None
                        
        content_type_map = {
            'image/jpeg': 'JPEG', 'image/jpg': 'JPEG', 'image/png': 'PNG',
            'image/gif': 'GIF', 'image/bmp': 'BMP', 'image/tiff': 'TIFF',
            'image/tif': 'TIFF', 'image/webp': 'WEBP',
        }
        
        mime_type = content_type.split(';')[0].strip()
        return content_type_map.get(mime_type)
    
    def _get_extension_from_format(self, image_format: str) -> str:
        """Lấy extension từ format"""
        format_map = {
            'JPEG': '.jpg', 'PNG': '.png', 'GIF': '.gif', 'BMP': '.bmp',
            'TIFF': '.tiff', 'TIF': '.tiff', 'WEBP': '.webp'
        }
        return format_map.get(image_format, '.jpg')
    
    def _save_image_to_temp(self, image: Image.Image, image_format: str) -> str:
        """
        Lưu image vào file tạm với format phù hợp
        
        Args:
            image: PIL Image (đã ở mode RGB)
            image_format: Format muốn lưu (JPEG, PNG, etc.)
            
        Returns:
            Đường dẫn file tạm
        """
        # Đảm bảo image ở mode RGB
        image = self._convert_to_rgb(image)
        
        # Xác định extension và format
        ext = self._get_extension_from_format(image_format)
        save_format = image_format if image_format in ('JPEG', 'PNG', 'GIF', 'BMP', 'TIFF', 'WEBP') else 'JPEG'
        
        # Tạo file tạm
        with tempfile.NamedTemporaryFile(delete=False, suffix=ext) as tmp_file:
            tmp_path = tmp_file.name
        
        try:
            image.save(tmp_path, save_format)
            return tmp_path
        except Exception as e:
            # Fallback về JPEG nếu lưu thất bại
            logger.warning(f"Failed to save as {save_format}, trying JPEG: {e}")
            tmp_path_jpg = tmp_path.rsplit('.', 1)[0] + '.jpg'
            image.save(tmp_path_jpg, 'JPEG')
            return tmp_path_jpg
    
    def read_mrz_from_url(self, image_url: str) -> Optional[Dict[str, Any]]:
        """Đọc MRZ từ URL ảnh"""
        try:
            image, content_type = self.download_image_from_url(image_url)
            
            # Detect format
            image_format = None
            if content_type:
                image_format = self._get_image_format_from_content_type(content_type)
            if not image_format:
                image_format = image.format or 'PNG'
            
            # Lưu vào file tạm
            tmp_path = self._save_image_to_temp(image, image_format)
            
            try:
                # Thử passporteye trước
                if self.use_passporteye:
                    result = self.read_mrz_with_passporteye(tmp_path)
                    if result:
                        return result
                
                # Fallback Tesseract
                if TESSERACT_AVAILABLE:
                    result = self.read_mrz_with_tesseract(image)
                    print(f"result: {result}")
                    if result:
                        return result
                
                    return None
            finally:
                try:
                    os.unlink(tmp_path)
                except Exception as e:
                    logger.warning(f"Failed to delete temp file: {e}")
        except Exception as e:
            logger.error(f"Error reading MRZ from URL: {str(e)}")
            return None
    
    def read_mrz_from_image(self, image: Union[Image.Image, bytes, str]) -> Optional[Dict[str, Any]]:
        """Đọc MRZ từ ảnh (hỗ trợ URL, file path, bytes, PIL Image)"""
        try:
            # Xử lý input khác nhau
            if isinstance(image, str):
                if image.startswith(('http://', 'https://')):
                    return self.read_mrz_from_url(image)
                elif os.path.exists(image):
                    pil_image = Image.open(image)
                else:
                    raise ValueError(f"Invalid image path or URL: {image}")
            elif isinstance(image, bytes):
                pil_image = Image.open(io.BytesIO(image))
            elif isinstance(image, Image.Image):
                pil_image = image
            else:
                raise ValueError(f"Unsupported image type: {type(image)}")
            
            # Convert và lưu file tạm
            image_format = pil_image.format or 'PNG'
            pil_image = self._convert_to_rgb(pil_image)
            tmp_path = self._save_image_to_temp(pil_image, image_format)
            
            try:
                # Thử passporteye trước
                if self.use_passporteye:
                    result = self.read_mrz_with_passporteye(tmp_path)
                    if result:
                        return result
                
                # Fallback Tesseract
                if TESSERACT_AVAILABLE:
                    result = self.read_mrz_with_tesseract(pil_image)
                    if result:
                        return result
                
                return None
            finally:
                try:
                    os.unlink(tmp_path)
                except Exception as e:
                    logger.warning(f"Failed to delete temp file: {e}")
        except Exception as e:
            logger.error(f"Error reading MRZ from image: {str(e)}")
            return None
    
    def _parse_mrz(self, mrz_data: str) -> Dict[str, Any]:
        """
        Parse MRZ data thành các trường thông tin
        
        MRZ format cho CCCD Việt Nam thường có 2-3 dòng:
        - Line 1: IDVNM + số căn cước (9-12 chữ số)
        - Line 2: Ngày sinh (YYMMDD) + giới tính + hạn (YYMMDD) + quốc tịch
        - Line 3: Họ tên (format: HO<<TEN<DEM)
        
        Ví dụ:
        IDVNM2000012341123400001234<<9
        0011209M4011227VNM<<<<<<<<<<<
        NGUYEN<<VAN<A<<<<<<<
        
        Args:
            mrz_data: MRZ string (2-3 dòng)
            
        Returns:
            Dict chứa các trường đã parse
        """
        try:
            lines = [line.strip() for line in mrz_data.strip().split('\n') if line.strip()]
            if len(lines) < 2:
                return {}
            
            parsed = {
                "line1": lines[0] if len(lines) > 0 else "",
                "line2": lines[1] if len(lines) > 1 else "",
                "line3": lines[2] if len(lines) > 2 else "",
            }
            
            line1 = lines[0]
            line2 = lines[1] if len(lines) > 1 else ""
            line3 = lines[2] if len(lines) > 2 else ""
            
            # Dòng 1: Extract số căn cước
            # Format: IDVNM + số căn cước (9-12 chữ số)
            if line1:
                # Tìm vị trí của << trong dòng
                double_lt_pos = line1.find('<<')
                
                if double_lt_pos > 0:
                    # Lấy phần trước <<
                    before_lt = line1[:double_lt_pos]
                    # Loại bỏ tất cả ký tự không phải số
                    digits_before_lt = re.sub(r'[^\d]', '', before_lt)
                    
                    if len(digits_before_lt) >= 12:
                        # Lấy 12 chữ số cuối cùng trước <<
                        parsed["id_number"] = digits_before_lt[-12:]
                    elif len(digits_before_lt) >= 9:
                        # Nếu không đủ 12 chữ số, lấy tất cả (9-11 chữ số)
                        parsed["id_number"] = digits_before_lt
                else:
                    # Nếu không có <<, fallback: lấy 12 chữ số cuối cùng của dòng
                    all_digits = re.sub(r'[^\d]', '', line1)
                    if len(all_digits) >= 12:
                        parsed["id_number"] = all_digits[-12:]
                    elif len(all_digits) >= 9:
                        parsed["id_number"] = all_digits
            
            # Dòng 2: Extract ngày sinh và hạn
            # Format: YYMMDD + giới tính + YYMMDD + quốc tịch
            # Ví dụ: 0011189M4011187VNM
            # Ngày sinh: 001118 → 18/11/2000 (YYMMDD)
            # Hạn: 401118 → 18/11/2040 (YYMMDD)
            if line2:
                # Tìm tất cả các chuỗi 6-7 chữ số (có thể có thêm 1 ký tự)
                date_matches = re.findall(r'\d{6,7}', line2)
                
                if len(date_matches) >= 2:
                    # Ngày sinh thường là số đầu tiên (6 chữ số)
                    dob = date_matches[0][:6]  # Lấy 6 chữ số đầu
                    if len(dob) == 6:
                        year_2digits = int(dob[0:2])
                        # Nếu năm < 50 thì là 20xx, ngược lại là 19xx
                        year = f"20{year_2digits:02d}" if year_2digits < 50 else f"19{year_2digits:02d}"
                        month = dob[2:4]
                        day = dob[4:6]
                        parsed["date_of_birth"] = f"{day}/{month}/{year}"
                    
                    # Hạn thường là số thứ 2 (6 chữ số)
                    # Format: YYMMDD → DD/MM/20YY
                    expiry = date_matches[1][:6]  # Lấy 6 chữ số đầu
                    if len(expiry) == 6:
                        year_2digits = int(expiry[0:2])
                        # Hạn luôn là năm tương lai (20xx)
                        year = f"20{year_2digits:02d}"
                        month = expiry[2:4]
                        day = expiry[4:6]
                        parsed["expiry_date"] = f"{day}/{month}/{year}"
                
                # Extract giới tính (M hoặc F)
                gender_match = re.search(r'[MF]', line2)
                if gender_match:
                    parsed["gender"] = "Nam" if gender_match.group() == "M" else "Nữ"
            
            # Dòng 3: Extract họ tên
            # Format: HO<<TEN<DEM hoặc HO<<TEN<DEM<<<<
            # Ví dụ: NGO<<QUANG<HIEU<<<<<<<<<c<cc<<
            if line3:
                # Loại bỏ các ký tự < và khoảng trắng ở đầu/cuối
                name_line = line3.strip()
                
                # Loại bỏ các ký tự lạ ở cuối (như c, cc, v.v.) - chỉ giữ lại chữ cái in hoa và <
                # Tìm vị trí cuối cùng có chữ cái in hoa hợp lệ
                # Pattern: tìm chuỗi chữ cái in hoa, sau đó có thể có < và ký tự lạ
                # Chỉ lấy phần từ đầu đến khi gặp ký tự lạ (chữ thường, số, ký tự đặc biệt không phải <)
                clean_name = re.match(r'^([A-Z<]+)', name_line)
                if clean_name:
                    name_line = clean_name.group(1)
                
                # Loại bỏ các ký tự < ở cuối
                name_line = name_line.rstrip('<').strip()
                
                # Tách bằng < hoặc <<
                name_parts = re.split(r'<+', name_line)
                
                # Lọc bỏ phần rỗng và chỉ giữ phần là chữ cái in hoa
                name_parts = [part.strip() for part in name_parts if part.strip() and part.strip().isalpha() and part.strip().isupper()]
                
                if name_parts:
                    # Họ thường là phần đầu tiên
                    if len(name_parts) >= 1:
                        parsed["last_name"] = name_parts[0]
                    # Tên đệm và tên
                    if len(name_parts) >= 2:
                        parsed["first_name"] = name_parts[-1]  # Tên là phần cuối
                        if len(name_parts) > 2:
                            parsed["middle_name"] = " ".join(name_parts[1:-1])  # Tên đệm ở giữa
                        # Full name
                        parsed["full_name"] = " ".join(name_parts)
            
            # Document type và country code từ dòng 1
            if line1.startswith("IDVNM"):
                parsed["document_type"] = "ID"
                parsed["country_code"] = "VNM"
            
            return parsed
            
        except Exception as e:
            logger.error(f"Error parsing MRZ: {str(e)}")
            return {}


# Tạo instance global
mrz_scan_service = MRZScanService()

