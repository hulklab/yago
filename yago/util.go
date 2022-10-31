package main

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

const (
	GOFILE    = "GOFILE"
	GOLINE    = "GOLINE"
	GOPACKAGE = "GOPACKAGE"
)

// GOFILE=ins.go GOLINE=132 GOPACKAGE=pgsqlservice

// -------------------- env lib --------------------------------

func getGoPath() string {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}
	return gopath
}

func isRunByGoGenerate() bool {
	gofile := os.Getenv(GOFILE)
	return len(gofile) > 0
}

func getFileByGoGenerate() (file string, err error) {
	gofile := os.Getenv(GOFILE)
	goline := os.Getenv(GOLINE)
	gopackage := os.Getenv(GOPACKAGE)

	if len(gofile) == 0 || len(goline) == 0 || len(gopackage) == 0 {
		err = fmt.Errorf("未获取到 GOFILE GOLINE GOPACKAGE 等环境变量，请确认是运行在 go:generate 下吗")
		return
	}

	err = filepath.Walk(pwd(), func(path string, f fs.FileInfo, err error) error {
		if f == nil {
			return err
		}

		if f.IsDir() {
			return nil
		}

		if filepath.Base(path) == gofile {
			file = path
			return filepath.SkipDir
		}

		return nil

	})

	return
}

func getGoLine() int {
	line := os.Getenv(GOLINE)

	i, _ := strconv.Atoi(line)
	return i
}

func getGoFile() string {
	file := os.Getenv(GOFILE)
	if len(file) > 0 {
		return getFileByFilename(os.Getenv(GOFILE))
	}
	return ""
}

//  ----------------- file lib -----------------------------

func pwd() string {
	dirPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err.Error())
	}

	return dirPath
}

func getPkgNameByFile(file string) string {
	return strings.ToLower(filepath.Base(filepath.Dir(file)))
}

func getModuleNameByFile(file string) string {
	rootPath := getProjectRootPathByFile(file)

	path := strings.TrimPrefix(strings.Replace(file, rootPath, "", 1), string(filepath.Separator))

	pathList := strings.Split(path, string(filepath.Separator))

	for i, v := range pathList {
		if v == "modules" {
			return pathList[i+1]
		}
	}

	// log.Fatalf("module name is not found in path %s", path)

	return ""
}

func getModulePathByFile(file string) string {
	rootPath := getProjectRootPathByFile(file)

	path := strings.TrimPrefix(strings.Replace(file, rootPath, "", 1), string(filepath.Separator))

	pathList := strings.Split(path, string(filepath.Separator))

	modulePath := []string{rootPath}
	hasModule := false
	for i, v := range pathList {
		modulePath = append(modulePath, v)
		if v == "modules" {
			modulePath = append(modulePath, pathList[i+1])
			hasModule = true
			break
		}
	}

	if !hasModule {
		log.Fatalf("can not found modules in the file path %s", file)
	}

	return filepath.Join(modulePath...)

}

func getFileNameByFile(file string) string {
	return strings.Replace(filepath.Base(file), ".go", "", 1)
}

var rootPath string
var once sync.Once

func getProjectRootPathByFile(file string) string {
	once.Do(func() {
		modPath := file
		if !isDir(file) {
			modPath = filepath.Dir(file)
		}

		var modFile string

		for {
			modFile = fmt.Sprintf("%s/go.mod", modPath)

			if fileExists(modFile) {
				break
			}

			modPath = filepath.Dir(modPath)
			if strings.Count(modPath, string(filepath.Separator)) == 0 {
				log.Fatalf("go.mod file is not found")
			}
		}

		rootPath = modPath

	})

	return rootPath
}

func getModNameByFile(file string) string {
	root := getProjectRootPathByFile(file)
	modFile := fmt.Sprintf("%s/go.mod", root)

	fi, err := os.Open(modFile)
	if err != nil {
		log.Fatalln(err)
	}

	defer func() { _ = fi.Close() }()

	br := bufio.NewReader(fi)

	spaceReg := regexp.MustCompile(`\s+`)

	for {
		line, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}

		lineStr := string(line)

		if strings.HasPrefix(lineStr, "module") {
			names := spaceReg.Split(lineStr, 3)
			if len(names) >= 2 {
				return names[1]
			}

			break
		}
	}

	log.Fatalln("get mod name failed")
	return ""

}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func isDir(filename string) bool {
	f, err := os.Stat(filename)

	return err == nil && f.IsDir()
}

func writeFileAppendOrCreate(filename, content string) {
	if !fileExists(filepath.Dir(filename)) {
		// 创建目录
		err := os.MkdirAll(filepath.Dir(filename), 0775)
		if err != nil {
			log.Fatalln(err)
		}
	}

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err = f.WriteString(content); err != nil {
		log.Fatalln(err)
	}
}

