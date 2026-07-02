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

// UploadImage 上传图片到阿里云 OSS
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
	var content []byte
	buf := make([]byte, 0, file.Size)
	for {
		tmp := make([]byte, 1024)
		n, err := uploadedFile.Read(tmp)
		if err != nil {
			break
		}
		buf = append(buf, tmp[:n]...)
	}
	content = buf

	// 5. 检测图片方向
	orientation := detectOrientation(content)

	// 6. 生成文件名和路径
	ext := filepath.Ext(file.Filename)
	if ext == "" {
		ext = ".webp"
	}
	// 确保扩展名包含点号
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	filename := uuid.New().String() + ext
	ossKey := config.C.OSSPrefix + filename

	// 7. 上传到 OSS
	bucket, err := getBucket()
	if err != nil {
		return nil, err
	}

	if err := bucket.PutObject(ossKey, bytes.NewReader(content)); err != nil {
		return nil, err
	}

	// 8. 构造访问 URL
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
// 使用 image.DecodeConfig 仅读取图片元数据（宽高），不加载完整图片到内存
func detectOrientation(content []byte) string {
	// 尝试不同格式，使用 DecodeConfig 只获取宽高
	decoders := []func(r io.Reader) (image.Config, error){
		jpeg.DecodeConfig,
		png.DecodeConfig,
		webp.DecodeConfig,
		gif.DecodeConfig,
	}

	for _, decode := range decoders {
		if config, err := decode(bytes.NewReader(content)); err == nil {
			if config.Width >= config.Height {
				return "landscape"
			}
			return "portrait"
		}
	}
	// 默认返回 landscape
	return "landscape"
}
