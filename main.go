package main

import (
	"bluebell/controller"
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/logger"
	"bluebell/pkg/snowflake"
	"bluebell/router"
	"bluebell/setting"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// 加载配置
	//Bash: `./myapp config.yaml debug`中：
	//os.Args[0] // "./myapp" —— 程序本身的名字
	//os.Args[1] // "config.yaml" —— 第一个命令行参数
	//os.Args[2] // "debug" 第二个命令行参数
	if err := setting.Init(os.Args[1]); err != nil {
		fmt.Printf("load config failed, err:%v\n", err)
		return
	}
	if err := logger.Init(setting.Conf.LogConfig, setting.Conf.Mode); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}
	if err := mysql.Init(setting.Conf.MySQLConfig); err != nil {
		fmt.Printf("init mysql failed, err:%v\n", err)
		return
	}
	defer mysql.Close() // 程序退出关闭数据库连接
	if err := redis.Init(setting.Conf.RedisConfig); err != nil {
		fmt.Printf("init redis failed, err:%v\n", err)
		return
	}
	defer redis.Close()
	// 初始化 snowflake
	// if err := snowflake.Init(setting.Conf.StartTime, setting.Conf.MachineID); err != nil {
	// 	fmt.Printf("init snowflake failed, err:%v\n", err)
	// 	return
	// }
	if err := snowflake.Init("2026-01-01", 1); err != nil {
		fmt.Printf("init snowflake failed, err:%v\n", err)
	}
	// 初始化自定义的validator
	if err := controller.InitValidator(); err != nil {
		fmt.Printf("init validaor failed, err:%v\n", err)
		return
	}
	// 注册路由
	r := router.Setup()
	//r := router.SetupRouter(setting.Conf.Mode)
	err := r.Run(fmt.Sprintf(":%d", setting.Conf.Port))
	if err != nil {
		fmt.Printf("run server failed, err:%v\n", err)
		return
	}
}
