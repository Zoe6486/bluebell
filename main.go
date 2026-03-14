package main

import (
	"bluebell/controller"
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/logger"
	"bluebell/pkg/snowflake"
	"bluebell/router"
	"bluebell/setting"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
)

func main() {
	// // 加载配置
	// //Bash: `./myapp config.yaml debug`中：
	// //os.Args[0] // "./myapp" —— 程序本身的名字
	// //os.Args[1] // "config.yaml" —— 第一个命令行参数
	// //os.Args[2] // "debug" 第二个命令行参数
	// if len(os.Args) < 2 {
	// 	fmt.Println("need config file.eg: bluebell config.yaml")
	// 	return
	// }
	// if err := setting.Init(os.Args[1]); err != nil {
	// 	fmt.Printf("load config failed, err:%v\n", err)
	// 	return
	// }

	// configPath := "./conf/dev.yaml" // 默认 dev
	// configPath := "./conf/config.yaml" // 部署就用config.yaml
	// if len(os.Args) >= 2 {
	// 	configPath = os.Args[1]
	// }

	var configPath string
	if len(os.Args) >= 2 {
		configPath = os.Args[1]
	} else {
		configPath = "./conf/config.yaml" // 默认值，文件不存在时会 warn 并 fallback 到 env
	}

	// 先尝试加载配置（现在 Init 内部已处理文件不存在的情况）
	if err := setting.Init(configPath); err != nil {
		// 只有真正致命的 Unmarshal 错误才会到这里
		// 建议在开发阶段打印详细错误，生产可改为 zap.Fatal
		fmt.Printf("Fatal: setting.Init failed: %v\n", err)
		os.Exit(1)
	}

	// 配置加载成功后，再初始化 logger（依赖 Conf.LogConfig）
	if err := logger.Init(setting.Conf.LogConfig, setting.Conf.Mode); err != nil {
		// logger 初始化失败是致命的
		fmt.Printf("Fatal: logger init failed: %v\n", err)
		os.Exit(1)
	}
	// 从这里开始，可以安全使用 zap.L()

	// 配置加载的总结日志（方便排查）
	zap.L().Info("Configuration loaded",
		zap.String("config_path_attempted", configPath),
		zap.Bool("using_env_fallback", true), // 可以加更多诊断信息
		zap.String("mode", setting.Conf.Mode),
		zap.Int("port", setting.Conf.Port),
	)
	// 初始化 snowflake（你已经硬编码了时间和 machine id，也可从 Conf 读）
	if err := snowflake.Init("2026-01-01", 1); err != nil {
		zap.L().Fatal("snowflake init failed", zap.Error(err))
	}

	// 初始化自定义 validator
	if err := controller.InitValidator(); err != nil {
		zap.L().Fatal("validator init failed", zap.Error(err))
	}

	// 初始化 MySQL
	db, err := mysql.Init(setting.Conf.MySQLConfig)
	if err != nil {
		zap.L().Fatal("mysql init failed", zap.Error(err))
	}
	defer func() {
		if err := db.Close(); err != nil {
			zap.L().Error("mysql db close failed", zap.Error(err))
		}
	}()

	// 初始化 Redis
	if err := redis.Init(setting.Conf.RedisConfig); err != nil {
		zap.L().Fatal("redis init failed", zap.Error(err))
	}
	defer redis.Close()

	// 注册路由
	r := router.Setup(db)

	// 启动 HTTP 服务 + 优雅关闭
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", setting.Conf.Port),
		Handler: r,
	}

	go func() {
		zap.L().Info("Starting server", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("listen failed", zap.Error(err))
		}
	}()

	// 等待退出信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zap.L().Info("Shutdown Server ...")

	// 优雅关闭（5秒超时）
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown failed", zap.Error(err))
	}

	zap.L().Info("Server exiting")
}
