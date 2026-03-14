package setting

import (
	"fmt"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var Conf = new(AppConfig)

type AppConfig struct {
	Name      string `mapstructure:"name"`
	Mode      string `mapstructure:"mode"`
	Version   string `mapstructure:"version"`
	StartTime string `mapstructure:"start_time"`
	MachineID int64  `mapstructure:"machine_id"`
	Port      int    `mapstructure:"port"`

	*LogConfig   `mapstructure:"log"`
	*MySQLConfig `mapstructure:"mysql"`
	*RedisConfig `mapstructure:"redis"`
}

type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DB           string `mapstructure:"dbname"`
	Port         int    `mapstructure:"port"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type RedisConfig struct {
	Host         string `mapstructure:"host"`
	Password     string `mapstructure:"password"`
	Port         int    `mapstructure:"port"`
	DB           int    `mapstructure:"db"`
	PoolSize     int    `mapstructure:"pool_size"`
	MinIdleConns int    `mapstructure:"min_idle_conns"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

func Init(filePath string) error {
	viper.SetConfigFile(filePath)

	// 所有字段设置默认值（生产安全 + fallback）
	viper.SetDefault("name", "bluebell")
	viper.SetDefault("mode", "release")
	viper.SetDefault("port", 8084)
	viper.SetDefault("version", "v0.0.1")
	viper.SetDefault("start_time", "2020-07-01")
	viper.SetDefault("machine_id", 1)

	viper.SetDefault("auth.jwt_expire", 8760)

	viper.SetDefault("log.level", "info") // 生产默认 info，避免 debug 泄露
	viper.SetDefault("log.filename", "/var/log/bluebell.log")
	viper.SetDefault("log.max_size", 200)
	viper.SetDefault("log.max_age", 30)
	viper.SetDefault("log.max_backups", 7)

	viper.SetDefault("mysql.host", "localhost")
	viper.SetDefault("mysql.port", 3306)
	viper.SetDefault("mysql.user", "")
	viper.SetDefault("mysql.password", "") // 必须 env 提供，否则启动失败
	viper.SetDefault("mysql.dbname", "bluebell")
	viper.SetDefault("mysql.max_open_conns", 200)
	viper.SetDefault("mysql.max_idle_conns", 50)

	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("redis.pool_size", 100)

	// 读文件（可选）
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Warning: No config file at %s: %v. Using defaults + env vars only.\n", filePath, err)
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // 下面那行不会自动映射嵌套字段，MYSQL_HOST，MYSQL_PORT不会映射到mysql.host，mysql.port
	viper.AutomaticEnv()                                   // 支持 MYSQL_HOST 或 BLUEBELL_MYSQL_HOST 等

	// 如果想统一前缀（AWS 上常用），可选打开
	// viper.SetEnvPrefix("APP")  // 环境变量变成 APP_MYSQL_HOST

	if err := viper.Unmarshal(Conf); err != nil {
		return fmt.Errorf("config unmarshal failed: %w", err)
	}

	// watch 配置变更（热更新可选，生产慎用）
	viper.WatchConfig()
	// viper.OnConfigChange(func(e fsnotify.Event) {
	// 	fmt.Println("Config changed:", e.Name)

	// 	if err := viper.Unmarshal(Conf); err != nil {
	// 		fmt.Printf("config reload failed: %v\n", err)
	// 	}
	// })
	// 更企业级
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config changed:", e.Name)

		if err := viper.Unmarshal(Conf); err != nil {
			panic(fmt.Errorf("config reload failed: %w", err))
		}
	})

	return nil
}
