// +build !windows

package yago

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func (a *App) startSignal() {
	pid := os.Getpid()
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP, syscall.SIGUSR2)

	for {
		s := <-signals
		log.Println("recv", s)

		switch s {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			log.Println("Graceful Shutdown...")

			a.Close()

			log.Println("Process", pid, "Exit OK")
			os.Exit(0)
		case syscall.SIGHUP:
			log.Println("Graceful restart...")

			a.Close()

			err := RestartApp()

			if err != nil {
				log.Fatalln("Process", pid, "Restart err")
			}

			log.Println("Process", pid, "Restart ok")
			os.Exit(0)
		case syscall.SIGUSR2:
			log.Println("Graceful reload...")
			err := ReloadConfig()
			if err != nil {
				log.Printf("reload config fail, err: %s", err)
			} else {
				log.Println("reload config ok...")
			}
		}
	}

}

func RestartApp() error {

	execSpec := &syscall.ProcAttr{
		Env:   os.Environ(),
		Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()},
	}

	newPid, err := syscall.ForkExec(os.Args[0], os.Args, execSpec)
	if err != nil {
		log.Println("restart app, ", "fork sub proc err:", err)
		return err
	}

	log.Println("fork sub proc pid:", newPid)

	return nil
}
