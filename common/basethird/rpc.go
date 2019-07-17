package basethird

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/hulklab/yago/libs/logger"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"sync"
	"time"
)

// 封装 rpc 的基础类
type RpcThird struct {
	c *grpc.ClientConn
	sync.Mutex
	Address string
	Timeout int
}

func (a *RpcThird) GetConn() (*grpc.ClientConn, error) {
	a.Lock()
	defer a.Unlock()

	var err error
	if a.c == nil {
		a.c, err = grpc.Dial(a.Address, grpc.WithInsecure())
	}

	return a.c, err
}

func (a *RpcThird) GetCtx() (context.Context, context.CancelFunc) {

	if a.Timeout == 0 {
		a.Timeout = 12
	}

	return context.WithTimeout(context.Background(), time.Duration(a.Timeout)*time.Second)
}

func (a *RpcThird) Call(f func(conn *grpc.ClientConn, ctx context.Context) (proto.Message, error), params interface{}) (proto.Message, error) {
	logInfo := logrus.Fields{
		"address":     a.Address,
		"timeout":     a.Timeout,
		"params":      params,
		"consume(ms)": 0,
		"error":       "",
		"result":      "",
		"category":    "third.rpc",
	}

	conn, err := a.GetConn()
	if err != nil {
		logInfo["error"] = err.Error()
		logger.Ins().WithFields(logInfo).Error()
		return nil, err
	}

	ctx, cancel := a.GetCtx()

	defer cancel()

	begin := time.Now()

	rep, err := f(conn, ctx)

	end := time.Now()
	consume := end.Sub(begin).Nanoseconds() / 1e6

	logInfo["consume"] = consume

	if rep != nil {
		logInfo["result"] = rep.String()
	}

	if err != nil {
		logInfo["error"] = err.Error()
		logger.Ins().WithFields(logInfo).Error()
	} else {
		logger.Ins().WithFields(logInfo).Info()
	}

	return rep, err
}
