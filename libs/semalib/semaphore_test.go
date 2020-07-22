package semalib

import (
	"errors"
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

	_ = sema.Wait()
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

	_ = sema.Wait()
	fmt.Printf("sema.AvailablePermits : %d \n", sema.AvailablePermits())

	fmt.Println("---------- TestTrySemaphore done ----------")
}

func TestErrRetrunSemaphore(t *testing.T) {
	sema := New(3)

	for i := 0; i < 3; i++ {
		sema.Add(func() error {
			// do some things
			fmt.Println("do some things")
			time.Sleep(time.Second)

			return nil
		})
	}

	sema.Add(func() error {
		// do some things
		fmt.Println("some error occur")
		time.Sleep(time.Millisecond * 200)
		return errors.New("occur error")
	})

	for i := 0; i < 3000; i++ {
		sema.Add(func() error {
			// do some things
			fmt.Println("do some things again")
			time.Sleep(time.Second)

			return nil
		})
	}

	err := sema.Wait()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("---------- TestErrRetrunSemaphore done ----------")
}
