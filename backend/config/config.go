package config

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// ConfigChangeCallback 配置变更回调函数类型
type ConfigChangeCallback func(key string, value interface{})

var (
	// Global 全局配置
	Global *Config
	v      *viper.Viper
	// 配置变更回调函数map
	callbacks = make(map[string]ConfigChangeCallback)
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
	Redis    RedisConfig    `yaml:"redis"`
	CORS     CORSConfig     `yaml:"cors"`
	Upload   UploadConfig   `yaml:"upload"`
	Log      LogConfig      `yaml:"log"`
	Security SecurityConfig `yaml:"security"`
}

type ServerConfig struct {
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"`
}

type DatabaseConfig struct {
	Driver          string `yaml:"driver"`
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	DBName          string `yaml:"dbname"`
	MaxIdleConns    int    `yaml:"max_idle_conns" mapstructure:"max_idle_conns"`
	MaxOpenConns    int    `yaml:"max_open_conns" mapstructure:"max_open_conns"`
	ConnMaxLifetime int    `yaml:"conn_max_lifetime" mapstructure:"conn_max_lifetime"`
}
type SecurityConfig struct {
	EncryptKey string `yaml:"encrypt_key" mapstructure:"encrypt_key"`
}
type JWTConfig struct {
	SecretKey  string `yaml:"secret_key" mapstructure:"secret_key"`
	ExpireTime int    `yaml:"expire_time" mapstructure:"expire_time"`
	Issuer     string `yaml:"issuer" mapstructure:"issuer"`
}

type RedisConfig struct {
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	Password    string `yaml:"password"`
	DB          int    `yaml:"db"`
	DefaultTTL  int    `yaml:"default_ttl" mapstructure:"default_ttl"`   // 改为 int
	LockTimeout int    `yaml:"lock_timeout" mapstructure:"lock_timeout"` // 改为 int
}

type CORSConfig struct {
	AllowedOrigins []string `yaml:"allowed_origins" mapstructure:"allowed_origins"`
	AllowedMethods []string `yaml:"allowed_methods" mapstructure:"allowed_methods"`
	AllowedHeaders []string `yaml:"allowed_headers" mapstructure:"allowed_headers"`
}

type UploadConfig struct {
	SavePath     string   `yaml:"save_path" mapstructure:"save_path"`
	AllowedTypes []string `yaml:"allowed_types" mapstructure:"allowed_types"`
	MaxSize      int64    `yaml:"max_size" mapstructure:"max_size"`
}

type LogConfig struct {
	Level      string `yaml:"level" mapstructure:"level"`
	Filename   string `yaml:"filename" mapstructure:"filename"`
	MaxSize    int    `yaml:"max_size" mapstructure:"max_size"`
	MaxAge     int    `yaml:"max_age" mapstructure:"max_age"`
	MaxBackups int    `yaml:"max_backups" mapstructure:"max_backups"`
	Compress   bool   `yaml:"compress" mapstructure:"compress"`
}

// InitConfig 初始化配置
func InitConfig(env string) error {
	v = viper.New()
	v.SetConfigType("yaml")

	// 获取当前文件的绝对路径
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("获取当前文件路径失败")
	}
	// 获取配置文件的目录
	configDir := filepath.Join(filepath.Dir(currentFile), "environments")

	// 设置配置文件路径
	v.AddConfigPath(configDir)

	// 加载基础配置
	v.SetConfigName("config")
	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("读取基础配置失败: %w", err)
	}

	// 加载环境配置
	if env != "" {
		envConfigFile := filepath.Join(configDir, "config."+env+".yaml")
		v.SetConfigFile(envConfigFile)
		if err := v.MergeInConfig(); err != nil {
			return fmt.Errorf("读取环境配置失败: %w", err)
		}
	}

	// 支持环境变量覆盖
	v.AutomaticEnv()
	v.SetEnvPrefix("APP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 解析配置到结构体
	Global = &Config{}
	if err := v.Unmarshal(Global); err != nil {
		return fmt.Errorf("解析配置失败: %w", err)
	}

	return nil
}

// GetString 获取字符串配置
func GetString(key string) string {
	return v.GetString(key)
}

// GetInt 获取整数配置
func GetInt(key string) int {
	return v.GetInt(key)
}

// GetBool 获取布尔配置
func GetBool(key string) bool {
	return v.GetBool(key)
}

// GetDuration 获取时间间隔配置
func GetDuration(key string) time.Duration {
	return v.GetDuration(key)
}

// GetStringSlice 获取字符串切片配置
func GetStringSlice(key string) []string {
	return v.GetStringSlice(key)
}

// GetStringMap 获取字符串映射配置
func GetStringMap(key string) map[string]interface{} {
	return v.GetStringMap(key)
}

// IsSet 检查配置是否存在
func IsSet(key string) bool {
	return v.IsSet(key)
}

// AllSettings 获取所有配置
func AllSettings() map[string]interface{} {
	return v.AllSettings()
}

// GetDSN 获取数据库连接字符串
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.DBName,
	)
}
