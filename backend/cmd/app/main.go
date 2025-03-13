package main

import (
	"fmt"
	"log"
	"normaladmin/backend/config"
	"normaladmin/backend/database"
	"normaladmin/backend/internal/services"
	"normaladmin/backend/pkg/cache"
	"normaladmin/backend/pkg/logger"
	"normaladmin/backend/pkg/rabbitmq"
	"normaladmin/backend/pkg/utils/encrypt"
	"os"
	"path/filepath"

	"normaladmin/backend/api/routes"
	_ "normaladmin/backend/cmd/app/docs" // swagger文档目录
	"normaladmin/backend/cmd/crons"
	_ "normaladmin/backend/database/migrations" // 导入迁移

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Backend API
// @version 1.0
// @description This is a backend server.

// @termsOfService http://swagger.io/terms/
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Bearer {token}

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemes http https
func main() {
	// 在 main 函数开始处添加
	rootDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("获取当前目录失败: %v", err)
	}
	// 因为是在 cmd/app 目录下，所以需要回退两级到项目根目录
	projectRoot := filepath.Join(rootDir, "../..")
	os.Setenv("PROJECT_ROOT", projectRoot)
	// 获取运行环境
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev" // 默认为开发环境
	}

	// 初始化配置
	if err := config.InitConfig(env); err != nil {
		log.Fatalf("初始化配置失败: %v", err)
	}
	// 初始化加密密钥
	encrypt.InitEncryptKey(config.Global.Security.EncryptKey) // 可以从配置或环境变量中获取

	// 初始化日志
	if err := logger.InitLogger(config.Global.Log); err != nil {
		log.Fatalf("初始化日志失败: %v", err)
	}

	// 初始化数据库
	if err := database.InitDB(config.Global.Database); err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	// 执行数据库迁移
	logger.Info("开始执行数据库迁移")
	if err := database.Migrate(); err != nil {
		logger.Fatal("数据库迁移失败",
			logger.Field("error", err),
		)
	}
	logger.Info("数据库迁移完成")

	// 初始化Redis
	if err := cache.InitRedis(config.Global.Redis); err != nil {
		log.Fatalf("初始化Redis失败: %v", err)
	}
	defer cache.Close() // 程序结束时关闭Redis连接
	// 初始化 RabbitMQ
	mq, err := rabbitmq.NewRabbitMQ("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer mq.Close()

	// 初始化系统监控服务
	systemMonitorService := services.NewSystemMonitorService(database.GetDB())

	// 启动系统监控定时任务
	monitorCron := crons.SetupSystemMonitorCron(systemMonitorService)
	defer monitorCron.Stop()

	gin.SetMode(config.Global.Server.Mode)

	// 创建 Gin 实例
	r := gin.Default()
	r.SetTrustedProxies([]string{"127.0.0.1", "localhost"})
	if env == "dev" {
		// Swagger配置
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// 注册路由
	routes.SetupRoutes(r, config.Global, mq)

	// 启动服务器
	addr := fmt.Sprintf(":%d", config.Global.Server.Port)
	logger.Info("服务器启动",
		logger.Field("env", env),
		logger.Field("port", config.Global.Server.Port),
		logger.Field("mode", config.Global.Server.Mode),
	)

	if err := r.Run(addr); err != nil {
		logger.Fatal("服务器启动失败",
			logger.Field("error", err),
		)
	}
}