func gofmt(file string) {
	cmd := exec.Command("gofmt", "-w", file)
	if err := cmd.Run(); err != nil {
		log.Printf("gofmt file %s err:%s", file, err.Error())
	}
}

func writeInLineOrCreate(file string, content string, line int) {
	if !fileExists(file) || line <= 0 {
		writeFileAppendOrCreate(file, content)
		return
	}

	f, err := os.Open(file)
	if err != nil {
		log.Fatalln(file, err)
	}

	defer f.Close()
	r := bufio.NewScanner(f)

	contents := []string{}
	for r.Scan() {
		row := r.Text()
		contents = append(contents, row)
	}

	contents = append(contents[0:line], append([]string{content}, contents[line:]...)...)

	err = ioutil.WriteFile(file, []byte(strings.Join(contents, "\n")), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

func addImports(file string, imports []string) {
	f, err := os.Open(file)
	if err != nil {
		log.Fatalln(file, err)
	}

	defer f.Close()
	r := bufio.NewScanner(f)

	importPos := 0
	contents := []string{}
	for r.Scan() {
		line := r.Text()
		if strings.Contains(line, "import (") {
			importPos = len(contents) + 1
		}

		contents = append(contents, line)
	}

	contents = append(contents[0:importPos], append(imports, contents[importPos:]...)...)

	err = ioutil.WriteFile(file, []byte(strings.Join(contents, "\n")), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

func addHttpRoute(file string, content string) {
	f, err := os.Open(file)
	if err != nil {
		log.Fatalln(file, err)
	}

	defer f.Close()
	r := bufio.NewScanner(f)

	initPos := 0
	initFlag := false
	contents := []string{}
	for r.Scan() {
		line := r.Text()
		if strings.Contains(line, "func init()") {
			initFlag = true
		}

		if initFlag && strings.HasPrefix(line, "}") {
			initPos = len(contents)
			initFlag = false
		}
		contents = append(contents, line)
		//fmt.Println(len(contents),line)
	}

	//fmt.Println("init pos",initPos)

	rs := []string{
		content,
	}

	contents = append(contents[0:initPos], append(rs, contents[initPos:]...)...)

	err = ioutil.WriteFile(file, []byte(strings.Join(contents, "\n")), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

func getFileByFilename(filename string) (file string) {
	if !strings.HasSuffix(filename, ".go") {
		filename = fmt.Sprintf("%s.go", filename)
	}
	file = filepath.Join(pwd(), filename)
	return
}

// ------------------------str lib -------------------------------

func capitalize(str string) string {
	var upperStr string
	vv := []rune(str) // 后文有介绍
	for i := 0; i < len(vv); i++ {
		if i == 0 {
			if vv[i] >= 97 && vv[i] <= 122 { // 后文有介绍
				vv[i] -= 32 // string的码表相差32位
				upperStr += string(vv[i])
			} else {
				return str
			}
		} else {
			upperStr += string(vv[i])
		}
	}
	return upperStr
}

func lcFirst(s string) string {

	if len(s) <= 1 {
		return strings.ToLower(s)
	}

	return fmt.Sprintf("%s%s", strings.ToLower(s[0:1]), s[1:])
}

func ucFirst(s string) string {
	if len(s) <= 1 {
		return strings.ToUpper(s)
	}

	return fmt.Sprintf("%s%s", strings.ToUpper(s[0:1]), s[1:])
}

func camelString(s string) string {
	data := make([]byte, 0, len(s))
	j := false
	k := false
	num := len(s) - 1
	for i := 0; i <= num; i++ {
		d := s[i]
		if !k && d >= 'A' && d <= 'Z' {
			k = true
		}
		if d >= 'a' && d <= 'z' && (j || !k) {
			d = d - 32
			j = false
			k = true
		}
		if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
			j = true
			continue
		}
		data = append(data, d)
	}
	return string(data[:])
}

func snakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		} else {
			j = false
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}

func lispString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '-')
		}
		if d != '-' {
			j = true
		}
		data = append(data, d)
	}

	return strings.ToLower(string(data[:]))
}

func If(b bool, s string) string {
	if b {
		return s
	}

	return ""
}

// -------------------------- ast -------------------------
type BaseGen struct {
	File             string
	Filename         string
	RootPath         string
	ModName          string
	ModulePath       string
	ModuleImportPath string
	ModuleName       string
	DaoInfo          *DaoInfo
}

