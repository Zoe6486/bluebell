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
	configPath := "./conf/config.yaml" // 部署就用config.yaml
	if len(os.Args) >= 2 {
		configPath = os.Args[1]
	}
	if err := setting.Init(configPath); err != nil {
		// fmt.Printf("load config failed, err:%v\n", err)
		// return
		panic(fmt.Sprintf("load config failed: %v", err))
	}

	if err := logger.Init(setting.Conf.LogConfig, setting.Conf.Mode); err != nil {
		// fmt.Printf("init logger failed, err:%v\n", err)
		// return
		panic(err)
	}

	// 初始化 snowflake
	// if err := snowflake.Init(setting.Conf.StartTime, setting.Conf.MachineID); err != nil {
	// 	fmt.Printf("init snowflake failed, err:%v\n", err)
	// 	return
	// }
	if err := snowflake.Init("2026-01-01", 1); err != nil {
		// fmt.Printf("init snowflake failed, err:%v\n", err)
		// return
		panic(err)
	}
	// 初始化自定义的validator
	if err := controller.InitValidator(); err != nil {
		// fmt.Printf("init validaor failed, err:%v\n", err)
		// return
		panic(err)
	}

	// 初始化 MySQL
	// if err := mysql.Init(setting.Conf.MySQLConfig); err != nil {
	// 	// fmt.Printf("init mysql failed, err:%v\n", err)
	// 	// return
	// 	panic(err)
	// }
	// defer mysql.Close() // 程序退出关闭数据库连接
	db, err := mysql.Init(setting.Conf.MySQLConfig)
	if err != nil {
		panic(err)
	}
	// defer db.Close()
	// 延迟关闭，并处理错误
	defer func() {
		if err := db.Close(); err != nil {
			zap.L().Error("mysql db close failed", zap.Error(err))
		}
	}()

	// 初始化 Redis
	if err := redis.Init(setting.Conf.RedisConfig); err != nil {
		// fmt.Printf("init redis failed, err:%v\n", err)
		// return
		panic(err)
	}
	defer redis.Close()

	// 4. 注册路由
	r := router.Setup(db)

	//
	// err := r.Run(fmt.Sprintf(":%d", setting.Conf.Port))
	// if err != nil {
	// 	fmt.Printf("run server failed, err:%v\n", err)
	// 	return
	// }

	// 5. 启动服务与优雅停机 (Standard Kiwi Enterprise Pattern)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", setting.Conf.Port),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("listen: ", zap.Error(err))
		}
	}()

	// 等待退出信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.L().Info("Shutdown Server ...")

	// 5秒内尝试处理完存量请求
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown:", zap.Error(err))
	}

	zap.L().Info("Server exiting")
}
