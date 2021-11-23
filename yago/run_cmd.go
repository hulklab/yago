package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
)

var (
	lastUpdateTime = time.Now().Unix()
	state          sync.Mutex
	cmd            *exec.Cmd
)

func readDir(dir string, files *[]string) error {
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return nil
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

		// fmt.Println(path, f.Name(), "----------------------------------------------")
		*files = append(*files, path)
		return nil
	})

	return err
}

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
					_ = watcher.Add(f)
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
		// go func() {
		// 	err := cmd.Run()
		// 	if err != nil && err.Error() != "signal: killed" {
		// 		log.Fatalln("[FATAL] start app error:", err)
		// 	}
		// }()
	}()
}
