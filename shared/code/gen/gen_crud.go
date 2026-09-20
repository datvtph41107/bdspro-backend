package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

var serviceName = "bdspro"
var path = fmt.Sprintf("../%s-service/v2", serviceName)
var domains = []string{
	"AssetSplitHistory",
	"AssetLegal",
	"AssetCost",
	"AssetCostType",
	"AssetExploitation",
	"AssetIncomeType",
	"IncomeDocument",
	"CostDocument",
}

// var domainName = "CostDocument"

// AssetSplitHistory,AssetLegal, AssetCost, AssetCostType, AssetExploitation,AssetIncomeType, IncomeDocument, CostDocument

// var api = "asset-cost"

func ToKebabCase(s string) string {
	// Chèn dấu "-" trước mỗi chữ in hoa (trừ chữ cái đầu tiên)
	re := regexp.MustCompile(`([a-z0-9])([A-Z])`)
	kebab := re.ReplaceAllString(s, `${1}-${2}`)
	return strings.ToLower(kebab)
}

func camelToSnake(s string) string {
	var result strings.Builder

	for i, r := range s {
		// Nếu là chữ in hoa
		if unicode.IsUpper(r) {
			// Nếu không phải ký tự đầu tiên thì thêm dấu "_"
			if i > 0 {
				result.WriteByte('_')
			}
			// Thêm chữ thường vào
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}

	return result.String()
}

// Viết hoa ký tự đầu tiên
// func upperFirst(s string) string {
// 	if len(s) == 0 {
// 		return s
// 	}
// 	return strings.ToUpper(s[:1]) + s[1:]
// }

func main() {
	// domainName :=
	fmt.Printf("GEN: %s", path)
	for i := 0; i < len(domains); i++ {
		GenerateDomain(domains[i])
	}
}

func GenerateDomain(domainName string) {
	// MakeDomain(domainName)
	// MakeRepo(domainName)
	// MakeUsecase(domainName)
	// MakeGormRepo(domainName)
	// MakeProtobuf(domainName)
	MakeService(domainName)
	// AppendYarmService(domainName)
	// Swagger(domainName)
	// AddServer(domainName)
	// AddConfig(domainName)
}

func AddServer(domainName string) {
	bindTag := "// @bind: register_server"
	dir := fmt.Sprintf("%s/../cmd/grpc", path)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		fmt.Printf("Lỗi tạo thư mục %s: %v\n", dir, err)
		return
	}

	// pkgName := camelToSnake(domainName)
	filename := fmt.Sprintf("%s/main.go", dir)

	// structName := domainName // fmt.Sprintf("%s", pkgName)

	content := `		pb_$1$.Register$2$ServiceServer(s, app.$2$Server)`

	content = strings.ReplaceAll(content, "$1$", serviceName)
	content = strings.ReplaceAll(content, "$2$", domainName)

	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Lỗi đọc file:", err)
		return
	}

	if strings.Contains(string(data), content) {
		fmt.Println("File chứa chuỗi 'abcd'")
	} else {
		input, err := os.ReadFile(filename)
		if err != nil {
			fmt.Printf("S:=%s", err)
			return
		}

		var output []string
		scanner := bufio.NewScanner(strings.NewReader(string(input)))
		fmt.Printf("S:=%s", output)

		for scanner.Scan() {
			line := scanner.Text()
			output = append(output, line)
			if strings.Contains(line, bindTag) {
				output = append(output, content) // Chèn sau dòng chứa bindTag
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Printf("S:=%s", err)
			return
		}

		// Ghi lại file
		os.WriteFile(filename, []byte(strings.Join(output, "\n")), 0644)
	}
}

