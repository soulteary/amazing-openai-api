package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	AoaModel "github.com/soulteary/amazing-openai-api/internal/model"
	AoaRouter "github.com/soulteary/amazing-openai-api/internal/router"
	"github.com/soulteary/amazing-openai-api/internal/version"
	"github.com/soulteary/amazing-openai-api/models/azure"
	"github.com/soulteary/amazing-openai-api/models/gemini"
	"github.com/soulteary/amazing-openai-api/models/yi"
	"github.com/soulteary/amazing-openai-api/pkg/logger"
)

const (
	_DEFAULT_PORT = 8080
	_DEFAULT_HOST = "0.0.0.0"
	_DEFAULT_TYPE = "azure"

	_ENV_KEY_NAME_PORT    = "AOA_PORT"
	_ENV_KEY_NAME_HOST    = "AOA_HOST"
	_ENV_KEY_SERVICE_TYPE = "AOA_TYPE"
)

// refs: https://github.com/soulteary/flare/blob/main/cmd/cmd.go
func Parse() {
	// 1. First try to get the environment variables
	flags := parseEnvVars()
	// 2. Then try to get the command line flags, overwrite the environment variables
	// flags := parseCLI(envs)

	log := logger.GetLogger()
	log.Println("程序启动中 🚀")
	log.Println("程序版本", version.Version)
	log.Println("程序构建日期", version.BuildDate)
	log.Println("程序 Git Commit", version.GitCommit)
	log.Println("程序服务地址", fmt.Sprintf("%s:%d", flags.Host, flags.Port))

	startDaemon(&flags)
}

// refs: https://github.com/soulteary/flare/blob/main/cmd/daemon.go
func startDaemon(flags *AoaModel.Flags) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	log := logger.GetLogger()

	router.Use(logger.Logger(log), gin.Recovery())

	AoaRouter.RegisterMiscRoute(router)

	switch flags.Type {
	case "azure":
		err := azure.Init()
		if err != nil {
			log.Fatalf("初始化 Azure OpenAI API 出错: %s\n", err)
		}
	case "yi":
		err := yi.Init()
		if err != nil {
			log.Fatalf("初始化 Yi API 出错: %s\n", err)
		}
	case "gemini":
		err := gemini.Init()
		if err != nil {
			log.Fatalf("初始化 Gemini API 出错: %s\n", err)
		}
	}
	AoaRouter.RegisterModelRoute(router, flags.Type)

	srv := &http.Server{
		Addr:              ":" + strconv.Itoa(flags.Port),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       5 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("程序启动出错: %s\n", err)
		}
	}()
	log.Println("程序已启动完毕 🚀")

	<-ctx.Done()

	stop()
	log.Println("程序正在关闭中，如需立即结束请按 CTRL+C")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("程序强制关闭: ", err)
	}

	log.Println("期待与你的再次相遇 ❤️")
}
