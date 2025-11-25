package client

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/liyue201/goqr"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"github.com/otiai10/gosseract/v2"
	"github.com/spf13/viper"
)

// OCRClient là client để OCR và đọc QR/MRZ
// @bind: assistant/internal/interface/provider.OCRProvider
type OCRClient struct {
	httpClient      *http.Client
	tesseractClient *gosseract.Client
	tempDir         string
}

// NewOCRClient khởi tạo OCR client
func NewOCRClient() *OCRClient {
	timeoutSeconds := viper.GetInt("ocr.timeout_seconds")
	if timeoutSeconds == 0 {
		timeoutSeconds = 60
	}

	tempDir := viper.GetString("ocr.temp_dir")
	if tempDir == "" {
		tempDir = os.TempDir() + "/assistant-ocr"
	}
	os.MkdirAll(tempDir, 0755)

	tesseractClient := gosseract.NewClient()
	tesseractClient.SetLanguage("vie", "eng")
	tesseractClient.SetPageSegMode(gosseract.PSM_AUTO)

	return &OCRClient{
		httpClient:      &http.Client{Timeout: time.Duration(timeoutSeconds) * time.Second},
		tesseractClient: tesseractClient,
		tempDir:         tempDir,
	}
}

// ExtractNationalCard trích xuất thông tin từ 2 ảnh CCCD
func (c *OCRClient) ExtractNationalCard(ctx context.Context, frontURL, backURL string) (map[string]interface{}, string, string, error) {
	startTime := time.Now()
	downloadStart := time.Now()

	var wg sync.WaitGroup
	var (
		frontData, backData               []byte
		qrFront, qrBack                   string
		mrzFront, mrzBack                 string
		mrzFromFrontText, mrzFromBackText string
		downloadFrontErr, downloadBackErr error
		qrFrontErr, qrBackErr             error
		mrzFrontErr, mrzBackErr           error
	)

	// Download 2 ảnh song song
	wg.Add(2)
	go func() {
		defer wg.Done()
		downloadFrontStart := time.Now()
		frontData, downloadFrontErr = c.downloadImage(ctx, frontURL)
		fmt.Printf("⏱️  [Download] Front image: %v\n", time.Since(downloadFrontStart))
	}()
	go func() {
		defer wg.Done()
		downloadBackStart := time.Now()
		backData, downloadBackErr = c.downloadImage(ctx, backURL)
		fmt.Printf("⏱️  [Download] Back image: %v\n", time.Since(downloadBackStart))
	}()
	wg.Wait()
	downloadDuration := time.Since(downloadStart)
	fmt.Printf("⏱️  [Download] Tổng thời gian tải về: %v\n", downloadDuration)

	if downloadFrontErr != nil {
		return nil, "", "", fmt.Errorf("download front image failed: %w", downloadFrontErr)
	}
	if downloadBackErr != nil {
		return nil, "", "", fmt.Errorf("download back image failed: %w", downloadBackErr)
	}

	// Đo thời gian infer - chạy TẤT CẢ song song: OCR, QR, MRZ cho cả 2 mặt
	inferStart := time.Now()
	wg.Add(4) // 2 OCR + 2 QR + 2 MRZ

	// OCR cả 2 ảnh song song
	// go func() {
	// 	defer wg.Done()
	// 	ocrFrontStart := time.Now()
	// 	frontText, ocrFrontErr = c.ocrBytes(frontData)
	// 	fmt.Printf("⏱️  [OCR] Front text: %v\n", time.Since(ocrFrontStart))
	// }()
	// go func() {
	// 	defer wg.Done()
	// 	ocrBackStart := time.Now()
	// 	backText, ocrBackErr = c.ocrBytes(backData)
	// 	fmt.Printf("⏱️  [OCR] Back text: %v\n", time.Since(ocrBackStart))
	// }()

	// Đọc QR code từ cả 2 mặt song song
	go func() {
		defer wg.Done()
		qrFrontStart := time.Now()
		qrFront, qrFrontErr = c.readQRCodeFromBytes(frontData)
		fmt.Printf("⏱️  [QR] Front: %v\n", time.Since(qrFrontStart))
	}()
	go func() {
		defer wg.Done()
		qrBackStart := time.Now()
		qrBack, qrBackErr = c.readQRCodeFromBytes(backData)
		fmt.Printf("⏱️  [QR] Back: %v\n", time.Since(qrBackStart))
	}()

	// Đọc MRZ từ cả 2 mặt song song
	go func() {
		defer wg.Done()
		mrzFrontStart := time.Now()
		mrzFront, mrzFrontErr = c.ocrMRZFromBytes(frontData)
		fmt.Printf("⏱️  [MRZ] Front: %v\n", time.Since(mrzFrontStart))
	}()
	go func() {
		defer wg.Done()
		mrzBackStart := time.Now()
		mrzBack, mrzBackErr = c.ocrMRZFromBytes(backData)
		fmt.Printf("⏱️  [MRZ] Back: %v\n", time.Since(mrzBackStart))
	}()

	wg.Wait()
	inferDuration := time.Since(inferStart)
	fmt.Printf("⏱️  [Infer] Tổng thời gian (QR + MRZ): %v\n", inferDuration)

	// Extract MRZ từ text song song (nếu cần fallback)
	extractMRZStart := time.Now()
	wg.Add(2)
	go func() {
		defer wg.Done()
		mrzFromBackText = c.extractMRZFromText(string(backData))
	}()
	go func() {
		defer wg.Done()
		mrzFromFrontText = c.extractMRZFromText(string(frontData))
	}()
	wg.Wait()
	fmt.Printf("⏱️  [Extract MRZ] Từ text: %v\n", time.Since(extractMRZStart))

	// Xử lý kết quả MRZ (ưu tiên mặt sau)
	var mrz string
	if mrzBackErr == nil && mrzBack != "" {
		mrz = mrzBack
	} else if mrzFrontErr == nil && mrzFront != "" {
		mrz = mrzFront
	} else {
		// fallback extract từ OCR text (ưu tiên backText vì MRZ thường ở mặt sau)
		if mrzFromBackText != "" {
			mrz = mrzFromBackText
		} else if mrzFromFrontText != "" {
			mrz = mrzFromFrontText
		}
	}

	// Xử lý kết quả QR code (ưu tiên mặt sau)
	var qr string
	if qrBackErr == nil && qrBack != "" {
		qr = qrBack
	} else if qrFrontErr == nil && qrFront != "" {
		qr = qrFront
	}

	// Đo thời gian parse
	parseStart := time.Now()
	// Parse thông tin - ưu tiên QR code trước
	info := make(map[string]interface{})
	var qrInfo, mrzInfo map[string]interface{}

	// Parse QR và MRZ song song (nếu có)
	parseQRMRZStart := time.Now()
	if qr != "" || mrz != "" {
		wg.Add(2)
		go func() {
			defer wg.Done()
			if qr != "" {
				parseQRStart := time.Now()
				qrInfo = c.parseQRCode(qr)
				fmt.Printf("⏱️  [Parse] QR code: %v\n", time.Since(parseQRStart))
			}
		}()
		go func() {
			defer wg.Done()
			if mrz != "" {
				parseMRZStart := time.Now()
				mrzInfo = c.parseMRZ(mrz)
				fmt.Printf("⏱️  [Parse] MRZ: %v\n", time.Since(parseMRZStart))
			}
		}()
		wg.Wait()
	}
	fmt.Printf("⏱️  [Parse] QR + MRZ tổng: %v\n", time.Since(parseQRMRZStart))

	// Merge QR code info (nếu có) - ưu tiên cho idNumber, fullName, dateOfBirth
	if qrInfo != nil {
		if idNumber, ok := qrInfo["idNumber"]; ok && idNumber != nil && idNumber != "" {
			info["idNumber"] = idNumber
		}
		if fullName, ok := qrInfo["fullName"]; ok && fullName != nil && fullName != "" {
			info["fullName"] = fullName
		}
		if dateOfBirth, ok := qrInfo["dateOfBirth"]; ok && dateOfBirth != nil && dateOfBirth != "" {
			info["dateOfBirth"] = dateOfBirth
		}
		// Merge các trường khác từ QR code
		for k, v := range qrInfo {
			if v != nil && v != "" && k != "idNumber" && k != "fullName" && k != "dateOfBirth" {
				info[k] = v
			}
		}
		info["qrCode"] = qr
	}

	// Merge MRZ info (nếu có) - chỉ lấy nếu QR code không có
	if mrzInfo != nil {
		if info["idNumber"] == nil || info["idNumber"] == "" {
			if idNumber, ok := mrzInfo["idNumber"]; ok && idNumber != nil && idNumber != "" {
				info["idNumber"] = idNumber
			}
		}
		if info["fullName"] == nil || info["fullName"] == "" {
			if fullName, ok := mrzInfo["fullName"]; ok && fullName != nil && fullName != "" {
				info["fullName"] = fullName
			}
		}
		if info["dateOfBirth"] == nil || info["dateOfBirth"] == "" {
			if dateOfBirth, ok := mrzInfo["dateOfBirth"]; ok && dateOfBirth != nil && dateOfBirth != "" {
				info["dateOfBirth"] = dateOfBirth
			}
		}
		// Merge các trường khác từ MRZ (như expiryDate)
		for k, v := range mrzInfo {
			if v != nil && v != "" && k != "idNumber" && k != "fullName" && k != "dateOfBirth" {
				if existing, exists := info[k]; !exists || existing == nil || existing == "" {
					info[k] = v
				}
			}
		}
	}

	// Parse thông tin từ OCR text (mặt trước và mặt sau) song song
	parseTextStart := time.Now()
	wg.Add(2)
	go func() {
		defer wg.Done()
		parseFrontStart := time.Now()
		c.parseFrontSide(string(frontData), info)
		fmt.Printf("⏱️  [Parse] Front side: %v\n", time.Since(parseFrontStart))
	}()
	go func() {
		defer wg.Done()
		parseBackStart := time.Now()
		c.parseBackSide(string(backData), info)
		fmt.Printf("⏱️  [Parse] Back side: %v\n", time.Since(parseBackStart))
	}()
	wg.Wait()
	fmt.Printf("⏱️  [Parse] Front + Back side tổng: %v\n", time.Since(parseTextStart))

	parseDuration := time.Since(parseStart)
	fmt.Printf("⏱️  [Parse] Tổng thời gian parse: %v\n", parseDuration)

	totalDuration := time.Since(startTime)
	fmt.Printf("⏱️  [TOTAL] Tổng thời gian xử lý: %v\n", totalDuration)

	return info, mrz, qr, nil
}

