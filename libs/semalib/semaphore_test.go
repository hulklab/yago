package semalib

import (
	"fmt"
	"testing"
	"time"
)

func TestSemaphore(t *testing.T) {
	sema := New(3)

	for i := 0; i < 10; i++ {
		fmt.Printf("sema.AvailablePermits : %d \n", sema.AvailablePermits())
		sema.Acquire() //数量不足，阻塞等待
		go func() {
			defer sema.Release()
			fmt.Println("sema")
			time.Sleep(time.Second)
		}()

	}

	sema.Wait()
	fmt.Printf("sema.AvailablePermits : %d \n", sema.AvailablePermits())

	fmt.Println("---------- TestSemaphore done ----------")
}

func TestTrySemaphore(t *testing.T) {
	sema := New(3)

	for i := 0; i < 10; i++ {
		fmt.Printf("sema.AvailablePermits : %d \n", sema.AvailablePermits())
		if sema.TryAcquire() { //不阻塞等待
			go func() {
				defer sema.Release()
				fmt.Println("sema")
				time.Sleep(time.Second)
			}()
		}
	}

	sema.Wait()
	fmt.Printf("sema.AvailablePermits : %d \n", sema.AvailablePermits())

	fmt.Println("---------- TestTrySemaphore done ----------")
}
