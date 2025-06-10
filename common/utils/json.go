package _utils

import (
	_routes "common/routes"
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/protoadapt"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ParseBody[T any](c *gin.Context) (*T, error) {
	var entity T

	// Đọc toàn bộ body từ request
	if err := c.ShouldBindJSON(&entity); err != nil {
		return nil, &_routes.Except{
			Code:    400,
			Message: err.Error(),
		}
	}

	return &entity, nil
}

func ParseBodyWithValidator[T any](c *gin.Context, target *T) error {
	if err := c.ShouldBindJSON(target); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			errorsMap := map[string]string{}
			for _, fe := range ve {
				field := strings.ToLower(fe.Field())

				typ := reflect.TypeOf(*target)
				strF := fe.StructField()
				f, ok := typ.FieldByName(strF)
				if ok {
					if tag := f.Tag.Get("json"); tag != "" {
						field = tag
					}
				}

				switch fe.Tag() {
				case "required":
					errorsMap[field] = "Trường " + field + " là bắt buộc"
				case "email":
					errorsMap[field] = "Trường " + field + " phải là email hợp lệ"
				case "min":
					errorsMap[field] = "Trường " + field + " phải có ít nhất " + fe.Param() + " ký tự"
				case "max":
					errorsMap[field] = "Trường " + field + " không được vượt quá " + fe.Param() + " ký tự"
				case "gte":
					errorsMap[field] = "Trường " + field + " phải lớn hơn hoặc bằng " + fe.Param()
				case "gt":
					errorsMap[field] = "Trường " + field + " phải lớn hơn " + fe.Param()
				case "lte":
					errorsMap[field] = "Trường " + field + " phải nhỏ hơn hoặc bằng " + fe.Param()
				case "lt":
					errorsMap[field] = "Trường " + field + " phải nhỏ hơn " + fe.Param()
				case "len":
					errorsMap[field] = "Trường " + field + " phải có độ dài chính xác là " + fe.Param()
				case "oneof":
					errorsMap[field] = "Trường " + field + " phải là một trong các giá trị: " + fe.Param()
				case "url":
					errorsMap[field] = "Trường " + field + " phải là URL hợp lệ"
				case "uuid":
					errorsMap[field] = "Trường " + field + " phải là UUID hợp lệ"
				case "numeric":
					errorsMap[field] = "Trường " + field + " phải là số"
				case "alpha":
					errorsMap[field] = "Trường " + field + " chỉ được chứa chữ cái"
				case "alphanum":
					errorsMap[field] = "Trường " + field + " chỉ được chứa chữ cái và số"
				case "startswith":
					errorsMap[field] = "Trường " + field + " phải bắt đầu với: " + fe.Param()
				case "endswith":
					errorsMap[field] = "Trường " + field + " phải kết thúc với: " + fe.Param()
				default:
					errorsMap[field] = "Trường " + field + " không hợp lệ"
				}
			}

			return &_routes.Except{
				Code:    400,
				Message: "Validation failed",
				Errors:  errorsMap,
			}
		}

		// Nếu không phải lỗi validator (ví dụ: JSON không hợp lệ)
		return &_routes.Except{
			Code:    400,
			Message: err.Error(),
		}
	}

	return nil
}

type fieldInfo struct {
	StructField reflect.StructField
	IndexPath   []int // đường dẫn tới field, dùng để gọi tValue.FieldByIndex
}

func getAllFields(t reflect.Type, indexPath []int) []fieldInfo {
	var fields []fieldInfo
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		path := append(indexPath, i)

		if f.Anonymous && f.Type.Kind() == reflect.Struct {
			fields = append(fields, getAllFields(f.Type, path)...)
		} else {
			fields = append(fields, fieldInfo{StructField: f, IndexPath: path})
		}
	}
	return fields
}