// downloadImage tải ảnh từ URL
func (c *OCRClient) downloadImage(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status=%d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}

// ocrBytes OCR từ bytes
// Tạo client riêng để có thể chạy song song an toàn
func (c *OCRClient) ocrBytes(data []byte) (string, error) {
	path := filepath.Join(c.tempDir, fmt.Sprintf("ocr_%d.jpg", time.Now().UnixNano()))
	defer os.Remove(path)
	if err := os.WriteFile(path, data, 0644); err != nil {
		return "", err
	}
	// Tạo client riêng để tránh lock khi chạy song song
	client := gosseract.NewClient()
	defer client.Close()
	client.SetLanguage("vie", "eng")
	client.SetPageSegMode(gosseract.PSM_AUTO)
	client.SetImage(path)
	return client.Text()
}

// ocrMRZFromBytes OCR MRZ từ bytes
func (c *OCRClient) ocrMRZFromBytes(data []byte) (string, error) {
	path := filepath.Join(c.tempDir, fmt.Sprintf("mrz_%d.png", time.Now().UnixNano()))
	defer os.Remove(path)
	if err := os.WriteFile(path, data, 0644); err != nil {
		return "", err
	}

	// Crop vùng MRZ
	cropPath, err := c.cropMRZRegion(path)
	// cropPath := path
	if err != nil {
		cropPath = path
	} else {
		defer os.Remove(cropPath)
	}

	client := gosseract.NewClient()
	defer client.Close()
	client.SetImage(cropPath)
	client.SetPageSegMode(gosseract.PSM_SINGLE_BLOCK)
	client.SetLanguage("eng")
	client.SetWhitelist("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789<")
	text, err := client.Text()
	if err != nil {
		return "", err
	}

	mrz := c.cleanMRZText(text)
	if c.isValidMRZ(mrz) {
		return mrz, nil
	}
	return "", fmt.Errorf("invalid MRZ")
}

