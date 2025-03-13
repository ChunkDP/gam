package sysconfig

// Get 获取字符串配置
func Get(key string, defaultValue string) string {
	return GetInstance().Get(key, defaultValue)
}

// GetInt 获取整数配置
func GetInt(key string, defaultValue int) int {
	return GetInstance().GetInt(key, defaultValue)
}

// GetBool 获取布尔配置
func GetBool(key string, defaultValue bool) bool {
	return GetInstance().GetBool(key, defaultValue)
}

// GetFloat 获取浮点数配置
func GetFloat(key string, defaultValue float64) float64 {
	return GetInstance().GetFloat(key, defaultValue)
}

// GetArray 获取数组配置
func GetArray(key string, defaultValue []string) []string {
	return GetInstance().GetArray(key, defaultValue)
}

// Set 设置配置
func Set(key string, value interface{}) error {
	return GetInstance().Set(key, value)
}

// Reload 重新加载配置
func Reload() error {
	return GetInstance().Reload()
}
