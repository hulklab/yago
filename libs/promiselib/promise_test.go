package promiselib

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestPromise(t *testing.T) {
	promise := New(3)

	for i := 0; i < 3; i++ {
		promise.Add(func() error {
			fmt.Println("do some things")
			// do some things
			time.Sleep(time.Second * 1)
			return nil
		})
	}

	promise.Add(func() error {
		// do some things , but occur err
		return errors.New("some errors occur")
	})

	for i := 0; i < 1000; i++ {
		promise.Add(func() error {
			fmt.Println("do some things")
			// do some things
			time.Sleep(time.Second * 1)
			return nil
		})
	}

	err := promise.Wait()
	if err != nil {
		fmt.Println("close all goruntine which is running")
	}

	fmt.Println("---------- TestPromise done ----------")
}
