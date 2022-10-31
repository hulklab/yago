package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type modelGen struct {
	BaseGen
	ModelName    string
	ModelFile    string
	ModelPackage string
}

func NewModelGen(conf *viper.Viper) *modelGen {
	s := new(modelGen)
	s.Init(conf)
	s.InitDaoInfo(conf.GetString("daoName"))

	// modelPackage
	s.ModelPackage = fmt.Sprintf("%smodel", s.ModuleName)
	// modelName
	if s.DaoInfo != nil {
		s.ModelName = strings.TrimSuffix(s.DaoInfo.DaoName, "Dao")
	} else {
		s.ModelName = camelString(s.Filename)
	}

	// modelFile
	s.ModelFile = fmt.Sprintf("%s/%smodel/%s.go", s.ModulePath, s.ModuleName, snakeString(s.ModelName))

	if !conf.GetBool("overwrite") && fileExists(s.ModelFile) {
		log.Fatalf("model 文件 %s 已存在,可以使用 -o 进行覆盖", s.ModelFile)
	}

	return s
}

func (m *modelGen) Gen() (err error) {
	lname := lcFirst(m.ModelName)

	data := ModelTmplData{
		Empty:   m.DaoInfo == nil,
		Package: m.ModelPackage,
		ModName: m.ModName,
		Lname:   lname,
		Name:    m.ModelName,
	}

	if m.DaoInfo != nil {
		data.DaoImportPath = m.DaoInfo.DaoImportPath
		data.DaoName = m.DaoInfo.DaoName
		data.DaoPackage = m.DaoInfo.DaoPackage
	}

	content := ExecuteTemplate(ModelTemplate, data)

	err = ioutil.WriteFile(m.ModelFile, []byte(content), 0644)
	if err != nil {
		return err
	}

	gofmt(m.ModelFile)

	return nil
}

func genModelCmd() *cobra.Command {
	// 定义二级命令: model
	var cmd = &cobra.Command{
		Use:   "gen-model",
		Short: "Gen model code",
		Run: func(cmd *cobra.Command, args []string) {
			s := NewModelGen(Conf)
			err := s.Gen()
			if err != nil {
				log.Fatalln(err)
			}
		},
	}

	cmd.Flags().StringP("daoName", "t", "", "Dao 名称,eg. UserDao")
	cmd.Flags().BoolP("overwrite", "o", false, "是否覆盖已存在文件")
	cmd.Flags().StringP("file", "f", getGoFile(), "file path,eg. ./ab_c.go")

	return cmd
}
