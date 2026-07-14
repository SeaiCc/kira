package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config 存储应用程序配置
type Config struct {
	// Server
	Port int `mapstructure:"PORT"`

	// Database
	DatabaseURL string `mapstructure:"DATABASE_URL"`

	// JWT
	JWTSecret      string `mapstructure:"JWT_SECRET"`
	JWTAlgorithm   string `mapstructure:"JWT_ALGORITHM"`
	JWTExpireHours int    `mapstructure:"JWT_EXPIRE_HOURS"`

	// CORS
	CORSOrigins []string `mapstructure:"CORS_ORIGINS"`

	// GitHub OAuth
	GitHubClientID     string `mapstructure:"GITHUB_CLIENT_ID"`
	GitHubClientSecret string `mapstructure:"GITHUB_CLIENT_SECRET"`

	// 存储类型: local(本地) 或 oss
	StorageType string `mapstructure:"STORAGE_TYPE"`

	// Aliyun OSS
	OSSAccessKeyID     string `mapstructure:"OSS_ACCESS_KEY_ID"`
	OSSAccessKeySecret string `mapstructure:"OSS_ACCESS_KEY_SECRET"`
	OSSBucket          string `mapstructure:"OSS_BUCKET"`
	OSSDomain          string `mapstructure:"OSS_DOMAIN"`
	OSSRegion          string `mapstructure:"OSS_REGION"`
	OSSPrefix          string `mapstructure:"OSS_PREFIX"`
}

// C 全局配置实例
var C *Config

// Load 加载配置文件到全局变量
func Load() {
	_ = godotenv.Load()

	C = &Config{
		// Server
		Port: getEnvInt("PORT", 8000),

		// Database
		DatabaseURL: getEnv("DATABASE_URL", "postgresql://postgres:123456@localhost:5432/kirameku_go?sslmode=disable&timezone=Asia/Shanghai"),

		// JWT
		JWTSecret:      getEnv("JWT_SECRET", "your-secret-key-here"),
		JWTAlgorithm:   getEnv("JWT_ALGORITHM", "HS256"),
		JWTExpireHours: getEnvInt("JWT_EXPIRE_HOURS", 72),

		// CORS
		CORSOrigins: strings.Split(getEnv("CORS_ORIGINS", "http://localhost:3000,http://localhost:5173"), ","),

		// GitHub OAuth
		GitHubClientID:     getEnv("GITHUB_CLIENT_ID", ""),
		GitHubClientSecret: getEnv("GITHUB_CLIENT_SECRET", ""),

		// Storage
		StorageType: getEnv("STORAGE_TYPE", "local"),

		// Aliyun OSS
		OSSAccessKeyID:     getEnv("OSS_ACCESS_KEY_ID", ""),
		OSSAccessKeySecret: getEnv("OSS_ACCESS_KEY_SECRET", ""),
		OSSBucket:          getEnv("OSS_BUCKET", ""),
		OSSDomain:          getEnv("OSS_DOMAIN", ""),
		OSSRegion:          getEnv("OSS_REGION", "oss-cn-hangzhou"),
		OSSPrefix:          getEnv("OSS_PREFIX", "uploads/"),
	}
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt 获取整型环境变量
func getEnvInt(key string, defaultValue int) int {
	value := getEnv(key, "")
	if value == "" {
		return defaultValue
	}

	result, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return result
}
