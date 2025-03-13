package handlers

import (
	"net/http"
	"normaladmin/backend/internal/models"
	"normaladmin/backend/internal/services"
	"strconv"
	"strings"

	"normaladmin/backend/pkg/utils/response"

	"github.com/gin-gonic/gin"
)

type MemberHandler struct {
	memberService services.MemberService
}

func NewMemberHandler(memberService services.MemberService) *MemberHandler {
	return &MemberHandler{memberService: memberService}
}

// GetMemberList godoc
// @Summary 获取会员列表
// @Description 获取会员列表,支持分页和筛选
// @Tags 会员管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param username query string false "用户名"
// @Param mobile query string false "手机号"
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param sortField query string false "排序字段，多个字段用逗号分隔" example("id,created_at")
// @Param sortOrder query string false "排序方式(asc/desc)，多个排序方式用逗号分隔" example("desc,asc")
// @Success 200 {object} response.ResponseData{data=object{members=[]models.Member,total=int64}} "成功"
// @Failure 401 {object} response.ResponseData "未授权"
// @Failure 500 {object} response.ResponseData "内部错误"
// @Router /gam/members [get]
func (h *MemberHandler) GetMemberList(c *gin.Context) {
	query := make(map[string]interface{})
	if username := c.Query("username"); username != "" {
		query["username"] = username
	}
	if mobile := c.Query("mobile"); mobile != "" {
		query["mobile"] = mobile
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

	members, total, err := h.memberService.List(query, page, pageSize, opts...)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get member list")
		return
	}

	response.Success(c, gin.H{"members": members, "total": total})
}

// CreateMember godoc
// @Summary 创建会员
// @Description 创建新的会员账号
// @Tags 会员管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param member body models.Member true "会员信息" example({"username":"test","mobile":"13800138000","password":"123456","status":1})
// @Success 200 {object} response.ResponseData{data=object{member=models.Member}} "成功"
// @Failure 400 {object} response.ResponseData "无效的请求参数"
// @Failure 401 {object} response.ResponseData "未授权"
// @Failure 500 {object} response.ResponseData "内部错误"
// @Router /gam/members [post]
func (h *MemberHandler) CreateMember(c *gin.Context) {
	var member models.Member
	if err := c.ShouldBindJSON(&member); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.memberService.Create(&member); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create member")
		return
	}

	response.Success(c, gin.H{"member": member})
}

// UpdateMember godoc
// @Summary 更新会员
// @Description 更新会员信息
// @Tags 会员管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "会员ID" minimum(1)
// @Param member body models.Member true "会员信息" example({"username":"test","mobile":"13800138000","status":1})
// @Success 200 {object} response.ResponseData{data=object{member=models.Member}} "成功"
// @Failure 400 {object} response.ResponseData "无效的请求参数"
// @Failure 401 {object} response.ResponseData "未授权"
// @Failure 404 {object} response.ResponseData "会员不存在"
// @Failure 500 {object} response.ResponseData "内部错误"
// @Router /gam/members/{id} [put]
func (h *MemberHandler) UpdateMember(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID format")
		return
	}

	var member models.Member
	if err := c.ShouldBindJSON(&member); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.memberService.Update(uint(id), &member); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update member")
		return
	}

	response.Success(c, gin.H{"member": member})
}

// DeleteMember godoc
// @Summary 删除会员
// @Description 根据ID删除会员
// @Tags 会员管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "会员ID" minimum(1)
// @Success 200 {object} response.ResponseData{data=object{message=string}} "成功"
// @Failure 400 {object} response.ResponseData "无效的会员ID"
// @Failure 401 {object} response.ResponseData "未授权"
// @Failure 404 {object} response.ResponseData "会员不存在"
// @Failure 500 {object} response.ResponseData "内部错误"
// @Router /gam/members/{id} [delete]
func (h *MemberHandler) DeleteMember(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID format")
		return
	}

	if err := h.memberService.Delete(uint(id), false); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete member")
		return
	}

	response.Success(c, gin.H{"message": "Member deleted successfully"})
}

// CheckMemberFieldUnique godoc
// @Summary 检查会员字段唯一性
// @Description 检查会员的某个字段值是否唯一（如用户名、手机号等）
// @Tags 会员管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param field query string true "字段名" example("username")
// @Param value query string true "字段值" example("test")
// @Param excludeId query int false "排除的ID" minimum(1)
// @Success 200 {object} response.ResponseData{data=object{unique=bool}} "成功"
// @Failure 400 {object} response.ResponseData "无效的请求参数"
// @Failure 401 {object} response.ResponseData "未授权"
// @Failure 500 {object} response.ResponseData "内部错误"
// @Router /gam/members/check-field [get]
func (h *MemberHandler) CheckMemberFieldUnique(c *gin.Context) {
	field := c.Query("field")
	value := c.Query("value")
	excludeID, _ := strconv.ParseUint(c.Query("excludeId"), 10, 32)

	if field == "" || value == "" {
		response.Error(c, http.StatusBadRequest, "Field and value are required")
		return
	}

	isUnique, err := h.memberService.CheckMemberFieldUnique(field, value, uint(excludeID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to check field uniqueness")
		return
	}

	response.Success(c, gin.H{"unique": isUnique})
}

// UpdateMemberStatus godoc
// @Summary 更新会员状态
// @Description 更新会员的启用/禁用状态（0：禁用，1：启用）
// @Tags 会员管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "会员ID" minimum(1)
// @Param status body object true "状态信息" example({"status":1})
// @Success 200 {object} response.ResponseData{data=object{message=string}} "成功"
// @Failure 400 {object} response.ResponseData "无效的请求参数"
// @Failure 401 {object} response.ResponseData "未授权"
// @Failure 404 {object} response.ResponseData "会员不存在"
// @Failure 500 {object} response.ResponseData "内部错误"
// @Router /gam/members/{id}/status [put]
func (h *MemberHandler) UpdateMemberStatus(c *gin.Context) {
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

	if err := h.memberService.Update(uint(id), map[string]interface{}{"status": req.Status}); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update member status")
		return
	}

	response.Success(c, gin.H{"message": "Member status updated successfully"})
}
