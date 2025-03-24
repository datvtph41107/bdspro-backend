package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
)

const wireFile = "../wire.go"
const rootPackage = "bdspro"

type NewFunc struct {
	FullPackage string
	Alias       string
	Func        string
}

func main() {
	newFuncs := findNewFunctions("../")
	existingFuncs := readExistingWire()

	updatedFuncs := mergeFunctions(existingFuncs, newFuncs)

	writeWireFile(updatedFuncs)
	fmt.Println("✅ wire.go updated successfully!")

	// cmd := exec.Command("cd .. && wire")
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", "cd .. && wire")
	} else {
		cmd = exec.Command("/bin/sh", "-c", "cd .. && wire")
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("wire generate failed: " + err.Error())
	} else {
		fmt.Println("✅ wire generate successfully!")
	}
}

// Tìm tất cả các hàm `NewXXX` trong các file Go và lấy package + đường dẫn đầy đủ
func findNewFunctions(root string) []NewFunc {
	var newFuncs []NewFunc
	re := regexp.MustCompile(`func\s+(New\w+)\(`)

	_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		fullPackage := getFullPackagePath(path)
		alias := generateAlias(fullPackage)
		if fullPackage == "" {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			fmt.Println("⚠️ Error opening file:", path, err)
			return nil
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())

			// Bỏ qua comment `// func NewXXX`
			if strings.HasPrefix(line, "// func New") {
				continue
			}

			if matches := re.FindStringSubmatch(line); len(matches) > 1 {
				newFuncs = append(newFuncs, NewFunc{FullPackage: fullPackage, Alias: alias, Func: matches[1]})
			}
		}

		return nil
	})

	sort.Slice(newFuncs, func(i, j int) bool {
		return newFuncs[i].Func < newFuncs[j].Func
	})
	return newFuncs
}

// Lấy đường dẫn import đầy đủ từ đường dẫn file, bỏ `../`
func getFullPackagePath(path string) string {
	relPath, err := filepath.Rel(".", filepath.Dir(path))
	if err != nil {
		return ""
	}
	if path == "../main.go" {
		return ""
	}
	return rootPackage + "/" + strings.TrimPrefix(strings.ReplaceAll(relPath, "\\", "/"), "../") // Bỏ `../`
}

// Tạo alias từ đường dẫn package (ví dụ: "mod/user/product" -> "mod_user_product")
func generateAlias(pkgPath string) string {
	return strings.ReplaceAll(pkgPath, "/", "_")
}

// Đọc các hàm đã có trong wire.go
func readExistingWire() []NewFunc {
	data, err := os.ReadFile(wireFile)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("📄 wire.go not found, creating a new one.")
			return nil
		}
		fmt.Println("⚠️ Error reading wire.go:", err)
		return nil
	}

	re := regexp.MustCompile(`(\S+)\.(New\w+)\(\)`)
	matches := re.FindAllStringSubmatch(string(data), -1)

	var funcs []NewFunc
	for _, match := range matches {
		funcs = append(funcs, NewFunc{FullPackage: match[1], Alias: generateAlias(match[1]), Func: match[2]})
	}

	sort.Slice(funcs, func(i, j int) bool {
		return funcs[i].Func < funcs[j].Func
	})
	return funcs
}

// Hợp nhất danh sách hàm cũ và mới, tránh trùng lặp
func mergeFunctions(existing, newFuncs []NewFunc) []NewFunc {
	funcSet := make(map[string]NewFunc)

	for _, f := range existing {
		funcSet[f.Func] = f
	}
	for _, f := range newFuncs {
		funcSet[f.Func] = f
	}

	var result []NewFunc
	for _, f := range funcSet {
		result = append(result, f)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Func < result[j].Func
	})
	return result
}

// Ghi danh sách các hàm vào wire.go
func writeWireFile(funcs []NewFunc) {
	file, err := os.Create(wireFile)
	if err != nil {
		fmt.Println("⚠️ Error creating wire.go:", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// Viết phần header
	writer.WriteString(`//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	_db "common/db"
`)

	// Viết danh sách import package với alias mới
	packageSet := make(map[string]string)
	for _, f := range funcs {
		packageSet[f.FullPackage] = f.Alias
	}

	var imports []string
	for pkg, alias := range packageSet {
		imports = append(imports, fmt.Sprintf("\t%s \"%s\"", alias, pkg))
	}
	sort.Strings(imports)

	for _, imp := range imports {
		writer.WriteString(imp + "\n")
	}

	writer.WriteString(")\n\n")

	// Viết danh sách WireSet
	writer.WriteString("// Inject dependencies\n")
	writer.WriteString(`var wireSet = wire.NewSet(
	_db.NewDB,
	NewConfigApp,
`)

	for _, f := range funcs {
		writer.WriteString(fmt.Sprintf("\t%s.%s,\n", f.Alias, f.Func))
	}

	// Đóng file
	writer.WriteString(`)
func InitializeApp() (*ConfigApp, error) {
	wire.Build(wireSet)
	return nil, nil
}`)
	writer.Flush()
}