// readQRCodeFromBytes đọc QR code từ bytes
func (c *OCRClient) readQRCodeFromBytes(data []byte) (string, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return "", err
	}

	// Chuyển sang BinaryBitmap để đọc QR
	luminanceSource := gozxing.NewLuminanceSourceFromImage(img)
	bitmap, _ := gozxing.NewBinaryBitmap(gozxing.NewHybridBinarizer(luminanceSource))

	// Tạo QRCodeReader
	qrReader := qrcode.NewQRCodeReader()
	result, err := qrReader.Decode(bitmap, nil)
	if err != nil {
		return "", err
	}
	return result.GetText(), nil
}

// Các hàm cropMRZRegion, cleanMRZText, isValidMRZ, extractMRZFromText, parseNationalCardInfo,
// parseQRCode vẫn giữ nguyên logic như file cũ, chỉ cần bỏ đi các read QR/MRZ lặp lại.

// parseNationalCardInfo chỉ parse MRZ và OCR text, QR code merge riêng sau

// cropMRZRegion - Crop ảnh để chỉ lấy phần MRZ (phần dưới 30% của ảnh)
func (c *OCRClient) cropMRZRegion(imagePath string) (string, error) {
	// Đọc ảnh
	file, err := os.Open(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to open image: %w", err)
	}
	defer file.Close()

	// Decode ảnh
	img, format, err := image.Decode(file)
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}

	// Lấy kích thước ảnh
	bounds := img.Bounds()

	cropY := bounds.Dy() * 6 / 10
	cropRect := image.Rect(bounds.Min.X, bounds.Min.Y+cropY, bounds.Max.X, bounds.Max.Y)

	croppedImg := img.(interface {
		SubImage(r image.Rectangle) image.Image
	}).SubImage(cropRect)

	// Lưu ảnh đã crop
	croppedPath := filepath.Join(c.tempDir, fmt.Sprintf("mrz_%d.png", time.Now().UnixNano()))
	croppedFile, err := os.Create(croppedPath)
	if err != nil {
		return "", fmt.Errorf("failed to create cropped image file: %w", err)
	}
	defer croppedFile.Close()

	// Encode và lưu ảnh
	switch format {
	case "jpeg", "jpg":
		err = jpeg.Encode(croppedFile, croppedImg, nil)
	case "png":
		err = png.Encode(croppedFile, croppedImg)
	default:
		// Default to JPEG
		err = jpeg.Encode(croppedFile, croppedImg, nil)
	}

	if err != nil {
		return "", fmt.Errorf("failed to encode cropped image: %w", err)
	}

	return croppedPath, nil
}

