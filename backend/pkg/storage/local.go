package storage

import (
	"fmt"
	"io"
	"normaladmin/backend/pkg/sysconfig"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	uploadPath string // 上传根目录
	urlPrefix  string // 访问URL前缀
}

func NewLocalStorage() Storage {
	uploadPath := sysconfig.Get("upload_path", "uploads")
	urlPrefix := sysconfig.Get("upload_url", "/uploads")

	// 确保上传目录存在
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		panic(fmt.Sprintf("创建上传目录失败: %v", err))
	}

	return &LocalStorage{
		uploadPath: uploadPath,
		urlPrefix:  urlPrefix,
	}
}

// Upload 上传文件到本地
func (s *LocalStorage) Upload(path string, file io.Reader) (string, error) {
	fullPath := filepath.Join(s.uploadPath, path)

	// 确保目录存在
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("创建目录失败: %w", err)
	}

	// 创建文件
	dst, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("创建文件失败: %w", err)
	}
	defer dst.Close()

	// 写入文件
	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("写入文件失败: %w", err)
	}

	// 返回可访问的URL
	return filepath.Join(s.urlPrefix, path), nil
}

// Delete 删除本地文件
func (s *LocalStorage) Delete(path string) error {
	fullPath := filepath.Join(s.uploadPath, path)
	if err := os.Remove(fullPath); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("删除文件失败: %w", err)
	}
	return nil
}

func (s *LocalStorage) GetType() string {
	return "local"
}
