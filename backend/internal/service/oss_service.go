package service

import (
	"bytes"
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/google/uuid"
	"golang.org/x/image/webp"

	"kira-go/internal/config"
)

// 允许的文件类型
var allowedTypes = map[string]bool{
	"image/jpeg":    true,
	"image/png":     true,
	"image/webp":    true,
	"image/gif":     true,
	"image/svg+xml": true,
}

// 最大文件大小 (10MB)
const maxFileSize = 10 * 1024 * 1024

// UploadResult 上传结果
type UploadResult struct {
	URL         string `json:"url"`
	Orientation string `json:"orientation"`
}

// UploadImage 上传图片
// file: 上传的文件
// 返回：上传 URL 和图片方向，错误
func UploadImage(file *multipart.FileHeader) (*UploadResult, error) {
	// 1. 校验文件类型
	contentType := file.Header.Get("Content-Type")
	if !allowedTypes[contentType] {
		return nil, errors.New("不支持的文件类型")
	}

	// 2. 校验文件大小
	if file.Size > maxFileSize {
		return nil, errors.New("文件大小不能超过 10MB")
	}

	// 3. 打开文件
	uploadedFile, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer uploadedFile.Close()

	// 4. 读取文件内容到内存（用于方向检测和上传）
	buf := make([]byte, 0, file.Size)
	for {
		tmp := make([]byte, 1024)
		n, readErr := uploadedFile.Read(tmp)
		if n > 0 {
			buf = append(buf, tmp[:n]...)
		}
		if readErr != nil {
			break
		}
	}
	content := buf

	// 5. 检测图片方向
	orientation := detectOrientation(content)

	// 6. 生成文件名
	ext := filepath.Ext(file.Filename)
	if ext == "" {
		ext = ".webp"
	}
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	filename := uuid.New().String() + ext

	// 7. 根据配置选择存储方式
	if config.C.StorageType == "oss" {
		return uploadToOSS(filename, content, orientation)
	}
	return uploadToLocal(filename, content, orientation)
}

// uploadToLocal 上传图片到本地服务器
func uploadToLocal(filename string, content []byte, orientation string) (*UploadResult, error) {
	// 确保上传目录存在
	uploadDir := "./uploads/images"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, err
	}

	// 写入文件
	filePath := filepath.Join(uploadDir, filename)
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		return nil, err
	}

	// 构造本地访问 URL
	url := "/uploads/images/" + filename

	return &UploadResult{
		URL:         url,
		Orientation: orientation,
	}, nil
}

// uploadToOSS 上传图片到阿里云 OSS
func uploadToOSS(filename string, content []byte, orientation string) (*UploadResult, error) {
	bucket, err := getBucket()
	if err != nil {
		return nil, err
	}

	ossKey := config.C.OSSPrefix + filename
	if err := bucket.PutObject(ossKey, bytes.NewReader(content)); err != nil {
		return nil, err
	}

	// 构造 OSS 访问 URL
	url := config.C.OSSDomain + "/" + ossKey

	return &UploadResult{
		URL:         url,
		Orientation: orientation,
	}, nil
}

// getBucket 获取 OSS Bucket 实例
func getBucket() (*oss.Bucket, error) {
	client, err := oss.New(
		config.C.OSSRegion,
		config.C.OSSAccessKeyID,
		config.C.OSSAccessKeySecret,
	)
	if err != nil {
		return nil, err
	}
	return client.Bucket(config.C.OSSBucket)
}

// detectOrientation 检测图片方向 (landscape 或 portrait)
func detectOrientation(content []byte) string {
	decoders := []func(r io.Reader) (image.Config, error){
		jpeg.DecodeConfig,
		png.DecodeConfig,
		webp.DecodeConfig,
		gif.DecodeConfig,
	}

	for i := 0; i < len(decoders); i++ {
		cfg, err := decoders[i](bytes.NewReader(content))
		if err == nil {
			if cfg.Width >= cfg.Height {
				return "landscape"
			}
			return "portrait"
		}
	}
	return "landscape"
}
