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

	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/kyson/e-shop-native/internal/user-srv/conf"
	"github.com/kyson/e-shop-native/internal/user-srv/data"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var flagconf string

func init() {
	//从终端读取config.yaml文件
	flag.StringVar(&flagconf, "conf", "./configs/config.yaml", "config path, eg: -conf config.yaml")

}

type App struct {
	grpc *grpc.Server
	http *http.Server
	conf_srv *conf.Server
	bc *conf.Bootstrap
	logger *zap.Logger
}

func NewApp(grpc *grpc.Server, http *http.Server, conf_server *conf.Server, bc *conf.Bootstrap, logger *zap.Logger) *App {
	return &App{
		grpc: grpc,
		http: http,
		conf_srv: conf_server,
		bc: bc,
		logger: logger,
	}
}

func (a *App) Run() []error {
	// 1. 创建一个可以在接收到信号时被取消的 context
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var srvErrs []error
	var wg sync.WaitGroup

	wg.Go(func() {
		log.Printf("GPRC servers starting: %s", a.conf_srv.GRPC.Addr)
		// 启动 gRPC 服务器
		lis, err := net.Listen("tcp", a.conf_srv.GRPC.Addr)
		if err != nil {
			srvErrs = append(srvErrs, fmt.Errorf("failed to listen for gRPC: %w", err))
		}
		// Serve 会在 GracefulStop 调用后返回错误，这是正常行为
		if err := a.grpc.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			// 我们只关心那些不是“服务器正常关闭”的错误
			srvErrs = append(srvErrs, fmt.Errorf("gRPC server failed to serve: %w", err))
		}
	})

	wg.Go(func() {
		<-ctx.Done() // 等待接收到终止信号
		log.Println("Shutting down GPRC servers...")

		// 关闭 gRPC 服务器
		a.grpc.GracefulStop()
	})

	wg.Go(func() {
		log.Printf("HTTP servers starting: %s", a.conf_srv.HTTP.Addr)
		// 启动 HTTP 服务器
		if err := a.http.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			srvErrs = append(srvErrs, fmt.Errorf("HTTP server failed to serve: %w", err))
		}

	})

	wg.Go(func() {
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

func main() {

	app, cleanup, err := InitializeApp()
	if err != nil {
		log.Printf("init app error: %v\n", err)
		panic(err)
	}
	defer cleanup()

	migrateDatabase(app.bc)

	srvErrs := app.Run()
	for _, err := range srvErrs {
		log.Printf("run app error: %v\n", err)
	}

}

func migrateDatabase(c *conf.Bootstrap) error {
	db, err := gorm.Open(mysql.Open(c.Data.MySQL.DSN), &gorm.Config{})
	if err != nil {
		return err
	}
	// 在这里同时迁移 UserPO 和 ProductPO
	return db.AutoMigrate(&data.UserPO{})
}
