package services

import (
	"fmt"
	"normaladmin/backend/internal/models"
	"normaladmin/backend/pkg/cache"
	"sync"
	"time"
)

// CacheBaseService 缓存装饰器
type CacheBaseService[T any] struct {
	next        BaseCRUD[T]
	cachePrefix string        // 缓存键前缀
	expiration  time.Duration // 缓存过期时间
	mutex       sync.Map      // 用于防止缓存击穿的互斥锁
	metrics     *CacheMetrics // 缓存统计
}

// CacheMetrics 缓存统计指标
type CacheMetrics struct {
	hits   int64      // 缓存命中次数
	misses int64      // 缓存未命中次数
	errors int64      // 错误次数
	mutex  sync.Mutex // 用于统计的互斥锁
}

// NewCacheBaseService 创建缓存服务实例
func NewCacheBaseService[T any](next BaseCRUD[T], prefix string, expiration ...time.Duration) BaseCRUD[T] {
	exp := 24 * time.Hour // 默认24小时
	if len(expiration) > 0 {
		exp = expiration[0]
	}
	return &CacheBaseService[T]{
		next:        next,
		cachePrefix: prefix,
		expiration:  exp,
		metrics:     &CacheMetrics{},
	}
}

// generateKey 生成缓存键
func (s *CacheBaseService[T]) generateKey(id uint) string {
	return fmt.Sprintf("%s:%d", s.cachePrefix, id)
}
func (s *CacheBaseService[T]) Create(entity *T) error {
	return s.next.Create(entity)
}

// List 方法直接调用下一个服务，不使用缓存
func (s *CacheBaseService[T]) List(query map[string]interface{}, page, pageSize string, opts ...models.QueryOption) ([]T, int64, error) {
	// 直接调用下一个服务的 List 方法
	return s.next.List(query, page, pageSize, opts...)
}

// GetByID 带防护机制的获取方法
func (s *CacheBaseService[T]) GetByID(id uint, opts ...models.QueryOption) (*T, error) {
	cacheKey := s.generateKey(id)

	// 1. 尝试从缓存获取
	var entity T
	err := cache.GetObject(cacheKey, &entity)
	if err == nil {
		s.metrics.recordHit()
		return &entity, nil
	}
	s.metrics.recordMiss()

	// 2. 防止缓存击穿：使用互斥锁
	mutex, _ := s.mutex.LoadOrStore(cacheKey, &sync.Mutex{})
	lock := mutex.(*sync.Mutex)
	lock.Lock()
	defer func() {
		lock.Unlock()
		s.mutex.Delete(cacheKey) // 使用完后删除锁
	}()

	// 双重检查，可能其他协程已经加载了缓存
	err = cache.GetObject(cacheKey, &entity)
	if err == nil {
		s.metrics.recordHit()
		return &entity, nil
	}

	// 3. 从数据库获取
	entity_ptr, err := s.next.GetByID(id, opts...)
	if err != nil {
		s.metrics.recordError()
		return nil, err
	}

	// 4. 防止缓存穿透：对空值也进行缓存，但过期时间较短
	if entity_ptr == nil {
		err = cache.SetObject(cacheKey, nil, 5*time.Minute) // 空值缓存5分钟
		s.metrics.recordError()
		return nil, nil
	}

	// 5. 设置缓存
	if err := cache.SetObject(cacheKey, entity_ptr, s.expiration); err != nil {
		s.metrics.recordError()
	}

	return entity_ptr, nil
}

// Update 更新实体（更新缓存）
func (s *CacheBaseService[T]) Update(id uint, data interface{}) error {
	if err := s.next.Update(id, data); err != nil {
		return err
	}

	// 删除缓存
	cacheKey := s.generateKey(id)
	if err := cache.Delete(cacheKey); err != nil {
		// TODO: 添加日志记录
	}

	return nil
}

// Delete 删除实体（删除缓存）
func (s *CacheBaseService[T]) Delete(id uint, hardDelete bool) error {
	if err := s.next.Delete(id, hardDelete); err != nil {
		return err
	}

	// 删除缓存
	cacheKey := s.generateKey(id)
	if err := cache.Delete(cacheKey); err != nil {
		// TODO: 添加日志记录
	}

	return nil
}

// BatchDelete 批量删除（批量删除缓存）
func (s *CacheBaseService[T]) BatchDelete(ids []uint, hardDelete bool) error {
	if err := s.next.BatchDelete(ids, hardDelete); err != nil {
		return err
	}

	// 批量删除缓存
	for _, id := range ids {
		cacheKey := s.generateKey(id)
		if err := cache.Delete(cacheKey); err != nil {
			// TODO: 添加日志记录
		}
	}

	return nil
}

// PrewarmCache 缓存预热
func (s *CacheBaseService[T]) PrewarmCache(ids []uint) error {
	for _, id := range ids {
		entity, err := s.next.GetByID(id)
		if err != nil {
			continue
		}
		cacheKey := s.generateKey(id)
		if err := cache.SetObject(cacheKey, entity, s.expiration); err != nil {
			s.metrics.recordError()
		}
	}
	return nil
}

// GetMetrics 获取缓存统计信息
func (s *CacheBaseService[T]) GetMetrics() map[string]interface{} {
	total := s.metrics.hits + s.metrics.misses
	hitRate := float64(0)
	if total > 0 {
		hitRate = float64(s.metrics.hits) / float64(total) * 100
	}

	return map[string]interface{}{
		"hits":     s.metrics.hits,
		"misses":   s.metrics.misses,
		"errors":   s.metrics.errors,
		"hit_rate": fmt.Sprintf("%.2f%%", hitRate),
	}
}

// 记录统计信息的方法
func (m *CacheMetrics) recordHit() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.hits++
}

func (m *CacheMetrics) recordMiss() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.misses++
}

func (m *CacheMetrics) recordError() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.errors++
}