func AddConfig(domainName string) {
	bindTag := "// @bind: register_server_config"
	bindTag2 := "// @bind: register_server_new"
	bindTag3 := "// @bind: register_server_param"
	dir := fmt.Sprintf("%s/../config", path)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		fmt.Printf("Lỗi tạo thư mục %s: %v\n", dir, err)
		return
	}

	// pkgName := camelToSnake(domainName)
	filename := fmt.Sprintf("%s/config.go", dir)

	// structName := domainName // fmt.Sprintf("%s", pkgName)

	content := `	$2$Server   *service.$2$Server`

	// content = strings.ReplaceAll(content, "$1$", serviceName)
	content = strings.ReplaceAll(content, "$2$", domainName)

	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Lỗi đọc file:", err)
		return
	}

	if strings.Contains(string(data), content) {
		fmt.Println("File chứa chuỗi 'abcd'")
	} else {
		input, err := os.ReadFile(filename)
		if err != nil {
			fmt.Printf("S:=%s", err)
			return
		}

		var output []string
		scanner := bufio.NewScanner(strings.NewReader(string(input)))
		fmt.Printf("S:=%s", output)

		for scanner.Scan() {
			line := scanner.Text()
			output = append(output, line)
			if strings.Contains(line, bindTag) {
				output = append(output, content) // Chèn sau dòng chứa bindTag
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Printf("S_1:=%s", err)
			return
		}

		// Ghi lại file
		os.WriteFile(filename, []byte(strings.Join(output, "\n")), 0644)
	}

	// -------
	content2 := `		$2$Server:   $2$Server,`
	content3 := `		$2$Server   *service.$2$Server,`

	// content = strings.ReplaceAll(content, "$1$", serviceName)
	content2 = strings.ReplaceAll(content2, "$2$", domainName)
	content3 = strings.ReplaceAll(content3, "$2$", domainName)

	data, err = os.ReadFile(filename)
	if err != nil {
		fmt.Println("Lỗi đọc file:", err)
		return
	}

	if strings.Contains(string(data), content2) {
		fmt.Println("File chứa chuỗi 'abcd'")
	} else {
		input, err := os.ReadFile(filename)
		if err != nil {
			fmt.Printf("S1:=%s", err)
			return
		}

		var output []string
		scanner := bufio.NewScanner(strings.NewReader(string(input)))
		fmt.Printf("S2:=%s", content2)

		for scanner.Scan() {
			line := scanner.Text()
			output = append(output, line)
			if strings.Contains(line, bindTag2) {
				output = append(output, content2) // Chèn sau dòng chứa bindTag
			}
			if strings.Contains(line, bindTag3) {
				output = append(output, content3) // Chèn sau dòng chứa bindTag
			}
		}
		// fmt.Printf("S3:=%s", output)

		if err := scanner.Err(); err != nil {
			fmt.Printf("S4:=%s", err)
			return
		}

		// Ghi lại file
		os.WriteFile(filename, []byte(strings.Join(output, "\n")), 0644)
	}
}

func MakeDomain(domainName string) {
	dir := fmt.Sprintf("%s/internal/domain", path)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		fmt.Printf("Lỗi tạo thư mục %s: %v\n", dir, err)
		return
	}

	pkgName := camelToSnake(domainName)
	fileName := fmt.Sprintf("%s.go", pkgName)

	// structName := upperFirst(pkgName) // fmt.Sprintf("%s", pkgName)

	content := fmt.Sprintf(`package domain
	
type %s struct {}`, domainName)
	// Tạo đường dẫn file
	fullPath := filepath.Join(dir, fileName)

	// Ghi nội dung vào file
	err = os.WriteFile(fullPath, []byte(content), 0644)
	if err != nil {
		fmt.Printf("Lỗi ghi file %s: %v\n", fullPath, err)
		return
	}

	fmt.Printf("Đã tạo file %s\n", fullPath)
}

func MakeRepo(domainName string) {
	dir := fmt.Sprintf("%s/repo", path)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		fmt.Printf("Lỗi tạo thư mục %s: %v\n", dir, err)
		return
	}

	pkgName := camelToSnake(domainName)
	fileName := fmt.Sprintf("%s.go", pkgName)

	// structName := domainName // fmt.Sprintf("%s", pkgName)

	content := `package repo

import (
	"$1$/v2/internal/domain"
	"context"
)

type $2$Repo interface {
	Create(c context.Context, entity *domain.$2$) error
	Update(c context.Context, id uint64, entity *domain.$2$) error
	Delete(c context.Context, id uint64) error
	GetData(c context.Context) ([]domain.$2$, error)
	GetByID(c context.Context,id uint64) (*domain.$2$, error)
}
`
	// Tạo đường dẫn file
	fullPath := filepath.Join(dir, fileName)

	content = strings.ReplaceAll(content, "$1$", serviceName)
	content = strings.ReplaceAll(content, "$2$", domainName)

	// Ghi nội dung vào file
	err = os.WriteFile(fullPath, []byte(content), 0644)
	if err != nil {
		fmt.Printf("Lỗi ghi file %s: %v\n", fullPath, err)
		return
	}

	fmt.Printf("Đã tạo file %s\n", fullPath)
}

