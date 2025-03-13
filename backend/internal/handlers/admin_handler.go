package handlers

import (
	"net/http"
	"normaladmin/backend/internal/models"
	"normaladmin/backend/internal/services"
	"normaladmin/backend/pkg/utils/response"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AdminHandler struct {
	adminService services.AdminService
}

func NewAdminHandler(adminService services.AdminService) *AdminHandler {
	return &AdminHandler{adminService: adminService}
}

// GetAdminList godoc
// @Summary 获取管理员列表
// @description 获取管理员列表，支持按用户名搜索、分页和排序
// @Tags 管理员管理
// @Accept json
// @Produce json
// @Param username query string false "用户名，用于搜索过滤"
// @Param page query int false "页码，默认值：1" default(1)
// @Param pageSize query int false "每页数量，默认值：10" default(10)
// @Param sortField query string false "排序字段，多个字段用逗号分隔" example("id,created_at")
// @Param sortOrder query string false "排序方式，多个排序方式用逗号分隔(asc/desc)" example("desc,asc")
// @Success 200 {object} response.ResponseData{data=object{admins=[]models.Admin,total=int}} "成功"
// @Failure 500 {object} response.ResponseData "内部错误"
// @Router /gam/admins [get]
func (h *AdminHandler) GetAdminList(c *gin.Context) {
	query := make(map[string]interface{})
	if username := c.Query("username"); username != "" {
		query["username"] = username
	}
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("pageSize", "10")
	// 处理排序
	var opts []models.QueryOption

	sortField := c.DefaultQuery("sortField", "")
	if sortField != "" {
		sortFields := strings.Split(sortField, ",")
		sortOrders := strings.Split(c.DefaultQuery("sortOrder", "desc"), ",")
		opts = append(opts, models.WithSort(sortFields, sortOrders))
	}
	admins, total, err := h.adminService.List(query, page, pageSize, opts...)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get admin list")
		return
	}

	response.Success(c, gin.H{
		"admins": admins,
		"total":  total,
	})
}

// GetAdmin godoc
// @Summary 获取单个管理员
// @Description 根据管理员ID获取详细信息
// @Tags 管理员管理
// @Accept json
// @Produce json
// @Param id path int true "管理员ID" minimum(1)
// @Success 200 {object} response.ResponseData{data=object{admin=models.Admin}} "成功"
// @Failure 400 {object} response.ResponseData "无效的管理员ID"
// @Failure 404 {object} response.ResponseData "管理员不存在"
// @Failure 500 {object} response.ResponseData "内部错误"
// @Router /gam/admins/{id} [get]
func (h *AdminHandler) GetAdmin(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID format")
		return
	}

	admin, err := h.adminService.GetByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get admin")
		return
	}

	response.Success(c, gin.H{"admin": admin})
}

// CreateAdmin godoc
// @Summary 创建管理员
// @Description 创建新的管理员账号，包括用户名、密码、角色等信息
// @Tags 管理员管理
// @Accept json
// @Produce json
// @Param admin body models.Admin true "管理员信息，包括用户名、密码、角色ID等" example({"username":"admin","password":"123456","role_id":1,"email":"admin@example.com"})
// @Success 200 {object} response.ResponseData{data=object{admin=models.Admin}} "成功"
// @Failure 400 {object} response.ResponseData "请求参数错误"
// @Failure 500 {object} response.ResponseData "内部错误"
// @Router /gam/admins [post]
func (h *AdminHandler) CreateAdmin(c *gin.Context) {
	var admin models.Admin
	if err := c.ShouldBindJSON(&admin); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to hash password")
		return
	}
	admin.Password = string(hashedPassword)
	if err := h.adminService.Create(&admin); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create admin")
		return
	}

	response.Success(c, gin.H{"admin": admin})
}

// UpdateAdmin godoc
// @Summary 更新管理员
// @Description 更新管理员信息，包括用户名、角色等（密码单独更新）
// @Tags 管理员管理
// @Accept json
// @Produce json
// @Param id path int true "管理员ID" minimum(1)
// @Param admin body models.Admin true "管理员信息" example({"username":"admin","role_id":1,"email":"admin@example.com"})
// @Success 200 {object} response.ResponseData{data=object{admin=models.Admin}} "成功"
// @Failure 400 {object} response.ResponseData "无效的请求参数"
// @Failure 404 {object} response.ResponseData "管理员不存在"
// @Failure 500 {object} response.ResponseData "内部错误"
// @Router /gam/admins/{id} [put]
func (h *AdminHandler) UpdateAdmin(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID format")
		return
	}

	var admin models.Admin
	if err := c.ShouldBindJSON(&admin); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if admin.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "Failed to hash password")
			return
		}
		admin.Password = string(hashedPassword)
	}
	if err := h.adminService.Update(uint(id), &admin); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update admin")
		return
	}

	response.Success(c, gin.H{"admin": admin})
}

