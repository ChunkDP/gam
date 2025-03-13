package validator

import (
	"fmt"
	"mime/multipart"
	"path"
	"strings"
)

// FileValidator 文件验证器
type FileValidator struct {
	MaxSize      int64    // 最大文件大小(字节)
	AllowedTypes []string // 允许的文件类型
}

// NewFileValidator 创建文件验证器
func NewFileValidator(maxSizeMB int, allowedTypes string) *FileValidator {
	return &FileValidator{
		MaxSize:      int64(maxSizeMB) * 1024 * 1024,
		AllowedTypes: strings.Split(allowedTypes, ","),
	}
}

// ValidateFile 验证文件
func (v *FileValidator) ValidateFile(file *multipart.FileHeader) error {
	// 验证文件大小
	if file.Size > v.MaxSize {
		return fmt.Errorf("文件大小超过限制: %d MB", v.MaxSize/1024/1024)
	}

	// 验证文件类型
	ext := strings.ToLower(strings.TrimPrefix(path.Ext(file.Filename), "."))
	if !v.isAllowedType(ext) {
		return fmt.Errorf("不支持的文件类型: %s", ext)
	}

	return nil
}

// isAllowedType 检查文件类型是否允许
func (v *FileValidator) isAllowedType(fileType string) bool {
	for _, t := range v.AllowedTypes {
		if strings.TrimSpace(t) == fileType {
			return true
		}
	}
	return false
}
