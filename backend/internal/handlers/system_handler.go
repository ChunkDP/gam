package handlers

import (
	"net/http"
	"normaladmin/backend/internal/services"
	"normaladmin/backend/pkg/utils/response"
	"time"

	"github.com/gin-gonic/gin"
)

type SystemHandler struct {
	systemService services.SystemService
}

func NewSystemHandler(systemService services.SystemService) *SystemHandler {
	return &SystemHandler{
		systemService: systemService,
	}
}

// GetSystemLogs godoc
// @Summary 获取系统日志列表
// @Description 分页获取系统操作日志
// @Tags 系统管理
// @Accept json
// @Produce json
// @Param page query string false "页码"
// @Param page_size query string false "每页数量"
// @Param username query string false "用户名"
// @Param module query string false "模块"
// @Param start_time query string false "开始时间"
// @Param end_time query string false "结束时间"
// @Success 200 {object} response.ResponseData{data=[]models.SystemLog} "成功"
// @Router /gam/system/logs [get]
func (h *SystemHandler) GetSystemLogs(c *gin.Context) {
	query := make(map[string]interface{})
	if username := c.Query("username"); username != "" {
		query["username"] = username
	}
	if module := c.Query("module"); module != "" {
		query["module"] = module
	}

	logs, total, err := h.systemService.GetLogList(query, c.Query("page"), c.Query("page_size"))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取日志失败")
		return
	}

	response.Success(c, gin.H{
		"list":  logs,
		"total": total,
	})
}

// GetSystemMonitor godoc
// @Summary 获取系统监控数据
// @Description 获取系统资源使用情况
// @Tags 系统管理
// @Accept json
// @Produce json
// @Success 200 {object} response.ResponseData{data=models.SystemMonitor} "成功"
// @Router /gam/system/monitor [get]
func (h *SystemHandler) GetSystemMonitor(c *gin.Context) {
	data, err := h.systemService.CollectSystemInfo()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取系统监控数据失败")
		return
	}
	response.Success(c, data)
}

// DeleteSystemLogs godoc
// @Summary 删除系统日志
// @Description 删除指定时间之前的系统日志
// @Tags 系统管理
// @Accept json
// @Produce json
// @Param before query string true "删除此时间之前的日志"
// @Success 200 {object} response.ResponseData{} "成功"
// @Router /gam/system/logs [delete]
func (h *SystemHandler) DeleteSystemLogs(c *gin.Context) {
	before := c.Query("before")
	if before == "" {
		response.Error(c, http.StatusBadRequest, "缺少时间参数")
		return
	}

	beforeTime, err := time.Parse(time.RFC3339, before)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "时间格式错误")
		return
	}

	if err := h.systemService.DeleteLogs(beforeTime); err != nil {
		response.Error(c, http.StatusInternalServerError, "删除日志失败")
		return
	}

	response.Success(c, nil)
}

// GetMonitorHistory godoc
// @Summary 获取监控历史数据
// @Description 获取指定时间范围内的监控历史数据
// @Tags 系统管理
// @Accept json
// @Produce json
// @Param duration query string true "时间范围，如1h, 24h, 7d"
// @Success 200 {object} response.ResponseData{data=[]models.SystemMonitor} "成功"
// @Router /gam/system/monitor/history [get]
func (h *SystemHandler) GetMonitorHistory(c *gin.Context) {
	duration := c.Query("duration")
	if duration == "" {
		response.Error(c, http.StatusBadRequest, "缺少时间范围参数")
		return
	}

	data, err := h.systemService.GetMonitorData(duration)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取监控历史数据失败")
		return
	}

	response.Success(c, data)
}