func MakeUsecase(domainName string) {
	dir := fmt.Sprintf("%s/usecases", path)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		fmt.Printf("Lỗi tạo thư mục %s: %v\n", dir, err)
		return
	}

	pkgName := camelToSnake(domainName)
	fileName := fmt.Sprintf("%s.go", pkgName)

	// structName := domainName // fmt.Sprintf("%s", pkgName)

	content := `package usecases

import (
	"$1$/v2/internal/domain"
	"$1$/v2/repo"
	"context"
)

type $2$Usecase struct {
	Repo repo.$2$Repo
}

func New$2$Uc(repo repo.$2$Repo) $2$Usecase {
	return $2$Usecase{
		Repo: repo,
	}
}

func (uc $2$Usecase) Create(c context.Context, data *domain.$2$) (*domain.$2$, error) {
	err := uc.Repo.Create(c, data)
	return data, err
}

func (uc $2$Usecase) Update(c context.Context, id uint64, data *domain.$2$) (*domain.$2$, error) {
	err := uc.Repo.Update(c, id, data)
	return data, err
}

func (uc $2$Usecase) Delete(c context.Context, id uint64) error {
	err := uc.Repo.Delete(c, id)
	return err
}

func (uc $2$Usecase) GetData(c context.Context, id uint64) ([]domain.$2$, error) {
	result, err := uc.Repo.GetData(c)
	return result, err
}

func (uc $2$Usecase) GetByID(c context.Context, id uint64) (*domain.$2$, error) {
	result, err := uc.Repo.GetByID(c, id)
	return result, err
}

`

	content = strings.ReplaceAll(content, "$1$", serviceName)
	content = strings.ReplaceAll(content, "$2$", domainName)
	// Tạo đường dẫn file
	fullPath := filepath.Join(dir, fileName)

	// Ghi nội dung vào file
	err = os.WriteFile(fullPath, []byte(content), 0644)
	if err != nil {
		fmt.Printf("Lỗi ghi file %s: %v\n", fullPath, err)
		return
	}

	fmt.Printf("Đã tạo file %s\n", fullPath)
}

func MakeGormRepo(domainName string) {
	dir := fmt.Sprintf("%s/external/repo", path)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		fmt.Printf("Lỗi tạo thư mục %s: %v\n", dir, err)
		return
	}

	pkgName := camelToSnake(domainName)
	fileName := fmt.Sprintf("%s.go", pkgName)

	structName := domainName // fmt.Sprintf("%s", pkgName)
	fmt.Print("packageName: ", structName)

	content := `package postgre

import (
	"$1$/v2/internal/domain"
	"context"
	"time"

	"gorm.io/gorm"
)

// @bind: $1$/v2/repo.$2$Repo
type Gorm$2$Repo struct {
	DB *gorm.DB
}

func NewGorm$2$Repo(DB *gorm.DB) *Gorm$2$Repo {
	return &Gorm$2$Repo{
		DB: DB,
	}
}

func (r *Gorm$2$Repo) Create(c context.Context, entity *domain.$2$) error {
	return r.DB.WithContext(c).Create(entity).Error
}

func (r *Gorm$2$Repo) Update(c context.Context, id uint64, entity *domain.$2$) error {
	return r.DB.WithContext(c).
		Model(entity).
		Where("id = ? and deleted_at is null", id).
		Updates(entity).Error
}

func (r *Gorm$2$Repo) Delete(c context.Context, id uint64) error {
	return r.DB.WithContext(c).
		Model(&domain.$2${}).
		Where("id = ?", id).
		Update("deleted_at", time.Now()).Error
}

func (r *Gorm$2$Repo) GetData(c context.Context) ([]domain.$2$, error) {
	var entities []domain.$2$
	err := r.DB.Where("deleted_at IS NULL").Find(&entities).Error
	return entities, err
}

func (r *Gorm$2$Repo) GetByID(c context.Context, id uint64) (*domain.$2$, error) {
	var entity domain.$2$
	err := r.DB.Where("id = ? AND deleted_at IS NULL", id).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}
`

	content = strings.ReplaceAll(content, "$1$", serviceName)
	content = strings.ReplaceAll(content, "$2$", domainName)
	// Tạo đường dẫn file
	fullPath := filepath.Join(dir, fileName)

	// Ghi nội dung vào file
	err = os.WriteFile(fullPath, []byte(content), 0644)
	if err != nil {
		fmt.Printf("Lỗi ghi file %s: %v\n", fullPath, err)
		return
	}

	fmt.Printf("Đã tạo file %s\n", fullPath)
}

