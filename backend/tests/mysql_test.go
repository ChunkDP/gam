package tests

import (
	"normaladmin/backend/pkg/database"
	"testing"
)

func TestMySQLConnection(t *testing.T) {
	// 获取数据库连接
	db := database.GetDB()

	// 测试ping
	err := db.DB().Ping()
	if err != nil {
		t.Errorf("MySQL connection test failed: %v", err)
		return
	}

	// 执行简单查询
	var result int
	err = db.Raw("SELECT 1").Scan(&result).Error
	if err != nil {
		t.Errorf("MySQL query test failed: %v", err)
		return
	}

	if result != 1 {
		t.Errorf("Expected 1, got %d", result)
		return
	}

	t.Log("MySQL connection test passed successfully")
}
