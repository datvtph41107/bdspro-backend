package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
)

type NewFunc struct {
	FullPackage string
	Alias       string
	Func        string
}

const packageName = "cmd"

func main() {
	// const rootPackage = "bdspro"
	rootPackage := os.Args[1]
	prefix := rootPackage + "/cmd"
	// path := fmt.Sprintf("../%s-service/cmd/http", rootPackage)
	wireFile := "wire/wire.go" //fmt.Sprint()
	// log.Println("wireFile: ", wireFile)
	folderFind := "./" //fmt.Sprintf("", rootPackage)
	// mod := os.Args[0]
	// log.Printf("MOD= %s; %s", mod, rootPackage)

	newFuncs := findNewFunctions(rootPackage, folderFind)
	existingFuncs := readExistingWire(wireFile)

	updatedFuncs := mergeFunctions(existingFuncs, newFuncs)
	binds := findBinds("./", rootPackage)
	// log.Println("binds: ", binds)

	moduleName := rootPackage
	// if moduleName == "notification" {
	// 	moduleName = "noti"
	// }
	if moduleName == "transaction" {
		moduleName = "tx"
	}
	writeWireFile(updatedFuncs, binds, rootPackage, prefix, wireFile, moduleName)
	fmt.Println("✅ wire.go updated successfully!")

	// cmd := exec.Command("cd .. && wire")
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", "cd .. && wire")
	} else {
		cmd = exec.Command("/bin/sh", "-c", fmt.Sprintf("cd wire && wire"))
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
func findNewFunctions(rootPackage, root string) []NewFunc {
	var newFuncs []NewFunc
	re := regexp.MustCompile(`func\s+(New\w+)\(`)

	_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		fullPackage := getFullPackagePath(rootPackage, path)
		// log.Printf("fullPackage: %s; %s", rootPackage, path)
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
func getFullPackagePath(rootPackage, path string) string {
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
func readExistingWire(wireFile string) []NewFunc {
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
func writeWireFile(inputs []NewFunc, binds []BindEntry, rootPackage, prefix string, wireFile string, moduleName string) {
	funcs := []NewFunc{}
	for _, f := range inputs {
		log.Println(f.FullPackage)
		if !strings.HasPrefix(f.FullPackage, prefix) {
			funcs = append(funcs, f)
		}
	}

	file, err := os.Create(wireFile)
	if err != nil {
		fmt.Println("⚠️ Error creating wire.go:", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// Viết phần header
	writer.WriteString(fmt.Sprintf(`//go:build wireinject
// +build wireinject
	
package wire

import (
	"github.com/google/wire"
	_db "common/db"
	"%s/initial"
`, moduleName))

	// Viết danh sách import package với alias mới
	packageSet := make(map[string]string)
	for _, f := range funcs {
		packageSet[f.FullPackage] = f.Alias
	}
	for _, b := range binds {
		packageSet[b.InterfacePath] = b.InterfaceAlias
		// writer.WriteString(fmt.Sprintf("\twire.Bind(new(%s.%s), new(*%s.%s)),\n",
		// 	b.InterfaceAlias, b.InterfaceName,
		// 	b.StructAlias, b.StructName))
	}

	var imports []string
	for pkg, alias := range packageSet {
		imports = append(imports, fmt.Sprintf("\t%s \"%s\"", alias, replacePrefix(pkg, moduleName)))

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
	_db.NewTransactionRepo,
`)

	// for _, f := range binds {
	// 	if f.IsBind {
	// 		// writer.WriteString(fmt.Sprintf("\twire.Bind(new(%s.%s), new(*%s.%s)),\n",
	// 		// 	f.Alias, f.BindTarget, f.Alias, f.Func))
	// 		writer.WriteString(fmt.Sprintf("\twire.Bind(new(%s.%s), new(*%s.%s)),\n",
	// 			f.Alias, f.Func, f.Alias, f.BindTarget))
	// 	}
	// }
	for _, b := range binds {
		writer.WriteString(fmt.Sprintf("\twire.Bind(new(%s.%s), new(*%s.%s)),\n",
			b.InterfaceAlias, b.InterfaceName,
			b.StructAlias, b.StructName))
	}

	for _, f := range funcs {
		writer.WriteString(fmt.Sprintf("\t%s.%s,\n", f.Alias, f.Func))
	}

	// Đóng file
	writer.WriteString(`)
func InitializeApp() (*initial.InitialApp, error) {
	wire.Build(wireSet)
	return nil, nil
}`)
	writer.Flush()
}

type BindEntry struct {
	InterfaceAlias string
	InterfacePath  string
	InterfaceName  string
	StructAlias    string
	StructPackage  string
	StructName     string
}

// func findBindStatements(root string) []BindStatement {
// 	var binds []BindStatement
// 	re := regexp.MustCompile(`@bind:\s+([\w./]+)\s*=>\s*([\w./]+)`)

// 	_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
// 		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".go") {
// 			return nil
// 		}

// 		file, err := os.Open(path)
// 		if err != nil {
// 			return nil
// 		}
// 		defer file.Close()

// 		scanner := bufio.NewScanner(file)
// 		for scanner.Scan() {
// 			line := strings.TrimSpace(scanner.Text())
// 			if matches := re.FindStringSubmatch(line); len(matches) == 3 {
// 				binds = append(binds, BindStatement{
// 					Interface: matches[1],
// 					Concrete:  matches[2],
// 				})
// 			}
// 		}
// 		return nil
// 	})

// 	return binds
// }

func findBinds(root string, rootPackage string) []BindEntry {
	var binds []BindEntry
	bindRe := regexp.MustCompile(`@bind:\s*([\w/]+)/(\w+)\.(\w+)`) // e.g., payment/internal/wallet/repo.AccountRepo

	_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return nil
		}
		defer file.Close()

		fullStructPackage := getFullPackagePath(rootPackage, path)
		structAlias := generateAlias(fullStructPackage)

		var pendingBind *BindEntry
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if matches := bindRe.FindStringSubmatch(line); len(matches) == 4 {
				interfacePkgPath := matches[1] + "/" + matches[2]
				interfaceAlias := generateAlias(interfacePkgPath)
				interfaceName := matches[3]

				pendingBind = &BindEntry{
					InterfaceAlias: interfaceAlias,
					InterfacePath:  interfacePkgPath,
					InterfaceName:  interfaceName,
					StructPackage:  fullStructPackage,
					StructAlias:    structAlias,
				}
				continue
			}

			// Expect: type Xxx struct {
			if pendingBind != nil && strings.HasPrefix(line, "type ") && strings.Contains(line, "struct") {
				words := strings.Fields(line)
				if len(words) >= 3 && strings.HasPrefix(words[2], "struct") {
					pendingBind.StructName = words[1]
					binds = append(binds, *pendingBind)
					pendingBind = nil
				}
			}
		}

		return nil
	})

	return binds
}

func replacePrefix(input, newPrefix string) string {
	idx := strings.Index(input, "/")
	if idx == -1 {
		return newPrefix // hoặc input nếu không muốn thay
	}
	return newPrefix + input[idx:]
}