// cleanMRZText - Làm sạch text MRZ, loại bỏ ký tự không hợp lệ và khoảng trắng
func (c *OCRClient) cleanMRZText(text string) string {
	lines := strings.Split(text, "\n")
	var cleanedLines []string

	mrzRegex := regexp.MustCompile(`[A-Z0-9<]{22,}`) // MRZ dòng thường có ít nhất 25 ký tự

	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Remove tất cả khoảng trắng trong MRZ
		line = strings.ReplaceAll(line, " ", "")
		// Chỉ giữ dòng có format MRZ (nhiều ký tự A-Z, 0-9, <)
		if len(line) >= 22 && mrzRegex.MatchString(line) && strings.Count(line, "<") > 1 {
			cleanedLines = append(cleanedLines, line)
		}
	}

	return strings.Join(cleanedLines, "\n")
}

// isValidMRZ - Kiểm tra xem text có phải MRZ hợp lệ không
func (c *OCRClient) isValidMRZ(text string) bool {
	lines := strings.Split(strings.TrimSpace(text), "\n")

	// MRZ thường có 2-3 dòng
	if len(lines) < 2 {
		return false
	}

	// Mỗi dòng phải có ít nhất 25 ký tự và chứa nhiều ký tự <
	validLines := 0
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) >= 22 && strings.Count(line, "<") > 1 {
			validLines++
		}
	}

	// Phải có ít nhất 2 dòng hợp lệ
	return validLines >= 2
}

// extractMRZFromText - Trích xuất MRZ từ text đã OCR (MRZ thường có 3 dòng)
// Format MRZ:
// - Dòng 1: IDVNM... (chứa ID và thông tin)
// - Dòng 2: 0011189M4011187VNM... (chứa DOB, gender, expiry)
// - Dòng 3: NGO<<QUANG<HIEU... (chứa tên)
func (c *OCRClient) extractMRZFromText(text string) string {
	lines := strings.Split(text, "\n")
	var mrzLines []string

	// Tìm các dòng có format MRZ
	// MRZ đặc trưng: chứa nhiều ký tự <, chữ số, và chữ cái in hoa
	mrzRegex := regexp.MustCompile(`[A-Z0-9<]{20,}`)

	// Tìm đoạn liên tiếp có 3 dòng MRZ
	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		// Remove khoảng trắng để dễ detect
		lineClean := strings.ReplaceAll(line, " ", "")

		// Kiểm tra nếu dòng này có format MRZ
		if len(lineClean) >= 20 && mrzRegex.MatchString(lineClean) && strings.Count(lineClean, "<") >= 2 {
			// Kiểm tra xem có phải là đoạn 3 dòng MRZ liên tiếp không
			if i+2 < len(lines) {
				line2 := strings.TrimSpace(lines[i+1])
				line2Clean := strings.ReplaceAll(line2, " ", "")
				line3 := strings.TrimSpace(lines[i+2])
				line3Clean := strings.ReplaceAll(line3, " ", "")

				// Kiểm tra cả 3 dòng đều có format MRZ
				if len(line2Clean) >= 20 && mrzRegex.MatchString(line2Clean) && strings.Count(line2Clean, "<") >= 2 &&
					len(line3Clean) >= 20 && mrzRegex.MatchString(line3Clean) && strings.Count(line3Clean, "<") >= 2 {
					// Tìm thấy đoạn 3 dòng MRZ
					mrzLines = []string{lineClean, line2Clean, line3Clean}
					break
				}
			}

			// Nếu không tìm thấy 3 dòng liên tiếp, thử tìm 2 dòng
			if len(mrzLines) == 0 && i+1 < len(lines) {
				line2 := strings.TrimSpace(lines[i+1])
				line2Clean := strings.ReplaceAll(line2, " ", "")
				if len(line2Clean) >= 20 && mrzRegex.MatchString(line2Clean) && strings.Count(line2Clean, "<") >= 2 {
					mrzLines = []string{lineClean, line2Clean}
				}
			}
		}
	}

	// Nếu tìm thấy 3 dòng, trả về 3 dòng
	if len(mrzLines) == 3 {
		return strings.Join(mrzLines, "\n")
	}

	// Nếu chỉ tìm thấy 2 dòng, trả về 2 dòng
	if len(mrzLines) >= 2 {
		return strings.Join(mrzLines, "\n")
	}

	return ""
}

