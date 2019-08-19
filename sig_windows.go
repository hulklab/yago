// +build windows

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
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	for {
		s := <-signals
		log.Println("recv", s)

		switch s {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			log.Println("Graceful Shutdown...")

			a.Close()

			log.Println("Process", pid, "Exit OK")
			os.Exit(0)
		}
	}
}
