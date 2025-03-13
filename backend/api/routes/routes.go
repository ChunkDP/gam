package routes

import (
	"log"
	v1 "normaladmin/backend/api/routes/v1"
	"normaladmin/backend/config"
	"normaladmin/backend/database"
	"normaladmin/backend/internal/handlers"
	"normaladmin/backend/internal/middleware"
	"normaladmin/backend/pkg/auth"
	"normaladmin/backend/pkg/rabbitmq"
	"normaladmin/backend/pkg/sysconfig"
	"os"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, conf *config.Config, mq *rabbitmq.RabbitMQ) {
	r.Use(middleware.CORS(config.Global.CORS))
	r.Use(middleware.AddHeaders)

	// 初始化数据库
	db := database.GetDB()

	// 获取项目根目录
	rootDir := os.Getenv("PROJECT_ROOT")
	if rootDir == "" {
		log.Fatal("PROJECT_ROOT environment variable is not set")
	}
	// 初始化 Casbin
	if _, err := auth.InitCasbin(db, rootDir); err != nil {
		log.Fatalf("初始化 Casbin 失败: %v", err)
	}

	// 初始化系统配置
	if err := sysconfig.Init(db); err != nil {
		log.Fatalf("初始化系统配置失败: %v", err)
	}

	// 文件访问路由
	files := r.Group("/uploads")
	files.Use(middleware.AntiLeech())
	{

		files.Static("/", sysconfig.Get("upload_path", "uploads"))
	}

	// 后台管理需要认证的路由
	gam := r.Group("/gam")
	gam.POST("/refresh-token", handlers.RefreshToken(conf.JWT))
	gam.POST("/login", handlers.Login(conf.JWT))

	gam.Use(middleware.JWTAuth(conf.JWT))
	gam.Use(middleware.RequestLoggerMiddleware(db))
	{
		// 认证相关路由
		gam.GET("/authmenus", handlers.GetAuthMenus) // 获取用户的菜单和权限信息

		gam.Use(middleware.CasbinMiddleware())

		// 注册路由
		v1.RegisterMenuRoutes(gam)
		v1.RegisterRoleRoutes(gam)
		v1.RegisterAdminRoutes(gam)
		v1.RegisterMemberRoutes(gam)
		v1.RegisterConfigRoutes(gam)
		v1.RegisterUploadRoutes(gam)
		v1.RegisterSystemRoutes(gam)
		v1.RegisterSystemMonitorRoutes(gam)
		v1.RegisterNotificationRoutes(gam, mq)
	}

	//前台路由才涉及api/v1的版本控制
	apiv1 := r.Group("/api")

	memberAuth := apiv1.Group("/member")
	{
		memberAuth.POST("/register", handlers.MemberRegister)
		memberAuth.POST("/login", handlers.MemberLogin(conf.JWT))
		memberAuth.POST("/refresh-token", handlers.RefreshToken(conf.JWT))
	}

	apiv1.Use(middleware.JWTAuth(conf.JWT))
	{

	}
}