func ParseQuery2[T any](c *gin.Context, target *T) error {
	tType := reflect.TypeOf(*target)
	tValue := reflect.ValueOf(target).Elem()

	query := c.Request.URL.Query()

	fields := getAllFields(tType, nil)

	for _, field := range fields {
		valueField := tValue.FieldByIndex(field.IndexPath)

		if !valueField.CanSet() {
			continue
		}

		queryKey := field.StructField.Tag.Get("form")
		if queryKey == "" {
			continue
		}

		raw := query.Get(queryKey)
		valid := field.StructField.Tag.Get("valid")
		if raw == "" && valid == "required" {
			return fmt.Errorf("%s là bắt buộc", queryKey)
		}
		if raw == "" {
			continue
		}

		parser := field.StructField.Tag.Get("parser")
		if parser == "uint64s" {
			rawValues := strings.Split(raw, ",")
			// slice := reflect.MakeSlice(uint64, 0, len(rawValues))
			slice := make([]uint64, len(rawValues)) // [0 0 0]

			for i, s := range rawValues {
				s = strings.TrimSpace(s)
				num, err := strconv.ParseUint(s, 10, 64)
				if err != nil {
					return fmt.Errorf("invalid uint64 in %s: %v", queryKey, err)
				}
				slice[i] = num
				// slice = reflect.Append(slice, reflect.ValueOf(num))
			}
			// values := reflect.ValueOf(&slice)
			valueField.Set(reflect.ValueOf(slice))
			continue
		}

		// Handle *time.Time
		if valueField.Kind() == reflect.Ptr && valueField.Type().Elem() == reflect.TypeOf(time.Time{}) {
			parsedTime, err := time.Parse(time.RFC3339, raw)
			if err != nil {
				return fmt.Errorf("invalid time for %s: %v", queryKey, err)
			}
			valueField.Set(reflect.ValueOf(&parsedTime))
			continue
		}

		if valueField.Kind() == reflect.Ptr && valueField.Type().Elem() == reflect.TypeOf(uint64(0)) {
			parsed, err := strconv.ParseUint(raw, 10, 64)
			if err != nil {
				return fmt.Errorf("invalid int for %s: %v", queryKey, err)
			}
			valueField.Set(reflect.ValueOf(&parsed))
			continue
		}

		if valueField.Kind() == reflect.Ptr && valueField.Type().Elem() == reflect.TypeOf(true) {
			parsed, err := strconv.ParseBool(raw)
			if err != nil {
				return fmt.Errorf("invalid int for %s: %v", queryKey, err)
			}
			valueField.Set(reflect.ValueOf(&parsed))
			continue
		}

		rawValues := strings.Split(raw, ",")
		kind := valueField.Kind()
		switch kind {
		case reflect.String:
			valueField.SetString(raw)

		case reflect.Float64:
			parsed, err := strconv.ParseFloat(raw, 64)
			if err != nil {
				return fmt.Errorf("invalid float for %s: %v", queryKey, err)
			}
			valueField.SetFloat(parsed)

		case reflect.Int, reflect.Int64:
			parsed, err := strconv.ParseInt(raw, 10, 64)
			if err != nil {
				return fmt.Errorf("invalid int for %s: %v", queryKey, err)
			}
			valueField.SetInt(parsed)
		case reflect.Uint, reflect.Uint64:
			parsed, err := strconv.ParseUint(raw, 10, 64)
			if err != nil {
				return fmt.Errorf("invalid int for %s: %v", queryKey, err)
			}
			valueField.SetUint(parsed)

		case reflect.Slice:
			elemKind := valueField.Type().Elem().Kind()
			slice := reflect.MakeSlice(valueField.Type(), 0, len(rawValues))
			for _, s := range rawValues {
				s = strings.TrimSpace(s)
				switch elemKind {
				case reflect.String:
					slice = reflect.Append(slice, reflect.ValueOf(s))

				case reflect.Uint64:
					num, err := strconv.ParseUint(s, 10, 64)
					if err != nil {
						return fmt.Errorf("invalid uint64 in %s: %v", queryKey, err)
					}
					slice = reflect.Append(slice, reflect.ValueOf(num))

				case reflect.Uint:
					num, err := strconv.ParseUint(s, 10, 64)
					if err != nil {
						return fmt.Errorf("invalid uint in %s: %v", queryKey, err)
					}
					slice = reflect.Append(slice, reflect.ValueOf(uint(num)))

				default:
					return fmt.Errorf("unsupported slice type in %s", queryKey)
				}
			}
			valueField.Set(slice)

		default:
			// Có thể thêm xử lý cho các loại khác nếu cần
		}
	}

	return nil
}

func ParseQuery[T any](c *gin.Context) (*T, error) {
	var entity T

	// Đọc toàn bộ body từ request
	if err := c.ShouldBindQuery(&entity); err != nil {
		return nil, &_routes.Except{
			Code:    400,
			Message: err.Error(),
		}
	}

	return &entity, nil
}

