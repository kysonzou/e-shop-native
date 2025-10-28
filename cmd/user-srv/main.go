package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/kyson/e-shop-native/internal/user-srv/conf"
	"github.com/kyson/e-shop-native/internal/user-srv/data"
	"github.com/kyson/e-shop-native/internal/user-srv/server"
)

var flagconf string

type App struct {
	*Server
	//conf_srv *conf.Server
	mysql_dsn string
	//data_srv*conf.Data
	logger *zap.Logger
}

type Server struct {
	grpc_srv  *grpc.Server
	http_srv  *http.Server
	admin_srv *http.Server
	grpc_addr string
}

func NewApp(grpc *server.BusinessGRPCServer,
	http *server.BusinessHTTPServer,
	conf_server *conf.Server,
	data_server *conf.Data,
	logger *zap.Logger,
	admin *server.AdminHTTPServer) *App {
	return &App{
		Server: &Server{
			grpc_srv:  grpc.Server,
			http_srv:  http.Server,
			admin_srv: admin.Server,
			grpc_addr: conf_server.GRPC.Addr,
		},
		//conf_srv: conf_server,
		//data_srv: data_server,
		mysql_dsn: data_server.MySQL.DSN,
		logger:    logger,
	}
}

func (a *Server) Run() []error {
	// 创建一个可以在接收到信号时被取消的 context
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var srvErrs []error
	var wg sync.WaitGroup

	// GRPC
	wg.Go(func() {
		log.Printf("GPRC servers starting: %s", a.grpc_addr)
		// 启动 gRPC 服务器
		lis, err := net.Listen("tcp", a.grpc_addr)
		if err != nil {
			srvErrs = append(srvErrs, fmt.Errorf("failed to listen for gRPC: %w", err))
		}
		// Serve 会在 GracefulStop 调用后返回错误，这是正常行为
		if err := a.grpc_srv.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			// 我们只关心那些不是“服务器正常关闭”的错误
			srvErrs = append(srvErrs, fmt.Errorf("gRPC server failed to serve: %w", err))
		}
	})

	wg.Go(func() {
		<-ctx.Done() // 等待接收到终止信号
		log.Println("Shutting down GPRC servers...")

		// 关闭 gRPC 服务器
		a.grpc_srv.GracefulStop()
	})

	// HTTP
	wg.Go(func() {
		log.Printf("HTTP servers starting: %s", a.http_srv.Addr)
		// 启动 HTTP 服务器
		if err := a.http_srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			srvErrs = append(srvErrs, fmt.Errorf("HTTP server failed to serve: %w", err))
		}

	})

	wg.Go(func() {
		<-ctx.Done() // 等待接收到终止信号
		log.Println("Shutting down HTTP servers...")
		//关闭HTTP
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := a.http_srv.Shutdown(shutdownCtx); err != nil {
			srvErrs = append(srvErrs, fmt.Errorf("HTTP server shutdown error: %w", err))
		}
	})

	// Amidn（mrteics）
	wg.Go(func() {
		log.Printf("Admin servers starting: %s", a.admin_srv.Addr)
		// 启动 HTTP 服务器
		if err := a.admin_srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			srvErrs = append(srvErrs, fmt.Errorf("admin server failed to serve: %w", err))
		}
	})

	wg.Go(func() {
		<-ctx.Done() // 等待接收到终止信号
		log.Println("Shutting down Admin servers...")
		// 关闭admin服务器
		ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancle()
		if err := a.admin_srv.Shutdown(ctx); err != nil {
			srvErrs = append(srvErrs, fmt.Errorf("admin server shutdown error: %w", err))
		}

	})
	wg.Wait()
	return srvErrs
}

func init() {
	//从终端读取config.yaml文件
	flag.StringVar(&flagconf, "conf", "./configs/config.yaml", "config path, eg: -conf config.yaml")

}

func main() {

	app, cleanup, err := InitializeApp()
	if err != nil {
		log.Printf("init app error: %v\n", err)
		panic(err)
	}
	defer cleanup()

	err = migrateDatabase(app.mysql_dsn)
	if err != nil {
		fmt.Printf("migrate database error: %v\n", err)
		panic(err)
	}

	srvErrs := app.Run()
	for _, err := range srvErrs {
		log.Printf("run app error: %v\n", err)
		panic(err)
	}

}

func migrateDatabase(dsn string) error {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	// 在这里同时迁移 UserPO 和 ProductPO
	return db.AutoMigrate(&data.UserPO{})
}
