package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type typeGen struct {
	BaseGen
	typeName string
}

func NewTypeGen(conf *viper.Viper) *typeGen {
	d := &typeGen{}
	d.Init(conf)

	var typeName = conf.GetString("typeName")
	if len(typeName) == 0 {
		// 从当前行寻找
		var b bool
		b, typeName = findTypeName(d.File)
		if !b {
			log.Fatalln("未找到类型")
		}
	}

	d.typeName = typeName

	return d
}

func (d *typeGen) genTypmethod() (t string, err error) {
	fileSet := token.NewFileSet()

	file, err := parser.ParseFile(fileSet, d.File, nil, parser.ParseComments)
	if err != nil {
		return
	}

	ast.Inspect(file, func(node ast.Node) bool {

		switch v := node.(type) {
		case *ast.TypeSpec:
			if v.Name == nil {
				return true
			}

			if v.Name.Name == d.typeName {
				t = v.Type.(*ast.Ident).Name
				return false
			}
		}

		return true
	})

	return
}

func (d *typeGen) genTypmap() (mapAssign string, err error) {
	fileSet := token.NewFileSet()

	// https://pkg.go.dev/go/parser#Mode
	file, err := parser.ParseFile(fileSet, d.File, nil, parser.ParseComments)
	if err != nil {
		return
	}

	ast.Inspect(file, func(n ast.Node) bool {
		switch v := n.(type) {
		case *ast.ValueSpec:
			if len(v.Names) == 0 {
				return true
			}

			if v.Comment == nil {
				return true
			}

			if vv, ok := v.Type.(*ast.Ident); !ok || vv.Name != d.typeName {
				return true
			}

			if v.Type.(*ast.Ident).Name != d.typeName {
				return true
			}

			// fmt.Printf("%T\n", v.Type)

			mapAssign += fmt.Sprintf("\t\t%s:\t\"%s\",\n", v.Names[0], strings.TrimSpace(v.Comment.Text()))

			return true
		}

		return true
	})

	// fmt.Println(mapAssign)

	return
}

func (d *typeGen) Gen() (err error) {

	var mapAssign string
	mapAssign, err = d.genTypmap()
	if err != nil {
		return err
	}

	tName, err := d.genTypmethod()
	if err != nil {
		return err
	}

	data := TypeTmplData{
		TypeName:  d.typeName,
		UType:     ucFirst(tName),
		Type:      tName,
		MapAssign: mapAssign,
	}
	content := ExecuteTemplate(TypeTemplate, data)

	fmt.Println(content)

	line := 0
	if getGoLine() > 0 {
		line = getGoLine() + 1
	}

	writeInLineOrCreate(d.File, content, line)

	gofmt(d.File)

	return nil
}

func genTypeCmd() *cobra.Command {
	// 定义二级命令: model
	var cmd = &cobra.Command{
		Use:   "gen-type",
		Short: "Gen custom type code",
		Run: func(cmd *cobra.Command, args []string) {
			s := NewTypeGen(Conf)
			err := s.Gen()
			if err != nil {
				log.Fatalln(err)
			}
		},
	}

	cmd.Flags().StringP("typeName", "t", "", "type 名称,eg. Status")
	cmd.Flags().StringP("file", "f", getGoFile(), "file path,eg. ./ab_c.go")

	return cmd
}
