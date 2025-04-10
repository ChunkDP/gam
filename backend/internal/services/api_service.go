package services

import (
	"encoding/json"
	"fmt"
	"normaladmin/backend/internal/models"

	"gorm.io/gorm"
)

// APIService 定义API管理相关的接口
type APIService interface {
	// 继承基础CRUD服务
	BaseCRUD[models.API]
	// 测试API
	TestAPI(id uint, params map[string]interface{}) (map[string]interface{}, error)
}

type apiService struct {
	BaseCRUD[models.API]
	db *gorm.DB
}

// NewAPIService 创建API服务实例
func NewAPIService(db *gorm.DB, base BaseCRUD[models.API]) APIService {

	return &apiService{
		BaseCRUD: base, // 使用装饰后的服务
		db:       db,   // 保存db实例
	}
}

// TestAPI 测试指定的API
func (s *apiService) TestAPI(id uint, params map[string]interface{}) (map[string]interface{}, error) {
	// 获取API信息
	api, err := s.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("API不存在: %v", err)
	}

	// 检查API状态
	if api.Status != 1 {
		return nil, fmt.Errorf("API已禁用")
	}

	// 验证参数
	if err := validateParams(api.Parameters, params); err != nil {
		return nil, fmt.Errorf("参数验证失败: %v", err)
	}

	// 构造测试响应
	var response map[string]interface{}
	if api.Response != "" {
		if err := json.Unmarshal([]byte(api.Response), &response); err != nil {
			return nil, fmt.Errorf("响应示例格式错误: %v", err)
		}
	} else {
		response = map[string]interface{}{
			"message": "API测试成功",
			"params":  params,
		}
	}

	return response, nil
}

// validateParams 验证请求参数
func validateParams(schemaStr string, params map[string]interface{}) error {
	if schemaStr == "" {
		return nil
	}

	var schema map[string]interface{}
	if err := json.Unmarshal([]byte(schemaStr), &schema); err != nil {
		return fmt.Errorf("参数模式格式错误: %v", err)
	}

	// 验证必需参数
	for key, val := range schema {
		if paramSchema, ok := val.(map[string]interface{}); ok {
			if required, ok := paramSchema["required"].(bool); ok && required {
				if _, exists := params[key]; !exists {
					return fmt.Errorf("缺少必需参数: %s", key)
				}
			}
		}
	}

	return nil
}
