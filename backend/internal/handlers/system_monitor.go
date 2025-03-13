package handlers

import (
	"net/http"
	"normaladmin/backend/internal/models"
	"normaladmin/backend/internal/services"
	"normaladmin/backend/pkg/utils/response"

	"github.com/gin-gonic/gin"
)

type SystemMonitorHandler struct {
	systemMonitorService services.SystemMonitorService
}

func NewSystemMonitorHandler(systemMonitorService services.SystemMonitorService) *SystemMonitorHandler {
	return &SystemMonitorHandler{
		systemMonitorService: systemMonitorService,
	}
}

// CollectSystemInfo 收集系统信息
// @Summary 收集系统信息
// @Description 收集当前系统的CPU、内存、磁盘等使用情况
// @Tags 系统监控
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=models.SystemMonitor}
// @Router /api/system/monitor/collect [post]
func (h *SystemMonitorHandler) CollectSystemInfo(c *gin.Context) {
	monitor, err := h.systemMonitorService.CollectSystemInfo()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "收集系统信息失败: "+err.Error())
		return
	}
	response.Success(c, monitor)
}

// GetSystemMonitors 获取系统监控列表
// @Summary 获取系统监控列表
// @Description 获取系统监控历史数据
// @Tags 系统监控
// @Accept json
// @Produce json
// @Param start_time query string false "开始时间"
// @Param end_time query string false "结束时间"
// @Param limit query int false "限制条数"
// @Success 200 {object} response.Response{data=[]models.SystemMonitor}
// @Router /api/system/monitor/list [get]
func (h *SystemMonitorHandler) GetSystemMonitors(c *gin.Context) {
	var query models.SystemMonitorQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, "无效的查询参数: "+err.Error())
		return
	}

	monitors, err := h.systemMonitorService.GetSystemMonitors(&query)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取系统监控列表失败: "+err.Error())
		return
	}
	response.Success(c, monitors)
}

// GetLatestSystemMonitor 获取最新的系统监控信息
// @Summary 获取最新的系统监控信息
// @Description 获取最新的系统监控数据
// @Tags 系统监控
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=models.SystemMonitor}
// @Router /api/system/monitor/latest [get]
func (h *SystemMonitorHandler) GetLatestSystemMonitor(c *gin.Context) {
	monitor, err := h.systemMonitorService.GetLatestSystemMonitor()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取最新系统监控信息失败: "+err.Error())
		return
	}
	response.Success(c, monitor)
}