// parseNationalCardInfo - Parse thông tin từ text đã OCR
// Hàm này không còn được dùng trực tiếp, logic đã được chuyển vào ExtractNationalCard
// Giữ lại để tương thích nếu có code khác gọi
func (c *OCRClient) parseNationalCardInfo(frontText, backText, mrz string) map[string]interface{} {
	info := make(map[string]interface{})

	// Parse MRZ trước (nếu có) vì MRZ chính xác hơn
	if mrz != "" {
		mrzInfo := c.parseMRZ(mrz)
		for k, v := range mrzInfo {
			info[k] = v
		}
	}

	// Parse thông tin từ mặt trước
	c.parseFrontSide(frontText, info)

	// Parse thông tin từ mặt sau
	c.parseBackSide(backText, info)

	return info
}

// parseFrontSide - Parse thông tin từ mặt trước căn cước
func (c *OCRClient) parseFrontSide(text string, info map[string]interface{}) {
	lines := strings.Split(text, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Tìm số CMND/CCCD (thường có format: số 9-12 chữ số)
		if idNumber := c.extractIDNumber(line); idNumber != "" && info["idNumber"] == nil {
			info["idNumber"] = idNumber
		}

		// Tìm họ tên (thường có từ khóa "Họ và tên" hoặc "Họ tên")
		// Chỉ set nếu chưa có từ QR code hoặc MRZ
		if (info["fullName"] == nil || info["fullName"] == "") &&
			strings.Contains(strings.ToLower(line), "họ") && strings.Contains(strings.ToLower(line), "tên") {
			// Extract tên sau từ khóa
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				name := strings.TrimSpace(parts[1])
				if name != "" {
					info["fullName"] = name
				}
			}
		}

		// Tìm ngày sinh
		// Chỉ set nếu chưa có từ QR code hoặc MRZ
		if (info["dateOfBirth"] == nil || info["dateOfBirth"] == "") &&
			(strings.Contains(strings.ToLower(line), "ngày sinh") || strings.Contains(strings.ToLower(line), "sinh ngày")) {
			dob := c.extractDate(line)
			if dob != "" {
				info["dateOfBirth"] = dob
			}
		}

		// Tìm giới tính
		if strings.Contains(strings.ToLower(line), "giới tính") || strings.Contains(strings.ToLower(line), "nam") || strings.Contains(strings.ToLower(line), "nữ") {
			gender := c.extractGender(line)
			if gender != "" {
				info["gender"] = gender
			}
		}

		// Tìm quốc tịch
		if strings.Contains(strings.ToLower(line), "quốc tịch") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				nationality := strings.TrimSpace(parts[1])
				if nationality != "" {
					info["nationality"] = nationality
				}
			}
		}

		// Tìm quê quán
		if strings.Contains(strings.ToLower(line), "quê quán") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				placeOfOrigin := strings.TrimSpace(parts[1])
				if placeOfOrigin != "" {
					info["placeOfOrigin"] = placeOfOrigin
				}
			}
		}

		// Tìm nơi thường trú
		if strings.Contains(strings.ToLower(line), "nơi thường trú") || strings.Contains(strings.ToLower(line), "thường trú") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				placeOfResidence := strings.TrimSpace(parts[1])
				if placeOfResidence != "" {
					info["placeOfResidence"] = placeOfResidence
				}
			}
		}

		// Tìm dân tộc
		if strings.Contains(strings.ToLower(line), "dân tộc") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				ethnic := strings.TrimSpace(parts[1])
				if ethnic != "" {
					info["ethnic"] = ethnic
				}
			}
		}

		// Tìm tôn giáo
		if strings.Contains(strings.ToLower(line), "tôn giáo") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				religion := strings.TrimSpace(parts[1])
				if religion != "" {
					info["religion"] = religion
				}
			}
		}
	}
}