func MakeProtobuf(domainName string) {
	dir := fmt.Sprintf("%s/../../protobuf/schema/%s", path, serviceName)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		fmt.Printf("Lỗi tạo thư mục %s: %v\n", dir, err)
		return
	}

	pkgName := camelToSnake(domainName)
	fileName := fmt.Sprintf("%s.proto", pkgName)

	structName := domainName // fmt.Sprintf("%s", pkgName)
	fmt.Print("packageName: ", structName)

	content := `syntax = "proto3";

package pb_$1$;
option go_package = "types/$1$";
import "google/protobuf/timestamp.proto";

message $2$DTO {
  uint64 id = 1;
  bool archived = 2;
}

message $2$ListDTO {
  repeated $2$DTO data = 1;
  uint32 totalElements = 2;
}

service $2$Service {
    rpc Get ($2$DTO) returns ($2$ListDTO);
    rpc Detail ($2$DTO) returns ($2$DTO);
    rpc Create ($2$DTO) returns ($2$DTO);
    rpc Update ($2$DTO) returns ($2$DTO);
    rpc Delete ($2$DTO) returns ($2$DTO);
}
`
	// Tạo đường dẫn file
	fullPath := filepath.Join(dir, fileName)
	content = strings.ReplaceAll(content, "$1$", serviceName)
	content = strings.ReplaceAll(content, "$2$", domainName)

	// Ghi nội dung vào file
	err = os.WriteFile(fullPath, []byte(content), 0644)
	if err != nil {
		fmt.Printf("Lỗi ghi file %s: %v\n", fullPath, err)
		return
	}

	fmt.Printf("Đã tạo file %s\n", fullPath)
}

func MakeService(domainName string) {
	dir := fmt.Sprintf("%s/external/service", path)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		fmt.Printf("Lỗi tạo thư mục %s: %v\n", dir, err)
		return
	}

	pkgName := camelToSnake(domainName)
	fileName := fmt.Sprintf("%s.go", pkgName)

	structName := domainName // fmt.Sprintf("%s", pkgName)
	fmt.Print("packageName: ", structName)

	content := `package service

import (
	"$1$/v2/internal/domain"
	"$1$/v2/usecases"
	"context"
	pb_$1$ "pb/types/$1$"

	"github.com/jinzhu/copier"
)

type $2$Server struct {
	pb_$1$.Unimplemented$2$ServiceServer
	UC usecases.$2$Usecase
}

func New$2$Server($2$Uc usecases.$2$Usecase) *$2$Server {
	return &$2$Server{
		UC: $2$Uc,
	}
}

func (s $2$Server) Get(c context.Context, dto *pb_$1$.$2$DTO) (*pb_$1$.$2$ListDTO, error) {
	result, err := s.UC.GetData(c, dto.Id)

	pbdatas := []*pb_$1$.$2$DTO{}
	for _, r := range result {
		pbdata := &pb_$1$.$2$DTO{}
		copier.Copy(&pbdata, &r)
		pbdatas = append(pbdatas, pbdata)
	}

	return &pb_$1$.$2$ListDTO{
		Data: pbdatas,
	}, err
}
func (s $2$Server) Detail(c context.Context, dto *pb_$1$.$2$DTO) (*pb_$1$.$2$DTO, error) {
	r, err := s.UC.GetByID(c, dto.Id)

	pbdata := &pb_$1$.$2$DTO{}
	copier.Copy(&pbdata, &r)

	return pbdata, err
}
func (s $2$Server) Create(c context.Context, dto *pb_$1$.$2$DTO) (*pb_$1$.$2$DTO, error) {
	d := &domain.$2${}
	copier.Copy(&d, &dto)

	r, err := s.UC.Create(c, d)

	copier.Copy(&dto, &r)
	return dto, err
}
func (s $2$Server) Update(c context.Context, dto *pb_$1$.$2$DTO) (*pb_$1$.$2$DTO, error) {
	d := &domain.$2${}
	copier.Copy(&d, &dto)

	r, err := s.UC.Update(c, dto.Id, d)

	copier.Copy(&dto, &r)
	return dto, err
}
func (s $2$Server) Delete(c context.Context, dto *pb_$1$.$2$DTO) (*pb_$1$.$2$DTO, error) {
	err := s.UC.Delete(c, dto.Id)
	return dto, err
}
`

	content = strings.ReplaceAll(content, "$1$", serviceName)
	content = strings.ReplaceAll(content, "$2$", domainName)
	// Tạo đường dẫn file
	fullPath := filepath.Join(dir, fileName)

	// Ghi nội dung vào file
	err = os.WriteFile(fullPath, []byte(content), 0644)
	if err != nil {
		fmt.Printf("Lỗi ghi file %s: %v\n", fullPath, err)
		return
	}

	fmt.Printf("Đã tạo file %s\n", fullPath)
}