func (g *BaseGen) Init(conf *viper.Viper) {
	var err error
	file := conf.GetString("file")
	if len(file) == 0 {
		log.Fatalln("file path is required")
	}

	f, err := filepath.Abs(file)
	if err != nil {
		log.Fatalf("get file %s abs err:%s", file, err.Error())
	}

	g.File = file
	g.Filename = strings.TrimSuffix(getFileNameByFile(f), ".go")

	g.RootPath = getProjectRootPathByFile(file)
	g.ModName = getModNameByFile(file)
	g.ModuleName = getModuleNameByFile(file)
	if len(g.ModuleName) > 0 {
		g.ModulePath = getModulePathByFile(file)
		g.ModuleImportPath = fmt.Sprintf("%s/%s", g.ModName, strings.TrimPrefix(strings.Replace(g.ModulePath, g.RootPath, "", 1), string(filepath.Separator)))
	}

}

func (g *BaseGen) InitDaoInfo(daoName string) {
	// 处理 dao
	if len(daoName) == 0 {
		var b bool
		b, daoName = FindDaoName(g.File)
		if b {
			g.DaoInfo = g.genDaoInfo(daoName)
		}
	} else {
		g.DaoInfo = g.genDaoInfo(daoName)
	}
}

func (g *BaseGen) genDaoInfo(daoName string) (dao *DaoInfo) {
	dao = &DaoInfo{
		DaoName: daoName,
		Gener:   g,
	}

	b := dao.genDaoInfo()
	if !b {
		log.Fatalf("Dao %s is not found", dao.DaoName)
	}

	return dao
}

func findTypeName(path string) (b bool, typeName string) {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, path, nil, parser.AllErrors)
	if err != nil {
		fmt.Println(err)
		return
	}

	line := getGoLine()

	ast.Inspect(file, func(node ast.Node) bool {
		switch v := node.(type) {
		case *ast.TypeSpec:
			typeLine := fileSet.Position(v.Pos()).Line

			if v.Name == nil {
				return true
			}

			if typeLine == line+1 {
				typeName = v.Name.Name
				return false
			}
		}

		return true
	})

	if len(typeName) == 0 {
		return
	}

	return true, typeName
}

func FindDaoName(path string) (b bool, daoName string) {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, path, nil, parser.AllErrors)
	if err != nil {
		fmt.Println(err)
		return
	}

	line := getGoLine()

	ast.Inspect(file, func(node ast.Node) bool {
		switch v := node.(type) {
		case *ast.TypeSpec:
			typeLine := fileSet.Position(v.Pos()).Line

			if v.Name == nil {
				return true
			}

			if typeLine == line+1 {
				daoName = v.Name.Name
				return false
			}
		}

		return true
	})

	if len(daoName) == 0 {
		return
	}

	return true, daoName
}

func FindServiceStructName(path string) (serviceName string) {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, path, nil, parser.AllErrors)
	if err != nil {
		fmt.Println(err)
		return
	}

	ast.Inspect(file, func(node ast.Node) bool {
		switch v := node.(type) {
		case *ast.TypeSpec:
			if v.Name == nil {
				return true
			}

			if strings.HasSuffix(v.Name.Name, "Service") {
				serviceName = v.Name.Name
				return false
			}

		}

		return true
	})

	if len(serviceName) == 0 {
		log.Fatalf("can't find SerivceStruct in file %s", path)
	}

	return
}

func isMethodExists(path string, methodName string) (b bool) {
	if !fileExists(path) {
		return
	}

	fileSet := token.NewFileSet()

	file, err := parser.ParseFile(fileSet, path, nil, parser.AllErrors)
	if err != nil {
		fmt.Println("parse dir err:", path, err)
		return
	}

	for _, decl := range file.Decls {
		funDecl, ok := decl.(*ast.FuncDecl)

		if !ok {
			continue
		}

		if funDecl.Name.Name == methodName {
			return true
		}
	}
	return
}

// 判断某个结构体是否已存在
func isStructExists(path string, structName string) (b bool) {
	var dir = path
	if !isDir(path) {
		dir = filepath.Dir(path)
	}

	fileSet := token.NewFileSet()

	pkgs, err := parser.ParseDir(fileSet, dir, nil, parser.AllErrors)
	if err != nil {
		fmt.Println("parse dir err:", dir, err)
		return
	}

	for _, pkg := range pkgs {
		for filename, file := range pkg.Files {
			fmt.Println("filename:", filename)
			if b {
				return
			}
			ast.Inspect(file, func(node ast.Node) bool {
				switch v := node.(type) {
				case *ast.TypeSpec:
					if v.Name == nil {
						return true
					}

					if v.Name.Name == structName {
						b = true
						return false
					}
				}

				return true
			})

		}
	}

	return
}
