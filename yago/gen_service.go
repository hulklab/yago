package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type serviceGen struct {
	BaseGen
	ServiceName    string
	ServiceLName   string
	ServiceFile    string
	ServicePackage string
	DtoFile        string
	DtoPackage     string
	ModelPackage   string
}

func NewServiceGen(conf *viper.Viper) *serviceGen {
	s := new(serviceGen)
	s.Init(conf)
	s.InitDaoInfo(conf.GetString("daoName"))

	s.ServicePackage = fmt.Sprintf("%sservice", s.ModuleName)
	s.DtoPackage = fmt.Sprintf("%sdto", s.ModuleName)
	s.ModelPackage = fmt.Sprintf("%smodel", s.ModuleName)

	if s.DaoInfo != nil {
		s.ServiceName = strings.TrimRight(s.DaoInfo.DaoName, "Dao")
	} else {
		s.ServiceName = camelString(s.Filename)
	}
	s.ServiceLName = lcFirst(s.ServiceName)

	s.ServiceFile = fmt.Sprintf("%s/%sservice/%s.go", s.ModulePath, s.ModuleName, snakeString(s.ServiceName))
	s.DtoFile = fmt.Sprintf("%s/%sdto/%s.go", s.ModulePath, s.ModuleName, snakeString(s.ServiceName))

	if !conf.GetBool("overwrite") && fileExists(s.ServiceFile) {
		log.Fatalf("service 文件 %s 已存在, 可以使用 -o 进行覆盖", s.ServiceFile)
	}

	return s
}

func (s *serviceGen) genService() (err error) {
	data := ServiceTmplData{
		Empty:   s.DaoInfo == nil,
		Lname:   s.ServiceLName,
		Name:    s.ServiceName,
		ModName: s.ModName,
		Package: s.ServicePackage,
	}

	if s.DaoInfo != nil {
		data.DaoImportPath = s.DaoInfo.DaoImportPath
		data.ModelPackage = s.ModelPackage
		data.DtoPackage = s.DtoPackage
		data.ModuleImportPath = s.ModuleImportPath
		data.FieldMap = s.DaoInfo.FieldMap
		data.DaoPackage = s.DaoInfo.DaoPackage
		data.DaoName = s.DaoInfo.DaoName
	}

	content := ExecuteTemplate(ServiceTemplate, data)

	err = ioutil.WriteFile(s.ServiceFile, []byte(content), 0644)
	if err != nil {
		return err
	}

	gofmt(s.ServiceFile)
	return nil
}

func (s *serviceGen) genDto() (err error) {
	if s.DaoInfo == nil {
		return nil
	}

	if isStructExists(s.DtoFile, fmt.Sprintf("%sAddReq", s.ServiceName)) {
		return nil
	}

	data := ServiceDtoTmplData{
		Package:  s.DtoPackage,
		Name:     s.ServiceName,
		FieldMap: s.DaoInfo.FieldMap,
	}

	content := ExecuteTemplate(ServiceDtoTemplate, data)

	writeFileAppendOrCreate(s.DtoFile, content)
	if err != nil {
		return err
	}

	gofmt(s.DtoFile)
	return nil
}
func (s *serviceGen) Gen() (err error) {
	err = s.genService()
	if err != nil {
		return err
	}

	err = s.genDto()

	return err
}

func genServiceCmd() *cobra.Command {
	var daoName string
	var overwrite bool

	// 定义二级命令: service
	var cmd = &cobra.Command{
		Use:   "gen-service",
		Short: "Gen service code",
		Run: func(cmd *cobra.Command, args []string) {
			s := NewServiceGen(Conf)
			err := s.Gen()
			if err != nil {
				log.Fatalln(err)
			}
		},
	}

	cmd.Flags().StringVarP(&daoName, "daoName", "t", "", "Dao 名称")
	cmd.Flags().BoolVarP(&overwrite, "overwrite", "o", false, "是否覆盖已存在文件")
	cmd.Flags().StringP("file", "f", getGoFile(), "file path,eg. ./ab_c.go")

	return cmd
}
