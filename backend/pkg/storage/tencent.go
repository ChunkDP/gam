package storage

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"normaladmin/backend/pkg/sysconfig"

	"github.com/tencentyun/cos-go-sdk-v5"
)

type TencentCOS struct {
	client *cos.Client
	domain string
}

func NewTencentCOS() Storage {
	secretId := sysconfig.Get("tencent_cos_key", "")
	secretKey := sysconfig.Get("tencent_cos_secret", "")
	region := sysconfig.Get("tencent_cos_region", "")
	bucketName := sysconfig.Get("tencent_cos_bucket", "")
	domain := sysconfig.Get("tencent_cos_domain", "")

	// 构建 COS 服务 URL
	bucketURL, err := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com", bucketName, region))
	if err != nil {
		panic(fmt.Sprintf("解析腾讯云COS URL失败: %v", err))
	}

	// 创建 COS 客户端
	client := cos.NewClient(&cos.BaseURL{BucketURL: bucketURL}, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  secretId,
			SecretKey: secretKey,
		},
	})

	return &TencentCOS{
		client: client,
		domain: domain,
	}
}

// Upload 上传文件到腾讯云COS
func (s *TencentCOS) Upload(path string, file io.Reader) (string, error) {
	// 上传文件
	_, err := s.client.Object.Put(context.Background(), path, file, nil)
	if err != nil {
		return "", fmt.Errorf("上传到COS失败: %w", err)
	}

	// 如果配置了自定义域名，使用自定义域名
	if s.domain != "" {
		return fmt.Sprintf("%s/%s", s.domain, path), nil
	}

	// 否则使用默认的COS域名
	return fmt.Sprintf("%s/%s", s.client.BaseURL.BucketURL.String(), path), nil
}

// Delete 从腾讯云COS删除文件
func (s *TencentCOS) Delete(path string) error {
	_, err := s.client.Object.Delete(context.Background(), path)
	if err != nil {
		return fmt.Errorf("从COS删除文件失败: %w", err)
	}
	return nil
}

func (s *TencentCOS) GetType() string {
	return "tencent"
}
