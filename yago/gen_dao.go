package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"log"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/Shelnutt2/db2struct"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type field struct {
	Name      string
	Type      string
	SnakeName string
	tagMap    map[string]string
}

func (f field) NotShowJson() bool {
	if len(f.tagMap) > 0 && f.tagMap["json"] == "-" {
		return true
	}

	return false
}

func (f field) NotShowXorm() bool {
	if len(f.tagMap) > 0 && f.tagMap["xorm"] == "-" {
		return true
	}

	return false
}

func (f field) IsOrmHook() bool {
	if len(f.tagMap) > 0 && (f.tagMap["xorm"] == "updated" || f.tagMap["xorm"] == "created" || f.tagMap["xorm"] == "deleted") {
		return true
	}

	return false
}

type DaoInfo struct {
	Filename      string
	DaoPackage    string
	DaoImportPath string
	FieldMap      map[string]field
	DaoName       string
	Gener         *BaseGen
}

func (d *DaoInfo) genDaoInfo() (b bool) {
	paths := d.listDaoFiles()

	for _, path := range paths {
		b, err := d.inspectDao(path)
		if !b {
			continue
		}
		if err != nil {
			fmt.Println("err:", err.Error())
			continue
		}

		return b
	}

	return false
}

func (d *DaoInfo) listDaoFiles() (paths []string) {
	daoPath := fmt.Sprintf("%s/app/modules/%s/%sdao", d.Gener.RootPath, d.Gener.ModuleName, d.Gener.ModuleName)

	paths = make([]string, 0)
	err := filepath.Walk(daoPath, func(path string, f fs.FileInfo, err error) error {
		if f == nil {
			return err
		}

		if strings.HasPrefix(f.Name(), ".") {
			return nil
		}

		if f.IsDir() {
			return nil
		}

		// 只取以 .go 结尾的文件
		if !strings.HasSuffix(f.Name(), ".go") {
			return nil
		}

		//fmt.Println(path, f.Name())
		paths = append(paths, path)
		return nil
	})

	if err != nil {
		log.Fatalln(err)
	}

	return paths
}

// @refer https://github.com/chai2010/go-ast-book/blob/master/ch7/readme.md
func (d *DaoInfo) inspectDao(path string) (b bool, err error) {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, path, nil, parser.AllErrors)
	if err != nil {
		fmt.Println(err)
		return
	}

	d.FieldMap = map[string]field{}
	d.Filename = filepath.Base(path)

	ast.Inspect(file, func(node ast.Node) bool {
		switch v := node.(type) {
		case *ast.TypeSpec:
			if v.Name == nil {
				return false
			}

			if v.Name.Name == d.DaoName {
				b = true
				stType, ok := v.Type.(*ast.StructType)
				if !ok {
					return false
				}

				if stType.Fields == nil {
					return false
				}

				for _, f := range stType.Fields.List {
					//fmt.Println(reflect.ValueOf(f.Type).Elem().Type())
					if len(f.Names) == 0 {
						// @todo 匿名struct
						continue
					}

					// @todo 处理组合类型(gmodel.Time)，结构体类型(struct)，自定义类型(Status)，引用外部包自定义类型（g.Status）
					var ty string
					ist, ok := f.Type.(*ast.Ident)
					if ok {
						ty = ist.Name
					}

					tagMap := map[string]string{}
					if f.Tag != nil {
						var tagStr = strings.Trim(f.Tag.Value, "`")
						tag := reflect.StructTag(tagStr)
						if jTag, ok := tag.Lookup("json"); ok {
							tagMap["json"] = jTag
						}

						if xTag, ok := tag.Lookup("xorm"); ok {
							tagMap["xorm"] = xTag
						}
					}

					for _, name := range f.Names {
						d.FieldMap[name.Name] = field{
							Name:      name.Name,
							Type:      ty,
							SnakeName: snakeString(name.Name),
							tagMap:    tagMap,
						}
					}

				}
				return false
			}

		case *ast.File:
			d.DaoPackage = v.Name.Name
			d.DaoImportPath = fmt.Sprintf("%s/%s", d.Gener.ModuleImportPath, v.Name.Name)
		}

		return true
	})

	return
}

