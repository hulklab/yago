package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func newCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "new",
		Short: "New module",
		Long:  `New a module named by input`,
		Run: func(cmd *cobra.Command, args []string) {
			if !fileExists("go.mod") {
				log.Println("current directory is not a go mod root path")
				return
			}

			curDir := strings.Split(pwd(), string(filepath.Separator))
			app := curDir[len(curDir)-1]

			module, _ := cmd.Flags().GetString("module")

			log.Println("create module", module)
			dirs := []string{"cmd", "dto", "dao", "http", "model", "service", "rpc", "task"}

			for _, d := range dirs {
				// dirPath := fmt.Sprintf("app/modules/%s/%s%s", module, module, d)
				dirPath := filepath.Join("app", "modules", module, module+d)
				if err := os.MkdirAll(dirPath, 0755); err != nil {
					log.Println(fmt.Sprintf("create module dir %s error:", dirPath), err.Error())
					return
				}
				// filePath := fmt.Sprintf("%s/%s.go", dirPath, module)
				filePath := filepath.Join(dirPath, module+".go")
				fileBody := fmt.Sprintf("package %s%s", module, d)
				if err := ioutil.WriteFile(filePath, []byte(fileBody), 0644); err != nil {
					log.Println(fmt.Sprintf("create module file %s error:", filePath), err.Error())
					return
				}
			}

			// routePath := "app/route/route.go"
			routePath := filepath.Join("app", "route", "route.go")
			routes := []string{"cmd", "http", "rpc", "task"}
			for _, d := range routes {
				// routePath := fmt.Sprintf("app/routes/%sroute/%s.go", d, d)
				var routeBody []byte
				var err error
				if routeBody, err = ioutil.ReadFile(routePath); err != nil {
					log.Println(fmt.Sprintf("read route file %s error:", routePath), err.Error())
					return
				}
				newRoute := fmt.Sprintf("\t_ \"%s/app/modules/%s/%s%s\"\n)", app, module, module, d)
				contentStr := strings.ReplaceAll(string(routeBody), ")", newRoute)
				if err = ioutil.WriteFile(routePath, []byte(contentStr), 0644); err != nil {
					log.Println(fmt.Sprintf("write route file %s error:", routePath), err.Error())
					return
				}
				cmd := exec.Command("gofmt", "-w", routePath)
				if err := cmd.Run(); err != nil {
					log.Println(fmt.Sprintf("gofmt route file %s error:", routePath), err.Error())
					return
				}
			}
		},
	}
	cmd.Flags().StringP("module", "m", "", "module name")
	_ = cmd.MarkFlagRequired("module")
	return cmd
}
