package main

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"unicode"

	"gopkg.in/yaml.v3"
)

type NewFunc struct {
	FullPackage    string
	Alias          string
	Func           string
	ReturnsCleanup bool
}

type BindEntry struct {
	InterfaceAlias string
	InterfacePath  string
	InterfaceName  string
	StructAlias    string
	StructPackage  string
	StructName     string
}

type InjectionConfig struct {
	Packages      []InjectionPackage            `yaml:"packages"`
	Exclude       map[string][]string           `yaml:"exclude"`
	ScanExclude   map[string][]string           `yaml:"scan_exclude"`
	ProcessInputs map[string]ProcessInputConfig `yaml:"process_inputs"`
}

type ProcessInputConfig struct {
	Database bool `yaml:"database"`
	Redis    bool `yaml:"redis"`
}

type InjectionPackage struct {
	Import       string   `yaml:"import"`
	Constructors []string `yaml:"constructors"`
}

func main() {
	if len(os.Args) < 2 || strings.TrimSpace(os.Args[1]) == "" {
		fmt.Fprintln(os.Stderr, "usage: go run gen_wire.go <service-module>")
		os.Exit(2)
	}

	rootPackage := strings.TrimSpace(os.Args[1])
	moduleName := rootPackage
	if moduleName == "transaction" {
		moduleName = "tx"
	}

	wireFile := filepath.Join("wire", "wire.go")
	injections, scanExclude, processInputs, err := loadGenerationConfig(rootPackage)
	if err != nil {
		fmt.Fprintln(os.Stderr, "load wire generation config:", err)
		os.Exit(1)
	}
	constructors, err := findNewFunctions(rootPackage, ".", scanExclude)
	if err != nil {
		fmt.Fprintln(os.Stderr, "scan constructors:", err)
		os.Exit(1)
	}

	cleanupEnabled := false
	for _, constructor := range constructors {
		if constructor.ReturnsCleanup {
			cleanupEnabled = true
			break
		}
	}

	binds, err := findBinds(".", rootPackage, scanExclude)
	if err != nil {
		fmt.Fprintln(os.Stderr, "scan binds:", err)
		os.Exit(1)
	}
	if err := writeWireFile(constructors, binds, injections, wireFile, moduleName, cleanupEnabled, processInputs); err != nil {
		fmt.Fprintln(os.Stderr, "write wire.go:", err)
		os.Exit(1)
	}

	fmt.Println("✅ wire.go updated successfully")

	cmd := exec.Command("wire", "gen")
	cmd.Dir = filepath.Dir(wireFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "wire generate failed:", err)
		os.Exit(1)
	}

	fmt.Println("✅ wire generate successfully")
}