func AppendYarmService(domainName string) {
	dir := fmt.Sprintf("%s/../../protobuf/schema/%s", path, serviceName)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		fmt.Printf("Lỗi tạo thư mục %s: %v\n", dir, err)
		return
	}

	// pkgName := camelToSnake(domainName)
	fileName := "service.yaml"
	api := ToKebabCase(domainName)

	structName := domainName // fmt.Sprintf("%s", pkgName)
	fmt.Print("packageName: ", structName)

	content := `
    - selector: pb_$1$.$2$Service.Get
      get: /v2/$1$/v2/$3$

    - selector: pb_$1$.$2$Service.Create
      post: /v2/$1$/v2/$3$
      body: "*"

    - selector: pb_$1$.$2$Service.Update
      put: /v2/$1$/v2/$3$/{id}
      body: "*"

    - selector: pb_$1$.$2$Service.Delete
      delete: /v2/$1$/v2/$3$/{id}

    - selector: pb_$1$.$2$Service.Detail
      get: /v2/$1$/v2/$3$/{id}

`
	content = strings.ReplaceAll(content, "$1$", serviceName)
	content = strings.ReplaceAll(content, "$2$", domainName)
	content = strings.ReplaceAll(content, "$3$", api)
	// Tạo đường dẫn file
	fullPath := filepath.Join(dir, fileName)

	data, err := os.ReadFile(fullPath)
	if err != nil {
		fmt.Println("Lỗi đọc file:", err)
		return
	}

	if strings.Contains(string(data), api) {
		fmt.Println("File chứa chuỗi 'abcd'")
	} else {
		f, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("Lỗi mở file:", err)
			return
		}
		defer f.Close()
		if _, err := f.WriteString(content); err != nil {
			fmt.Println("Lỗi ghi file:", err)
		}
	}
}

func Swagger(domainName string) {
	dir := fmt.Sprintf("%s/swagger", path)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		fmt.Printf("Lỗi tạo thư mục %s: %v\n", dir, err)
		return
	}

	// pkgName := camelToSnake(domainName)
	pkgName := camelToSnake(domainName)
	fileName := fmt.Sprintf("%s.go", pkgName)

	api := ToKebabCase(domainName)

	structName := domainName // fmt.Sprintf("%s", pkgName)
	fmt.Print("packageName: ", structName)

	content := `
package swagger_gen

import (
	pb_$1$ "pb/types/$1$"
)

// @Summary Lấy danh sách
// @Tags $2$
// @Produce json
// @Param page query int false "Page"
// @Param size query int false "Size"
// @Param query query pb_$1$.$2$DTO true "Size"
// @Security BearerAuth
// @Router /v2/$1$/v2/$3$ [get]
func Get$2$() {}

// @Summary Xem chi tiết
// @Tags $2$
// @Produce json
// @Param id path int true "ID"
// @Security BearerAuth
// @Router /v2/$1$/v2/$3$/{id} [get]
func Detail$2$() {}

// @Summary Tạo mới
// @Tags $2$
// @Produce json
// @Param body body pb_$1$.$2$DTO true "Body"
// @Security BearerAuth
// @Router /v2/$1$/v2/$3$ [post]
func Create$2$(pb_$1$.$2$DTO) {}

// @Summary Lấy theo ID
// @Tags $2$
// @Produce json
// @Param id path int true "ID"
// @Param body body pb_$1$.$2$DTO true "Body"
// @Security BearerAuth
// @Router /v2/$1$/v2/$3$/{id} [put]
func Update$2$() {}

// @Summary Xóa
// @Tags $2$
// @Produce json
// @Param id path int true "ID"
// @Security BearerAuth
// @Router /v2/$1$/v2/$3$/{id} [delete]
func Delete$2$() {}

`
	// Tạo đường dẫn file
	fullPath := filepath.Join(dir, fileName)
	content = strings.ReplaceAll(content, "$1$", serviceName)
	content = strings.ReplaceAll(content, "$2$", domainName)
	content = strings.ReplaceAll(content, "$3$", api)

	// data, err := os.ReadFile(fullPath)
	// if err != nil {
	// 	fmt.Println("Lỗi đọc file:", err)
	// 	return
	// }

	// Ghi nội dung vào file
	err = os.WriteFile(fullPath, []byte(content), 0644)
	if err != nil {
		fmt.Printf("Lỗi ghi file %s: %v\n", fullPath, err)
		return
	}
}
