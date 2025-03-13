package models

import "gorm.io/gorm"

// UploadFile 上传文件记录
type UploadFile struct {
	gorm.Model `swaggerignore:"true"`
	FileName   string `gorm:"type:varchar(255);not null" json:"file_name"` // 文件名
	FilePath   string `gorm:"type:varchar(255);not null" json:"file_path"` // 文件路径
	FileType   string `gorm:"type:varchar(50);not null" json:"file_type"`  // 文件类型
	FileSize   int64  `gorm:"not null" json:"file_size"`                   // 文件大小
	FileExt    string `gorm:"type:varchar(50);not null" json:"file_ext"`   // 文件扩展名
	FileUrl    string `gorm:"type:varchar(255);not null" json:"file_url"`  //

}

// UploadConfig 上传配置
type UploadConfig struct {
	MaxSize      int      `json:"max_size"`      // 最大文件大小(MB)
	AllowedTypes []string `json:"allowed_types"` // 允许的文件类型
	Driver       string   `json:"driver"`        // 存储驱动
	UploadPath   string   `json:"upload_path"`   // 上传路径(本地存储时使用)
	AccessURL    string   `json:"access_url"`    // 访问URL前缀(本地存储时使用)
}
