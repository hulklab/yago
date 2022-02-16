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
	"testing"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"github.com/robfig/cron"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

// support custom app before run
type AppInitHook func(app *App) error

var appInitHooks = make([]AppInitHook, 0)

func AddAppInitHook(hs ...AppInitHook) {
	appInitHooks = append(appInitHooks, hs...)
}

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

	// com close chan
	comCloseDoneChan chan int
}

var (
	errHttpRouteEmpty = errors.New("http router is empty")
	errTaskRouteEmpty = errors.New("task router is empty")
)

type httpStaticPath struct {
	Route string `json:"route"`
	Path  string `json:"path"`
}

func NewApp() *App {
	// new app
	app := new(App)

	app.DebugMode = Config.GetBool("app.debug")

	// init http
	app.HttpEnable = Config.GetBool("app.http_enable")
	if app.HttpEnable {
		if app.DebugMode {
			app.HttpRunMode = gin.DebugMode
		} else {
			app.HttpRunMode = gin.ReleaseMode
		}
		gin.SetMode(app.HttpRunMode)
		app.httpEngine = gin.New()
		// use logger
		if app.DebugMode {
			app.httpEngine.Use(gin.Logger())
		} else {
			app.httpEngine.Use(gin.Recovery())
		}
		app.httpCloseChan = make(chan int, 1)
		app.httpCloseDoneChan = make(chan int, 1)

		app.HttpViewRender = Config.GetBool("app.http_view_render")
		if app.HttpViewRender {
			app.HttpViewPath = Config.GetString("app.http_view_path")
			if app.HttpViewPath != "" {
				app.httpEngine.LoadHTMLGlob(app.HttpViewPath)
			}

			if Config.IsSet("app.http_static_paths") {
				hsp := Config.Get("app.http_static_paths")
				httpStaticPaths := make([]httpStaticPath, 0)
				err := mapstructure.Decode(hsp, &httpStaticPaths)
				if err != nil {
					log.Fatalln("parse http static paths err:", err.Error())
				}

				for _, staticPath := range httpStaticPaths {
					app.httpEngine.Static(staticPath.Route, staticPath.Path)
				}
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

	app.comCloseDoneChan = make(chan int, 1)

	return app
}

// 此 init 最先执行，配置文件此处初始化
func init() {
	// avoid go test error
	testing.Init()

	initConfig()

	log.SetFlags(log.LstdFlags)

	initGrpcServer()
}

func (a *App) HttpEngine() *gin.Engine {
	return a.httpEngine
}

func (a *App) Run() {
	if len(appInitHooks) > 0 {
		for _, f := range appInitHooks {
			err := f(a)
			if err != nil {
				log.Fatalf("init err:%s", err.Error())
			}
		}
	}

	if a.TaskEnable {
		// 开启 task
		go a.runTask()
	}

	if a.RpcEnable {
		// 开启 rpc
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
	pidFile, ok := getPidFile()
	if !ok {
		return
	}

	pf, err := os.Create(pidFile)
	if err != nil {
		log.Fatalf("pidfile check err:%v\n", err)
		return
	}

	defer pf.Close()

	newPid := os.Getpid()
	_, err = pf.Write([]byte(fmt.Sprintf("%d", newPid)))
	if err != nil {
		log.Fatalf("write pid err:%v\n", err)
		return
	}

	debug("app is running with pid:", newPid)
}

func (a *App) registerHttpRouter(g *HttpGroupRouter) {
	for _, r := range g.HttpRouterList {
		method := strings.ToUpper(r.Method)
		actions := r.Actions

		var handlers []gin.HandlerFunc

		for _, handler := range actions {
			do := handler
			handlers = append(handlers, func(c *gin.Context) {
				ctx, _ := getCtxFromGin(c)
				do(ctx)
			})
		}

		name := runtime.FuncForPC(reflect.ValueOf(actions[len(actions)-1]).Pointer()).Name()
		debugf("[HTTP] %-6s %-25s --> %s\n", method, r.Url(), strings.NewReplacer("(", "", ")", "", "*", "").Replace(name))

		switch method {
		case http.MethodGet:
			g.GinGroup.GET(r.Path, handlers...)
		case http.MethodPost:
			g.GinGroup.POST(r.Path, handlers...)
		case http.MethodDelete:
			g.GinGroup.DELETE(r.Path, handlers...)
		case http.MethodPut:
			g.GinGroup.PUT(r.Path, handlers...)
		case http.MethodOptions:
			g.GinGroup.OPTIONS(r.Path, handlers...)
		case http.MethodPatch:
			g.GinGroup.PATCH(r.Path, handlers...)
		case http.MethodHead:
			g.GinGroup.HEAD(r.Path, handlers...)
		default:
			g.GinGroup.Any(r.Path, handlers...)
		}
	}
}

func (a *App) registerHttpGroupRouter(group map[string]*HttpGroupRouter) {
	for _, g := range group {
		if g.Parent == nil {
			g.GinGroup = a.httpEngine.Group(g.Prefix)
		} else {
			g.GinGroup = g.Parent.GinGroup.Group(g.Prefix)
		}

		if len(g.Middlewares) > 0 {
			for _, m := range g.Middlewares {
				handler := m
				g.GinGroup.Use(func(c *gin.Context) {
					ctx, err := getCtxFromGin(c)
					if err != nil {
						log.Println(err)
						return
					}
					handler(ctx)
				})
			}
		}

		a.registerHttpRouter(g)

		a.registerHttpGroupRouter(g.Children)
	}
}

func (a *App) loadHttpRouter() error {
	if len(httpGroupRouterMap) == 0 {
		return errHttpRouteEmpty
	}

	// cors
	if a.HttpCorsAllowAllOrigins || len(a.HttpCorsAllowOrigins) != 0 {
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

	// no route handler
	if httpNoRouterHandler != nil {
		a.httpEngine.NoRoute(func(c *gin.Context) {
			ctx := newCtx(c)
			httpNoRouterHandler(ctx)
		})
	}

	// register global middleware
	for _, m := range httpGlobalMiddleware {
		handler := m
		a.httpEngine.Use(func(c *gin.Context) {
			ctx, err := getCtxFromGin(c)
			if err != nil {
				log.Println(err)
				return
			}
			handler(ctx)
		})
	}

	a.registerHttpGroupRouter(httpGroupRouterMap)

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

	// defend slow dos attack
	if Config.IsSet("app.http_read_timeout") {
		srv.ReadTimeout = Config.GetDuration("app.http_read_timeout")
	}

	if Config.IsSet("app.http_read_header_timeout") {
		srv.ReadTimeout = Config.GetDuration("app.http_read_header_timeout")
	}

	if a.HttpSslOn {
		go func() {
			// service connections
			debugf("https listen on: %s", srv.Addr)

			if err := srv.ListenAndServeTLS(a.HttpCertFile, a.HttpKeyFile); err != nil && err != http.ErrServerClosed {
				log.Fatalf("listen: %s\n", err)
			} else {
			}
		}()
	} else {
		go func() {
			// service connections
			debugf("http listen on: %s", srv.Addr)

			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("listen: %s\n", err)
			}
		}()
	}

	<-a.httpCloseChan
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(Config.GetInt64("app.http_stop_time_wait"))*time.Second,
	)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			debug("http server already closed")
		} else if errors.Is(err, context.DeadlineExceeded) {
			debug("http server gracefully shutdown timeout")
		} else {
			debug("http server gracefully shutdown error", err)
		}
	} else {
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
		name := runtime.FuncForPC(reflect.ValueOf(action).Pointer()).Name()
		name = strings.NewReplacer("(", "", ")", "", "*", "").Replace(name)
		if router.Spec == "@loop" {
			wg.Add(1)
			go func() {
				defer wg.Done()
				action()
				debugf("[TASK] %-32s --> %s\n", "stop", name)
			}()
		} else {
			err := c.AddFunc(router.Spec, func() {
				wg.Add(1)
				defer wg.Done()
				action()
			})
			if err != nil {
				continue
			}
		}
		debugf("[TASK] %-32s --> %s\n", router.Spec, name)
	}

	c.Start()

	<-a.taskCloseChan
	go func() {
		c.Stop()
		wg.Wait()
		a.taskCloseDoneChan <- 1
	}()
	time.Sleep(time.Duration(Config.GetInt64("app.task_stop_time_wait")) * time.Second)
	a.taskCloseDoneChan <- 1

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
		cred, err := credentials.NewServerTLSFromFile(certFile, keyFile)
		if err != nil {
			log.Fatalf("Failed to generate credentials %v", err)
		}

		// 实例化 grpc Server, 并开启 TSL 认证
		RpcServer = grpc.NewServer(grpc.Creds(cred))

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
			debugf("[TASK] %-32s --> %s\n", k, method.Name)
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
		} else {
			debugf("rpc listen on: %s", rpcAddr)
		}
	}()

	<-a.rpcCloseChan
	a.rpcEngine.GracefulStop()
	a.rpcCloseDoneChan <- 1
}

func (a *App) Close() {
	close(StopChan)

	if a.TaskEnable {
		//close(TaskCloseChan)
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

	go func() {
		Component.Close()
		a.comCloseDoneChan <- 1
	}()

	var comCloseTimeWait time.Duration
	if Config.IsSet("app.com_close_time_wait") {
		comCloseTimeWait = time.Duration(Config.GetInt64("app.com_stop_time_wait")) * time.Second
	} else {
		comCloseTimeWait = 10 * time.Second
	}
	select {
	case <-a.comCloseDoneChan:
		log.Println("Components Close OK")
	case <-time.After(comCloseTimeWait):
		log.Println("Components Close Timeout")
	}
}

//var TaskCloseChan = make(chan int)
var StopChan = make(chan struct{})
