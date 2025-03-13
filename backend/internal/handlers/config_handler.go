package handlers

import (
	"fmt"
	"net/http"
	"normaladmin/backend/internal/services"
	"strconv"

	"normaladmin/backend/pkg/utils/response"

	"github.com/gin-gonic/gin"
)

type ConfigHandler struct {
	configService services.ConfigService
}

func NewConfigHandler(configService services.ConfigService) *ConfigHandler {
	return &ConfigHandler{
		configService: configService,
	}
}

// GetConfigGroups godoc
// @Summary 获取配置组列表
// @Description 获取所有系统配置组信息，包括组名称、标识符等
// @Tags 系统配置
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.ResponseData "成功"
// @Failure 401 {object} response.ResponseData "未授权或Token无效"
// @Failure 500 {object} response.ResponseData "内部错误"
// @Router /gam/configs/groups [get]
func (h *ConfigHandler) GetConfigGroups(c *gin.Context) {
	groups, err := h.configService.GetConfigGroups()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取配置组失败")
		return
	}
	response.Success(c, gin.H{"groups": groups})
}

// GetConfigItems godoc
// @Summary 获取配置项列表
// @Description 获取指定配置组下的所有配置项，包括配置键、值、描述等
// @Tags 系统配置
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param groupId path int true "配置组ID" minimum(1)
// @Success 200 {object} response.ResponseData "成功"
// @Failure 400 {object} response.ResponseData "无效的配置组ID"
// @Failure 401 {object} response.ResponseData "未授权或Token无效"
// @Failure 404 {object} response.ResponseData "配置组不存在"
// @Failure 500 {object} response.ResponseData "内部错误"
// @Router /gam/configs/items/{groupId} [get]
func (h *ConfigHandler) GetConfigItems(c *gin.Context) {
	groupID, err := strconv.ParseUint(c.Param("groupId"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的配置组ID")
		return
	}

	items, err := h.configService.GetConfigItems(int64(groupID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取配置项失败")
		return
	}
	response.Success(c, gin.H{"items": items})
}

// UpdateConfigValue godoc
// @Summary 更新配置值
// @Description 更新单个配置项的值
// @Tags 系统配置
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body object true "配置更新信息" example({"group_id":1,"item_key":"site_name","value":"My Website"})
// @Success 200 {object} response.ResponseData{} "成功"
// @Failure 400 {object} response.ResponseData "无效的请求参数"
// @Failure 401 {object} response.ResponseData "未授权或Token无效"
// @Failure 404 {object} response.ResponseData "配置项不存在"
// @Failure 500 {object} response.ResponseData "内部错误"
// @Router /gam/configs/value [put]
func (h *ConfigHandler) UpdateConfigValue(c *gin.Context) {
	var req struct {
		GroupID int64  `json:"group_id" binding:"required"`
		ItemKey string `json:"item_key" binding:"required"`
		Value   string `json:"value"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.configService.UpdateConfigValue(req.GroupID, req.ItemKey, req.Value)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "更新配置失败")
		return
	}
	response.Success(c, nil)
}

// BatchUpdateConfigs godoc
// @Summary 批量更新配置
// @Description 批量更新指定配置组下的多个配置项值
// @Tags 系统配置
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body object true "批量配置更新信息" example({"group_id":1,"configs":{"site_name":"My Website","site_description":"My Website Description"}})
// @Success 200 {object} response.ResponseData{} "成功"
// @Failure 400 {object} response.ResponseData "无效的请求参数"
// @Failure 401 {object} response.ResponseData "未授权或Token无效"
// @Failure 404 {object} response.ResponseData "配置组不存在"
// @Failure 500 {object} response.ResponseData "内部错误"
// @Router /gam/configs/batch [put]
func (h *ConfigHandler) BatchUpdateConfigs(c *gin.Context) {
	var req struct {
		GroupID int64             `json:"group_id" binding:"required"`
		Configs map[string]string `json:"configs" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	fmt.Println(req)
	err := h.configService.BatchUpdateConfigs(req.GroupID, req.Configs)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "批量更新配置失败")
		return
	}
	response.Success(c, nil)
}
