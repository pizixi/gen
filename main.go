package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"gen/config"
	"gen/log"
	"gen/models"
	"gen/router"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	// 定时任务
	// gincron "gen/cron"

	"github.com/gin-gonic/gin"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "conf", "app.ini", "config file path")
	flag.Parse()

	// 加载配置
	cfg := config.InitConfig(configFile)
	err := cfg.Load()
	if err != nil {
		panic(fmt.Sprintf("load config failed, file: %s, error: %s", configFile, err))
	}

	// 初始化日志---原作者使用的zap
	// log.Init(cfg)
	// defer func() {
	// 	if err := log.Logger.Sync(); err != nil {
	// 		fmt.Printf("Failed to close log: %s\n", err)
	// 	}
	// }()

	// 初始化日志。个人封装
	log.Logger.SetAsync()
	defer log.Logger.Flush()

	//定时任务.个人封装
	// c := gincron.NewCron()
	// c.AddRoute(&gincron.Route{})
	// c.Start()
	// defer c.Stop()

	// 初始化数据库
	err = models.InitDB(cfg)
	if err != nil {
		panic(fmt.Sprintf("init db failed, error: %s", err))
	}

	// 启动Web服务
	err = startServer(cfg)
	if err != nil {
		panic(fmt.Sprintf("Server started failed: %s", err))
	}
}

func startServer(cfg *config.AppConfig) error {
	server := &http.Server{
		Addr:    cfg.HttpAddr + ":" + cfg.HttpPort,
		Handler: getEngine(cfg),
	}
	ctx, cancel := context.WithCancel(context.Background())
	go listenToSystemSignals(cancel)

	go func() {
		<-ctx.Done()
		if err := server.Shutdown(context.Background()); err != nil {
			log.Logger.Error(fmt.Sprintf("Failed to shutdown server: %s", err))
		}
	}()
	log.Logger.Debug("Server started success")
	err := server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		log.Logger.Debug("Server was shutdown gracefully")
		return nil
	}
	return err
}

func getEngine(cfg *config.AppConfig) *gin.Engine {
	gin.SetMode(func() string {
		if cfg.IsDevEnv() {
			return gin.DebugMode
		}
		return gin.ReleaseMode
	}())
	engine := gin.New()
	engine.Use(gin.CustomRecovery(func(c *gin.Context, err interface{}) {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "服务器内部错误，请稍后再试！",
		})
	}))
	router.RegisterRoutes(engine)
	return engine
}

func listenToSystemSignals(cancel context.CancelFunc) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	for {
		select {
		case <-signalChan:
			cancel()
			return
		}
	}
}
