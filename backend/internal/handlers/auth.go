package handlers

import (
	"fmt"
	"net/http"
	"normaladmin/backend/config"
	"normaladmin/backend/database"
	"normaladmin/backend/internal/middleware"
	"normaladmin/backend/internal/models"
	"normaladmin/backend/pkg/utils/response"

	"github.com/gin-gonic/gin"
)

// LoginRequest 登录请求参数
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RefreshTokenRequest 刷新令牌请求参数
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// MemberLoginRequest 会员登录请求参数
type MemberLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// MemberRegisterRequest 会员注册请求参数
type MemberRegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"required"`
}

// Login godoc
// @Summary 用户登录
// @Description 管理员登录接口，验证用户名密码并返回访问令牌
// @Tags 认证管理
// @Accept json
// @Produce json
// @Param req body LoginRequest true "登录信息" example({"username":"admin","password":"123456"})
// @Success 200 {object} response.ResponseData{data=object{token=string,refreshToken=string,user=object}} "成功"
// @Failure 400 {object} response.ResponseData "请求参数错误"
// @Failure 401 {object} response.ResponseData "用户名或密码错误"
// @Failure 500 {object} response.ResponseData "内部错误"
// @Router /gam/login [post]
func Login(cfg config.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			fmt.Println(err, req)
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		// 验证用户名密码
		var user models.Admin
		if err := database.GetDB().Where("username = ?", req.Username).First(&user).Error; err != nil {
			response.Error(c, http.StatusUnauthorized, "Invalid username or password")
			return
		}

		if !user.CheckPassword(req.Password) {
			response.Error(c, http.StatusUnauthorized, "Invalid username or password")
			return
		}

		// 生成token
		accessToken, refreshToken, err := middleware.GenerateToken(user.ID, user.RoleID, "admin", user.Username, cfg)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "Failed to generate token")
			return
		}

		response.Success(c, gin.H{
			"token":        accessToken,
			"refreshToken": refreshToken,
			"user": gin.H{
				"id":       user.ID,
				"username": user.Username,
				"email":    user.Email,
				"role_id":  user.RoleID,
			},
		})
	}
}

// RefreshToken godoc
// @Summary 刷新访问令牌
// @Description 使用刷新令牌获取新的访问令牌和刷新令牌
// @Tags 认证管理
// @Accept json
// @Produce json
// @Param req body RefreshTokenRequest true "刷新令牌" example({"refresh_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."})
// @Success 200 {object} response.ResponseData{data=object{token=string,refreshToken=string}} "成功"
// @Failure 400 {object} response.ResponseData "请求参数错误"
// @Failure 401 {object} response.ResponseData "无效的刷新令牌"
// @Failure 500 {object} response.ResponseData "内部错误"
// @Router /gam/refresh-token [post]
func RefreshToken(cfg config.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RefreshTokenRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		// 验证刷新令牌
		claims := &middleware.Claims{}
		token, err := middleware.ParseToken(req.RefreshToken, cfg.SecretKey)
		if err != nil || !token.Valid {
			response.Error(c, http.StatusUnauthorized, "Invalid refresh token")
			return
		}

		claims, ok := token.Claims.(*middleware.Claims)
		if !ok {
			response.Error(c, http.StatusUnauthorized, "Invalid token claims")
			return
		}

		// 生成新的访问令牌
		accessToken, refreshToken, err := middleware.GenerateToken(claims.UserID, claims.RoleID, claims.UserType, claims.Username, cfg)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "Failed to generate token")
			return
		}

		response.Success(c, gin.H{
			"token":        accessToken,
			"refreshToken": refreshToken,
		})
	}
}

// MemberLogin godoc
// @Summary 会员登录
// @Description 小程序会员登录接口，验证用户名密码并返回访问令牌
// @Tags 会员认证
// @Accept json
// @Produce json
// @Param req body MemberLoginRequest true "登录信息" example({"username":"member","password":"123456"})
// @Success 200 {object} response.ResponseData{data=object{token=string,refreshToken=string,member=object}} "成功"
// @Failure 400 {object} response.ResponseData "请求参数错误"
// @Failure 401 {object} response.ResponseData "用户名或密码错误"
// @Failure 500 {object} response.ResponseData "内部错误"
// @Router /api/member/login [post]
func MemberLogin(cfg config.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req MemberLoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		// 验证用户名密码
		var member models.Member
		if err := database.GetDB().Where("username = ?", req.Username).First(&member).Error; err != nil {
			response.Error(c, http.StatusUnauthorized, "用户名错误")
			return
		}

		if !member.CheckPassword(req.Password) {
			response.Error(c, http.StatusUnauthorized, "密码错误")
			return
		}

		// 生成token
		accessToken, refreshToken, err := middleware.GenerateToken(member.ID, member.LevelID, "member", member.Username, cfg)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "生成token失败")
			return
		}

		response.Success(c, gin.H{
			"token":        accessToken,
			"refreshToken": refreshToken,
			"member": gin.H{
				"id":       member.ID,
				"username": member.Username,
				"email":    member.Email,
				"mobile":   member.Mobile,
				"level_id": member.LevelID,
			},
		})
	}
}

// MemberRegister godoc
// @Summary 会员注册
// @Description 小程序会员注册接口，创建新会员账号
// @Tags 会员认证
// @Accept json
// @Produce json
// @Param req body MemberRegisterRequest true "注册信息" example({"username":"newmember","password":"123456","email":"member@example.com","phone":"13800138000"})
// @Success 200 {object} response.ResponseData{data=object{message=string}} "成功"
// @Failure 400 {object} response.ResponseData "请求参数错误或用户名已存在"
// @Failure 500 {object} response.ResponseData "内部错误"
// @Router /api/member/register [post]
func MemberRegister(c *gin.Context) {
	var req MemberRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// 检查用户名是否已存在
	var count int64
	if err := database.GetDB().Model(&models.Member{}).Where("username = ?", req.Username).Count(&count).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "系统错误")
		return
	}
	if count > 0 {
		response.Error(c, http.StatusBadRequest, "用户名已存在")
		return
	}

	// 创建新会员
	member := models.Member{
		Username: req.Username,
		Email:    req.Email,
		Mobile:   req.Phone,
		LevelID:  1, // 默认等级
	}
	member.SetPassword(req.Password)

	if err := database.GetDB().Create(&member).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "注册失败")
		return
	}

	response.Success(c, gin.H{
		"message": "注册成功",
	})
}

// GetAuthMenus godoc
// @Summary 获取用户菜单权限
// @Description 获取当前登录用户的菜单列表和按钮权限
// @Tags 认证管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.ResponseData{data=object{menus=[]models.Menu,permissions=[]string}} "成功"
// @Failure 401 {object} response.ResponseData "未授权或Token无效"
// @Failure 500 {object} response.ResponseData "内部错误"
// @Router /gam/authmenus [get]
func GetAuthMenus(c *gin.Context) {
	roleID, exists := c.Get("role_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Role ID not found")
		return
	}

	// 首先获取所有菜单（type = 'menu'）
	var menus []models.Menu
	if err := database.GetDB().
		Select("menus.*").
		Table("menus").
		Joins("INNER JOIN role_menus ON menus.id = role_menus.menu_id").
		Where("role_menus.role_id = ? AND menus.status = 1 AND menus.type = 'menu'", roleID).
		Order("menus.sort ASC").
		Find(&menus).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch menus")
		return
	}

	// 获取所有按钮权限（type = 'button'）
	var buttons []string
	if err := database.GetDB().
		Select("menus.permission").
		Table("menus").
		Joins("INNER JOIN role_menus ON menus.id = role_menus.menu_id").
		Where("role_menus.role_id = ? AND menus.status = 1 AND menus.type = 'button'", roleID).
		Find(&buttons).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch button permissions")
		return
	}

	response.Success(c, gin.H{
		"menus":       menus,
		"permissions": buttons,
	})
}
