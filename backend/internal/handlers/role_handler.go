package handlers

import (
	"net/http"
	"normaladmin/backend/internal/models"
	"normaladmin/backend/internal/services"
	"normaladmin/backend/pkg/utils/response"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	roleService services.RoleService
}

func NewRoleHandler(roleService services.RoleService) *RoleHandler {
	return &RoleHandler{
		roleService: roleService,
	}
}

// GetRoleList godoc
// @Summary 获取角色列表
// @Description 获取角色列表，支持按名称搜索、分页和排序
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param name query string false "角色名称，用于搜索过滤"
// @Param page query int false "页码，默认值：1" default(1)
// @Param pageSize query int false "每页数量，默认值：10" default(10)
// @Param sortField query string false "排序字段，多个字段用逗号分隔" example("id,created_at")
// @Param sortOrder query string false "排序方式，多个排序方式用逗号分隔(asc/desc)" example("desc,asc")
// @Success 200 {object} response.ResponseData{data=object{roles=[]models.Role,total=int}} "成功"
// @Failure 500 {object} response.ResponseData "内部错误"
// @Router /gam/roles [get]
func (h *RoleHandler) GetRoleList(c *gin.Context) {
	query := make(map[string]interface{})
	if name := c.Query("name"); name != "" {
		query["name"] = name
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

	// 添加其他查询选项
	// if c.Query("with_profile") == "true" {
	// 	opts = append(opts, services.WithPreload("Profile"))
	// }

	roles, total, err := h.roleService.List(query, page, pageSize, opts...)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get role list")
		return
	}

	response.Success(c, gin.H{"roles": roles, "total": total})
}

// GetRole godoc
// @Summary 获取单个角色
// @Description 根据角色ID获取角色的详细信息
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param id path int true "角色ID" minimum(1)
// @Success 200 {object} response.ResponseData{data=object{role=models.Role}} "成功"
// @Failure 400 {object} response.ResponseData "无效的角色ID"
// @Failure 404 {object} response.ResponseData "角色不存在"
// @Failure 500 {object} response.ResponseData "内部错误"
// @Router /gam/roles/{id} [get]
func (h *RoleHandler) GetRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID format")
		return
	}

	role, err := h.roleService.GetByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get role")
		return
	}

	response.Success(c, gin.H{"role": role})
}

// CreateRole godoc
// @Summary 创建角色
// @Description 创建新的角色，包括角色名称、描述等信息
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param role body models.Role true "角色信息，包括名称、描述等"
// @Success 200 {object} response.ResponseData{data=object{role=models.Role}} "成功"
// @Failure 400 {object} response.ResponseData "请求参数错误"
// @Failure 500 {object} response.ResponseData "内部错误"
// @Router /gam/roles [post]
func (h *RoleHandler) CreateRole(c *gin.Context) {
	var role models.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.roleService.Create(&role); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create role")
		return
	}

	response.Success(c, gin.H{"role": role})
}

// UpdateRole godoc
// @Summary 更新角色
// @Description 根据角色ID更新角色信息
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param id path int true "角色ID" minimum(1)
// @Param role body models.Role true "需要更新的角色信息"
// @Success 200 {object} response.ResponseData{data=object{role=models.Role}} "成功"
// @Failure 400 {object} response.ResponseData "无效的请求参数"
// @Failure 404 {object} response.ResponseData "角色不存在"
// @Failure 500 {object} response.ResponseData "内部错误"
// @Router /gam/roles/{id} [put]
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID format")
		return
	}

	var role models.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.roleService.Update(uint(id), &role); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update role")
		return
	}

	response.Success(c, gin.H{"role": role})
}

// DeleteRole godoc
// @Summary 删除角色
// @Description 根据角色ID删除指定角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param id path int true "角色ID" minimum(1)
// @Success 200 {object} response.ResponseData{data=object{message=string}} "成功"
// @Failure 400 {object} response.ResponseData "无效的角色ID"
// @Failure 404 {object} response.ResponseData "角色不存在"
// @Failure 500 {object} response.ResponseData "内部错误"
// @Router /gam/roles/{id} [delete]
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID format")
		return
	}

	if err := h.roleService.Delete(uint(id), false); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete role")
		return
	}

	response.Success(c, gin.H{"message": "Role deleted successfully"})
}

