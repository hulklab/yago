package promiselib

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestPromise(t *testing.T) {
	promise := New(3)

	for i := 0; i < 5; i++ {
		if i == 1 {
			promise.Add(func() error {
				// do some things , but occur err
				return errors.New("some errors occur")
			})
			continue
		}

		promise.Add(func() error {
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
