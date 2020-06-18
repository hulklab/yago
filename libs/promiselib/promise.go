package promiselib

import "context"

type promise struct {
	ctx      context.Context
	cancel   context.CancelFunc
	usedSize int        // nums of the goruntine
	error    chan error // log result for goruntine
}

func New(concurrencyNum int) *promise {
	w := new(promise)
	w.error = make(chan error, concurrencyNum)
	w.ctx, w.cancel = context.WithCancel(context.Background())

	return w
}

// add goruntine
func (w *promise) Add(f func() error) {
	w.usedSize ++

	go func() {
		select {
		case <-w.ctx.Done():
			return
		default:
			w.error <- f()
		}
	}()
}

// wait for result
func (w *promise) Wait() error {
	for i := 0; i < w.usedSize; i ++ {
		if err := <-w.error; err != nil {
			w.cancel()

			return err
		}
	}

	return nil
}