"""
HTTP Server - Đọc QR code và MRZ từ ảnh căn cước
Port: 8218
Logic được port từ assistant-service/infra/client/national_card.go
"""

import os
import re
import time
import tempfile
import asyncio
import io
from typing import Optional, Dict, Any, Tuple
from fastapi import FastAPI, File, UploadFile, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
from pyzbar.pyzbar import decode
import uvicorn
import httpx

# MRZ dependencies
try:
    import pytesseract
    TESSERACT_AVAILABLE = True
except ImportError:
    TESSERACT_AVAILABLE = False

try:
    from passporteye import read_mrz
    PASSPORTEYE_AVAILABLE = True
except ImportError:
    PASSPORTEYE_AVAILABLE = False

try:
    from PIL import Image
    PIL_AVAILABLE = True
except ImportError:
    PIL_AVAILABLE = False


class NationalCardRequest(BaseModel):
    frontImageUrl: str
    backImageUrl: str


class NationalCardExtractResponse(BaseModel):
    code: int  # 0 = success, khác 0 = error
    message: str  # "success" hoặc message lỗi
    data: Dict[str, Any]  # Chứa info, mrz, qr, và các thông tin khác


app = FastAPI(title="National Card Extraction", version="2.0.0")
app.add_middleware(CORSMiddleware, allow_origins=["*"], allow_credentials=True, allow_methods=["*"], allow_headers=["*"])


def _parse_int(s: str) -> int:
    """Parse string thành int (helper)"""
    result = 0
    for c in s:
        if '0' <= c <= '9':
            result = result * 10 + int(c)
    return result


def _read_qr_code_from_bytes(image_data: bytes) -> str:
    """Đọc QR code từ image bytes"""
    try:
        if not PIL_AVAILABLE:
            return ""
        
        image = Image.open(io.BytesIO(image_data))
        results = decode(image)
        if not results:
            return ""
        
        # Trả về QR code đầu tiên
        return results[0].data.decode('utf-8')
    except Exception as e:
        print(f"Error reading QR code: {e}")
        return ""


def _crop_mrz_region(image_path: str) -> Optional[str]:
    """Crop ảnh để chỉ lấy phần MRZ (phần dưới 40% của ảnh)"""
    try:
        if not PIL_AVAILABLE:
            return None
        
        image = Image.open(image_path)
        width, height = image.size
        
        # Crop phần dưới 40% (từ 60% xuống)
        crop_y = int(height * 0.6)
        cropped = image.crop((0, crop_y, width, height))
        
        # Lưu vào file tạm
        cropped_path = tempfile.mktemp(suffix='.png')
        cropped.save(cropped_path)
        return cropped_path
    except Exception as e:
        print(f"Error cropping MRZ region: {e}")
        return None


def _clean_mrz_text(text: str) -> str:
    """Làm sạch text MRZ, loại bỏ ký tự không hợp lệ và khoảng trắng"""
    lines = text.split('\n')
    cleaned_lines = []
    
    mrz_regex = re.compile(r'[A-Z0-9<]{22,}')
    
    for line in lines:
        line = line.strip()
        line = line.replace(' ', '')
        # Chỉ giữ dòng có format MRZ (nhiều ký tự A-Z, 0-9, <)
        if len(line) >= 22 and mrz_regex.match(line) and line.count('<') > 1:
            cleaned_lines.append(line)
    
    return '\n'.join(cleaned_lines)


def _is_valid_mrz(text: str) -> bool:
    """Kiểm tra xem text có phải MRZ hợp lệ không"""
    lines = text.strip().split('\n')
    
    if len(lines) < 2:
        return False
    
    valid_lines = 0
    for line in lines:
        line = line.strip()
        if len(line) >= 22 and line.count('<') > 1:
            valid_lines += 1
    
    return valid_lines >= 2