// parseBackSide - Parse thông tin từ mặt sau căn cước
func (c *OCRClient) parseBackSide(text string, info map[string]interface{}) {
	lines := strings.Split(text, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Tìm ngày cấp
		// Chỉ set nếu chưa có từ QR code
		if (info["issueDate"] == nil || info["issueDate"] == "") &&
			(strings.Contains(strings.ToLower(line), "ngày cấp") || strings.Contains(strings.ToLower(line), "cấp ngày")) {
			issueDate := c.extractDate(line)
			if issueDate != "" {
				info["issueDate"] = issueDate
			}
		}

		// Tìm nơi cấp
		if strings.Contains(strings.ToLower(line), "nơi cấp") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				issuePlace := strings.TrimSpace(parts[1])
				if issuePlace != "" {
					info["issuePlace"] = issuePlace
				}
			}
		}

		// Tìm đặc điểm nhận dạng
		if strings.Contains(strings.ToLower(line), "đặc điểm") || strings.Contains(strings.ToLower(line), "nhận dạng") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				features := strings.TrimSpace(parts[1])
				if features != "" {
					info["features"] = features
				}
			}
		}

		// Tìm địa chỉ
		if strings.Contains(strings.ToLower(line), "địa chỉ") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				address := strings.TrimSpace(parts[1])
				if address != "" {
					info["address"] = address
				}
			}
		}
	}
}

