package main

import (
	"io/ioutil"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type taskGen struct {
	BaseGen
}

func NewTaskGen(conf *viper.Viper) *taskGen {
	d := &taskGen{}
	d.Init(conf)

	return d
}

func (d *taskGen) Gen() (err error) {

	name := camelString(d.Filename)
	if len(name) < 2 {
		log.Fatalln("the filename length is greater than at least 2")
	}

	data := TaskTmplData{
		Package: getPkgNameByFile(d.File),
		Name:    name,
	}

	content := ExecuteTemplate(TaskTemplate, data)

	err = ioutil.WriteFile(d.File, []byte(content), 0644)
	if err != nil {
		return err
	}

	gofmt(d.File)

	return nil
}

func genTaskCmd() *cobra.Command {

	var cmd = &cobra.Command{
		Use:   "gen-task",
		Short: "Gen task code",
		Run: func(cmd *cobra.Command, args []string) {
			s := NewTaskGen(Conf)
			err := s.Gen()
			if err != nil {
				log.Fatalln(err)
			}
		},
	}

	cmd.Flags().StringP("file", "f", getGoFile(), "file path,eg. ./ab_c.go")

	return cmd
}