def _ocr_mrz_from_bytes(image_data: bytes) -> str:
    """OCR MRZ từ bytes"""
    try:
        if not TESSERACT_AVAILABLE:
            return ""
        
        # Lưu vào file tạm
        temp_path = tempfile.mktemp(suffix='.png')
        with open(temp_path, 'wb') as f:
            f.write(image_data)
        
        try:
            # Crop vùng MRZ
            crop_path = _crop_mrz_region(temp_path)
            if crop_path:
                mrz_image_path = crop_path
            else:
                mrz_image_path = temp_path
            
            # OCR với Tesseract
            text = pytesseract.image_to_string(
                Image.open(mrz_image_path),
                lang='eng',
                config='--psm 6 -c tessedit_char_whitelist=ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789<'
            )
            
            mrz = _clean_mrz_text(text)
            if _is_valid_mrz(mrz):
                return mrz
            
            return ""
        finally:
            try:
                os.unlink(temp_path)
                if crop_path and crop_path != temp_path:
                    os.unlink(crop_path)
            except:
                pass
    except Exception as e:
        print(f"Error OCR MRZ: {e}")
        return ""


def _extract_mrz_from_text(text: str) -> str:
    """Trích xuất MRZ từ text đã OCR"""
    lines = text.split('\n')
    mrz_lines = []
    
    mrz_regex = re.compile(r'[A-Z0-9<]{20,}')
    
    # Tìm đoạn liên tiếp có 3 dòng MRZ
    for i in range(len(lines)):
        line = lines[i].strip()
        line_clean = line.replace(' ', '')
        
        if len(line_clean) >= 20 and mrz_regex.match(line_clean) and line_clean.count('<') >= 2:
            # Kiểm tra xem có phải là đoạn 3 dòng MRZ liên tiếp không
            if i + 2 < len(lines):
                line2 = lines[i + 1].strip()
                line2_clean = line2.replace(' ', '')
                line3 = lines[i + 2].strip()
                line3_clean = line3.replace(' ', '')
                
                if (len(line2_clean) >= 20 and mrz_regex.match(line2_clean) and line2_clean.count('<') >= 2 and
                    len(line3_clean) >= 20 and mrz_regex.match(line3_clean) and line3_clean.count('<') >= 2):
                    mrz_lines = [line_clean, line2_clean, line3_clean]
                    break
            
            # Nếu không tìm thấy 3 dòng liên tiếp, thử tìm 2 dòng
            if len(mrz_lines) == 0 and i + 1 < len(lines):
                line2 = lines[i + 1].strip()
                line2_clean = line2.replace(' ', '')
                if len(line2_clean) >= 20 and mrz_regex.match(line2_clean) and line2_clean.count('<') >= 2:
                    mrz_lines = [line_clean, line2_clean]
    
    if len(mrz_lines) >= 2:
        return '\n'.join(mrz_lines)
    
    return ""


def _parse_qr_code(qr_text: str) -> Dict[str, Any]:
    """Parse QR code string thành thông tin
    Format QR code: ID||FullName|DOB|Gender|Address|IssueDate|...
    DOB format: DDMMYYYY (18112000 -> 18/11/2000)
    IssueDate format: DDMMYYYY (17092025 -> 17/09/2025)
    """
    result = {}
    
    # Remove brackets nếu có
    qr_text = qr_text.strip()
    if qr_text.startswith('[') and qr_text.endswith(']'):
        qr_text = qr_text[1:-1]
    
    # Split by |
    parts = qr_text.split('|')
    if len(parts) < 5:
        return result
    
    # [0]: ID number
    if len(parts) > 0 and parts[0]:
        result['idNumber'] = parts[0]
    
    # [1]: Empty (||)
    # [2]: Full name
    if len(parts) > 2 and parts[2]:
        result['fullName'] = parts[2]
    
    # [3]: Date of birth (format: DDMMYYYY hoặc YYMMDD -> dd/mm/yyyy)
    if len(parts) > 3 and parts[3]:
        dob = parts[3]
        if len(dob) == 8:
            # Format: DDMMYYYY
            day = dob[0:2]
            month = dob[2:4]
            year = dob[4:8]
            result['dateOfBirth'] = f"{day}/{month}/{year}"
        elif len(dob) == 6:
            # Format: YYMMDD
            yy = dob[0:2]
            month = dob[2:4]
            day = dob[4:6]
            yy_int = _parse_int(yy)
            year = f"20{yy}" if yy_int < 50 else f"19{yy}"
            result['dateOfBirth'] = f"{day}/{month}/{year}"
    
    # [4]: Gender
    if len(parts) > 4 and parts[4]:
        result['gender'] = parts[4]
    
    # [5]: Address
    if len(parts) > 5 and parts[5]:
        result['address'] = parts[5]
        result['placeOfResidence'] = parts[5]
    
    # [6]: Issue date (format: DDMMYYYY hoặc YYMMDD -> dd/mm/yyyy)
    if len(parts) > 6 and parts[6]:
        issue_date = parts[6]
        if len(issue_date) == 8:
            day = issue_date[0:2]
            month = issue_date[2:4]
            year = issue_date[4:8]
            result['issueDate'] = f"{day}/{month}/{year}"
        elif len(issue_date) == 6:
            yy = issue_date[0:2]
            month = issue_date[2:4]
            day = issue_date[4:6]
            year = f"20{yy}"
            result['issueDate'] = f"{day}/{month}/{year}"
    
    return result