func CompareEqual[T1, T2 any](t1 *T1, t2 *T2) bool {
	// Kiểm tra nếu cả hai đều là nil
	if t1 == nil && t2 == nil {
		return true
	}

	// Nếu cả hai đều nil, trả về false
	if reflect.TypeOf(t1) == reflect.TypeOf(t2) {
		return reflect.DeepEqual(t1, t2) // So sánh khác nhau
	}

	return false
}

func ParseQueryToStruct[T any](query url.Values) (*T, error) {
	var t T
	v := reflect.ValueOf(&t).Elem()
	tType := v.Type()

	for i := 0; i < tType.NumField(); i++ {
		field := tType.Field(i)
		valueField := v.Field(i)

		if !valueField.CanSet() {
			continue
		}

		// Lấy tên query từ tag
		queryKey := field.Tag.Get("form")
		if queryKey == "" {
			queryKey = field.Tag.Get("parser")
		}
		if queryKey == "" {
			continue
		}

		raw := query.Get(queryKey)
		if raw == "" {
			continue
		}
		// Handle *time.Time
		if valueField.Kind() == reflect.Ptr && valueField.Type().Elem() == reflect.TypeOf(time.Time{}) {
			parsedTime, err := time.Parse(time.RFC3339, raw)
			if err != nil {
				return nil, fmt.Errorf("invalid time for %s: %v", queryKey, err)
			}
			valueField.Set(reflect.ValueOf(&parsedTime))
			continue
		}

		rawValues := strings.Split(raw, ",")

		switch valueField.Kind() {
		case reflect.String:
			valueField.SetString(raw)

		case reflect.Float64:
			parsed, err := strconv.ParseFloat(raw, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid float for %s: %v", queryKey, err)
			}
			valueField.SetFloat(parsed)

		case reflect.Uint, reflect.Uint64:
			parsed, err := strconv.ParseUint(raw, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid int for %s: %v", queryKey, err)
			}
			valueField.SetUint(parsed)

		case reflect.Slice:
			elemKind := valueField.Type().Elem().Kind()
			slice := reflect.MakeSlice(valueField.Type(), 0, len(rawValues))
			for _, s := range rawValues {
				s = strings.TrimSpace(s)
				switch elemKind {
				case reflect.String:
					slice = reflect.Append(slice, reflect.ValueOf(s))

				case reflect.Uint64:
					num, err := strconv.ParseUint(s, 10, 64)
					if err != nil {
						return nil, fmt.Errorf("invalid int64 in %s: %v", queryKey, err)
					}
					slice = reflect.Append(slice, reflect.ValueOf(num))

				case reflect.Uint:
					num, err := strconv.ParseUint(s, 10, 64)
					if err != nil {
						return nil, fmt.Errorf("invalid uint in %s: %v", queryKey, err)
					}
					slice = reflect.Append(slice, reflect.ValueOf(uint(num)))

				default:
					return nil, fmt.Errorf("unsupported slice type in %s", queryKey)
				}
			}
			valueField.Set(slice)

			// default:
			// 	return nil, fmt.Errorf("unsupported field kind %s in %s", valueField.Kind(), field.Name)
		}
	}

	return &t, nil
}

var vietnameseNumbers = map[string]int{
	"một":  1,
	"hai":  2,
	"ba":   3,
	"bốn":  4,
	"năm":  5,
	"sáu":  6,
	"bảy":  7,
	"tám":  8,
	"chín": 9,
	"mười": 10,
}

func ParseNumber(str string) *int {
	str = strings.TrimSpace(strings.ToLower(str))
	if num, err := strconv.Atoi(str); err == nil {
		return &num
	}
	if val, exists := vietnameseNumbers[str]; exists {
		return &val
	}
	return nil
}

func normalizeNumberWords(input string) string {
	input = strings.ToLower(input)
	for word, digit := range vietnameseNumbers {
		input = strings.ReplaceAll(input, word, strconv.FormatInt(int64(digit), 10))
	}
	input = strings.ReplaceAll(input, "phẩy", ".")
	input = strings.ReplaceAll(input, "chấm", ".")
	input = strings.ReplaceAll(input, ",", ".")
	input = strings.ReplaceAll(input, "tỷ", "")
	return strings.TrimSpace(input)
}