// parseMRZ - Parse MRZ string thành thông tin
func (c *OCRClient) parseMRZ(mrz string) map[string]interface{} {
	result := make(map[string]interface{})

	// MRZ format cho căn cước Việt Nam có thể có 2 hoặc 3 dòng
	lines := strings.Split(strings.TrimSpace(mrz), "\n")
	if len(lines) < 2 {
		return result
	}

	// Parse line 1, line 2, và line 3 (nếu có) của MRZ
	// Remove khoảng trắng từ mỗi dòng để đảm bảo parse chính xác
	line1 := strings.ReplaceAll(strings.TrimSpace(lines[0]), " ", "")
	line2 := strings.ReplaceAll(strings.TrimSpace(lines[1]), " ", "")
	var line3 string
	if len(lines) >= 3 {
		line3 = strings.ReplaceAll(strings.TrimSpace(lines[2]), " ", "")
	}

	// Extract ID number từ MRZ line 1
	// Tìm vị trí << trong line1, lấy 12 số cuối trước đó
	if idx := strings.Index(line1, "<<"); idx > 0 {
		// Lấy phần trước <<
		before := line1[:idx]
		// Tìm tất cả số trong phần đó
		digits := regexp.MustCompile(`\d+`).FindAllString(before, -1)
		if len(digits) > 0 {
			// Lấy số cuối cùng và kiểm tra nếu có 12 chữ số
			lastDigits := digits[len(digits)-1]
			if len(lastDigits) >= 12 {
				// Lấy 12 số cuối
				result["idNumber"] = lastDigits[len(lastDigits)-12:]
			} else if len(lastDigits) == 12 {
				result["idNumber"] = lastDigits
			}
		}
	}
	// Fallback: nếu không tìm thấy <<, tìm số 12 chữ số bắt đầu bằng 0
	if result["idNumber"] == nil {
		idRegex := regexp.MustCompile(`\d{12}`)
		matches := idRegex.FindAllString(line1, -1)
		if len(matches) > 0 {
			for _, match := range matches {
				if match[0] == '0' {
					result["idNumber"] = match
					break
				}
			}
			if result["idNumber"] == nil {
				result["idNumber"] = matches[0]
			}
		}
	}

	// Extract date of birth từ line 2 (format: YYMMDD)
	// Ví dụ: 0011189M4011187VNM -> tìm 001118 (số đầu tiên 6 chữ số) -> 18/11/2000
	dobRegex := regexp.MustCompile(`\d{6}`)
	if matches := dobRegex.FindAllString(line2, -1); len(matches) > 0 {
		// Lấy số đầu tiên (date of birth)
		match := matches[0]
		if len(match) == 6 {
			yy := match[0:2]
			month := match[2:4]
			day := match[4:6]
			monthInt := parseInt(month)
			dayInt := parseInt(day)

			if monthInt >= 1 && monthInt <= 12 && dayInt >= 1 && dayInt <= 31 {
				// Xác định century: YY < 50 -> 20YY, YY >= 50 -> 19YY
				var year string
				if parseInt(yy) < 50 {
					year = "20" + yy
				} else {
					year = "19" + yy
				}
				result["dateOfBirth"] = fmt.Sprintf("%s/%s/%s", day, month, year)
			}
		}
	}

	// Extract full name từ line 3 (nếu có)
	// Ví dụ: NGUYEN<<VAN<A<<<<<<< -> NGUYEN VAN A
	if line3 != "" {
		// Loại bỏ các ký tự không phải chữ cái và < ở cuối
		// Tìm phần chứa tên (trước khi có ký tự lạ như c, cc)
		nameRegex := regexp.MustCompile(`^([A-Z<]+)`)
		if nameMatch := nameRegex.FindString(line3); nameMatch != "" {
			// Loại bỏ các < ở cuối
			nameMatch = strings.TrimRight(nameMatch, "<")
			// Tách bằng < và lọc các phần có chữ cái
			nameParts := strings.Split(nameMatch, "<")
			var cleanParts []string
			for _, part := range nameParts {
				part = strings.TrimSpace(part)
				// Chỉ lấy phần chứa chữ cái in hoa (loại bỏ phần rỗng)
				if len(part) > 0 && regexp.MustCompile(`^[A-Z]+$`).MatchString(part) {
					cleanParts = append(cleanParts, part)
				}
			}
			if len(cleanParts) > 0 {
				result["fullName"] = strings.Join(cleanParts, " ")
			}
		}
	}

	// Extract expiry date từ line 2 (format: YYMMDD)
	// Ví dụ: 0011189M4011187VNM -> tìm 401118 (số thứ 2, sau M) -> 18/11/2040
	// Trong MRZ line 2, expiry date là số 6 chữ số thứ 2 (sau date of birth và M)
	expiryRegex := regexp.MustCompile(`\d{6}`)
	expiryMatches := expiryRegex.FindAllString(line2, -1)
	if len(expiryMatches) >= 2 {
		// Lấy số thứ 2 (expiry date)
		match := expiryMatches[1]
		if len(match) == 6 {
			yy := match[0:2]
			month := match[2:4]
			day := match[4:6]

			// Validate month và day
			monthInt := parseInt(month)
			dayInt := parseInt(day)

			if monthInt >= 1 && monthInt <= 12 && dayInt >= 1 && dayInt <= 31 {
				// Expiry date luôn là năm tương lai, nên dùng 20YY
				year := "20" + yy
				result["expiryDate"] = fmt.Sprintf("%s/%s/%s", day, month, year)
			}
		}
	} else if len(expiryMatches) == 1 {
		// Nếu chỉ có 1 số, kiểm tra xem có phải expiry date không (YY >= 30)
		match := expiryMatches[0]
		if len(match) == 6 {
			yy := match[0:2]
			month := match[2:4]
			day := match[4:6]
			monthInt := parseInt(month)
			dayInt := parseInt(day)
			yyInt := parseInt(yy)

			if monthInt >= 1 && monthInt <= 12 && dayInt >= 1 && dayInt <= 31 {
				// Nếu YY >= 30, có thể là expiry date
				if yyInt >= 30 {
					year := "20" + yy
					result["expiryDate"] = fmt.Sprintf("%s/%s/%s", day, month, year)
				}
			}
		}
	}

	return result
}

// Helper functions để extract thông tin từ text
func (c *OCRClient) extractIDNumber(text string) string {
	// Tìm số CMND/CCCD: 9-12 chữ số liên tiếp
	idRegex := regexp.MustCompile(`\d{9,12}`)
	matches := idRegex.FindAllString(text, -1)
	if len(matches) > 0 {
		return matches[0]
	}
	return ""
}

func (c *OCRClient) extractDate(text string) string {
	// Tìm ngày tháng năm: dd/mm/yyyy hoặc dd-mm-yyyy
	dateRegex := regexp.MustCompile(`(\d{1,2})[/-](\d{1,2})[/-](\d{4})`)
	matches := dateRegex.FindStringSubmatch(text)
	if len(matches) == 4 {
		return fmt.Sprintf("%s/%s/%s", matches[1], matches[2], matches[3])
	}
	return ""
}

func (c *OCRClient) extractGender(text string) string {
	textLower := strings.ToLower(text)
	if strings.Contains(textLower, "nam") {
		return "Nam"
	}
	if strings.Contains(textLower, "nữ") {
		return "Nữ"
	}
	return ""
}