def _parse_mrz(mrz: str) -> Dict[str, Any]:
    """Parse MRZ string thành thông tin"""
    result = {}
    
    lines = mrz.strip().split('\n')
    if len(lines) < 2:
        return result
    
    line1 = lines[0].strip().replace(' ', '')
    line2 = lines[1].strip().replace(' ', '')
    line3 = lines[2].strip().replace(' ', '') if len(lines) >= 3 else ""
    
    # Extract ID number từ MRZ line 1
    if '<<' in line1:
        idx = line1.index('<<')
        before = line1[:idx]
        digits = re.findall(r'\d+', before)
        if digits:
            last_digits = digits[-1]
            if len(last_digits) >= 12:
                result['idNumber'] = last_digits[-12:]
            elif len(last_digits) == 12:
                result['idNumber'] = last_digits
    
    # Fallback: tìm số 12 chữ số bắt đầu bằng 0
    if 'idNumber' not in result:
        id_matches = re.findall(r'\d{12}', line1)
        if id_matches:
            for match in id_matches:
                if match[0] == '0':
                    result['idNumber'] = match
                    break
            if 'idNumber' not in result:
                result['idNumber'] = id_matches[0]
    
    # Extract date of birth từ line 2 (format: YYMMDD)
    dob_matches = re.findall(r'\d{6}', line2)
    if dob_matches:
        match = dob_matches[0]
        if len(match) == 6:
            yy = match[0:2]
            month = match[2:4]
            day = match[4:6]
            month_int = _parse_int(month)
            day_int = _parse_int(day)
            
            if 1 <= month_int <= 12 and 1 <= day_int <= 31:
                yy_int = _parse_int(yy)
                year = f"20{yy}" if yy_int < 50 else f"19{yy}"
                result['dateOfBirth'] = f"{day}/{month}/{year}"
    
    # Extract full name từ line 3 (nếu có)
    if line3:
        name_match = re.match(r'^([A-Z<]+)', line3)
        if name_match:
            name_line = name_match.group(1).rstrip('<')
            name_parts = [p for p in name_line.split('<') if p.strip() and p.strip().isalpha() and p.strip().isupper()]
            if name_parts:
                result['fullName'] = ' '.join(name_parts)
    
    # Extract expiry date từ line 2
    expiry_matches = re.findall(r'\d{6}', line2)
    if len(expiry_matches) >= 2:
        match = expiry_matches[1]
        if len(match) == 6:
            yy = match[0:2]
            month = match[2:4]
            day = match[4:6]
            month_int = _parse_int(month)
            day_int = _parse_int(day)
            
            if 1 <= month_int <= 12 and 1 <= day_int <= 31:
                year = f"20{yy}"
                result['expiryDate'] = f"{day}/{month}/{year}"
    elif len(expiry_matches) == 1:
        match = expiry_matches[0]
        if len(match) == 6:
            yy = match[0:2]
            month = match[2:4]
            day = match[4:6]
            month_int = _parse_int(month)
            day_int = _parse_int(day)
            yy_int = _parse_int(yy)
            
            if 1 <= month_int <= 12 and 1 <= day_int <= 31 and yy_int >= 30:
                year = f"20{yy}"
                result['expiryDate'] = f"{day}/{month}/{year}"
    
    return result


