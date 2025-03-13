package services

import (
	"fmt"
	"mime/multipart"
	"normaladmin/backend/internal/models"
	"normaladmin/backend/pkg/storage"
	"normaladmin/backend/pkg/sysconfig"
	"normaladmin/backend/pkg/utils/validator"
	"path"
	"strings"
	"time"

	"log"

	"gorm.io/gorm"
)

type UploadService interface {
	UploadFile(file *multipart.FileHeader, fileType string) (*models.UploadFile, error)
	GetUploadConfig() (*models.UploadConfig, error)
	DeleteFile(filePath string) error
	BatchUploadFiles(files []*multipart.FileHeader, fileType string) ([]map[string]interface{}, error)
}

type uploadService struct {
	db      *gorm.DB
	storage storage.Storage
}

func NewUploadService(db *gorm.DB) UploadService {
	// 根据配置初始化存储驱动
	storageDriver := sysconfig.Get("upload_driver", "local")
	var store storage.Storage

	switch storageDriver {
	case "aliyun":
		store = storage.NewAliyunOSS()
	case "tencent":
		store = storage.NewTencentCOS()
	case "qiniu":
		store = storage.NewQiniuStorage()
	default:
		store = storage.NewLocalStorage()
	}

	return &uploadService{
		db:      db,
		storage: store,
	}
}

// UploadFile 上传文件
func (s *uploadService) UploadFile(file *multipart.FileHeader, fileType string) (*models.UploadFile, error) {
	// 创建文件验证器
	validator := validator.NewFileValidator(
		sysconfig.GetInt("upload_max_size", 10),
		sysconfig.Get("upload_mime_types", ""),
	)

	// 验证文件
	if err := validator.ValidateFile(file); err != nil {
		return nil, err
	}

	// 生成文件路径
	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), strings.ToLower(path.Ext(file.Filename)))
	filePath := fmt.Sprintf("%s/%s/%s", time.Now().Format("2006/01/02"), fileType, fileName)

	// 上传文件
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	url, err := s.storage.Upload(filePath, src)
	if err != nil {
		return nil, err
	}

	// 如果是本地存储，使用配置的访问URL前缀
	if s.storage.GetType() == "local" {
		accessURL := sysconfig.Get("upload_url", "/uploads")
		url = fmt.Sprintf("%s/%s", strings.TrimRight(accessURL, "/"), filePath)
	}

	// 保存文件记录
	uploadFile := &models.UploadFile{
		FileName: file.Filename,
		FilePath: filePath,
		FileType: fileType,
		FileSize: file.Size,
		FileExt:  strings.ToLower(path.Ext(file.Filename)),
		FileUrl:  url,
	}

	if err := s.db.Create(uploadFile).Error; err != nil {
		return nil, err
	}

	return uploadFile, nil
}

// GetUploadConfig 获取上传配置
func (s *uploadService) GetUploadConfig() (*models.UploadConfig, error) {
	return &models.UploadConfig{
		MaxSize:      sysconfig.GetInt("upload_max_size", 10),
		AllowedTypes: strings.Split(sysconfig.Get("upload_mime_types", ""), ","),
		Driver:       sysconfig.Get("upload_driver", "local"),
		UploadPath:   sysconfig.Get("upload_path", "uploads"),
		AccessURL:    sysconfig.Get("upload_url", "/uploads"),
	}, nil
}

// DeleteFile 删除文件
func (s *uploadService) DeleteFile(filePath string) error {
	if err := s.storage.Delete(filePath); err != nil {
		return err
	}

	return s.db.Where("file_path = ?", filePath).Delete(&models.UploadFile{}).Error
}

// 辅助函数：检查切片是否包含某个值
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if strings.TrimSpace(s) == item {
			return true
		}
	}
	return false
}

// BatchUploadFiles 批量上传文件
func (s *uploadService) BatchUploadFiles(files []*multipart.FileHeader, fileType string) ([]map[string]interface{}, error) {
	results := make([]map[string]interface{}, 0)
	validator := validator.NewFileValidator(
		sysconfig.GetInt("upload_max_size", 10),
		sysconfig.Get("upload_mime_types", ""),
	)
	for _, file := range files {
		// 检查文件类型
		if err := validator.ValidateFile(file); err != nil {
			return nil, err
		}

		// 上传单个文件
		uploadFile, err := s.UploadFile(file, fileType)
		if err != nil {
			log.Printf("上传文件失败: %v", err)
			continue
		}

		results = append(results, map[string]interface{}{
			"id":   uploadFile.ID,
			"name": uploadFile.FileName,
			"url":  uploadFile.FileUrl,
		})
	}

	return results, nil
}
