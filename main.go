package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/zommage/livemanager/common"
	"github.com/zommage/livemanager/conf"
	. "github.com/zommage/livemanager/logs"
	models "github.com/zommage/livemanager/models"
	liveRouter "github.com/zommage/livemanager/router"
)

var (
	confInfo *conf.Config

	confPath = flag.String("config", "./conf/app.dev.json", "profilePath")

	httpSrv  *http.Server
	httpsSrv *http.Server
)

func Init() error {
	flag.Parse()

	var err error = nil
	confInfo, err = conf.InitConfig(confPath)
	if err != nil {
		return fmt.Errorf("init config is err: %v", err)
	}

	err = InitLog(confInfo.LogConf.LogPath, confInfo.LogConf.LogLevel)
	if err != nil {
		return fmt.Errorf("init log is err: %v", err)
	}

	err = models.InitDb()
	if err != nil {
		return fmt.Errorf("init db err: %v", err)
	}

	// 初始化 rsa 加密
	err = common.InitRsaKey()
	if err != nil {
		return fmt.Errorf("init rsa key err: %v", err)
	}

	return nil
}

func main() {
	//catch global panic
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("panic err: ", err)
			Log.Errorf("Panic err: %v", err)
		}
	}()

	err := Init()
	if err != nil {
		fmt.Println("main init err: ", err)
		return
	}

	router := gin.Default()

	// 解决跨域问题
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AddAllowHeaders("token")
	router.Use(cors.New(corsConfig))

	liveRouter.ApiRouter(router)

	// http services
	httpSrv = &http.Server{
		Addr:    ":" + confInfo.BaseConf.HttpPort,
		Handler: router,
	}

	// https services
	httpsSrv = &http.Server{
		Addr:    ":" + confInfo.BaseConf.HttpsPort,
		Handler: router,
	}

	if confInfo.BaseConf.Env == "dev" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	go func() {
		// http service connections
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			Log.Errorf("http listen: %v", err)
			panic(err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds. os.Interrupt==syscall.SIGINT
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	go handleSignal(quit)

	fmt.Println("livemanager server is start")
	Log.Info("livemanager server is start")
	if err := httpsSrv.ListenAndServeTLS(confInfo.BaseConf.SslCrt, confInfo.BaseConf.SslKey); err != nil && err != http.ErrServerClosed {
		Log.Errorf("https listen: %v", err)
		panic(err)
	}
}

func handleSignal(c chan os.Signal) {
	switch <-c {
	case syscall.SIGQUIT:
		fmt.Println("Shutdown quickly, bye...")
		Log.Info("Shutdown quickly, bye...")
	case os.Interrupt, syscall.SIGTERM: // os.Interrupt==syscall.SIGINT
		fmt.Println("Shutdown gracefully, bye...")
		Log.Info("Shutdown gracefully, bye...")
		// do graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := httpsSrv.Shutdown(ctx); err != nil {
			Log.Error("https Server Shutdown err:", err)
		}
		if err := httpSrv.Shutdown(ctx); err != nil {
			Log.Error("http Server Shutdown err:", err)
		}
	}
	os.Exit(0)
}