def _extract_id_number(text: str) -> str:
    """Tìm số CMND/CCCD: 9-12 chữ số liên tiếp"""
    matches = re.findall(r'\d{9,12}', text)
    if matches:
        return matches[0]
    return ""


def _extract_date(text: str) -> str:
    """Tìm ngày tháng năm: dd/mm/yyyy hoặc dd-mm-yyyy"""
    match = re.search(r'(\d{1,2})[/-](\d{1,2})[/-](\d{4})', text)
    if match:
        return f"{match.group(1)}/{match.group(2)}/{match.group(3)}"
    return ""


def _extract_gender(text: str) -> str:
    """Extract giới tính từ text"""
    text_lower = text.lower()
    if 'nam' in text_lower:
        return "Nam"
    if 'nữ' in text_lower:
        return "Nữ"
    return ""


async def _download_image(url: str, timeout: int = 60) -> bytes:
    """Download ảnh từ URL"""
    async with httpx.AsyncClient(timeout=timeout) as client:
        try:
            _clean_url = url.strip()
            response = await client.get(_clean_url)
            response.raise_for_status()
            if response.status_code != 200:
                raise Exception(f"Failed to download image: status={response.status_code}")
            return response.content
        except httpx.HTTPError as e:
            raise Exception(f"Failed to download image from {url}: {str(e)}")


def _ocr_text_from_bytes(image_data: bytes) -> str:
    """OCR text từ bytes"""
    try:
        if not TESSERACT_AVAILABLE:
            return ""
        
        # Lưu vào file tạm
        temp_path = tempfile.mktemp(suffix='.png')
        with open(temp_path, 'wb') as f:
            f.write(image_data)
        
        try:
            # OCR với Tesseract (tiếng Việt + tiếng Anh)
            text = pytesseract.image_to_string(
                Image.open(temp_path),
                lang='vie+eng',
                config='--psm 6'
            )
            return text
        finally:
            try:
                os.unlink(temp_path)
            except:
                pass
    except Exception as e:
        print(f"Error OCR text: {e}")
        return ""


def _parse_front_side(text: str, info: Dict[str, Any]):
    """Parse thông tin từ mặt trước căn cước"""
    lines = text.split('\n')
    
    for line in lines:
        line = line.strip()
        if not line:
            continue
        
        # Tìm số CMND/CCCD
        id_number = _extract_id_number(line)
        if id_number and 'idNumber' not in info:
            info['idNumber'] = id_number
        
        # Tìm họ tên
        if ('fullName' not in info or not info['fullName']) and 'họ' in line.lower() and 'tên' in line.lower():
            parts = line.split(':')
            if len(parts) > 1:
                name = parts[1].strip()
                if name:
                    info['fullName'] = name
        
        # Tìm ngày sinh
        if ('dateOfBirth' not in info or not info['dateOfBirth']) and ('ngày sinh' in line.lower() or 'sinh ngày' in line.lower()):
            dob = _extract_date(line)
            if dob:
                info['dateOfBirth'] = dob
        
        # Tìm giới tính
        if 'giới tính' in line.lower() or 'nam' in line.lower() or 'nữ' in line.lower():
            gender = _extract_gender(line)
            if gender:
                info['gender'] = gender
        
        # Tìm quốc tịch
        if 'quốc tịch' in line.lower():
            parts = line.split(':')
            if len(parts) > 1:
                nationality = parts[1].strip()
                if nationality:
                    info['nationality'] = nationality
        
        # Tìm quê quán
        if 'quê quán' in line.lower():
            parts = line.split(':')
            if len(parts) > 1:
                place_of_origin = parts[1].strip()
                if place_of_origin:
                    info['placeOfOrigin'] = place_of_origin
        
        # Tìm nơi thường trú
        if 'nơi thường trú' in line.lower() or 'thường trú' in line.lower():
            parts = line.split(':')
            if len(parts) > 1:
                place_of_residence = parts[1].strip()
                if place_of_residence:
                    info['placeOfResidence'] = place_of_residence
        
        # Tìm dân tộc
        if 'dân tộc' in line.lower():
            parts = line.split(':')
            if len(parts) > 1:
                ethnic = parts[1].strip()
                if ethnic:
                    info['ethnic'] = ethnic
        
        # Tìm tôn giáo
        if 'tôn giáo' in line.lower():
            parts = line.split(':')
            if len(parts) > 1:
                religion = parts[1].strip()
                if religion:
                    info['religion'] = religion


