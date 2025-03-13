package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"normaladmin/backend/pkg/sysconfig"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

type QiniuStorage struct {
	mac      *qbox.Mac
	bucket   string
	domain   string
	uploader *storage.FormUploader
}

func NewQiniuStorage() Storage {
	accessKey := sysconfig.Get("qiniu_access_key", "")
	secretKey := sysconfig.Get("qiniu_secret_key", "")
	bucket := sysconfig.Get("qiniu_bucket", "")
	domain := sysconfig.Get("qiniu_domain", "")

	mac := qbox.NewMac(accessKey, secretKey)
	cfg := storage.Config{
		Zone:          &storage.ZoneHuadong, // 默认使用华东区域，可以根据需要修改
		UseCdnDomains: false,
	}

	return &QiniuStorage{
		mac:      mac,
		bucket:   bucket,
		domain:   domain,
		uploader: storage.NewFormUploader(&cfg),
	}
}

// Upload 上传文件到七牛云
func (s *QiniuStorage) Upload(path string, file io.Reader) (string, error) {
	// 生成上传凭证
	putPolicy := storage.PutPolicy{
		Scope: s.bucket,
	}
	upToken := putPolicy.UploadToken(s.mac)

	// 读取文件内容
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, file); err != nil {
		return "", fmt.Errorf("读取文件失败: %w", err)
	}

	// 上传文件
	ret := storage.PutRet{}
	err := s.uploader.Put(context.Background(), &ret, upToken, path, bytes.NewReader(buf.Bytes()), int64(buf.Len()), nil)
	if err != nil {
		return "", fmt.Errorf("上传到七牛云失败: %w", err)
	}

	// 如果配置了自定义域名，使用自定义域名
	if s.domain != "" {
		return fmt.Sprintf("%s/%s", s.domain, path), nil
	}

	// 否则返回七牛云默认域名
	return fmt.Sprintf("http://%s.qiniudn.com/%s", s.bucket, path), nil
}

// Delete 从七牛云删除文件
func (s *QiniuStorage) Delete(path string) error {
	bucketManager := storage.NewBucketManager(s.mac, nil)
	err := bucketManager.Delete(s.bucket, path)
	if err != nil {
		return fmt.Errorf("从七牛云删除文件失败: %w", err)
	}
	return nil
}

func (s *QiniuStorage) GetType() string {
	return "qiniu"
}
