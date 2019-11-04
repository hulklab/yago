package basethird

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/hulklab/yago/coms/logger"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"sync"
	"time"
)

// 封装 rpc 的基础类
type RpcThird struct {
	c *grpc.ClientConn
	sync.Mutex
	Address          string
	Timeout          int
	MaxRecvMsgsizeMb int
	MaxSendMsgsizeMb int
	SslOn            bool
	CertFile         string
	Hostname         string
	logInfoOff       bool
}

func (a *RpcThird) GetConn() (*grpc.ClientConn, error) {
	a.Lock()
	defer a.Unlock()

	var err error
	if a.c == nil {
		if a.MaxRecvMsgsizeMb == 0 {
			a.MaxRecvMsgsizeMb = 4
		}
		if a.MaxSendMsgsizeMb == 0 {
			a.MaxSendMsgsizeMb = 4
		}

		if !a.SslOn {
			a.c, err = grpc.Dial(
				a.Address,
				grpc.WithInsecure(),
				grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(a.MaxRecvMsgsizeMb*1024*1024)),
				grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(a.MaxSendMsgsizeMb*1024*1024)),
			)

		} else {
			if a.CertFile == "" {
				log.Fatalln("server cert file is required when ssl on")
			}
			creds, err := credentials.NewClientTLSFromFile(a.CertFile, a.Hostname)
			if err != nil {
				log.Fatalf("failed to create TLS credentials %v", err)
			}

			a.c, err = grpc.Dial(
				a.Address,
				grpc.WithTransportCredentials(creds),
				grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(a.MaxRecvMsgsizeMb*1024*1024)),
				grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(a.MaxSendMsgsizeMb*1024*1024)),
			)
		}
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
		// 默认是日志没关
		if !a.logInfoOff {
			logger.Ins().WithFields(logInfo).Info()
		}
	}

	return rep, err
}

// 设置是否要关闭 info 日志
func (a *RpcThird) SetLogInfoFlag(on bool) {
	if on {
		a.logInfoOff = false
	} else {
		a.logInfoOff = true
	}
}
