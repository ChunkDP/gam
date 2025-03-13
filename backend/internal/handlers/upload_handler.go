package handlers

import (
	"net/http"
	"normaladmin/backend/internal/services"
	"normaladmin/backend/pkg/utils/response"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	uploadService services.UploadService
}

func NewUploadHandler(uploadService services.UploadService) *UploadHandler {
	return &UploadHandler{
		uploadService: uploadService,
	}
}

// UploadFile godoc
// @Summary 通用文件上传
// @Description 上传单个文件，支持多种文件类型
// @Tags 文件管理
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "要上传的文件"
// @Success 200 {object} response.ResponseData{data=object{url=string,name=string,path=string}} "成功"
// @Failure 400 {object} response.ResponseData "获取上传文件失败或上传过程出错"
// @Failure 500 {object} response.ResponseData "内部错误"
// @Router /gam/upload/file [post]
func (h *UploadHandler) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "获取上传文件失败")
		return
	}

	uploadFile, err := h.uploadService.UploadFile(file, "file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, gin.H{
		"url":  uploadFile.FileUrl,
		"name": uploadFile.FileName,
		"path": uploadFile.FilePath,
	})
}

// UploadImage godoc
// @Summary 图片文件上传
// @Description 上传单个图片文件，仅支持图片格式（jpg、jpeg、png、gif等）
// @Tags 文件管理
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "要上传的图片文件"
// @Success 200 {object} response.ResponseData{data=object{url=string,name=string,path=string}} "成功"
// @Failure 400 {object} response.ResponseData "获取上传文件失败或非图片格式"
// @Failure 500 {object} response.ResponseData "内部错误"
// @Router /gam/upload/image [post]
func (h *UploadHandler) UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "获取上传文件失败")
		return
	}

	uploadFile, err := h.uploadService.UploadFile(file, "image")
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, gin.H{
		"url":  uploadFile.FileUrl,
		"name": uploadFile.FileName,
		"path": uploadFile.FilePath,
	})
}

// GetUploadConfig godoc
// @Summary 获取上传配置
// @Description 获取文件上传相关的配置信息，包括文件大小限制、支持的文件类型等
// @Tags 文件管理
// @Accept json
// @Produce json
// @Success 200 {object} response.ResponseData{data=object{}} "成功"
// @Failure 400 {object} response.ResponseData "获取配置失败"
// @Failure 500 {object} response.ResponseData "内部错误"
// @Router /gam/upload/config [get]
func (h *UploadHandler) GetUploadConfig(c *gin.Context) {
	config, err := h.uploadService.GetUploadConfig()
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, config)
}

// DeleteFile godoc
// @Summary 删除文件
// @Description 根据文件路径删除已上传的文件
// @Tags 文件管理
// @Accept json
// @Produce json
// @Param request body object true "文件路径信息" example({"file_path": "uploads/file/example.jpg"})
// @Success 200 {object} response.ResponseData{} "成功"
// @Failure 400 {object} response.ResponseData "参数错误或删除失败"
// @Failure 500 {object} response.ResponseData "内部错误"
// @Router /gam/upload/delete [post]
func (h *UploadHandler) DeleteFile(c *gin.Context) {
	var req struct {
		FilePath string `json:"file_path" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	if err := h.uploadService.DeleteFile(req.FilePath); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, nil)
}

// BatchUpload godoc
// @Summary 批量上传文件
// @Description 同时上传多个文件
// @Tags 文件管理
// @Accept multipart/form-data
// @Produce json
// @Param files[] formData file true "要上传的文件数组"
// @Success 200 {object} response.ResponseData{data=object{url=string,name=string,path=string}} "成功"
// @Failure 400 {object} response.ResponseData "获取上传文件失败或上传过程出错"
// @Failure 500 {object} response.ResponseData "内部错误"
// @Router /gam/upload/batch [post]
func (h *UploadHandler) BatchUpload(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		response.Error(c, http.StatusBadRequest, "获取上传文件失败")
		return
	}

	files := form.File["files[]"]
	if len(files) == 0 {
		response.Error(c, http.StatusBadRequest, "未找到上传文件")
		return
	}

	results, err := h.uploadService.BatchUploadFiles(files, "file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, results)
}
