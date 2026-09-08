package gen_mapper

// import (
// 	"bytes"
// 	"fmt"
// 	"go/types"
// 	"io/fs"
// 	"os"
// 	"path/filepath"
// 	"strings"

// 	"golang.org/x/tools/go/packages"
// )

// // aliasReplace: dùng để tìm struct trong source code thật
// var aliasReplace = map[string]string{
// 	"pb":   "../protobuf",
// 	"auth": "../..", // thêm dòng này
// }

// type PkgObj struct {
// 	ImportPath string // để import trong code gen
// 	LoadPath   string // để packages.Load tìm
// 	Obj        string // struct name
// 	Alias      string // alias trong import
// }

// // --- utils ---

// func getModuleName(serviceName string) string {
// 	// đi ngược lên 3 folder tìm go.mod (theo targetDir trong main)
// 	data, err := os.ReadFile("../../../" + serviceName + "-service/go.mod")
// 	if err != nil {
// 		fmt.Println("⚠️ cannot read go.mod:", err)
// 		return ""
// 	}
// 	for _, line := range strings.Split(string(data), "\n") {
// 		if strings.HasPrefix(line, "module ") {
// 			return strings.TrimSpace(strings.TrimPrefix(line, "module "))
// 		}
// 	}
// 	return ""
// }

// func resolveForLoad(pkgPath string, serviceName string) string {
// 	// alias pb => ../protobuf
// 	for k, v := range aliasReplace {
// 		if strings.HasPrefix(pkgPath, k+"/") {
// 			pkgPath = strings.Replace(pkgPath, k, v, 1)
// 		}
// 	}

// 	// ../protobuf/... => protobuf/...
// 	if strings.HasPrefix(pkgPath, "../") {
// 		pkgPath = strings.TrimPrefix(pkgPath, "../")
// 	}

// 	// thêm module name nếu thiếu
// 	modName := getModuleName(serviceName)
// 	if modName != "" && !strings.HasPrefix(pkgPath, modName) {
// 		pkgPath = modName + "/" + pkgPath
// 	}
// 	return pkgPath
// }

// func resolveForImport(pkgPath string, serviceName string) string {
// 	// giữ nguyên import path như comment
// 	return pkgPath
// }

// func splitPkgObj(full string, serviceName string) PkgObj {
// 	parts := strings.Split(full, ".")
// 	if len(parts) < 2 {
// 		panic("invalid mapper path: " + full)
// 	}
// 	pkgPath := strings.Join(parts[:len(parts)-1], ".")
// 	obj := parts[len(parts)-1]

// 	importPath := resolveForImport(pkgPath, serviceName)
// 	loadPath := resolveForLoad(pkgPath, serviceName)
// 	alias := strings.ReplaceAll(importPath, "/", "_")

// 	return PkgObj{
// 		ImportPath: importPath,
// 		LoadPath:   loadPath,
// 		Obj:        obj,
// 		Alias:      alias,
// 	}
// }

// func loadStruct(pkgPath, obj string) *types.Struct {
// 	cfg := &packages.Config{Mode: packages.NeedTypes | packages.NeedSyntax | packages.NeedDeps}
// 	pkgs, err := packages.Load(cfg, pkgPath)
// 	if err != nil || len(pkgs) == 0 {
// 		panic(fmt.Sprintf("cannot load package %s: %v", pkgPath, err))
// 	}
// 	scope := pkgs[0].Types.Scope()
// 	if scope == nil {
// 		panic("no scope in package: " + pkgPath)
// 	}
// 	typ := scope.Lookup(obj)
// 	if typ == nil {
// 		panic(fmt.Sprintf("struct %s not found in %s", obj, pkgPath))
// 	}
// 	named, ok := typ.Type().(*types.Named)
// 	if !ok {
// 		panic(fmt.Sprintf("%s is not a named type", obj))
// 	}
// 	st, ok := named.Underlying().(*types.Struct)
// 	if !ok {
// 		panic(fmt.Sprintf("%s is not a struct", obj))
// 	}
// 	return st
// }

// // --- main ---

