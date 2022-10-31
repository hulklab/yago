package main

import (
	"io/ioutil"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type apiGen struct {
	BaseGen
}

func NewApiGen(conf *viper.Viper) *apiGen {
	g := new(apiGen)
	g.Init(conf)
	return g
}

func (d *apiGen) Gen() (err error) {

	name := camelString(d.Filename)
	if len(name) < 2 {
		log.Fatalln("the filename length is greater than at least 2")
	}

	data := ApiTmplData{
		Package: getPkgNameByFile(d.File),
		LName:   lcFirst(name),
		OName:   d.Filename,
	}

	content := ExecuteTemplate(ApiTemplate, data)

	err = ioutil.WriteFile(d.File, []byte(content), 0644)
	if err != nil {
		return err
	}

	gofmt(d.File)

	return nil
}

func genApiCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "gen-api",
		Short: "Gen api code",
		Run: func(cmd *cobra.Command, args []string) {
			s := NewApiGen(Conf)
			err := s.Gen()
			if err != nil {
				log.Fatalln(err)
			}
		},
	}

	cmd.Flags().StringP("file", "f", getGoFile(), "file path,eg. ./ab_c.go")

	return cmd
}