func findNewFunctions(rootPackage, root string, excluded []string) ([]NewFunc, error) {
	functions := make(map[string]NewFunc)
	fset := token.NewFileSet()

	err := filepath.Walk(root, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if info.IsDir() {
			name := info.Name()
			if name == ".git" || name == "vendor" || name == "node_modules" {
				return filepath.SkipDir
			}

			relative, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			relativeSlash := filepath.ToSlash(relative)
			if relativeSlash == "wire" || isScanExcluded(relativeSlash, excluded) {
				return filepath.SkipDir
			}

			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		file, err := parser.ParseFile(fset, path, nil, parser.SkipObjectResolution)
		if err != nil {
			return fmt.Errorf("parse %s: %w", path, err)
		}
		fullPackage := getFullPackagePath(rootPackage, path)
		if fullPackage == "" || strings.HasPrefix(fullPackage, rootPackage+"/cmd") {
			return nil
		}

		for _, declaration := range file.Decls {
			function, ok := declaration.(*ast.FuncDecl)
			if !ok || function.Recv != nil || !strings.HasPrefix(function.Name.Name, "New") {
				continue
			}
			item := NewFunc{
				FullPackage:    fullPackage,
				Alias:          generateAlias(fullPackage),
				Func:           function.Name.Name,
				ReturnsCleanup: hasWireCleanupResult(function),
			}
			functions[item.FullPackage+"\x00"+item.Func] = item
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	result := make([]NewFunc, 0, len(functions))
	for _, item := range functions {
		result = append(result, item)
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Func == result[j].Func {
			return result[i].FullPackage < result[j].FullPackage
		}
		return result[i].Func < result[j].Func
	})
	return result, nil
}

func hasWireCleanupResult(function *ast.FuncDecl) bool {
	if function == nil || function.Type.Results == nil {
		return false
	}

	resultTypes := make([]ast.Expr, 0, len(function.Type.Results.List))
	for _, field := range function.Type.Results.List {
		count := len(field.Names)
		if count == 0 {
			count = 1
		}
		for i := 0; i < count; i++ {
			resultTypes = append(resultTypes, field.Type)
		}
	}

	if len(resultTypes) != 3 {
		return false
	}

	cleanup, ok := resultTypes[1].(*ast.FuncType)
	if !ok {
		return false
	}
	if cleanup.Params != nil && len(cleanup.Params.List) != 0 {
		return false
	}
	if cleanup.Results != nil && len(cleanup.Results.List) != 0 {
		return false
	}

	errType, ok := resultTypes[2].(*ast.Ident)
	return ok && errType.Name == "error"
}

func findBinds(root, rootPackage string, excluded []string) ([]BindEntry, error) {
	bindPattern := regexp.MustCompile(`@bind:\s*([A-Za-z0-9_./-]+)`)
	binds := make(map[string]BindEntry)
	fset := token.NewFileSet()

	err := filepath.Walk(root, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if info.IsDir() {
			name := info.Name()
			if name == ".git" || name == "vendor" || name == "node_modules" {
				return filepath.SkipDir
			}
			relative, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			if isScanExcluded(filepath.ToSlash(relative), excluded) {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		file, err := parser.ParseFile(fset, path, nil, parser.ParseComments|parser.SkipObjectResolution)
		if err != nil {
			return fmt.Errorf("parse %s: %w", path, err)
		}
		structPackage := getFullPackagePath(rootPackage, path)
		structAlias := generateAlias(structPackage)

		for _, declaration := range file.Decls {
			gen, ok := declaration.(*ast.GenDecl)
			if !ok || gen.Tok != token.TYPE {
				continue
			}
			for _, specification := range gen.Specs {
				typeSpec, ok := specification.(*ast.TypeSpec)
				if !ok {
					continue
				}
				if _, ok := typeSpec.Type.(*ast.StructType); !ok {
					continue
				}

				comment := commentText(gen.Doc) + "\n" + commentText(typeSpec.Doc)
				matches := bindPattern.FindAllStringSubmatch(comment, -1)
				for _, match := range matches {
					if len(match) != 2 {
						continue
					}
					target := strings.TrimSuffix(strings.TrimSpace(match[1]), ".go")
					separator := strings.LastIndex(target, ".")
					if separator <= 0 || separator == len(target)-1 {
						return fmt.Errorf("%s: invalid @bind target %q", path, target)
					}
					interfacePath := target[:separator]
					interfaceName := target[separator+1:]
					if !ast.IsExported(typeSpec.Name.Name) {
						return fmt.Errorf("%s: @bind concrete struct %s must be exported", path, typeSpec.Name.Name)
					}

					item := BindEntry{
						InterfaceAlias: generateAlias(interfacePath),
						InterfacePath:  interfacePath,
						InterfaceName:  interfaceName,
						StructAlias:    structAlias,
						StructPackage:  structPackage,
						StructName:     typeSpec.Name.Name,
					}
					key := item.InterfacePath + "." + item.InterfaceName + "\x00" + item.StructPackage + "." + item.StructName
					binds[key] = item
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	result := make([]BindEntry, 0, len(binds))
	for _, item := range binds {
		result = append(result, item)
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].InterfacePath == result[j].InterfacePath {
			return result[i].InterfaceName < result[j].InterfaceName
		}
		return result[i].InterfacePath < result[j].InterfacePath
	})
	return result, nil
}

func commentText(group *ast.CommentGroup) string {
	if group == nil {
		return ""
	}
	return group.Text()
}

func getFullPackagePath(rootPackage, path string) string {
	relativeDirectory, err := filepath.Rel(".", filepath.Dir(path))
	if err != nil {
		return ""
	}
	relativeDirectory = filepath.ToSlash(relativeDirectory)
	if relativeDirectory == "." || relativeDirectory == "" {
		return rootPackage
	}
	return rootPackage + "/" + strings.TrimPrefix(relativeDirectory, "../")
}

func generateAlias(packagePath string) string {
	var output strings.Builder
	for _, character := range packagePath {
		if unicode.IsLetter(character) || unicode.IsDigit(character) || character == '_' {
			output.WriteRune(character)
			continue
		}
		output.WriteByte('_')
	}
	return strings.Trim(output.String(), "_")
}

func scriptDir() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "."
	}
	return filepath.Dir(filename)
}

func loadGenerationConfig(service string) ([]InjectionPackage, []string, ProcessInputConfig, error) {
	configPath := filepath.Join(scriptDir(), "injection_clients.yml")
	content, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil, ProcessInputConfig{}, nil
		}
		return nil, nil, ProcessInputConfig{}, fmt.Errorf("read %s: %w", configPath, err)
	}

	var config InjectionConfig
	if err := yaml.Unmarshal(content, &config); err != nil {
		return nil, nil, ProcessInputConfig{}, fmt.Errorf("parse %s: %w", configPath, err)
	}

	service = strings.TrimSpace(service)
	return filterInjectionPackages(
		config.Packages,
		config.Exclude[service],
	), normalizeScanExclude(config.ScanExclude[service]), config.ProcessInputs[service], nil
}

func normalizeScanExclude(items []string) []string {
	seen := make(map[string]struct{}, len(items))
	result := make([]string, 0, len(items))
	for _, item := range items {
		item = strings.Trim(strings.TrimSpace(filepath.ToSlash(item)), "/")
		if item == "" || item == "." {
			continue
		}
		if _, exists := seen[item]; exists {
			continue
		}
		seen[item] = struct{}{}
		result = append(result, item)
	}
	sort.Strings(result)
	return result
}

func isScanExcluded(relative string, excluded []string) bool {
	relative = strings.Trim(strings.TrimSpace(filepath.ToSlash(relative)), "/")
	if relative == "" || relative == "." {
		return false
	}
	for _, item := range excluded {
		if relative == item || strings.HasPrefix(relative, item+"/") {
			return true
		}
	}
	return false
}

func filterInjectionPackages(
	packages []InjectionPackage,
	excluded []string,
) []InjectionPackage {
	if len(excluded) == 0 {
		return packages
	}

	excludedSet := make(map[string]struct{}, len(excluded))
	for _, item := range excluded {
		item = strings.TrimSpace(item)
		if item != "" {
			excludedSet[item] = struct{}{}
		}
	}

	filtered := make([]InjectionPackage, 0, len(packages))
	for _, injection := range packages {
		packagePath := strings.TrimSpace(injection.Import)
		constructors := make([]string, 0, len(injection.Constructors))

		for _, constructor := range injection.Constructors {
			constructor = strings.TrimSpace(constructor)
			if constructor == "" {
				continue
			}

			key := packagePath + "." + constructor
			if _, skip := excludedSet[key]; skip {
				continue
			}

			constructors = append(constructors, constructor)
		}

		if len(constructors) == 0 {
			continue
		}

		filtered = append(filtered, InjectionPackage{
			Import:       packagePath,
			Constructors: constructors,
		})
	}

	return filtered
}

func writeWireFile(functions []NewFunc, binds []BindEntry, injections []InjectionPackage, wireFile, moduleName string, cleanupEnabled bool, processInputs ProcessInputConfig) error {
	if err := os.MkdirAll(filepath.Dir(wireFile), 0o755); err != nil {
		return err
	}
	file, err := os.Create(wireFile)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	initialPackage := moduleName + "/initial"
	initialAlias := generateAlias(initialPackage)

	if _, err := fmt.Fprintf(writer, `//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	_db "common/db"
`); err != nil {
		return err
	}
	if !processInputs.Redis {
		if _, err := writer.WriteString("\t_provider \"common/provider\"\n"); err != nil {
			return err
		}
	}
	if processInputs.Database {
		if _, err := writer.WriteString("\t\"gorm.io/gorm\"\n"); err != nil {
			return err
		}
	}
	if processInputs.Redis {
		if _, err := writer.WriteString("\t_redis \"common/redis\"\n"); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintf(writer, "\t%s %q\n", initialAlias, initialPackage); err != nil {
		return err
	}

	packages := make(map[string]string)
	for _, function := range functions {
		packages[function.FullPackage] = function.Alias
	}
	for _, bind := range binds {
		packages[bind.InterfacePath] = bind.InterfaceAlias
		packages[bind.StructPackage] = bind.StructAlias
	}
	for _, injection := range injections {
		packagePath := strings.TrimSpace(injection.Import)
		if packagePath == "" {
			continue
		}
		packages[packagePath] = generateAlias(packagePath)
	}

	imports := make([]string, 0, len(packages))
	for packagePath, alias := range packages {
		if packagePath == initialPackage || packagePath == "common/db" || packagePath == "common/provider" {
			continue
		}
		imports = append(imports, fmt.Sprintf("\t%s %q", alias, resolveImportPath(packagePath, moduleName)))
	}
	sort.Strings(imports)
	for _, item := range imports {
		if _, err := writer.WriteString(item + "\n"); err != nil {
			return err
		}
	}
	if _, err := writer.WriteString(")\n\n// Inject dependencies\nvar wireSet = wire.NewSet(\n"); err != nil {
		return err
	}
	if !processInputs.Database {
		if _, err := writer.WriteString("\t_db.NewDB,\n"); err != nil {
			return err
		}
	}
	if _, err := writer.WriteString("\t_db.NewTransactionRepo,\n"); err != nil {
		return err
	}
	if !processInputs.Redis {
		if _, err := writer.WriteString("\t_provider.SyncProviderSet,\n"); err != nil {
			return err
		}
	}

	for _, bind := range binds {
		if _, err := fmt.Fprintf(writer, "\twire.Bind(new(%s.%s), new(*%s.%s)),\n", bind.InterfaceAlias, bind.InterfaceName, bind.StructAlias, bind.StructName); err != nil {
			return err
		}
	}
	for _, injection := range injections {
		alias := packages[strings.TrimSpace(injection.Import)]
		for _, constructor := range injection.Constructors {
			constructor = strings.TrimSpace(constructor)
			if constructor == "" {
				continue
			}
			if _, err := fmt.Fprintf(writer, "\t%s.%s,\n", alias, constructor); err != nil {
				return err
			}
		}
	}
	for _, function := range functions {
		if _, err := fmt.Fprintf(writer, "\t%s.%s,\n", function.Alias, function.Func); err != nil {
			return err
		}
	}
	parameters := make([]string, 0, 2)
	if processInputs.Database {
		parameters = append(parameters, "database *gorm.DB")
	}
	if processInputs.Redis {
		parameters = append(parameters, "redisService *_redis.RedisService")
	}
	parameterList := strings.Join(parameters, ", ")

	injectorTemplate := `)

func InitializeApp(%s) (*%s.InitialApp, error) {
	wire.Build(wireSet)
	return nil, nil
}
`
	if cleanupEnabled {
		injectorTemplate = `)

func InitializeApp(%s) (*%s.InitialApp, func(), error) {
	wire.Build(wireSet)
	return nil, nil, nil
}
`
	}

	if _, err := fmt.Fprintf(writer, injectorTemplate, parameterList, initialAlias); err != nil {
		return err
	}
	if err := writer.Flush(); err != nil {
		return err
	}

	format := exec.Command("gofmt", "-w", wireFile)
	format.Stdout = os.Stdout
	format.Stderr = os.Stderr
	return format.Run()
}

func resolveImportPath(packagePath, moduleName string) string {
	if strings.HasPrefix(packagePath, moduleName+"/") || packagePath == moduleName {
		return replacePrefix(packagePath, moduleName)
	}
	if !strings.Contains(packagePath, "/") {
		return moduleName
	}
	return packagePath
}

func replacePrefix(input, newPrefix string) string {
	index := strings.Index(input, "/")
	if index == -1 {
		return newPrefix
	}
	return newPrefix + input[index:]
}
