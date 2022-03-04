package main

import (
	"io/ioutil"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type rpcGen struct {
	BaseGen
}

func NewRpcGen(conf *viper.Viper) *rpcGen {
	d := new(rpcGen)
	d.Init(conf)

	return d
}

func (d *rpcGen) Gen() (err error) {

	name := camelString(d.Filename)
	if len(name) < 2 {
		log.Fatalln("the filename length is greater than at least 2")
	}

	data := RpcTmplData{
		Package: getPkgNameByFile(d.File),
		Name:    name,
	}

	content := ExecuteTemplate(RpcTemplate, data)

	err = ioutil.WriteFile(d.File, []byte(content), 0644)
	if err != nil {
		return err
	}

	gofmt(d.File)

	return nil
}

func genRpcCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "gen-rpc",
		Short: "Gen rpc code",
		Run: func(cmd *cobra.Command, args []string) {
			s := NewRpcGen(Conf)
			err := s.Gen()
			if err != nil {
				log.Fatalln(err)
			}
		},
	}

	cmd.Flags().StringP("file", "f", getGoFile(), "file path,eg. ./ab_c.go")

	return cmd
}