def _parse_back_side(text: str, info: Dict[str, Any]):
    """Parse thông tin từ mặt sau căn cước"""
    lines = text.split('\n')
    
    for line in lines:
        line = line.strip()
        if not line:
            continue
        
        # Tìm ngày cấp
        if ('issueDate' not in info or not info['issueDate']) and ('ngày cấp' in line.lower() or 'cấp ngày' in line.lower()):
            issue_date = _extract_date(line)
            if issue_date:
                info['issueDate'] = issue_date
        
        # Tìm nơi cấp
        if 'nơi cấp' in line.lower():
            parts = line.split(':')
            if len(parts) > 1:
                issue_place = parts[1].strip()
                if issue_place:
                    info['issuePlace'] = issue_place
        
        # Tìm đặc điểm nhận dạng
        if 'đặc điểm' in line.lower() or 'nhận dạng' in line.lower():
            parts = line.split(':')
            if len(parts) > 1:
                features = parts[1].strip()
                if features:
                    info['features'] = features
        
        # Tìm địa chỉ
        if 'địa chỉ' in line.lower():
            parts = line.split(':')
            if len(parts) > 1:
                address = parts[1].strip()
                if address:
                    info['address'] = address


async def _extract_national_card_from_data(front_image_data: bytes, back_image_data: bytes) -> Tuple[Dict[str, Any], str, str]:
    """
    Trích xuất thông tin từ 2 ảnh CCCD (từ file data)
    Trả về: (info, mrz, qr)
    """
    start_time = time.time()
    
    # Đọc QR code và MRZ song song từ cả 2 mặt
    qr_front_task = asyncio.create_task(asyncio.to_thread(_read_qr_code_from_bytes, front_image_data))
    qr_back_task = asyncio.create_task(asyncio.to_thread(_read_qr_code_from_bytes, back_image_data))
    mrz_front_task = asyncio.create_task(asyncio.to_thread(_ocr_mrz_from_bytes, front_image_data))
    mrz_back_task = asyncio.create_task(asyncio.to_thread(_ocr_mrz_from_bytes, back_image_data))
    
    qr_front, qr_back, mrz_front, mrz_back = await asyncio.gather(
        qr_front_task, qr_back_task, mrz_front_task, mrz_back_task
    )
    
    print(f"⏱️  [Infer] QR Front: {qr_front[:50] if qr_front else 'None'}")
    print(f"⏱️  [Infer] QR Back: {qr_back[:50] if qr_back else 'None'}")
    
    # Extract MRZ từ text song song (nếu cần fallback)
    mrz_from_front_text = _extract_mrz_from_text(str(front_image_data))
    mrz_from_back_text = _extract_mrz_from_text(str(back_image_data))
    
    # Xử lý kết quả MRZ (ưu tiên mặt sau)
    mrz = ""
    if mrz_back:
        mrz = mrz_back
    elif mrz_front:
        mrz = mrz_front
    else:
        if mrz_from_back_text:
            mrz = mrz_from_back_text
        elif mrz_from_front_text:
            mrz = mrz_from_front_text
    
    # Xử lý kết quả QR code (ưu tiên mặt sau)
    qr = ""
    if qr_back:
        qr = qr_back
    elif qr_front:
        qr = qr_front
    
    # Parse thông tin - ưu tiên QR code trước
    info = {}
    qr_info = {}
    mrz_info = {}
    
    # Parse QR và MRZ song song (nếu có)
    if qr:
        qr_info = _parse_qr_code(qr)
    if mrz:
        mrz_info = _parse_mrz(mrz)
    
    # Merge QR code info (nếu có) - ưu tiên cho idNumber, fullName, dateOfBirth
    if qr_info:
        if qr_info.get('idNumber'):
            info['idNumber'] = qr_info['idNumber']
        if qr_info.get('fullName'):
            info['fullName'] = qr_info['fullName']
        if qr_info.get('dateOfBirth'):
            info['dateOfBirth'] = qr_info['dateOfBirth']
        # Merge các trường khác từ QR code
        for k, v in qr_info.items():
            if v and k not in ['idNumber', 'fullName', 'dateOfBirth']:
                info[k] = v
        # info['qrCode'] = qr
    
    # Merge MRZ info (nếu có) - chỉ lấy nếu QR code không có
    if mrz_info:
        if 'idNumber' not in info or not info['idNumber']:
            if mrz_info.get('idNumber'):
                info['idNumber'] = mrz_info['idNumber']
        if 'fullName' not in info or not info['fullName']:
            if mrz_info.get('fullName'):
                info['fullName'] = mrz_info['fullName']
        if 'dateOfBirth' not in info or not info['dateOfBirth']:
            if mrz_info.get('dateOfBirth'):
                info['dateOfBirth'] = mrz_info['dateOfBirth']
        # Merge các trường khác từ MRZ (như expiryDate)
        for k, v in mrz_info.items():
            if v and k not in ['idNumber', 'fullName', 'dateOfBirth']:
                if k not in info or not info[k]:
                    info[k] = v
    
    # Parse thông tin từ OCR text (mặt trước và mặt sau) song song
    if TESSERACT_AVAILABLE:
        front_text_task = asyncio.create_task(asyncio.to_thread(_ocr_text_from_bytes, front_image_data))
        back_text_task = asyncio.create_task(asyncio.to_thread(_ocr_text_from_bytes, back_image_data))
        front_text, back_text = await asyncio.gather(front_text_task, back_text_task)
        
        # Parse front và back side song song
        _parse_front_side(front_text, info)
        _parse_back_side(back_text, info)
    
    total_duration = time.time() - start_time
    print(f"⏱️  [TOTAL] Tổng thời gian xử lý: {total_duration:.2f}s")
    
    return info, mrz, qr