// DeleteAdmin godoc
// @Summary 删除管理员
// @Description 根据ID删除指定管理员账号
// @Tags 管理员管理
// @Accept json
// @Produce json
// @Param id path int true "管理员ID" minimum(1)
// @Success 200 {object} response.ResponseData{data=object{message=string}} "成功"
// @Failure 400 {object} response.ResponseData "无效的管理员ID"
// @Failure 404 {object} response.ResponseData "管理员不存在"
// @Failure 500 {object} response.ResponseData "内部错误"
// @Router /gam/admins/{id} [delete]
func (h *AdminHandler) DeleteAdmin(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID format")
		return
	}

	if err := h.adminService.Delete(uint(id), false); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete admin")
		return
	}

	response.Success(c, gin.H{"message": "Admin deleted successfully"})
}

// UpdateAdminStatus godoc
// @Summary 更新管理员状态
// @Description 更新管理员的启用/禁用状态（0：禁用，1：启用）
// @Tags 管理员管理
// @Accept json
// @Produce json
// @Param id path int true "管理员ID" minimum(1)
// @Param status body object true "状态信息" example({"status":1})
// @Success 200 {object} response.ResponseData{data=object{message=string}} "成功"
// @Failure 400 {object} response.ResponseData "无效的请求参数"
// @Failure 404 {object} response.ResponseData "管理员不存在"
// @Failure 500 {object} response.ResponseData "内部错误"
// @Router /gam/admins/{id}/status [put]
func (h *AdminHandler) UpdateAdminStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID format")
		return
	}

	var req struct {
		Status int `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.adminService.Update(uint(id), map[string]interface{}{"status": req.Status}); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update admin status")
		return
	}

	response.Success(c, gin.H{"message": "Admin status updated successfully"})
}

// UpdatePassword godoc
// @Summary 更新管理员密码
// @Description 更新管理员的登录密码，需要提供旧密码
// @Tags 管理员管理
// @Accept json
// @Produce json
// @Param id path int true "管理员ID" minimum(1)
// @Param password body object true "密码信息" example({"old_password":"123456","new_password":"654321"})
// @Success 200 {object} response.ResponseData{data=object{message=string}} "成功"
// @Failure 400 {object} response.ResponseData "无效的请求参数"
// @Failure 401 {object} response.ResponseData "旧密码错误"
// @Failure 404 {object} response.ResponseData "管理员不存在"
// @Failure 500 {object} response.ResponseData "内部错误"
// @Router /gam/admins/{id}/password [put]
func (h *AdminHandler) UpdatePassword(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID format")
		return
	}

	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.adminService.UpdatePassword(uint(id), req.OldPassword, req.NewPassword); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update password")
		return
	}

	response.Success(c, gin.H{"message": "Password updated successfully"})
}

// CheckAdminFieldUnique godoc
// @Summary 检查管理员字段唯一性
// @Description 检查指定字段的值在管理员表中是否唯一
// @Tags 管理员管理
// @Accept json
// @Produce json
// @Param field query string true "字段名称" example("username")
// @Param value query string true "字段值" example("admin")
// @Param excludeId query int false "需要排除的管理员ID" minimum(1)
// @Success 200 {object} response.ResponseData{data=object{unique=bool}} "成功"
// @Failure 400 {object} response.ResponseData "无效的请求参数"
// @Failure 500 {object} response.ResponseData "内部错误"
// @Router /gam/admins/check-field [get]
func (h *AdminHandler) CheckAdminFieldUnique(c *gin.Context) {
	field := c.Query("field")
	value := c.Query("value")
	excludeID, _ := strconv.ParseUint(c.Query("excludeId"), 10, 32)

	if field == "" || value == "" {
		response.Error(c, http.StatusBadRequest, "Field and value are required")
		return
	}

	isUnique, err := h.adminService.CheckAdminFieldUnique(field, value, uint(excludeID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to check field uniqueness")
		return
	}

	response.Success(c, gin.H{"unique": isUnique})
}
