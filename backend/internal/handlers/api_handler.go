package handlers

import (
	"net/http"
	"normaladmin/backend/internal/models"
	"normaladmin/backend/internal/services"
	"normaladmin/backend/pkg/utils/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// APIHandler API管理处理器
type APIHandler struct {
	apiService services.APIService
}

// NewAPIHandler 创建API处理器实例
func NewAPIHandler(apiService services.APIService) *APIHandler {
	return &APIHandler{
		apiService: apiService,
	}
}

// GetAPIList godoc
// @Summary 获取API列表
// @Description 获取API列表，支持分页和搜索
// @Tags API管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param name query string false "API名称，用于搜索"
// @Param group query string false "API分组"
// @Success 200 {object} response.ResponseData{data=map[string]interface{}} "成功"
// @Failure 500 {object} response.ResponseData "内部错误"
// @Router /gam/apis [get]
func (h *APIHandler) GetAPIList(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("pageSize", "10")
	// 构建查询条件
	query := make(map[string]interface{})
	if name := c.Query("name"); name != "" {
		query["name"] = name
	}
	if group := c.Query("group"); group != "" {
		query["group"] = group
	}

	// 获取API列表
	apis, total, err := h.apiService.List(query, page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取API列表失败")
		return
	}

	response.Success(c, gin.H{
		"list":  apis,
		"total": total,
	})
}

// GetAPI godoc
// @Summary 获取API详情
// @Description 获取指定ID的API详情
// @Tags API管理
// @Accept json
// @Produce json
// @Param id path int true "API ID"
// @Success 200 {object} response.ResponseData{data=models.API} "成功"
// @Failure 400 {object} response.ResponseData "无效的ID"
// @Failure 404 {object} response.ResponseData "API不存在"
// @Router /gam/apis/{id} [get]
func (h *APIHandler) GetAPI(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的API ID")
		return
	}

	api, err := h.apiService.GetByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "API不存在")
		return
	}

	response.Success(c, api)
}

// CreateAPI godoc
// @Summary 创建API
// @Description 创建新的API
// @Tags API管理
// @Accept json
// @Produce json
// @Param api body models.API true "API信息"
// @Success 200 {object} response.ResponseData{data=models.API} "成功"
// @Failure 400 {object} response.ResponseData "无效的请求参数"
// @Failure 500 {object} response.ResponseData "创建失败"
// @Router /gam/apis [post]
func (h *APIHandler) CreateAPI(c *gin.Context) {
	var api models.API
	if err := c.ShouldBindJSON(&api); err != nil {
		response.Error(c, http.StatusBadRequest, "无效的请求参数")
		return
	}

	if err := h.apiService.Create(&api); err != nil {
		response.Error(c, http.StatusInternalServerError, "创建API失败")
		return
	}

	response.Success(c, api)
}

// UpdateAPI godoc
// @Summary 更新API
// @Description 更新指定ID的API信息
// @Tags API管理
// @Accept json
// @Produce json
// @Param id path int true "API ID"
// @Param api body models.API true "API信息"
// @Success 200 {object} response.ResponseData{} "成功"
// @Failure 400 {object} response.ResponseData "无效的请求参数"
// @Failure 500 {object} response.ResponseData "更新失败"
// @Router /gam/apis/{id} [put]
func (h *APIHandler) UpdateAPI(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的API ID")
		return
	}

	var api models.API
	if err := c.ShouldBindJSON(&api); err != nil {
		response.Error(c, http.StatusBadRequest, "无效的请求参数")
		return
	}

	if err := h.apiService.Update(uint(id), &api); err != nil {
		response.Error(c, http.StatusInternalServerError, "更新API失败")
		return
	}

	response.Success(c, nil)
}

// DeleteAPI godoc
// @Summary 删除API
// @Description 删除指定ID的API
// @Tags API管理
// @Accept json
// @Produce json
// @Param id path int true "API ID"
// @Success 200 {object} response.ResponseData{} "成功"
// @Failure 400 {object} response.ResponseData "无效的ID"
// @Failure 500 {object} response.ResponseData "删除失败"
// @Router /gam/apis/{id} [delete]
func (h *APIHandler) DeleteAPI(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的API ID")
		return
	}

	if err := h.apiService.Delete(uint(id), false); err != nil {
		response.Error(c, http.StatusInternalServerError, "删除API失败")
		return
	}

	response.Success(c, nil)
}

// TestAPI godoc
// @Summary 测试API
// @Description 测试指定ID的API
// @Tags API管理
// @Accept json
// @Produce json
// @Param id path int true "API ID"
// @Param params body map[string]interface{} true "测试参数"
// @Success 200 {object} response.ResponseData{} "成功"
// @Failure 400 {object} response.ResponseData "无效的请求参数"
// @Failure 500 {object} response.ResponseData "测试失败"
// @Router /gam/apis/{id}/test [post]
func (h *APIHandler) TestAPI(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的API ID")
		return
	}

	var params map[string]interface{}
	if err := c.ShouldBindJSON(&params); err != nil {
		response.Error(c, http.StatusBadRequest, "无效的请求参数")
		return
	}

	result, err := h.apiService.TestAPI(uint(id), params)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}
