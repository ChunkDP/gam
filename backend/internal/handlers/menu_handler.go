package handlers

import (
	"fmt"
	"net/http"
	"normaladmin/backend/database"
	"normaladmin/backend/internal/models"
	"normaladmin/backend/pkg/cache"
	"normaladmin/backend/pkg/utils/response"
	"strconv"

	"normaladmin/backend/internal/services"

	"github.com/gin-gonic/gin"
)

type MenuHandler struct {
	menuService services.MenuService
}

func NewMenuHandler(menuService services.MenuService) *MenuHandler {
	return &MenuHandler{menuService: menuService}
}

// GetMenuTree godoc
// @Summary 获取菜单树
// @Description 获取完整的菜单层级结构树
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.ResponseData{data=object{menus=[]models.Menu}} "成功"
// @Failure 401 {object} response.ResponseData "未授权或Token无效"
// @Failure 500 {object} response.ResponseData "内部错误"
// @Router /gam/menus/tree [get]
func (h *MenuHandler) GetMenuTree(c *gin.Context) {
	menus, err := h.menuService.GetMenuTree()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get menu tree")
		return
	}
	response.Success(c, gin.H{"menuTree": menus})
}

// GetMenuList godoc
// @Summary 获取菜单列表
// @Description 获取菜单列表，支持按名称搜索和排序
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param name query string false "菜单名称，用于搜索过滤"
// @Param sortField query string false "排序字段" example("sort")
// @Param sortOrder query string false "排序方式(asc/desc)" example("asc")
// @Success 200 {object} response.ResponseData{data=object{menus=[]models.Menu}} "成功"
// @Failure 401 {object} response.ResponseData "未授权或Token无效"
// @Failure 500 {object} response.ResponseData "内部错误"
// @Router /gam/menus [get]
func (h *MenuHandler) GetMenuList(c *gin.Context) {
	query := make(map[string]interface{})
	if title := c.Query("title"); title != "" {
		query["title"] = title
	}
	if name := c.Query("name"); name != "" {
		query["name"] = name
	}

	menus, err := h.menuService.GetMenus(query)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get menus")
		return
	}

	response.Success(c, gin.H{"menus": menus})
}

// CreateMenu godoc
// @Summary 创建菜单
// @Description 创建新的菜单项，支持设置父级菜单
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param menu body models.Menu true "菜单信息" example({"name":"系统管理","path":"/system","component":"Layout","sort":1,"parent_id":0,"meta":{"title":"系统管理","icon":"setting"}})
// @Success 200 {object} response.ResponseData{data=object{menu=models.Menu}} "成功"
// @Failure 400 {object} response.ResponseData "无效的请求参数"
// @Failure 401 {object} response.ResponseData "未授权或Token无效"
// @Failure 500 {object} response.ResponseData "内部错误"
// @Router /gam/menus [post]
func (h *MenuHandler) CreateMenu(c *gin.Context) {
	var menu models.Menu
	if err := c.ShouldBindJSON(&menu); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.menuService.CreateMenu(&menu); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create menu")
		return
	}

	response.Success(c, gin.H{"menu": menu})
}

// UpdateMenu godoc
// @Summary 更新菜单
// @Description 更新现有菜单的信息，包括名称、路径、组件等
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "菜单ID" minimum(1)
// @Param menu body models.Menu true "菜单信息" example({"name":"系统设置","path":"/settings","component":"Layout","sort":1,"meta":{"title":"系统设置","icon":"setting"}})
// @Success 200 {object} response.ResponseData{data=object{menu=models.Menu}} "成功"
// @Failure 400 {object} response.ResponseData "无效的请求参数"
// @Failure 401 {object} response.ResponseData "未授权或Token无效"
// @Failure 404 {object} response.ResponseData "菜单不存在"
// @Failure 500 {object} response.ResponseData "内部错误"
// @Router /gam/menus/{id} [put]
func (h *MenuHandler) UpdateMenu(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID format")
		return
	}

	var menu models.Menu
	if err := c.ShouldBindJSON(&menu); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.menuService.UpdateMenu(uint(id), &menu); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update menu")
		return
	}

	response.Success(c, gin.H{"menu": menu})
}

