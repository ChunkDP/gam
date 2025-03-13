package v1

import (
	"normaladmin/backend/database"
	"normaladmin/backend/internal/handlers"
	"normaladmin/backend/internal/models"
	"normaladmin/backend/internal/services"

	"github.com/gin-gonic/gin"
)

// RegisterMemberRoutes 注册会员相关路由
func RegisterMemberRoutes(r *gin.RouterGroup) {
	db := database.GetDB()
	base := services.NewBaseCRUDService[models.Member](db)
	cache := services.NewCacheBaseService(base, "member")
	log := services.NewLogBaseService(cache, "member", db)
	memberService := services.NewMemberService(db, log)
	h := handlers.NewMemberHandler(memberService)

	// 会员管理
	members := r.Group("/members")
	{
		members.GET("", h.GetMemberList)
		members.POST("", h.CreateMember)
		members.PUT("/:id", h.UpdateMember)
		members.DELETE("/:id", h.DeleteMember)
		//members.PUT("/:id/status", h.UpdateMemberStatus)
		members.GET("/check-field", h.CheckMemberFieldUnique) // 添加这一行
	}

}