// readQRCodeFromImage - Đọc QR code từ image bytes
func (c *OCRClient) readQRCodeFromImage(imageData []byte) (string, error) {
	// Decode ảnh từ bytes
	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		// Nếu decode trực tiếp thất bại, lưu vào file tạm rồi decode
		imagePath := filepath.Join(c.tempDir, fmt.Sprintf("qr_%d.png", time.Now().UnixNano()))
		defer os.Remove(imagePath)

		if err := os.WriteFile(imagePath, imageData, 0644); err != nil {
			return "", fmt.Errorf("failed to save image: %w", err)
		}

		file, err := os.Open(imagePath)
		if err != nil {
			return "", fmt.Errorf("failed to open image: %w", err)
		}
		defer file.Close()

		img, _, err = image.Decode(file)
		if err != nil {
			return "", fmt.Errorf("failed to decode image: %w", err)
		}
	}

	// Convert sang grayscale
	gray := image.NewGray(img.Bounds())
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			gray.Set(x, y, color.GrayModel.Convert(img.At(x, y)))
		}
	}

	// Decode QR code
	codes, err := goqr.Recognize(gray)
	if err != nil {
		return "", fmt.Errorf("failed to recognize QR code: %w", err)
	}

	if len(codes) == 0 {
		return "", fmt.Errorf("no QR code found")
	}

	// Trả về text của QR code đầu tiên
	return string(codes[0].Payload), nil
}

// parseQRCode - Parse QR code string thành thông tin
// Format QR code: ID||FullName|DOB|Gender|Address|IssueDate|...
// DOB format: DDMMYYYY (18112000 -> 18/11/2000)
// IssueDate format: DDMMYYYY (17092025 -> 17/09/2025)
func (c *OCRClient) parseQRCode(qrText string) map[string]interface{} {
	result := make(map[string]interface{})

	// Remove brackets nếu có
	qrText = strings.TrimSpace(qrText)
	if strings.HasPrefix(qrText, "[") && strings.HasSuffix(qrText, "]") {
		qrText = qrText[1 : len(qrText)-1]
	}

	// Split by |
	parts := strings.Split(qrText, "|")
	if len(parts) < 5 {
		return result // QR code không đúng format
	}

	// Parse các trường theo thứ tự
	// [0]: ID number
	if len(parts) > 0 && parts[0] != "" {
		result["idNumber"] = parts[0]
	}

	// [1]: Empty (||)
	// [2]: Full name
	if len(parts) > 2 && parts[2] != "" {
		result["fullName"] = parts[2]
	}

	// [3]: Date of birth (format: DDMMYYYY hoặc YYMMDD -> dd/mm/yyyy)
	if len(parts) > 3 && parts[3] != "" {
		dob := parts[3]
		if len(dob) == 8 {
			// Format: DDMMYYYY (ví dụ: 18112000 -> 18/11/2000)
			day := dob[0:2]
			month := dob[2:4]
			year := dob[4:8]
			result["dateOfBirth"] = fmt.Sprintf("%s/%s/%s", day, month, year)
		} else if len(dob) == 6 {
			// Format: YYMMDD
			yy := dob[0:2]
			month := dob[2:4]
			day := dob[4:6]
			// Xác định century: nếu YY < 50 thì là 20YY, ngược lại là 19YY
			var year string
			if yyInt := parseInt(yy); yyInt < 50 {
				year = "20" + yy
			} else {
				year = "19" + yy
			}
			result["dateOfBirth"] = fmt.Sprintf("%s/%s/%s", day, month, year)
		}
	}

	// [4]: Gender
	if len(parts) > 4 && parts[4] != "" {
		result["gender"] = parts[4]
	}

	// [5]: Address
	if len(parts) > 5 && parts[5] != "" {
		result["address"] = parts[5]
		result["placeOfResidence"] = parts[5] // Có thể là nơi thường trú
	}

	// [6]: Issue date (format: DDMMYYYY hoặc YYMMDD -> dd/mm/yyyy)
	if len(parts) > 6 && parts[6] != "" {
		issueDate := parts[6]
		if len(issueDate) == 8 {
			// Format: DDMMYYYY (ví dụ: 17092025 -> 17/09/2025)
			day := issueDate[0:2]
			month := issueDate[2:4]
			year := issueDate[4:8]
			result["issueDate"] = fmt.Sprintf("%s/%s/%s", day, month, year)
		} else if len(issueDate) == 6 {
			// Format: YYMMDD
			yy := issueDate[0:2]
			month := issueDate[2:4]
			day := issueDate[4:6]
			// Issue date thường là năm gần đây, nên giả định 20YY
			year := "20" + yy
			result["issueDate"] = fmt.Sprintf("%s/%s/%s", day, month, year)
		}
	}

	return result
}

// parseInt - Helper để parse string thành int
func parseInt(s string) int {
	var result int
	for _, c := range s {
		if c >= '0' && c <= '9' {
			result = result*10 + int(c-'0')
		}
	}
	return result
}
