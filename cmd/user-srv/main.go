package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"

	//"os"
	"os/signal"
	"syscall"
	"time"

	kratosconfig "github.com/go-kratos/kratos/v2/config"
	kratosconfigfile "github.com/go-kratos/kratos/v2/config/file"
	"github.com/kyson/e-shop-native/internal/user-srv/conf"
	"google.golang.org/grpc"
)


var flagconf string

func init(){
	//从终端读取config.yaml文件
	flag.StringVar(&flagconf, "conf", "./configs/config.yaml", "config path, eg: -conf config.yaml")
	
}

type App struct{
	grpc *grpc.Server
	http *http.Server
	conf *conf.Server
}

func NewApp(grpc *grpc.Server, http *http.Server, conf *conf.Server) *App{
	return &App{
		grpc: grpc,
		http: http,
		conf: conf,
	}
}

func (a *App) Run() []error {
	// 1. 创建一个可以在接收到信号时被取消的 context
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var srvErrs []error
	var wg sync.WaitGroup
	
	wg.Go(func ()  {
		log.Printf("GPRC servers starting: %s", a.conf.GRPC.Address)
		// 启动 gRPC 服务器
		lis, err := net.Listen("tcp", a.conf.GRPC.Address)
		if err != nil {
			srvErrs = append(srvErrs, fmt.Errorf("failed to listen for gRPC: %w", err))
		}
		// Serve 会在 GracefulStop 调用后返回错误，这是正常行为
		if err := a.grpc.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
            // 我们只关心那些不是“服务器正常关闭”的错误
			srvErrs = append(srvErrs, fmt.Errorf("gRPC server failed to serve: %w", err))
		}
	})

	wg.Go(func ()  {
		<-ctx.Done() // 等待接收到终止信号
		log.Println("Shutting down GPRC servers...")
		
		// 关闭 gRPC 服务器
		a.grpc.GracefulStop()
	})

	wg.Go(func ()  {	
		log.Printf("HTTP servers starting: %s", a.conf.HTTP.Address)
		// 启动 HTTP 服务器
		if err := a.http.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			srvErrs = append(srvErrs, fmt.Errorf("HTTP server failed to serve: %w", err))
		}
	
	})

	wg.Go(func ()  {
		<-ctx.Done() // 等待接收到终止信号
		log.Println("Shutting down HTTP servers...")
		//关闭HTTP
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()	
		if err := a.http.Shutdown(shutdownCtx); err != nil {
			srvErrs = append(srvErrs, fmt.Errorf("HTTP server shutdown error: %w", err))
		}
	})
	wg.Wait()
	return srvErrs
}

func main(){
	bc, err := newConfig()
	if err != nil {
		log.Printf("init config error: %v\n", err)
		panic(err)
	} 

	app, cleanup, err := InitializeApp(bc.Server, bc.Data)
	if err != nil {
		log.Printf("init app error: %v\n", err)
		panic(err)
	}
	defer cleanup()

	srvErrs := app.Run()
	for _, err := range srvErrs {
		log.Printf("run app error: %v\n", err)
	}	
}

func newConfig() (*conf.Bootstrap , error){
	flag.Parse()
	// 1. 读取配置文件
	config := kratosconfig.New(
		kratosconfig.WithSource(
			kratosconfigfile.NewSource(flagconf),
		),
	)
	defer config.Close()

	// 2. 加载配置文件
	if err := config.Load(); err != nil {
		return nil, err
	}

	// 3. 反射到结构体
	var bc conf.Bootstrap	
	if err := config.Scan(&bc); err != nil {
		return nil, err
	}
	return &bc, nil
}