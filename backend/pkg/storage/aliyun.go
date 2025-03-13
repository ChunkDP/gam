package storage

import (
	"fmt"
	"io"
	"normaladmin/backend/pkg/sysconfig"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type AliyunOSS struct {
	client *oss.Client
	bucket *oss.Bucket
	domain string
}

func NewAliyunOSS() Storage {
	accessKey := sysconfig.Get("aliyun_oss_key", "")
	secretKey := sysconfig.Get("aliyun_oss_secret", "")
	endpoint := sysconfig.Get("aliyun_oss_endpoint", "")
	bucketName := sysconfig.Get("aliyun_oss_bucket", "")
	domain := sysconfig.Get("aliyun_oss_domain", "")

	client, err := oss.New(endpoint, accessKey, secretKey)
	if err != nil {
		panic(fmt.Sprintf("初始化阿里云OSS客户端失败: %v", err))
	}

	bucket, err := client.Bucket(bucketName)
	if err != nil {
		panic(fmt.Sprintf("获取Bucket失败: %v", err))
	}

	return &AliyunOSS{
		client: client,
		bucket: bucket,
		domain: domain,
	}
}

// Upload 上传文件到阿里云OSS
func (s *AliyunOSS) Upload(path string, file io.Reader) (string, error) {
	// 上传文件
	if err := s.bucket.PutObject(path, file); err != nil {
		return "", fmt.Errorf("上传到OSS失败: %w", err)
	}

	// 如果配置了自定义域名，使用自定义域名
	if s.domain != "" {
		return fmt.Sprintf("%s/%s", s.domain, path), nil
	}

	// 否则使用默认的OSS域名
	return fmt.Sprintf("https://%s.%s/%s", s.bucket.BucketName, s.client.Config.Endpoint, path), nil
}

// Delete 从阿里云OSS删除文件
func (s *AliyunOSS) Delete(path string) error {
	if err := s.bucket.DeleteObject(path); err != nil {
		return fmt.Errorf("从OSS删除文件失败: %w", err)
	}
	return nil
}

func (s *AliyunOSS) GetType() string {
	return "aliyun"
}
