package models

import (
	"strconv"

	"gorm.io/gorm"
)

// QueryOption 定义查询选项的函数类型
type QueryOption func(*gorm.DB) *gorm.DB

// Paginate 分页方法
func Paginate(page, pageSize string) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		offset, limit := ParsePagination(page, pageSize)
		return db.Offset(offset).Limit(limit)
	}
}
func ParsePagination(page, pageSize string) (offset, limit int) {
	pageNum, _ := strconv.Atoi(page)
	if pageNum < 1 {
		pageNum = 1
	}

	limitNum, _ := strconv.Atoi(pageSize)
	if limitNum < 1 {
		limitNum = 10
	}

	offset = (pageNum - 1) * limitNum
	return offset, limitNum
}

// WithPreload 创建预加载查询选项
func WithPreload(relations ...string) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		for _, relation := range relations {
			db = db.Preload(relation)
		}
		return db
	}
}

// WithSort 创建排序查询选项
func WithSort(sortFields []string, sortOrders []string) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		// 如果排序字段为空或只包含空字符串，则不进行排序
		if len(sortFields) == 0 || (len(sortFields) == 1 && sortFields[0] == "") {
			return db
		}

		for i, field := range sortFields {
			// 跳过空的排序字段
			if field == "" {
				continue
			}

			order := "ASC"
			if i < len(sortOrders) && sortOrders[i] == "desc" {
				order = "DESC"
			}
			db = db.Order(field + " " + order)
		}
		return db
	}
}

// WithJoin 创建连接查询选项
func WithJoin(joins ...string) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		for _, join := range joins {
			db = db.Joins(join)
		}
		return db
	}
}