func ExtractPriceInSentence(sentence string) (int64, error) {
	// Tìm cụm có dạng "... tỷ ..."
	re := regexp.MustCompile(`([\p{L}\d\s\.,]+?)\s*tỷ`)
	matches := re.FindAllStringSubmatch(sentence, -1)

	if len(matches) == 0 {
		return 0, fmt.Errorf("Không tìm thấy giá trị tỷ")
	}

	// Xử lý cụm đầu tiên (nếu có nhiều thì có thể duyệt thêm)
	for _, match := range matches {
		raw := normalizeNumberWords(match[1])
		raw = strings.ReplaceAll(raw, " ", "")
		if raw == "" {
			continue
		}
		if val, err := strconv.ParseFloat(raw, 64); err == nil {
			return int64(val * 1_000_000_000), nil
		}
	}

	return 0, fmt.Errorf("Không parse được số")
}

type Rule struct {
	Field    string
	Pattern  *regexp.Regexp
	GetValue func([]string) interface{}
}

func ParseWithRules(sentence string, rules []Rule) map[string]interface{} {
	result := make(map[string]interface{})
	for _, rule := range rules {
		match := rule.Pattern.FindStringSubmatch(sentence)
		if len(match) > 0 {
			value := rule.GetValue(match)
			if value != nil {
				result[rule.Field] = value
			}
		}
	}
	return result
}

func isZeroValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Int, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Bool:
		return !v.Bool()
	default:
		return v.IsZero()
	}
}

func MapWithValidation[S any, D any](src *S, dest *D) *_routes.Except {
	srcVal := reflect.ValueOf(*src)
	srcType := reflect.TypeOf(*src)

	// srcVal := reflect.ValueOf(*src)
	// srcType := reflect.TypeOf(*src)

	destVal := reflect.ValueOf(dest).Elem()
	destType := destVal.Type()
	var timestampType = reflect.TypeOf(&timestamppb.Timestamp{})

	for i := 0; i < srcType.NumField(); i++ {
		srcField := srcType.Field(i)
		srcValue := srcVal.Field(i)

		// Map theo tên trường nếu tồn tại trong dest
		if destField, ok := destType.FieldByName(srcField.Name); ok {
			destFieldValue := destVal.FieldByName(destField.Name)

			// Kiểm tra khả năng set
			if destFieldValue.CanSet() {
				srcKind := srcValue.Kind()
				destKind := destFieldValue.Kind()

				// Cho phép auto-convert int64 → uint64 nếu có thể
				switch {
				case srcValue.Type() == timestampType && destFieldValue.Type() == reflect.TypeOf((*time.Time)(nil)):
					// *timestamppb.Timestamp → *time.Time
					if !srcValue.IsNil() {
						t := srcValue.Interface().(*timestamppb.Timestamp).AsTime()
						destFieldValue.Set(reflect.ValueOf(&t))
					}
				case srcKind == reflect.Int64 && destKind == reflect.Uint64:
					destFieldValue.SetUint(uint64(srcValue.Int()))
				case srcKind == destKind:
					destFieldValue.Set(srcValue)
					// Bạn có thể thêm các case khác nếu muốn hỗ trợ thêm kiểu
				}
			}
		}
	}

	errorsMap := map[string]string{}

	for i := 0; i < destType.NumField(); i++ {
		field := destType.Field(i)
		value := destVal.Field(i)

		// Lấy tên json để báo lỗi
		jsonName := field.Name
		if jsonTag := field.Tag.Get("json"); jsonTag != "" {
			jsonName = strings.Split(jsonTag, ",")[0]
		}

		// Kiểm tra binding:"required"
		if bindingTag := field.Tag.Get("binding"); strings.Contains(bindingTag, "required") {
			if isZeroValue(value) {
				errorsMap[jsonName] = "Trường " + jsonName + " là bắt buộc"
			}
		}
	}

	if len(errorsMap) > 0 {
		return &_routes.Except{
			Code:    400,
			Message: "Validation failed",
			Errors:  errorsMap,
		}
	}

	return nil
}

func ErrorField(err *_routes.Except, detail protoadapt.MessageV1) error {
	st := status.New(codes.Code(err.Code), err.Message)
	stWithDetails, _ := st.WithDetails(detail)
	return stWithDetails.Err()
}

func StructToJSONString(p interface{}) *string {
	jsonBytes, err := json.Marshal(p)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	jsonString := string(jsonBytes)
	return &jsonString
}

func ShortenContent(s string, limit int) string {
	if len(s) <= limit {
		return s
	}
	return s[:limit] + "..."
}
