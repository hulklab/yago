package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type httpGen struct {
	BaseGen
	HttpName     string
	HttpFile     string
	HttpPackage  string
	HttpLispName string
	entry        string
}

func NewHttpGen(conf *viper.Viper) *httpGen {
	h := new(httpGen)
	h.Init(conf)
	h.InitDaoInfo(conf.GetString("daoName"))
	h.entry = conf.GetString("entry")

	if h.DaoInfo == nil {
		h.HttpName = camelString(h.Filename)
	} else {
		h.HttpName = strings.TrimRight(h.DaoInfo.DaoName, "Dao")
	}

	if len(h.entry) == 0 {
		h.HttpFile = fmt.Sprintf("%s/%shttp/%s.go", h.ModulePath, h.ModuleName, snakeString(h.HttpName))
		h.HttpPackage = fmt.Sprintf("%shttp", h.ModuleName)
	} else {
		h.HttpFile = fmt.Sprintf("%s/%shttp/%s/%s.go", h.ModulePath, h.ModuleName, h.entry, snakeString(h.HttpName))
		h.HttpPackage = h.entry
	}

	h.HttpLispName = lispString(h.HttpName)

	if !conf.GetBool("overwrite") && fileExists(h.HttpFile) {
		log.Fatalf("http 文件 %s 已存在, 使用 -o 进行覆盖", h.HttpFile)
	}

	return h
}

func (h *httpGen) Gen() (err error) {
	data := HttpTmplData{
		Empty:          h.DaoInfo == nil,
		Package:        h.HttpPackage,
		Name:           h.HttpName,
		DtoPackage:     fmt.Sprintf("%sdto", h.ModuleName),
		ServicePackage: fmt.Sprintf("%sservice", h.ModuleName),
		ModName:        h.ModName,
		ModuleName:     h.ModuleName,
		Entry:          h.entry,
	}

	if h.DaoInfo != nil {
		data.AddRoute = h.genHttpRoutesByMethod("Add")
		data.DelRoute = h.genHttpRoutesByMethod("Delete")
		data.UpdateRoute = h.genHttpRoutesByMethod("Update")
		data.ListRoute = h.genHttpRoutesByMethod("List")
		data.DetailRoute = h.genHttpRoutesByMethod("Detail")

	}

	content := ExecuteTemplate(HttpTemplate, data)

	writeFileAppendOrCreate(h.HttpFile, content)

	gofmt(h.HttpFile)

	return nil
}

func (h *httpGen) genHttpRoutesByMethod(method string) string {
	lispMethod := lispString(method)
	data := HttpRouteTmplData{
		ModuleName: h.ModuleName,
		LispName:   h.HttpLispName,
		LispMethod: lispMethod,
		Method:     method,
	}

	return ExecuteTemplate(HttpRouteTemplate, data)
}

func genHttpCmd() *cobra.Command {
	// 定义二级命令: http
	var cmd = &cobra.Command{
		Use:   "gen-http",
		Short: "Gen http code",
		Run: func(cmd *cobra.Command, args []string) {
			s := NewHttpGen(Conf)
			err := s.Gen()
			if err != nil {
				log.Fatalln(err)
			}
		},
	}

	cmd.Flags().StringP("daoName", "t", "", "Dao 名称")
	cmd.Flags().StringP("entry", "e", "", "入口,admin front api open")
	cmd.Flags().BoolP("overwrite", "o", false, "是否覆盖已存在文件")
	cmd.Flags().StringP("file", "f", getGoFile(), "file path,eg. ./ab_c.go")

	return cmd
}
