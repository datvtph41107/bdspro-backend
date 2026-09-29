"""
Service để đọc QR code từ ảnh URL
Hỗ trợ download ảnh từ URL và decode QR code
"""

import io
import logging
from typing import Optional, List, Dict, Any
from urllib.request import urlopen

try:
    from PIL import Image
    PIL_AVAILABLE = True
except ImportError:
    PIL_AVAILABLE = False

try:
    from pyzbar.pyzbar import decode
    try:
        from pyzbar.pyzbar import ZBarSymbol
        ZBAR_SYMBOL_AVAILABLE = True
    except ImportError:
        ZBarSymbol = None
        ZBAR_SYMBOL_AVAILABLE = False
    PYZBAR_AVAILABLE = True
except ImportError:
    PYZBAR_AVAILABLE = False
    ZBarSymbol = None
    ZBAR_SYMBOL_AVAILABLE = False

logger = logging.getLogger(__name__)


class EasyOCRService:
    """Service để đọc QR code từ ảnh URL"""
    
    def __init__(self):
        """Khởi tạo EasyOCR service"""
        if not PIL_AVAILABLE:
            logger.warning("PIL/Pillow not available. QR code reading may not work.")
        if not PYZBAR_AVAILABLE:
            logger.warning("pyzbar not available. QR code reading may not work.")
    
    def download_image_from_url(self, image_url: str) -> Image.Image:
        """
        Download ảnh từ URL và convert sang PIL Image
        
        Args:
            image_url: URL của ảnh
            
        Returns:
            PIL Image object
            
        Raises:
            ValueError: Nếu không thể download hoặc decode ảnh
        """
        try:
            with urlopen(image_url) as response:
                image_data = response.read()
                image = Image.open(io.BytesIO(image_data))
                
                # Convert sang RGB nếu cần
                if image.mode != 'RGB':
                    if image.mode == 'RGBA':
                        # Convert RGBA sang RGB với background trắng
                        rgb_image = Image.new('RGB', image.size, (255, 255, 255))
                        rgb_image.paste(image, mask=image.split()[3])
                        return rgb_image
                    else:
                        return image.convert('RGB')
                
                return image
        except Exception as e:
            raise ValueError(f"Failed to download image from {image_url}: {str(e)}")
    
    def preprocess_image_for_qr(self, image: Image.Image) -> List[Image.Image]:
        """
        Preprocess ảnh để tăng khả năng đọc QR code
        Tạo nhiều version của ảnh với các xử lý khác nhau
        
        Args:
            image: PIL Image
            
        Returns:
            List các PIL Image đã preprocess
        """
        processed_images = [image]  # Thêm ảnh gốc
        
        try:
            # 1. Resize nếu ảnh quá nhỏ (QR code cần độ phân giải tối thiểu)
            width, height = image.size
            if width < 200 or height < 200:
                scale_factor = max(200 / width, 200 / height)
                new_size = (int(width * scale_factor), int(height * scale_factor))
                resized = image.resize(new_size, Image.Resampling.LANCZOS)
                processed_images.append(resized)
            
            # 2. Convert sang grayscale (thường tốt hơn cho QR code)
            if image.mode != 'L':
                gray = image.convert('L')
                processed_images.append(gray)
            
            # 3. Enhance contrast
            try:
                from PIL import ImageEnhance
                enhancer = ImageEnhance.Contrast(image)
                enhanced = enhancer.enhance(2.0)  # Tăng contrast 2 lần
                processed_images.append(enhanced)
            except:
                pass
            
            # 4. Grayscale với contrast enhanced
            try:
                from PIL import ImageEnhance
                gray_enhanced = image.convert('L')
                enhancer = ImageEnhance.Contrast(gray_enhanced)
                gray_enhanced = enhancer.enhance(2.0)
                processed_images.append(gray_enhanced)
            except:
                pass
            
        except Exception as e:
            logger.warning(f"Image preprocessing failed: {e}")
        
        return processed_images
    
    def read_qr_code_from_url(self, image_url: str) -> Optional[List[Dict[str, Any]]]:
        """
        Đọc QR code từ URL ảnh
        
        Args:
            image_url: URL của ảnh chứa QR code
            
        Returns:
            List các dict chứa thông tin QR code đã decode, hoặc None nếu không tìm thấy
            Format: [{"data": "text", "type": "QRCODE", "rect": (x, y, width, height)}]
        """
        try:
            if not PYZBAR_AVAILABLE:
                print("pyzbar not available. Cannot read QR code.")
                return None
            
            # Download ảnh từ URL
            image = self.download_image_from_url(image_url)
            print(f"Downloaded image: size={image.size}, mode={image.mode}")
            
            # Thử decode với nhiều cách khác nhau
            decoded_objects = None
            
            # 1. Thử với ảnh gốc (RGB) - chỉ tìm QRCODE
            try:
                if ZBAR_SYMBOL_AVAILABLE and ZBarSymbol:
                    decoded_objects = decode(image, symbols=[ZBarSymbol.QRCODE])
                else:
                    decoded_objects = decode(image)
            except:
                decoded_objects = decode(image)
            if decoded_objects:
                logger.debug("Found QR code in original RGB image")
            
            # 2. Nếu không tìm thấy, thử với grayscale
            if not decoded_objects:
                gray_image = image.convert('L')
                try:
                    if ZBAR_SYMBOL_AVAILABLE and ZBarSymbol:
                        decoded_objects = decode(gray_image, symbols=[ZBarSymbol.QRCODE])
                    else:
                        decoded_objects = decode(gray_image)
                except:
                    decoded_objects = decode(gray_image)
                if decoded_objects:
                    logger.debug("Found QR code in grayscale image")
            
            # 3. Nếu vẫn không tìm thấy, thử với các ảnh đã preprocess
            if not decoded_objects:
                logger.debug("No QR code found in original image, trying preprocessed versions...")
                processed_images = self.preprocess_image_for_qr(image)
                
                for i, processed_img in enumerate(processed_images):
                    if processed_img is image:  # Đã thử ảnh gốc rồi
                        continue
                    logger.debug(f"Trying preprocessed image #{i}: size={processed_img.size}, mode={processed_img.mode}")
                    decoded_objects = decode(processed_img)
                    if decoded_objects:
                        logger.info(f"Found QR code using preprocessed image #{i}")
                        break
            
            if not decoded_objects:
                logger.warning(f"No QR code found in image from {image_url} after all preprocessing attempts")
                return None
            
            # Format kết quả và loại bỏ duplicate
            results = []
            seen_data = set()
            for obj in decoded_objects:
                try:
                    data = obj.data.decode('utf-8')
                    if data in seen_data:
                        continue
                    seen_data.add(data)
                    
                    result = {
                        "data": data,
                        "type": obj.type,
                        "rect": {
                            "x": obj.rect.left,
                            "y": obj.rect.top,
                            "width": obj.rect.width,
                            "height": obj.rect.height
                        }
                    }
                    results.append(result)
                except UnicodeDecodeError:
                    logger.warning(f"Failed to decode QR code data as UTF-8: {obj.data}")
                    continue
            
            logger.info(f"Found {len(results)} unique QR code(s) in image from {image_url}")
            return results if results else None
            
        except Exception as e:
            logger.error(f"Error reading QR code from URL {image_url}: {str(e)}", exc_info=True)
            return None
    
    def read_qr_code_from_image(self, image: Image.Image) -> Optional[List[Dict[str, Any]]]:
        """
        Đọc QR code từ PIL Image object
        
        Args:
            image: PIL Image object
            
        Returns:
            List các dict chứa thông tin QR code đã decode, hoặc None nếu không tìm thấy
        """
        try:
            if not PYZBAR_AVAILABLE:
                logger.error("pyzbar not available. Cannot read QR code.")
                return None
            
            # Convert sang RGB nếu cần
            if image.mode != 'RGB':
                if image.mode == 'RGBA':
                    rgb_image = Image.new('RGB', image.size, (255, 255, 255))
                    rgb_image.paste(image, mask=image.split()[3])
                    image = rgb_image
                else:
                    image = image.convert('RGB')
            
            logger.debug(f"Processing image: size={image.size}, mode={image.mode}")
            
            # Thử decode với nhiều cách khác nhau
            decoded_objects = None
            
            # 1. Thử với ảnh gốc (RGB) - chỉ tìm QRCODE
            try:
                if ZBAR_SYMBOL_AVAILABLE and ZBarSymbol:
                    decoded_objects = decode(image, symbols=[ZBarSymbol.QRCODE])
                else:
                    decoded_objects = decode(image)
            except:
                decoded_objects = decode(image)
            if decoded_objects:
                logger.debug("Found QR code in original RGB image")
            
            # 2. Nếu không tìm thấy, thử với grayscale
            if not decoded_objects:
                gray_image = image.convert('L')
                try:
                    if ZBAR_SYMBOL_AVAILABLE and ZBarSymbol:
                        decoded_objects = decode(gray_image, symbols=[ZBarSymbol.QRCODE])
                    else:
                        decoded_objects = decode(gray_image)
                except:
                    decoded_objects = decode(gray_image)
                if decoded_objects:
                    logger.debug("Found QR code in grayscale image")
            
            # 3. Nếu vẫn không tìm thấy, thử với các ảnh đã preprocess
            if not decoded_objects:
                logger.debug("No QR code found in original image, trying preprocessed versions...")
                processed_images = self.preprocess_image_for_qr(image)
                
                for i, processed_img in enumerate(processed_images):
                    if processed_img is image:  # Đã thử ảnh gốc rồi
                        continue
                    logger.debug(f"Trying preprocessed image #{i}: size={processed_img.size}, mode={processed_img.mode}")
                    decoded_objects = decode(processed_img)
                    if decoded_objects:
                        logger.info(f"Found QR code using preprocessed image #{i}")
                        break
            
            if not decoded_objects:
                logger.warning("No QR code found in image after all preprocessing attempts")
                return None
            
            # Format kết quả và loại bỏ duplicate
            results = []
            seen_data = set()
            for obj in decoded_objects:
                try:
                    data = obj.data.decode('utf-8')
                    if data in seen_data:
                        continue
                    seen_data.add(data)
                    
                    result = {
                        "data": data,
                        "type": obj.type,
                        "rect": {
                            "x": obj.rect.left,
                            "y": obj.rect.top,
                            "width": obj.rect.width,
                            "height": obj.rect.height
                        }
                    }
                    results.append(result)
                except UnicodeDecodeError:
                    logger.warning(f"Failed to decode QR code data as UTF-8: {obj.data}")
                    continue
            
            logger.info(f"Found {len(results)} unique QR code(s) in image")
            return results if results else None
            
        except Exception as e:
            logger.error(f"Error reading QR code from image: {str(e)}", exc_info=True)
            return None


# Tạo instance global
easyocr_service = EasyOCRService()

