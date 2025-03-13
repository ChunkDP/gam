package services

import (
	"fmt"
	"normaladmin/backend/internal/models"
	"normaladmin/backend/pkg/logger"
	"time"

	"gorm.io/gorm"
)

type LogService[T any] interface {
	BaseCRUD[T]
	logOperation(operation string, args ...interface{})
}

// LogBaseService 日志装饰器
type LogBaseService[T any] struct {
	next        BaseCRUD[T]
	serviceName string // 服务名称，用于日志标识
	db          *gorm.DB
}

// NewLogBaseService 创建日志服务实例
func NewLogBaseService[T any](next BaseCRUD[T], serviceName string, db *gorm.DB) BaseCRUD[T] {
	return &LogBaseService[T]{
		next:        next,
		serviceName: serviceName,
		db:          db,
	}
}

// logOperation 记录操作日志到数据库
func (s *LogBaseService[T]) logOperation(operation string, args ...interface{}) {
	// 构建日志记录
	log := models.SystemLog{
		Module:    s.serviceName,
		Action:    operation,
		Method:    "API",
		Params:    fmt.Sprintf("%v", args),
		Status:    200, // 默认成功状态
		CreatedAt: time.Now(),
	}

	// 异步保存日志到数据库
	go func(log models.SystemLog) {
		if err := s.db.Create(&log).Error; err != nil {
			// 如果数据库记录失败，则使用 zap 记录错误
			logger.Error("保存操作日志失败",
				logger.Field("error", err),
				logger.Field("log", log),
			)
		}
	}(log)
}

// Create 创建实体（带日志）
func (s *LogBaseService[T]) Create(entity *T) error {
	startTime := time.Now()
	err := s.next.Create(entity)
	duration := time.Since(startTime).Milliseconds()

	// 记录操作日志
	status := 200
	result := "成功"
	if err != nil {
		status = 500
		result = err.Error()
	}

	// 构建并保存日志
	s.logOperationWithResult("Create", status, result, duration, entity)

	return err
}

// GetByID 获取单个实体（带日志）
func (s *LogBaseService[T]) GetByID(id uint, opts ...models.QueryOption) (*T, error) {
	startTime := time.Now()
	result, err := s.next.GetByID(id, opts...)
	duration := time.Since(startTime).Milliseconds()

	// 记录操作日志
	status := 200
	resultMsg := "成功"
	if err != nil {
		status = 500
		resultMsg = err.Error()
	}

	s.logOperationWithResult("GetByID", status, resultMsg, duration, id)
	return result, err
}

// List 获取列表（带日志）
func (s *LogBaseService[T]) List(query map[string]interface{}, page, pageSize string, opts ...models.QueryOption) ([]T, int64, error) {
	startTime := time.Now()
	result, total, err := s.next.List(query, page, pageSize, opts...)
	duration := time.Since(startTime).Milliseconds()

	// 记录操作日志
	status := 200
	resultMsg := fmt.Sprintf("成功，获取到 %d 条记录", total)
	if err != nil {
		status = 500
		resultMsg = err.Error()
	}

	s.logOperationWithResult("List", status, resultMsg, duration, query, page, pageSize)
	return result, total, err
}

// Update 更新实体（带日志）
func (s *LogBaseService[T]) Update(id uint, data interface{}) error {
	startTime := time.Now()
	err := s.next.Update(id, data)
	duration := time.Since(startTime).Milliseconds()

	// 记录操作日志
	status := 200
	result := "成功"
	if err != nil {
		status = 500
		result = err.Error()
	}

	s.logOperationWithResult("Update", status, result, duration, id, data)
	return err
}

// Delete 删除实体（带日志）
func (s *LogBaseService[T]) Delete(id uint, hardDelete bool) error {
	startTime := time.Now()
	err := s.next.Delete(id, hardDelete)
	duration := time.Since(startTime).Milliseconds()

	// 记录操作日志
	status := 200
	result := "成功"
	if err != nil {
		status = 500
		result = err.Error()
	}

	s.logOperationWithResult("Delete", status, result, duration, id, hardDelete)
	return err
}

// BatchDelete 批量删除（带日志）
func (s *LogBaseService[T]) BatchDelete(ids []uint, hardDelete bool) error {
	startTime := time.Now()
	err := s.next.BatchDelete(ids, hardDelete)
	duration := time.Since(startTime).Milliseconds()

	// 记录操作日志
	status := 200
	result := "成功"
	if err != nil {
		status = 500
		result = err.Error()
	}

	s.logOperationWithResult("BatchDelete", status, result, duration, ids, hardDelete)
	return err
}

// logOperationWithResult 记录带结果的操作日志
func (s *LogBaseService[T]) logOperationWithResult(operation string, status int, result string, duration int64, args ...interface{}) {
	// 构建日志记录
	log := models.SystemLog{
		Module:    s.serviceName,
		Action:    operation,
		Method:    "API",
		Params:    fmt.Sprintf("%v", args),
		Result:    result,
		Status:    status,
		Duration:  duration,
		CreatedAt: time.Now(),
	}

	// 异步保存日志到数据库
	go func(log models.SystemLog) {
		if err := s.db.Create(&log).Error; err != nil {
			// 如果数据库记录失败，则使用 zap 记录错误
			logger.Error("保存操作日志失败",
				logger.Field("error", err),
				logger.Field("log", log),
			)
		}
	}(log)
}
