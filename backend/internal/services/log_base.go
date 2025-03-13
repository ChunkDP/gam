package services

import (
	"context"
	"normaladmin/backend/internal/models"
	"normaladmin/backend/pkg/logger"
	"time"

	"go.uber.org/zap"
)

type LogService[T any] interface {
	BaseCRUDService[T]
	logOperation(ctx context.Context, operation string, args ...interface{})
}

// LogBaseService 日志装饰器
type LogBaseService[T any] struct {
	next        BaseCRUD[T]
	serviceName string // 服务名称，用于日志标识
}

// NewLogBaseService 创建日志服务实例
func NewLogBaseService[T any](next BaseCRUD[T], serviceName string) BaseCRUD[T] {
	return &LogBaseService[T]{
		next:        next,
		serviceName: serviceName,
	}
}

// logOperation 记录操作日志
func (s *LogBaseService[T]) logOperation(operation string, args ...interface{}) {
	// 构建基础日志字段

	fields := []zap.Field{
		logger.Field("service", s.serviceName),
		logger.Field("operation", operation),
		logger.Field("timestamp", time.Now()),
	}

	// 添加参数信息
	if len(args) > 0 {
		fields = append(fields, logger.Field("args", args))
	}

	// 记录日志
	logger.Info("操作日志", fields...)
}

func (s *LogBaseService[T]) Create(entity *T) error {
	s.logOperation("Create", entity)
	err := s.next.Create(entity)
	if err != nil {
		// 错误日志记录
		logger.Error("创建失败",
			logger.Field("service", s.serviceName),
			logger.Field("operation", "Create"),
			logger.Field("entity", entity),
			logger.Field("error", err),
		)
	}
	return err
}

// GetByID 获取单个实体（带日志）
func (s *LogBaseService[T]) GetByID(id uint, opts ...models.QueryOption) (*T, error) {
	s.logOperation("GetByID", id)
	result, err := s.next.GetByID(id, opts...)
	if err != nil {
		logger.Error("查询失败",
			logger.Field("service", s.serviceName),
			logger.Field("operation", "GetByID"),
			logger.Field("id", id),
			logger.Field("error", err),
		)
	}
	return result, err
}

// List 获取列表（带日志）
func (s *LogBaseService[T]) List(query map[string]interface{}, page, pageSize string, opts ...models.QueryOption) ([]T, int64, error) {
	s.logOperation("List", query, page, pageSize)
	result, total, err := s.next.List(query, page, pageSize, opts...)
	if err != nil {
		logger.Error("列表查询失败",
			logger.Field("service", s.serviceName),
			logger.Field("operation", "List"),
			logger.Field("query", query),
			logger.Field("page", page),
			logger.Field("pageSize", pageSize),
			logger.Field("error", err),
		)
	} else {
		logger.Info("列表查询成功",
			logger.Field("service", s.serviceName),
			logger.Field("operation", "List"),
			logger.Field("total", total),
			logger.Field("resultCount", len(result)),
		)
	}
	return result, total, err
}

// Update 更新实体（带日志）
func (s *LogBaseService[T]) Update(id uint, data interface{}) error {
	s.logOperation("Update", id, data)
	err := s.next.Update(id, data)
	if err != nil {
		logger.Error("更新失败",
			logger.Field("service", s.serviceName),
			logger.Field("operation", "Update"),
			logger.Field("id", id),
			logger.Field("data", data),
			logger.Field("error", err),
		)
	}
	return err
}

// Delete 删除实体（带日志）
func (s *LogBaseService[T]) Delete(id uint, hardDelete bool) error {
	s.logOperation("Delete", id, hardDelete)
	err := s.next.Delete(id, hardDelete)
	if err != nil {
		logger.Error("删除失败",
			logger.Field("service", s.serviceName),
			logger.Field("operation", "Delete"),
			logger.Field("id", id),
			logger.Field("hardDelete", hardDelete),
			logger.Field("error", err),
		)
	}
	return err
}

func (s *LogBaseService[T]) BatchDelete(ids []uint, hardDelete bool) error {
	s.logOperation("BatchDelete", ids, hardDelete)
	err := s.next.BatchDelete(ids, hardDelete)
	if err != nil {
		logger.Error("批量删除失败",
			logger.Field("service", s.serviceName),
			logger.Field("operation", "BatchDelete"),
			logger.Field("ids", ids),
			logger.Field("hardDelete", hardDelete),
			logger.Field("error", err),
		)
	}
	return err
}
