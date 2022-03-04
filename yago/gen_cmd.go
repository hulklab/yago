package main

import (
	"io/ioutil"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type cmdGen struct {
	BaseGen
}

func NewCmdGen(conf *viper.Viper) *cmdGen {
	g := new(cmdGen)
	g.Init(conf)
	return g
}

func (d *cmdGen) Gen() (err error) {

	name := camelString(d.Filename)
	if len(name) < 2 {
		log.Fatalln("the filename length is greater than at least 2")
	}

	data := CmdTmplData{
		Package: getPkgNameByFile(d.File),
		Name:    name,
	}

	content := ExecuteTemplate(CmdTemplate, data)

	err = ioutil.WriteFile(d.File, []byte(content), 0644)
	if err != nil {
		return err
	}

	gofmt(d.File)

	return nil
}

func genCommandCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "gen-cmd",
		Short: "Gen cmd code",
		Run: func(cmd *cobra.Command, args []string) {
			s := NewCmdGen(Conf)
			err := s.Gen()
			if err != nil {
				log.Fatalln(err)
			}
		},
	}

	cmd.Flags().StringP("file", "f", getGoFile(), "file path,eg. ./ab_c.go")

	return cmd
}
