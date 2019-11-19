package basethird

import (
	"context"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc/metadata"

	"github.com/hulklab/yago/coms/logger"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type UnaryClientInterceptor func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error
type StreamClientInterceptor func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, opts ...grpc.CallOption) error

// 封装 rpc 的基础类
type RpcThird struct {
	c *grpc.ClientConn
	sync.Mutex
	Address                       string
	Timeout                       int
	MaxRecvMsgsizeMb              int
	MaxSendMsgsizeMb              int
	SslOn                         bool
	CertFile                      string
	Hostname                      string
	logInfoOff                    bool
	beforeUnaryClientInterceptor  UnaryClientInterceptor
	afterUnaryClientInterceptor   UnaryClientInterceptor
	beforeStreamClientInterceptor StreamClientInterceptor
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

		dialOptions := []grpc.DialOption{
			grpc.WithUnaryInterceptor(a.unaryClientInterceptor),
			grpc.WithStreamInterceptor(a.streamClientInterceptor),
			grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(a.MaxRecvMsgsizeMb * 1024 * 1024)),
			grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(a.MaxSendMsgsizeMb * 1024 * 1024)),
		}

		if !a.SslOn {
			dialOptions = append(dialOptions, grpc.WithInsecure())

		} else {
			if a.CertFile == "" {
				log.Fatalln("server cert file is required when ssl on")
			}
			creds, err := credentials.NewClientTLSFromFile(a.CertFile, a.Hostname)
			if err != nil {
				log.Fatalf("failed to create TLS credentials %v", err)
			}

			dialOptions = append(dialOptions, grpc.WithTransportCredentials(creds))
		}

		a.c, err = grpc.Dial(
			a.Address,
			dialOptions...,
		)
	}

	return a.c, err
}

func (a *RpcThird) GetCtx() (context.Context, context.CancelFunc) {

	if a.Timeout == 0 {
		a.Timeout = 12
	}

	return context.WithTimeout(context.Background(), time.Duration(a.Timeout)*time.Second)
}

// 设置是否要关闭 info 日志
func (a *RpcThird) SetLogInfoFlag(on bool) {
	if on {
		a.logInfoOff = false
	} else {
		a.logInfoOff = true
	}
}

func (a *RpcThird) SetBeforeUnaryClientInterceptor(unary UnaryClientInterceptor) {
	a.beforeUnaryClientInterceptor = unary
}

func (a *RpcThird) SetAfterUnaryClientInterceptor(unary UnaryClientInterceptor) {
	a.afterUnaryClientInterceptor = unary
}

func (a *RpcThird) SetBeforeStreamClientInterceptor(stream StreamClientInterceptor) {
	a.beforeStreamClientInterceptor = stream
}

func (a *RpcThird) unaryClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	if a.beforeUnaryClientInterceptor != nil {
		err := a.beforeUnaryClientInterceptor(ctx, method, req, reply, cc, opts...)
		if err != nil {
			return err
		}
	}

	logInfo := logrus.Fields{
		"address":  a.Address,
		"timeout":  a.Timeout,
		"method":   method,
		"params":   req,
		"consume":  0,
		"category": "third.rpc",
	}

	md, ok := metadata.FromOutgoingContext(ctx)
	if ok {
		logInfo["metadata"] = md
	}
	//log.Printf("before invoker. method: %+v, request:%+v", method, req)
	begin := time.Now()

	err := invoker(ctx, method, req, reply, cc, opts...)

	end := time.Now()
	consume := end.Sub(begin).Nanoseconds() / 1e6
	logInfo["consume"] = consume

	if err != nil {
		logInfo["hint"] = err.Error()
		logger.Ins().WithFields(logInfo).Error()
	} else {
		// 默认是日志没关
		if !a.logInfoOff {
			logInfo["result"] = reply
		}

		logger.Ins().WithFields(logInfo).Info()
	}
	if err != nil {
		return err
	}

	if a.afterUnaryClientInterceptor != nil {
		return a.afterUnaryClientInterceptor(ctx, method, req, reply, cc, opts...)
	}
	return nil
}
func (a *RpcThird) streamClientInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	if a.beforeStreamClientInterceptor != nil {
		err := a.beforeStreamClientInterceptor(ctx, desc, cc, method, opts...)
		if err != nil {
			return nil, err
		}
	}

	logInfo := logrus.Fields{
		"address":  a.Address,
		"timeout":  a.Timeout,
		"method":   method,
		"consume":  0,
		"category": "third.rpc",
	}

	md, ok := metadata.FromOutgoingContext(ctx)
	if ok {
		logInfo["metadata"] = md
	}

	begin := time.Now()

	clientStream, err := streamer(ctx, desc, cc, method, opts...)

	end := time.Now()
	consume := end.Sub(begin).Nanoseconds() / 1e6
	logInfo["consume"] = consume

	if err != nil {
		logInfo["hint"] = err.Error()
		logger.Ins().WithFields(logInfo).Error()
	} else {
		logger.Ins().WithFields(logInfo).Info()
	}

	// 此时只是打开了 stream 通道，还未开始传输数据，只有 stream 得到一个 EOF 的错误时才算传输完成
	return clientStream, err
}
