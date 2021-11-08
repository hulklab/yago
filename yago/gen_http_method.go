package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type httpMethodGen struct {
	BaseGen
	Entry      string
	name       string
	httpFile   string
	methodInfo *ServiceMethodInfo
}

func NewHttpMethodGen(conf *viper.Viper) *httpMethodGen {
	s := new(httpMethodGen)
	s.Init(conf)
	s.Entry = conf.GetString("entry")

	methodInfo := findServiceMethod(s.File)
	methodInfo.Import = fmt.Sprintf("%s/%s", s.ModuleImportPath, methodInfo.Package)
	s.methodInfo = methodInfo

	// @todo support param
	s.name = capitalize(strings.Replace(s.methodInfo.StructName, "Service", "", 1))
	// s.filename = SnakeString(s.name)

	if len(s.Entry) == 0 {
		s.httpFile = filepath.Join(s.RootPath, "app", "modules", s.ModuleName, s.ModuleName+"http", s.Filename+".go")
	} else {
		s.httpFile = filepath.Join(s.RootPath, "app", "modules", s.ModuleName, s.ModuleName+"http", s.Entry, s.Filename+".go")
	}

	return s
}

type Param struct {
	Package string
	Name    string
	Import  string
}

type ServiceMethodInfo struct {
	Name       string
	Package    string
	Import     string
	StructName string
	Req        *Param
	Resp       *Param
	HasErr     bool
}

func findServiceMethod(path string) (methodInfo *ServiceMethodInfo) {

	fileSet := token.NewFileSet()

	file, err := parser.ParseFile(fileSet, path, nil, parser.AllErrors)
	if err != nil {
		fmt.Println("parse file err:", path, err)
		return
	}

	// fmt.Println("==============")
	// ast.Print(nil, file)

	line := getGoLine()

	imports := map[string]string{}

	for _, s := range file.Imports {
		path := strings.Trim(s.Path.Value, `"`)
		var name string
		if s.Name == nil {
			pathList := strings.Split(path, "/")
			name = pathList[len(pathList)-1]
		} else {
			name = s.Name.Name
		}
		imports[name] = s.Path.Value
	}

	// fmt.Printf("imports: %+v", imports)

	methodInfo = &ServiceMethodInfo{}
	methodInfo.Package = file.Name.Name
	// fmt.Println("package:", file.Name.Name)

	for _, decl := range file.Decls {
		funDecl, ok := decl.(*ast.FuncDecl)

		if !ok {
			continue
		}

		typeLine := fileSet.Position(decl.Pos()).Line

		if typeLine == line+1 {
			// var methodName string
			// var hasReq,hasResp,hasErr bool
			// var reqName string

			// methodName = funDecl.Name.Name

			// 结构体
			if funDecl.Recv != nil {
				recv := funDecl.Recv.List[0]
				if v, ok := recv.Type.(*ast.StarExpr); ok {
					if vv, ok := v.X.(*ast.Ident); ok {
						methodInfo.StructName = vv.Name
					}
				}
			}

			methodInfo.Name = funDecl.Name.Name

			if len(funDecl.Type.Params.List) > 0 {
				param := funDecl.Type.Params.List[0]

				// 判断是否为指针
				if v, ok := param.Type.(*ast.StarExpr); ok {
					// 判断是否为外部包
					if vv, ok := v.X.(*ast.SelectorExpr); ok {
						pkg := vv.X.(*ast.Ident).Name
						req := &Param{
							Name:    vv.Sel.Name,
							Package: pkg,
							Import:  imports[pkg],
						}

						methodInfo.Req = req
					}
				}
			}

			if funDecl.Type.Results != nil && len(funDecl.Type.Results.List) > 0 {
				for _, result := range funDecl.Type.Results.List {
					// fmt.Println("result:", i)
					if v, ok := result.Type.(*ast.StarExpr); ok {
						// 判断是否为外部包
						if vv, ok := v.X.(*ast.SelectorExpr); ok {
							pkg := vv.X.(*ast.Ident).Name
							resp := &Param{
								Name:    vv.Sel.Name,
								Package: pkg,
								Import:  imports[pkg],
							}

							methodInfo.Resp = resp
						}
						continue
					}

					if v, ok := result.Type.(*ast.Ident); ok {
						if v.Name == "error" {
							methodInfo.HasErr = true
						}
					}
				}

			}

			break
		}

	}

	if len(methodInfo.Name) == 0 || len(methodInfo.StructName) == 0 {
		log.Fatalln("get method info err")
	}

	return

}

func (s *httpMethodGen) Gen() (err error) {
	// 获取 method-code
	hContent, hImports := s.genHttpMethod()
	fmt.Println(hContent, hImports)

	rContent := s.genHttpRouteMethod()

	fmt.Println(rContent)

	// 再写 http 文件
	if !fileExists(s.httpFile) {
		conf := viper.New()
		conf.Set("file", s.httpFile)
		conf.Set("entry", s.Entry)
		hg := NewHttpGen(conf)

		err = hg.Gen()
		if err != nil {
			log.Fatalln("生成 http 文件失败", err)
		}
	}

	writeFileAppendOrCreate(s.httpFile, hContent)
	if len(hImports) > 0 {
		addImports(s.httpFile, hImports)
	}

	if len(rContent) > 0 {
		addHttpRoute(s.httpFile, rContent)
	}
	gofmt(s.httpFile)

	return nil

}

func (s *httpMethodGen) genHttpMethod() (content string, imports []string) {

	actionName := fmt.Sprintf("%sAction", s.methodInfo.Name)
	if isMethodExists(s.httpFile, actionName) {
		return "", nil
	}

	data := HttpMethodTmplData{
		HasReq:     s.methodInfo.Req != nil,
		HasResp:    s.methodInfo.Resp != nil,
		HasErr:     s.methodInfo.HasErr,
		CamelName:  s.name,
		Method:     s.methodInfo.Name,
		ModuleName: s.ModuleName,
		ReqPackage: If(s.methodInfo.Req != nil, s.methodInfo.Req.Package),
		ReqName:    If(s.methodInfo.Req != nil, s.methodInfo.Req.Name),
	}
	content = ExecuteTemplate(HttpMethodTemplate, data)

	imports = []string{
		fmt.Sprintf(`"%s/%s"`, s.ModuleImportPath, s.methodInfo.Package),
	}

	if s.methodInfo.Req != nil {
		imports = append(imports, s.methodInfo.Req.Import)
	}

	return
}

func (s *httpMethodGen) genHttpRouteMethod() string {
	data := HttpRouteTmplData{
		ModuleName: s.ModuleName,
		LispName:   lispString(s.name),
		LispMethod: lispString(s.methodInfo.Name),
		Method:     s.methodInfo.Name,
	}

	// @todo Entry -> Group

	content := ExecuteTemplate(HttpRouteTemplate, data)
	return content
}

func genHttpMethodCmd() *cobra.Command {
	var entry string

	var cmd = &cobra.Command{
		Use:   "gen-http-method",
		Short: "Gen http-method code",
		Run: func(cmd *cobra.Command, args []string) {
			s := NewHttpMethodGen(Conf)
			err := s.Gen()
			if err != nil {
				log.Fatalln(err)
			}
		},
	}

	cmd.Flags().StringVarP(&entry, "entry", "e", "", "入口, admin front api open")
	cmd.Flags().StringP("file", "f", getGoFile(), "file path,eg. ./ab_c.go")

	return cmd
}
