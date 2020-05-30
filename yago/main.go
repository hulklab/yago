package main

import (
	"errors"
	"fmt"
	"go/build"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
)

func GenDir(srcPath string, destPath, app string) error {
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
			//path := strings.Replace(path, "\\", "/", -1)
			destNewPath := strings.Replace(path, srcPath, destPath, -1)
			if err := GenFile(path, destNewPath, app); err != nil {
				log.Println(fmt.Sprintf("create file %s error:", destNewPath), err.Error())
				return err
			}
		}
		return nil
	})
	return err
}

func GenFile(src, dest, app string) (err error) {
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
			b, _ := pathExists(destSplitPath)
			if b == false {
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

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func getCurDir() string {
	dir, _ := filepath.Abs(filepath.Dir("."))
	return filepath.Clean(dir)
	//return strings.Replace(dir, "\\", "/", -1)
}

func getGoPath() string {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}
	return gopath
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init app",
	Long:  `Init a app named by input`,
	Run: func(cmd *cobra.Command, args []string) {
		useMod, _ := cmd.Flags().GetBool("mod")
		app, _ := cmd.Flags().GetString("app")

		log.Println("create app", app)
		if err := os.MkdirAll(app, 0755); err != nil {
			log.Println("create app dir error:", err.Error())
		}
		var src string
		if useMod {
			src = filepath.Join(getGoPath(), "pkg", "mod", "github.com", "hulklab", "yago@"+Version, "example")
		} else {
			src = filepath.Join(getGoPath(), "src", "github.com", "hulklab", "yago", "example")
		}
		dest := app

		if err := GenDir(src, dest, app); err != nil {
			log.Println("create app error:", err.Error())
		}
	},
}

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "New module",
	Long:  `New a module named by input`,
	Run: func(cmd *cobra.Command, args []string) {
		if exist, _ := pathExists("go.mod"); !exist {
			log.Println("current directory is not a go mod root path")
			return
		}

		curDir := strings.Split(getCurDir(), string(filepath.Separator))
		app := curDir[len(curDir)-1]

		module, _ := cmd.Flags().GetString("module")
		advance, _ := cmd.Flags().GetBool("advance")

		log.Println("create module", module)
		dirs := []string{"cmd", "dao", "http", "model", "rpc", "task"}
		if advance {
			dirs = append(dirs, "service")
		}

		for _, d := range dirs {
			//dirPath := fmt.Sprintf("app/modules/%s/%s%s", module, module, d)
			dirPath := filepath.Join("app", "modules", module, module+d)
			if err := os.MkdirAll(dirPath, 0755); err != nil {
				log.Println(fmt.Sprintf("create module dir %s error:", dirPath), err.Error())
				return
			}
			//filePath := fmt.Sprintf("%s/%s.go", dirPath, module)
			filePath := filepath.Join(dirPath, module+".go")
			fileBody := fmt.Sprintf("package %s%s", module, d)
			if err := ioutil.WriteFile(filePath, []byte(fileBody), 0644); err != nil {
				log.Println(fmt.Sprintf("create module file %s error:", filePath), err.Error())
				return
			}
		}

		//routePath := "app/route/route.go"
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

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version",
	Long:  `Print version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("yago version", Version)
	},
}

func camelString(s string) (string, error) {
	s = strings.Trim(s, "_")

	r, _ := regexp.Compile("^[a-zA-Z][a-zA-Z_0-9]*[a-zA-Z0-9]$")
	if !r.MatchString(s) {
		return "", errors.New("only support a-z A-Z _")
	}

	ss := strings.Split(s, "_")
	for k, v := range ss {
		ss[k] = strings.Title(v)
	}
	return strings.Join(ss, ""), nil
}

func genFileByTemplate(filename, template string) {
	pkgName := strings.ToLower(filepath.Base(getCurDir()))
	name, err := camelString(filename)
	if err != nil {
		log.Fatalln(err)
	}

	if len(name) < 2 {
		log.Fatalln("the filename length is greater than at least 2")
	}

	lname := fmt.Sprintf("%s%s", strings.ToLower(name[0:1]), name[1:])

	r := strings.NewReplacer("{{PACKAGE}}", pkgName, "{{NAME}}", name, "{{LNAME}}", lname, "{{ONAME}}", filename)
	content := r.Replace(template)

	filePath := filepath.Join(getCurDir(), filename+".go")
	if err := ioutil.WriteFile(filePath, []byte(content), 0644); err != nil {
		log.Println(fmt.Sprintf("create file %s error:", filePath), err.Error())
		return
	} else {
		log.Println(fmt.Sprintf("create file %s succ", filePath))
	}
}

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "auto generate file",
	Long:  `auto generate http, rpc, task, cmd, service, model, api file`,
	Run: func(cmd *cobra.Command, args []string) {
		p, err := cmd.Flags().GetString("http")
		if err != nil {
			log.Fatalln(err)
		}

		if len(p) > 0 {
			genFileByTemplate(p, HttpTemplate)
			return
		}

		r, err := cmd.Flags().GetString("rpc")
		if err != nil {
			log.Fatalln(err)
		}

		if len(r) > 0 {
			genFileByTemplate(r, RpcTemplate)
			return
		}

		c, err := cmd.Flags().GetString("cmd")
		if err != nil {
			log.Fatalln(err)
		}

		if len(c) > 0 {
			genFileByTemplate(c, CmdTemplate)
			return
		}

		t, err := cmd.Flags().GetString("task")
		if err != nil {
			log.Fatalln(err)
		}

		if len(t) > 0 {
			genFileByTemplate(t, TaskTemplate)
			return
		}

		s, err := cmd.Flags().GetString("service")
		if err != nil {
			log.Fatalln(err)
		}

		if len(s) > 0 {
			genFileByTemplate(s, ServiceTemplate)
			return
		}

		m, err := cmd.Flags().GetString("model")
		if err != nil {
			log.Fatalln(err)
		}

		if len(m) > 0 {
			genFileByTemplate(m, ModelTemplate)
			return
		}

		a, err := cmd.Flags().GetString("api")
		if err != nil {
			log.Fatalln(err)
		}

		if len(a) > 0 {
			genFileByTemplate(a, ApiTemplate)
			return
		}
	},
}

var (
	lastUpdateTime = time.Now().Unix()
	state          sync.Mutex
	cmd            *exec.Cmd
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Hot run app",
	Long:  "Hot build and run app",
	Run: func(cmd *cobra.Command, args []string) {
		pwd, _ := os.Getwd()
		appName := filepath.Base(pwd)

		files := make([]string, 0)

		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatalln(err)
		}

		defer watcher.Close()

		go func() {
			for {
				select {
				case event := <-watcher.Events:
					if !strings.HasSuffix(event.Name, ".go") {
						continue
					}

					now := time.Now().Unix()
					// 3 秒不重复编译
					if now-lastUpdateTime < 3 {
						continue
					}

					lastUpdateTime = now

					autoBuildApp(appName)

					// change file
					if event.Op&fsnotify.Write == fsnotify.Write {
						log.Println("[INFO] modified file: ", event.Name, " - ", event.String())
					}

				case err := <-watcher.Errors:
					log.Fatalln("[FATAL] watch error:", err)
				}

			}
		}()

		// 定时刷新监听的文件
		go func() {
			for {
				err := readDir(pwd, &files)
				if err != nil {
					log.Fatalln("[FATAL] read dir err", err)
				}

				for _, f := range files {
					watcher.Add(f)

				}

				time.Sleep(30 * time.Second)
			}
		}()

		// 先启动一次
		autoBuildApp(appName)

		select {}
	},
}

func autoBuildApp(appName string) {
	state.Lock()
	defer state.Unlock()

	log.Println("[INFO] rebuild app start ...")
	if runtime.GOOS == "windows" {
		appName += ".exe"
	}

	bcmd := exec.Command("go", "build", "-o", appName)
	bcmd.Stderr = os.Stderr
	bcmd.Stdout = os.Stdout
	err := bcmd.Run()
	if err != nil {
		log.Println("[ERROR] rebuild app error: ", err)
		return
	}

	restartApp(appName)

	log.Println("[INFO] rebuild app success.")
}

func restartApp(appName string) {
	log.Println("[INFO] restart app ", appName, "...")

	defer func() {
		if e := recover(); e != nil {
			log.Println("[ERROR] restart app error: ", e)
		}
	}()

	// 杀掉原先的 app
	if cmd != nil && cmd.Process != nil {
		err := cmd.Process.Kill()
		if err != nil {
			log.Fatalln("[FATAL] stop app error: ", err)
			return
		}
	}

	// 重启新的 app
	go func() {
		cmd = exec.Command("./" + appName)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		go cmd.Run()
		//go func() {
		//	err := cmd.Run()
		//	if err != nil && err.Error() != "signal: killed" {
		//		log.Fatalln("[FATAL] start app error:", err)
		//	}
		//}()
	}()
}

func readDir(dir string, files *[]string) error {
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
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
		*files = append(*files, path)
		return nil
	})

	return err
}

var rootCmd = &cobra.Command{}

func main() {

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(genCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Println("cmd run error:", err.Error())
		os.Exit(1)
	}
}

func init() {
	// init cmd
	initCmd.Flags().BoolP("mod", "", true, "use go mod ? only for dev user")
	// init cmd
	initCmd.Flags().StringP("app", "a", "", "app name")
	_ = initCmd.MarkFlagRequired("app")

	// module cmd
	newCmd.Flags().StringP("module", "m", "", "module name")
	newCmd.Flags().BoolP("advance", "a", false, "gen advance module which include service")
	_ = newCmd.MarkFlagRequired("module")

	// gen cmd
	genCmd.Flags().StringP("http", "p", "", "http file name")
	genCmd.Flags().StringP("rpc", "r", "", "rpc file name")
	genCmd.Flags().StringP("cmd", "c", "", "cmd file name")
	genCmd.Flags().StringP("task", "t", "", "task file name")
	genCmd.Flags().StringP("service", "s", "", "service file name")
	genCmd.Flags().StringP("model", "m", "", "model file name")
	genCmd.Flags().StringP("api", "a", "", "api file name, with no _api postfix")
}