// DeleteMenu godoc
// @Summary 删除菜单
// @Description 删除指定的菜单项，如果存在子菜单则不允许删除
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "菜单ID" minimum(1)
// @Success 200 {object} response.ResponseData{data=object{message=string}} "成功"
// @Failure 400 {object} response.ResponseData "无效的菜单ID"
// @Failure 401 {object} response.ResponseData "未授权或Token无效"
// @Failure 404 {object} response.ResponseData "菜单不存在"
// @Failure 409 {object} response.ResponseData "存在子菜单，无法删除"
// @Failure 500 {object} response.ResponseData "内部错误"
// @Router /gam/menus/{id} [delete]
func (h *MenuHandler) DeleteMenu(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID format")
		return
	}

	if err := h.menuService.DeleteMenu(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete menu")
		return
	}

	response.Success(c, gin.H{"message": "Menu deleted successfully"})
}

// UpdateMenuStatus godoc
// @Summary 更新菜单状态
// @Description 更新菜单的启用/禁用状态（0：禁用，1：启用）
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "菜单ID" minimum(1)
// @Param status body object true "状态信息" example({"status":1})
// @Success 200 {object} response.ResponseData{data=object{message=string}} "成功"
// @Failure 400 {object} response.ResponseData "无效的请求参数"
// @Failure 401 {object} response.ResponseData "未授权或Token无效"
// @Failure 404 {object} response.ResponseData "菜单不存在"
// @Failure 500 {object} response.ResponseData "内部错误"
// @Router /gam/menus/{id}/status [put]
func (h *MenuHandler) UpdateMenuStatus(c *gin.Context) {
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

	if err := h.menuService.Update(uint(id), map[string]interface{}{"status": req.Status}); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update menu status")
		return
	}

	response.Success(c, gin.H{"message": "Menu status updated successfully"})
}

// UpdateMenuSort godoc
// @Summary 更新菜单排序
// @Description 更新菜单的显示顺序
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "菜单ID" minimum(1)
// @Param sort body object true "排序信息" example({"sort":1})
// @Success 200 {object} response.ResponseData{data=object{message=string}} "成功"
// @Failure 400 {object} response.ResponseData "无效的请求参数"
// @Failure 401 {object} response.ResponseData "未授权或Token无效"
// @Failure 404 {object} response.ResponseData "菜单不存在"
// @Failure 500 {object} response.ResponseData "内部错误"
// @Router /gam/menus/{id}/sort [put]
func (h *MenuHandler) UpdateMenuSort(c *gin.Context) {
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

	if err := h.menuService.Update(uint(id), map[string]interface{}{"sort": req.Sort}); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update menu sort")
		return
	}

	response.Success(c, gin.H{"message": "Menu sort updated successfully"})
}

// UpdateMenuHidden 更新菜单显示状态
// @Summary 更新菜单显示状态
// @Description 更新菜单的显示/隐藏状态
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Param id path int true "菜单ID"
// @Param hidden body object true "显示状态信息"
// @Success 200 {object} response.ResponseData "成功"
// @Failure 400 {object} response.ResponseData "请求错误"
// @Failure 500 {object} response.ResponseData "内部错误"
// @Router /gam/menus/{id}/hidden [put]
func (h *MenuHandler) UpdateMenuHidden(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID format")
		return
	}

	var req struct {
		IsHidden bool `json:"is_hidden"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.menuService.Update(uint(id), map[string]interface{}{"is_hidden": req.IsHidden}); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update menu hidden status")
		return
	}

	response.Success(c, gin.H{"message": "Menu hidden status updated successfully"})
}

// GetUserMenus 获取用户菜单
// @Summary 获取用户菜单
// @Description 获取当前用户的菜单列表
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Success 200 {object} response.ResponseData{data=[]models.Menu} "成功"
// @Failure 500 {object} response.ResponseData "内部错误"
// @Router /api/user/menus [get]
func GetUserMenus(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		return
	}
	cacheKey := fmt.Sprintf("user:menu:%d", userID)

	// 尝试从缓存获取
	var menus []models.Menu
	if err := cache.GetObject(cacheKey, &menus); err == nil {
		response.Success(c, gin.H{"menus": menus})
		return
	}

	// 缓存未命中，从数据库获取
	if err := database.GetDB().Where("status = ?", 1).Find(&menus).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get menus")
		return
	}

	// 设置缓存
	if err := cache.SetObject(cacheKey, menus); err != nil {
		// 这里只记录日志，不影响返回结果
	}

	response.Success(c, gin.H{"menus": menus})
}