@app.post("/detect/national-card", response_model=NationalCardExtractResponse)
async def detect_national_card(request: NationalCardRequest):
    """Đọc QR code và MRZ từ ảnh căn cước (mặt trước và mặt sau) từ URL"""
    try:
        # Validate URLs
        if not request.frontImageUrl or not request.backImageUrl:
            return NationalCardExtractResponse(
                code=400,
                message="frontImageUrl and backImageUrl are required",
                data={}
            )
        
        # Download 2 ảnh song song (giống Go code)
        download_start = time.time()
        front_task = asyncio.create_task(_download_image(request.frontImageUrl))
        back_task = asyncio.create_task(_download_image(request.backImageUrl))
        
        front_image_data, back_image_data = await asyncio.gather(front_task, back_task)
        download_duration = time.time() - download_start
        print(f"⏱️  [Download] Tổng thời gian tải về: {download_duration:.2f}s")
        
        # Extract thông tin
        info, mrz, qr = await _extract_national_card_from_data(front_image_data, back_image_data)
        
        # Build data object
        data = {
            "info": info,
            "mrz": mrz,
            "qr": qr
        }
        
        return NationalCardExtractResponse(
            code=0,
            message="success",
            data=data
        )
        
    except Exception as e:
        # Trả về format lỗi với code != 0
        error_msg = str(e)
        if "download" in error_msg.lower():
            return NationalCardExtractResponse(
                code=400,
                message=f"Lỗi khi tải ảnh: {error_msg}",
                data={}
            )
        return NationalCardExtractResponse(
            code=500,
            message=f"Lỗi khi xử lý ảnh: {error_msg}",
            data={}
        )


if __name__ == "__main__":
    port = int(os.getenv("HTTP_PORT", 8218))
    uvicorn.run("http_server:app", host="0.0.0.0", port=port, reload=True)
