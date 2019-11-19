package yago

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

type App struct {
	// 是否开启debug模式
	DebugMode bool
	// http web 引擎
	httpEngine *gin.Engine
	// http run mode
	HttpRunMode string
	// 开启http服务
	HttpEnable bool
	// http路由配置
	HttpRouters []*HttpRouter
	// http close chan
	httpCloseChan     chan int
	httpCloseDoneChan chan int
	// https 证书配置
	HttpSslOn    bool
	HttpCertFile string
	HttpKeyFile  string
	// http html 模版配置
	HttpViewRender bool
	HttpViewPath   string
	HttpStaticPath string
	// http cors 跨域配置
	HttpCorsAllowAllOrigins  bool
	HttpCorsAllowOrigins     []string
	HttpCorsAllowMethods     []string
	HttpCorsAllowHeaders     []string
	HttpCorsExposeHeaders    []string
	HttpCorsAllowCredentials bool
	HttpCorsMaxAge           time.Duration
	// http gzip 压缩
	HttpGzipOn    bool
	HttpGzipLevel int
	// http pprof
	HttpPprof bool

	// 开启task服务
	TaskEnable bool
	// task路由配置
	TaskRouters []*TaskRouter
	// task close chan
	taskCloseChan     chan int
	taskCloseDoneChan chan int

	// rpc
	RpcEnable bool
	rpcEngine *grpc.Server

	// rpc close chan
	rpcCloseChan     chan int
	rpcCloseDoneChan chan int
}

var (
	errHttpRouteEmpty = errors.New("http router is empty")
	errTaskRouteEmpty = errors.New("task router is empty")
)

func NewApp() *App {
	// new app
	app := new(App)

	app.DebugMode = Config.GetBool("app.debug")

	// init http
	app.HttpEnable = Config.GetBool("app.http_enable")
	if app.HttpEnable {
		if app.DebugMode == true {
			app.HttpRunMode = gin.DebugMode
		} else {
			app.HttpRunMode = gin.ReleaseMode
		}
		gin.SetMode(app.HttpRunMode)
		app.httpEngine = gin.New()
		// use logger
		app.httpEngine.Use(gin.Logger())
		app.httpCloseChan = make(chan int, 1)
		app.httpCloseDoneChan = make(chan int, 1)

		app.HttpViewRender = Config.GetBool("app.http_view_render")
		if app.HttpViewRender {
			app.HttpViewPath = Config.GetString("app.http_view_path")
			if app.HttpViewPath != "" {
				app.httpEngine.LoadHTMLGlob(app.HttpViewPath)
			}
			app.HttpStaticPath = Config.GetString("app.http_static_path")
			if app.HttpStaticPath != "" {
				app.httpEngine.Static("/static", app.HttpStaticPath)
			}
		}

		app.HttpSslOn = Config.GetBool("app.http_ssl_on")
		if app.HttpSslOn {
			app.HttpCertFile = Config.GetString("app.http_cert_file")
			app.HttpKeyFile = Config.GetString("app.http_key_file")
		}

		app.HttpCorsAllowAllOrigins = true
		app.HttpCorsAllowOrigins = []string{}
		app.HttpCorsAllowHeaders = []string{"Origin", "Content-Length", "Content-Type"}
		app.HttpCorsExposeHeaders = []string{}
		app.HttpCorsAllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"}
		app.HttpCorsAllowCredentials = false
		app.HttpCorsMaxAge = time.Duration(12) * time.Hour

		if Config.IsSet("app.http_cors_allow_all_origins") {
			app.HttpCorsAllowAllOrigins = Config.GetBool("app.http_cors_allow_all_origins")
		}
		if Config.IsSet("app.http_cors_allow_origins") {
			app.HttpCorsAllowOrigins = Config.GetStringSlice("app.http_cors_allow_origins")
		}
		if Config.IsSet("app.http_cors_allow_headers") {
			app.HttpCorsAllowHeaders = Config.GetStringSlice("app.http_cors_allow_headers")
		}
		if Config.IsSet("app.http_cors_expose_headers") {
			app.HttpCorsExposeHeaders = Config.GetStringSlice("app.http_cors_expose_headers")
		}
		if Config.IsSet("app.http_cors_allow_methods") {
			app.HttpCorsAllowMethods = Config.GetStringSlice("app.http_cors_allow_methods")
		}
		if Config.IsSet("app.http_cors_allow_credentials") {
			app.HttpCorsAllowCredentials = Config.GetBool("app.http_cors_allow_credentials")
		}
		if Config.IsSet("app.http_cors_max_age") {
			app.HttpCorsMaxAge = Config.GetDuration("app.http_cors_max_age")
		}

		if Config.IsSet("app.http_gzip_on") {
			app.HttpGzipOn = Config.GetBool("app.http_gzip_on")
		} else {
			app.HttpGzipOn = true
		}
		if app.HttpGzipOn {
			switch Config.GetInt("app.http_gzip_level") {
			case 1:
				app.HttpGzipLevel = gzip.DefaultCompression
			case 2:
				app.HttpGzipLevel = gzip.BestSpeed
			case 3:
				app.HttpGzipLevel = gzip.BestCompression
			default:
				app.HttpGzipLevel = gzip.DefaultCompression
			}
		}

		app.HttpPprof = Config.GetBool("app.http_pprof_on")
	}

	// init rpc
	app.RpcEnable = Config.GetBool("app.rpc_enable")
	if app.RpcEnable {
		app.rpcCloseChan = make(chan int, 1)
		app.rpcCloseDoneChan = make(chan int, 1)
	}

	// init task
	app.TaskEnable = Config.GetBool("app.task_enable")
	if app.TaskEnable {
		app.taskCloseChan = make(chan int, 1)
		app.taskCloseDoneChan = make(chan int, 1)
	}

	return app
}

