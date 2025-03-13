package storage

import "io"

// Storage 存储接口
type Storage interface {
	Upload(path string, file io.Reader) (string, error)
	Delete(path string) error
	GetType() string // 获取存储类型
}