type daoGen struct {
	BaseGen
	DaoFile    string
	DaoPackage string
	db         string
	tableName  string
}

func NewDaoGen(conf *viper.Viper) *daoGen {
	g := new(daoGen)
	g.Init(conf)
	g.db = conf.GetString("db")
	g.tableName = conf.GetString("tableName")
	g.DaoFile = g.File
	g.DaoPackage = fmt.Sprintf("%sdao", g.ModuleName)

	return g
}

func (d *daoGen) Gen() (err error) {
	confPath := fmt.Sprintf("%s/app.toml", d.RootPath)
	if !fileExists(confPath) {
		return fmt.Errorf("%s 文件不存在", confPath)
	}

	cfg := viper.New()
	cfg.SetConfigFile(confPath)
	err = cfg.ReadInConfig()
	if err != nil {
		return err
	}

	// Generate struct string based on columnDataTypes
	// 根据 db 获取配置
	conf := cfg.GetStringMap(d.db)
	if conf == nil {
		log.Fatalln("Error, conf", d.db, "is not exists")
		return
	}

	var dbName, userName, passName, hostName string
	var portName int
	database, ok := conf["database"]
	if ok {
		dbName = database.(string)
	}

	user, ok := conf["user"]
	if ok {
		userName = user.(string)
	}

	pass, ok := conf["password"]
	if ok {
		passName = pass.(string)
	}

	host, ok := conf["host"]
	if ok {
		hostName = host.(string)
	}

	port, ok := conf["port"]
	if ok {
		portName, _ = strconv.Atoi(port.(string))
	}

	tableName := d.tableName
	if tableName == "" {
		fmt.Println("Error, please input table name")
		return
	}

	structName := camelString(tableName) + "Dao"

	columnDataTypes, columnsSorted, err := db2struct.GetColumnsFromMysqlTable(userName, passName, hostName, portName, dbName, tableName)

	if err != nil {
		fmt.Println("Error in selecting column data information from mysql information schema")
		return
	}

	pkgName := d.DaoPackage

	struc, err := db2struct.Generate(*columnDataTypes, columnsSorted, tableName, structName, pkgName, true, false, false)

	if err != nil {
		fmt.Println("Error in creating struct from json: " + err.Error())
		return
	}

	reg := regexp.MustCompile(`(_?U?ID\s+)int`)

	strucStr := reg.ReplaceAllString(string(struc), "${1}int64")

	r := strings.NewReplacer(
		"ID", "Id",
		"UID", "Uid",
		"URL", "Url",
		"time.Time", "string",
		`json:"ctime"`, `json:"ctime" xorm:"created"`,
		`json:"utime"`, `json:"utime" xorm:"updated"`,
		`json:"created_at"`, `json:"created_at" xorm:"created"`,
		`json:"updated_at"`, `json:"updated_at" xorm:"updated"`,
		`json:"id"`, `json:"id" xorm:"autoincr pk"`,
		"package "+pkgName, "", // @todo 新建不替换
	)

	tn := `
func (t *%s) TableName() string {
	return "%s"
}
`
	tnStr := fmt.Sprintf(tn, structName, tableName)

	daoContent := strings.TrimSpace(fmt.Sprintf("%s\n\n%s", r.Replace(strucStr), tnStr))

	fmt.Println(daoContent)

	writeInLineOrCreate(d.DaoFile, daoContent, getGoLine())

	gofmt(d.DaoFile)

	return nil
}

func genDaoCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "gen-dao",
		Short: "Gen dao code",
		Run: func(cmd *cobra.Command, args []string) {
			d := NewDaoGen(Conf)
			err := d.Gen()
			if err != nil {
				log.Fatalln(err)
			}
		},
	}

	cmd.Flags().StringP("tableName", "t", "", "表名称")
	cmd.Flags().StringP("db", "d", "db", "数据库组件配置")
	cmd.Flags().StringP("file", "f", getGoFile(), "file path,eg. ./ab_c.go")

	_ = cmd.MarkFlagRequired("tableName")

	return cmd

}