// 此 init 最先执行，配置文件此处初始化
func init() {
	initConfig()

	log.SetFlags(log.LstdFlags)

	initGrpcServer()
}

func (a *App) Run() {
	if a.TaskEnable {
		// 开启 task
		go a.runTask()
	}

	if a.RpcEnable {
		// 开启 rpc 服务
		go a.runRpc()
	}

	if a.HttpEnable {
		// 开启 http
		go a.runHttp()
	}

	// 生成 pid
	a.genPid()

	a.startSignal()
}

func (a *App) genPid() {
	pidfile, ok := getPidFile()
	if !ok {
		return
	}

	pf, err := os.Create(pidfile)

	defer pf.Close()

	if err != nil {
		log.Fatalf("pidfile check err:%v\n", err)
		return

	} else {
		newPid := os.Getpid()
		_, err := pf.Write([]byte(fmt.Sprintf("%d", newPid)))
		if err != nil {
			log.Fatalf("write pid err:%v\n", err)
			return
		}

		log.Println("app is running with pid:", newPid)
		return
	}
	return
}

func (a *App) loadHttpRouter() error {
	if len(HttpRouterMap) == 0 {
		return errHttpRouteEmpty
	}

	// cors
	if a.HttpCorsAllowAllOrigins == true || len(a.HttpCorsAllowOrigins) != 0 {
		a.httpEngine.Use(cors.New(cors.Config{
			AllowAllOrigins:        a.HttpCorsAllowAllOrigins,
			AllowOrigins:           a.HttpCorsAllowOrigins,
			AllowMethods:           a.HttpCorsAllowMethods,
			AllowHeaders:           a.HttpCorsAllowHeaders,
			ExposeHeaders:          a.HttpCorsExposeHeaders,
			AllowCredentials:       a.HttpCorsAllowCredentials,
			MaxAge:                 a.HttpCorsMaxAge,
			AllowWebSockets:        true,
			AllowBrowserExtensions: true,
			AllowWildcard:          true,
		}))

	}

	// gzip
	if a.HttpGzipOn {
		a.httpEngine.Use(gzip.Gzip(a.HttpGzipLevel))
	}
	// pprof
	if a.HttpPprof {
		pprof.Register(a.httpEngine)
	}

	// params
	a.httpEngine.Use(func(c *gin.Context) {
		req := c.Request

		query := req.URL.Query()

		for k, v := range query {
			c.Set(k, v[0])
		}

		switch c.ContentType() {
		case "application/x-www-form-urlencoded":
			err := req.ParseForm()
			if err != nil {
				log.Println("parse form", err.Error())
				return
			}
			for k, v := range req.PostForm {
				c.Set(k, v[0])
			}
		case "multipart/form-data":
			err := req.ParseMultipartForm(a.httpEngine.MaxMultipartMemory)
			if err != nil {
				log.Println("parse multi form", err.Error())
				return
			} else if req.MultipartForm != nil {
				for k, v := range req.MultipartForm.Value {
					c.Set(k, v[0])
				}
			}
		}
	})

	// no route handler
	if httpNoRouterHandler != nil {
		a.httpEngine.NoRoute(func(c *gin.Context) {
			ctx := NewCtx(c)
			httpNoRouterHandler(ctx)
		})
	}

	for _, r := range HttpRouterMap {
		method := strings.ToUpper(r.Method)
		action := r.Action
		controller := r.h
		handler := func(c *gin.Context) {
			ctx := NewCtx(c)
			if e := controller.BeforeAction(ctx); e.HasErr() {
				ctx.SetError(e)
			} else {
				if err := ctx.Validate(); err != nil {
					ctx.SetError(ErrParam, err.Error())
				} else {
					action(ctx)
				}
			}

			controller.AfterAction(ctx)
		}

		log.Println("[HTTP]", r.Url, runtime.FuncForPC(reflect.ValueOf(action).Pointer()).Name())
		switch method {
		case http.MethodGet:
			a.httpEngine.GET(r.Url, handler)
		case http.MethodPost:
			a.httpEngine.POST(r.Url, handler)
		case http.MethodDelete:
			a.httpEngine.DELETE(r.Url, handler)
		case http.MethodPut:
			a.httpEngine.PUT(r.Url, handler)
		case http.MethodOptions:
			a.httpEngine.OPTIONS(r.Url, handler)
		case http.MethodPatch:
			a.httpEngine.PATCH(r.Url, handler)
		case http.MethodHead:
			a.httpEngine.HEAD(r.Url, handler)
		default:
			a.httpEngine.Any(r.Url, handler)
		}
	}
	return nil
}

