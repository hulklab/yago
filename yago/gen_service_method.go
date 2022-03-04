package main

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type serviceMethodGen struct {
	BaseGen
	Method     string
	Entry      string
	StructName string
	name       string
	dtoFile    string
	httpFile   string
	hasReq     bool
	hasResp    bool
	hasErr     bool
	hasHttp    bool
}

func NewServiceMethodGen(conf *viper.Viper) *serviceMethodGen {
	s := &serviceMethodGen{
		Entry:   conf.GetString("entry"),
		Method:  conf.GetString("name"),
		hasReq:  conf.GetBool("req"),
		hasResp: conf.GetBool("resp"),
		hasErr:  conf.GetBool("err"),
		hasHttp: conf.GetBool("http"),
	}
	s.Init(conf)

	fmt.Println("req:", conf.GetBool("req"), "entry", conf.GetString("entry"))

	structName := conf.GetString("structName")
	if len(structName) == 0 {
		structName = FindServiceStructName(s.File)
	}

	s.StructName = structName

	// 获取当前模块信息
	s.name = capitalize(strings.Replace(s.StructName, "Service", "", 1))
	s.dtoFile = fmt.Sprintf("%s/app/modules/%s/%sdto/%s.go", s.RootPath, s.ModuleName, s.ModuleName, s.Filename)
	if len(s.Entry) == 0 {
		s.httpFile = filepath.Join(s.RootPath, "app", "modules", s.ModuleName, s.ModuleName+"http", s.Filename+".go")
	} else {
		s.httpFile = filepath.Join(s.RootPath, "app", "modules", s.ModuleName, s.ModuleName+"http", s.Entry, s.Filename+".go")
	}

	//fmt.Println(s.modName, s.moduleName, s.name, s.dtoFile, s.httpFile)

	return s
}

func (s *serviceMethodGen) Gen() (err error) {
	// 获取 method-code
	content, imports := s.genServiceMethod()
	fmt.Println(content)

	// 生成 base-code
	var dContent string
	dContents := make([]string, 0)
	if s.hasReq {
		reqDto := s.genReqDto()
		if len(reqDto) > 0 {
			dContents = append(dContents, reqDto)
		}
	}

	if s.hasResp {
		respDto := s.genRespDto()
		if len(respDto) > 0 {
			dContents = append(dContents, respDto)
		}
	}

	if len(dContents) > 0 && !fileExists(s.dtoFile) {
		dContents = append([]string{fmt.Sprintf("package %sdto", s.ModuleName)}, dContents...)
	}

	if len(dContents) > 0 {
		dContent = strings.Join(dContents, "\n")
	}

	fmt.Println(dContent)

	hContent, hImports := s.genHttpMethod()
	fmt.Println(hContent)

	rContent := s.genHttpRouteMethod()

	fmt.Println(rContent)

	// 先写 base 文件
	if len(dContent) > 0 {
		writeFileAppendOrCreate(s.dtoFile, dContent)
		gofmt(s.dtoFile)
	}

	// 再写 service 文件
	writeInLineOrCreate(s.File, content, getGoLine())
	if len(imports) > 0 {
		addImports(s.File, imports)
	}
	gofmt(s.File)

	// 再写 http 文件
	if s.hasHttp && len(hContent) > 0 {
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
	}

	return nil

}

func (s *serviceMethodGen) genServiceMethod() (content string, imports []string) {
	imports = []string{}

	data := MethodTmplData{
		StructName: s.StructName,
		ModuleName: s.ModuleName,
		Method:     s.Method,
		HasReq:     s.hasReq,
		HasResp:    s.hasResp,
		HasErr:     s.hasErr,
	}

	content = ExecuteTemplate(MethodTemplate, data)

	if s.hasReq || s.hasResp {
		imports = append(imports, fmt.Sprintf(`"%s/%sdto"`, s.ModuleImportPath, s.ModuleName))
	}

	return
}

func (s *serviceMethodGen) genReqDto() string {
	structName := fmt.Sprintf("%sReq", s.Method)
	if isStructExists(s.dtoFile, structName) {
		return ""
	}

	content := ExecuteTemplate(DtoReqTemplate, map[string]string{"Method": s.Method})

	return content
}

func (s *serviceMethodGen) genRespDto() string {
	structName := fmt.Sprintf("%sResp", s.Method)
	if isStructExists(s.dtoFile, structName) {
		return ""
	}

	content := ExecuteTemplate(DtoRespTemplate, map[string]string{"Method": s.Method})
	return content
}

func (s *serviceMethodGen) genHttpMethod() (content string, imports []string) {
	if !s.hasHttp {
		return
	}

	actionName := fmt.Sprintf("%sAction", s.Method)
	if isMethodExists(s.httpFile, actionName) {
		return "", nil
	}

	data := HttpMethodTmplData{
		HasReq:     s.hasReq,
		HasResp:    s.hasResp,
		HasErr:     s.hasErr,
		CamelName:  s.name,
		Method:     s.Method,
		ModuleName: s.ModuleName,
		ReqPackage: fmt.Sprintf("%sdto", s.ModuleName),
		ReqName:    fmt.Sprintf("%sReq", s.Method),
	}
	content = ExecuteTemplate(HttpMethodTemplate, data)

	imports = []string{
		fmt.Sprintf(`"%s/%sservice"`, s.ModuleImportPath, s.ModuleName),
	}

	if s.hasReq {
		imports = append(imports, fmt.Sprintf(`"%s/%sdto"`, s.ModuleImportPath, s.ModuleName))
	}

	return
}

func (s *serviceMethodGen) genHttpRouteMethod() string {
	if !s.hasHttp {
		return ""
	}

	data := HttpRouteTmplData{
		ModuleName: s.ModuleName,
		LispName:   lispString(s.name),
		LispMethod: lispString(s.Method),
		Method:     s.Method,
	}

	content := ExecuteTemplate(HttpRouteTemplate, data)
	return content
}

func genServiceMethodCmd() *cobra.Command {
	// 定义二级命令: service-method
	var cmd = &cobra.Command{
		Use:   "gen-service-method",
		Short: "Gen service-method code",
		Run: func(cmd *cobra.Command, args []string) {
			s := NewServiceMethodGen(Conf)
			err := s.Gen()
			if err != nil {
				log.Fatalln(err)
			}
		},
	}

	cmd.Flags().StringP("name", "n", "", "方法名称")
	cmd.Flags().StringP("structName", "s", "", "结构体名称")
	cmd.Flags().StringP("entry", "e", "", "入口,admin front api openapi")
	cmd.Flags().Bool("req", true, "是否需要 request 参数")
	cmd.Flags().Bool("resp", true, "是否需要返回 resp 参数")
	cmd.Flags().Bool("err", true, "是否需要返回 error 参数")
	cmd.Flags().Bool("http", true, "是否需要 http 方法")
	cmd.Flags().StringP("file", "f", getGoFile(), "file path,eg. ./ab_c.go")

	_ = cmd.MarkFlagRequired("name")

	return cmd
}
