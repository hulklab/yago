package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func initAppDir(srcPath string, destPath, app string) error {
	if srcInfo, err := os.Stat(srcPath); err != nil {
		return err
	} else {
		if !srcInfo.IsDir() {
			return errors.New("src path is not a correct directory！")
		}
	}
	if destInfo, err := os.Stat(destPath); err != nil {
		return err
	} else {
		if !destInfo.IsDir() {
			return errors.New("dest path is not a correct directory！")
		}
	}

	err := filepath.Walk(srcPath, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if !f.IsDir() {
			// path := strings.Replace(path, "\\", "/", -1)
			destNewPath := strings.Replace(path, srcPath, destPath, -1)
			if err := initAppFile(path, destNewPath, app); err != nil {
				log.Println(fmt.Sprintf("create file %s error:", destNewPath), err.Error())
				return err
			}
		}
		return nil
	})
	return err
}

func initAppFile(src, dest, app string) (err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destSplitPathDirs := strings.Split(dest, string(filepath.Separator))

	destSplitPath := ""
	for index, dir := range destSplitPathDirs {
		if index < len(destSplitPathDirs)-1 {
			destSplitPath = filepath.Join(destSplitPath, dir)
			b := fileExists(destSplitPath)
			if !b {
				err := os.Mkdir(destSplitPath, os.ModePerm)
				if err != nil {
					return err
				}
			}
		}
	}
	dstFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	srcFileInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	content := make([]byte, srcFileInfo.Size())
	if _, err := srcFile.Read(content); err != nil {
		return err
	}

	contentStr := strings.ReplaceAll(string(content), "github.com/hulklab/yago/example", app)

	if _, err := dstFile.WriteString(contentStr); err != nil {
		return err
	}
	return nil
}

func initCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Init app",
		Long:  `Init a app named by input`,
		Run: func(cmd *cobra.Command, args []string) {
			app, _ := cmd.Flags().GetString("app")
			src, _ := cmd.Flags().GetString("src")

			log.Println("create app", app)
			if err := os.MkdirAll(app, 0755); err != nil {
				log.Println("create app dir error:", err.Error())
			}

			if len(src) == 0 {
				src = filepath.Join(getGoPath(), "pkg", "mod", "github.com", "hulklab", "yago@"+Version, "example")
			}

			dest := app

			if err := initAppDir(src, dest, app); err != nil {
				log.Println("create app error:", err.Error())
			}
		},
	}
	cmd.Flags().StringP("app", "a", "", "app name")
	cmd.Flags().StringP("src", "s", "", "[optional]src yago example dir")
	_ = cmd.MarkFlagRequired("app")

	return cmd
}

func init() {
	// init cmd
}