// UpdateRoleStatus godoc
// @Summary 更新角色状态
// @Description 更新指定角色的启用/禁用状态
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param id path int true "角色ID" minimum(1)
// @Param status body object true "状态信息" example({"status": 1})
// @Success 200 {object} response.ResponseData{data=object{message=string}} "成功"
// @Failure 400 {object} response.ResponseData "无效的请求参数"
// @Failure 404 {object} response.ResponseData "角色不存在"
// @Failure 500 {object} response.ResponseData "内部错误"
// @Router /gam/roles/{id}/status [put]
func (h *RoleHandler) UpdateRoleStatus(c *gin.Context) {
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

	if err := h.roleService.Update(uint(id), map[string]interface{}{"status": req.Status}); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update role status")
		return
	}

	response.Success(c, gin.H{"message": "Role status updated successfully"})
}

// UpdateRoleSort godoc
// @Summary 更新角色排序
// @Description 更新指定角色的排序值
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param id path int true "角色ID" minimum(1)
// @Param sort body object true "排序信息" example({"sort": 1})
// @Success 200 {object} response.ResponseData{data=object{message=string}} "成功"
// @Failure 400 {object} response.ResponseData "无效的请求参数"
// @Failure 404 {object} response.ResponseData "角色不存在"
// @Failure 500 {object} response.ResponseData "内部错误"
// @Router /gam/roles/{id}/sort [put]
func (h *RoleHandler) UpdateRoleSort(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID format")
		return
	}

	var req struct {
		Sort int `json:"sort"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.roleService.Update(uint(id), map[string]interface{}{"sort": req.Sort}); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update role sort")
		return
	}

	response.Success(c, gin.H{"message": "Role sort updated successfully"})
}

// CheckRoleFieldUnique godoc
// @Summary 检查角色字段唯一性
// @Description 检查指定字段的值在角色表中是否唯一
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param field query string true "字段名称" example("name")
// @Param value query string true "字段值" example("admin")
// @Param excludeId query int false "需要排除的角色ID" minimum(1)
// @Success 200 {object} response.ResponseData{data=object{unique=bool}} "成功"
// @Failure 400 {object} response.ResponseData "无效的请求参数"
// @Failure 500 {object} response.ResponseData "内部错误"
// @Router /gam/roles/check-field [get]
func (h *RoleHandler) CheckRoleFieldUnique(c *gin.Context) {
	field := c.Query("field")
	value := c.Query("value")
	excludeID, _ := strconv.ParseUint(c.Query("excludeId"), 10, 32)

	if field == "" || value == "" {
		response.Error(c, http.StatusBadRequest, "Field and value are required")
		return
	}

	isUnique, err := h.roleService.CheckRoleFieldUnique(field, value, uint(excludeID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to check field uniqueness")
		return
	}

	response.Success(c, gin.H{"unique": isUnique})
}

// GetRoleMenus godoc
// @Summary 获取角色菜单权限
// @Description 获取指定角色的菜单权限树和已选中的菜单列表
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param roleId path int true "角色ID" minimum(1)
// @Success 200 {object} response.ResponseData{data=object{menuTree=[]models.Menu,checkedMenus=[]uint}} "成功"
// @Failure 400 {object} response.ResponseData "无效的角色ID"
// @Failure 404 {object} response.ResponseData "角色不存在"
// @Failure 500 {object} response.ResponseData "内部错误"
// @Router /gam/roles/{roleId}/menus [get]
func (h *RoleHandler) GetRoleMenus(c *gin.Context) {
	roleID, err := strconv.ParseUint(c.Param("roleId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid role ID format")
		return
	}

	menuTree, checkedMenus, err := h.roleService.GetRoleMenus(uint(roleID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get role menus")
		return
	}

	response.Success(c, gin.H{"menuTree": menuTree, "checkedMenus": checkedMenus})
}

// UpdateRoleMenus godoc
// @Summary 更新角色菜单权限
// @Description 更新指定角色的菜单权限
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param roleId path int true "角色ID" minimum(1)
// @Param menuIds body object true "菜单ID列表" example({"menuIds": [1,2,3]})
// @Success 200 {object} response.ResponseData{data=object{message=string}} "成功"
// @Failure 400 {object} response.ResponseData "无效的请求参数"
// @Failure 404 {object} response.ResponseData "角色不存在"
// @Failure 500 {object} response.ResponseData "内部错误"
// @Router /gam/roles/{roleId}/menus [put]
func (h *RoleHandler) UpdateRoleMenus(c *gin.Context) {
	roleID, err := strconv.ParseUint(c.Param("roleId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid role ID format")
		return
	}

	var req struct {
		MenuIDs []uint `json:"menuIds" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.roleService.UpdateRoleMenus(uint(roleID), req.MenuIDs); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update role menus")
		return
	}

	response.Success(c, gin.H{"message": "Role menus updated successfully"})
}