func (a *App) runHttp() {
	// load router
	if err := a.loadHttpRouter(); err != nil {
		a.httpCloseDoneChan <- 1
		return
	}

	// listen and serve
	srv := &http.Server{
		Addr:    Config.GetString("app.http_addr"),
		Handler: a.httpEngine,
	}

	if a.HttpSslOn {
		go func() {
			// service connections
			if err := srv.ListenAndServeTLS(a.HttpCertFile, a.HttpKeyFile); err != nil && err != http.ErrServerClosed {
				log.Fatalf("listen: %s\n", err)
			}
		}()

	} else {
		go func() {
			// service connections
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("listen: %s\n", err)
			}
		}()
	}

	select {
	case <-a.httpCloseChan:
		ctx, cancel := context.WithTimeout(
			context.Background(),
			time.Duration(Config.GetInt64("app.http_stop_time_wait"))*time.Second,
		)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			// log
		}
		a.httpCloseDoneChan <- 1
	}
}

func (a *App) loadTaskRouter() error {
	if len(TaskRouterList) == 0 {
		return errTaskRouteEmpty
	}

	return nil
}

func (a *App) runTask() {
	if !Config.GetBool("app.task_enable") {
		return
	}

	if err := a.loadTaskRouter(); err != nil {
		a.taskCloseDoneChan <- 1
		return
	}

	c := cron.New()
	wg := sync.WaitGroup{}
	for _, router := range TaskRouterList {
		action := router.Action
		if router.Spec == "@loop" {
			go func() {
				wg.Add(1)
				action()
				log.Println("[TASK]", "stop", runtime.FuncForPC(reflect.ValueOf(action).Pointer()).Name())
				wg.Done()
			}()
		} else {
			err := c.AddFunc(router.Spec, func() {
				wg.Add(1)
				action()
				wg.Done()
			})
			if err != nil {
				continue
			}
		}
		log.Println("[TASK]", router.Spec, runtime.FuncForPC(reflect.ValueOf(router.Action).Pointer()).Name())
	}

	c.Start()

	select {
	case <-a.taskCloseChan:
		go func() {
			c.Stop()
			wg.Wait()
			a.taskCloseDoneChan <- 1
		}()
		time.Sleep(time.Duration(Config.GetInt64("app.task_stop_time_wait")) * time.Second)
		a.taskCloseDoneChan <- 1
	}
}

var RpcServer *grpc.Server

func initGrpcServer() {
	isSslOn := Config.GetBool("app.rpc_ssl_on")
	if isSslOn {
		certFile := Config.GetString("app.rpc_cert_file")
		keyFile := Config.GetString("app.rpc_key_file")
		if certFile == "" || keyFile == "" {
			log.Fatalln("rpc ssl cert file or key file is required when rpc ssl on")
		}
		creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
		if err != nil {
			log.Fatalf("Failed to generate credentials %v", err)
		}

		// 实例化 grpc Server, 并开启 TSL 认证
		RpcServer = grpc.NewServer(grpc.Creds(creds))

	} else {
		RpcServer = grpc.NewServer()
	}
}

// rpc
func (a *App) runRpc() {
	rpcAddr := Config.GetString("app.rpc_addr")
	lis, err := net.Listen("tcp", rpcAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	a.rpcEngine = RpcServer

	for k, v := range a.rpcEngine.GetServiceInfo() {
		for _, method := range v.Methods {
			log.Println("[GRPC]", k, method.Name)
		}
	}

	// open rpc reflection, then you can use gpc_cli
	rpcReflectOn := Config.GetBool("app.rpc_reflect_on")
	if rpcReflectOn {
		reflection.Register(RpcServer)
	}

	go func() {
		if err := a.rpcEngine.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	select {
	case <-a.rpcCloseChan:
		a.rpcEngine.GracefulStop()
		a.rpcCloseDoneChan <- 1
	}
}

func (a *App) Close() {
	if a.TaskEnable {
		close(TaskCloseChan)
		a.taskCloseChan <- 1
		<-a.taskCloseDoneChan
		log.Println("Task Server Stop OK")
	}
	if a.HttpEnable {
		a.httpCloseChan <- 1
		<-a.httpCloseDoneChan
		log.Println("Http Server Stop OK")
	}

	if a.RpcEnable {
		a.rpcCloseChan <- 1
		<-a.rpcCloseDoneChan
		log.Println("Rpc Server Stop OK")
	}

	// broad close chan
	close(GlobalCloseChan)
}

var TaskCloseChan = make(chan int)

var GlobalCloseChan = make(chan int)