// func main() {
// 	if len(os.Args) < 2 {
// 		fmt.Println("Usage: go run map_gen.go <serviceName>")
// 		os.Exit(1)
// 	}
// 	serviceName := strings.TrimSpace(os.Args[1])
// 	targetDir, err := filepath.Abs("../../../" + serviceName + "-service/infra/mapper")
// 	if err != nil {
// 		panic(err)
// 	}

// 	// quét mapper folder
// 	err = filepath.WalkDir(targetDir, func(path string, d fs.DirEntry, err error) error {
// 		if d.IsDir() || !strings.HasSuffix(d.Name(), ".go") {
// 			return nil
// 		}
// 		processFile(path, serviceName)
// 		return nil
// 	})
// 	if err != nil {
// 		panic(err)
// 	}
// }

// func processFile(path string, serviceName string) {
// 	data, err := os.ReadFile(path)
// 	if err != nil {
// 		panic(err)
// 	}
// 	lines := strings.Split(string(data), "\n")

// 	var mappers []struct {
// 		src, dst PkgObj
// 	}

// 	for _, line := range lines {
// 		line = strings.TrimSpace(line)
// 		if strings.HasPrefix(line, "// @mapper:") {
// 			line = strings.TrimPrefix(line, "// @mapper:")
// 			parts := strings.Split(strings.TrimSpace(line), "-")
// 			if len(parts) != 2 {
// 				continue
// 			}
// 			src := splitPkgObj(strings.TrimSpace(parts[0]), serviceName)
// 			dst := splitPkgObj(strings.TrimSpace(parts[1]), serviceName)
// 			mappers = append(mappers, struct {
// 				src, dst PkgObj
// 			}{src, dst})
// 		}
// 	}

// 	if len(mappers) == 0 {
// 		return
// 	}

// 	var buf bytes.Buffer
// 	buf.WriteString("package mapper\n\n")
// 	buf.WriteString("import (\n")
// 	imported := map[string]bool{}
// 	for _, m := range mappers {
// 		for _, pkg := range []PkgObj{m.src, m.dst} {
// 			if !imported[pkg.ImportPath] {
// 				buf.WriteString(fmt.Sprintf("\t%s \"%s\"\n", pkg.Alias, pkg.ImportPath))
// 				imported[pkg.ImportPath] = true
// 			}
// 		}
// 	}
// 	buf.WriteString(")\n\n")

// 	for _, m := range mappers {
// 		srcStruct := loadStruct(m.src.LoadPath, m.src.Obj)
// 		dstStruct := loadStruct(m.dst.LoadPath, m.dst.Obj)

// 		fnName := fmt.Sprintf("Map_%s_To_%s", m.src.Obj, m.dst.Obj)
// 		buf.WriteString(fmt.Sprintf("func %s(src *%s.%s) *%s.%s {\n",
// 			fnName, m.src.Alias, m.src.Obj, m.dst.Alias, m.dst.Obj))
// 		buf.WriteString("\tif src == nil { return nil }\n")
// 		buf.WriteString(fmt.Sprintf("\tdst := &%s.%s{}\n", m.dst.Alias, m.dst.Obj))

// 		srcFields := map[string]types.Type{}
// 		for i := 0; i < srcStruct.NumFields(); i++ {
// 			f := srcStruct.Field(i)
// 			if f.Exported() {
// 				srcFields[f.Name()] = f.Type()
// 			}
// 		}
// 		for i := 0; i < dstStruct.NumFields(); i++ {
// 			f := dstStruct.Field(i)
// 			if f.Exported() {
// 				if t, ok := srcFields[f.Name()]; ok && types.Identical(f.Type(), t) {
// 					buf.WriteString(fmt.Sprintf("\tdst.%s = src.%s\n", f.Name(), f.Name()))
// 				}
// 			}
// 		}
// 		buf.WriteString("\treturn dst\n")
// 		buf.WriteString("}\n\n")
// 	}

// 	// append vào cuối file mapper
// 	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer f.Close()
// 	if _, err := f.Write(buf.Bytes()); err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("Generated mapper code in:", path)
// }
