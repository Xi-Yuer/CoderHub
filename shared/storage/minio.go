package storage

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"path/filepath"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/nfnt/resize"
)

type Minio struct {
	Client    *minio.Client
	Bucket    string
	Region    string
	Endpoint  string
	AccessKey string
	SecretKey string
	UseSSL    bool
}

func NewMinio(endpoint, accessKey, secretKey, bucket, region string, useSSL bool) *Minio {
	return &Minio{
		Endpoint:  endpoint,
		AccessKey: accessKey,
		SecretKey: secretKey,
		Bucket:    bucket,
		Region:    region,
		UseSSL:    useSSL,
	}
}

func (m *Minio) Connect() error {
	// 初始化 Minio 客户端
	client, err := minio.New(m.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(m.AccessKey, m.SecretKey, ""),
		Secure: m.UseSSL,
	})
	if err != nil {
		return err
	}

	m.Client = client

	// 检查并创建 bucket
	exists, err := client.BucketExists(context.Background(), m.Bucket)
	if err != nil {
		return fmt.Errorf("检查 bucket 失败: %v", err)
	}

	if !exists {
		// 创建 bucket
		err = client.MakeBucket(context.Background(), m.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("创建 bucket 失败: %v", err)
		}
	}

	// 设置 bucket 公开读取策略
	policy := `{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": {"AWS": ["*"]},
				"Action": ["s3:GetObject"],
				"Resource": ["arn:aws:s3:::` + m.Bucket + `/*"]
			}
		]
	}`

	err = client.SetBucketPolicy(context.Background(), m.Bucket, policy)
	if err != nil {
		return fmt.Errorf("设置 bucket 策略失败: %v", err)
	}

	return nil
}

// 创建桶
func (m *Minio) CreateBucket(bucketName string) error {
	return m.Client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{Region: m.Region})
}

// 上传文件
func (m *Minio) UploadFile(bucketName, objectName string, reader io.Reader, objectSize int64, contentType string) (minio.UploadInfo, error) {
	// 判断桶是否存在
	exists, err := m.Client.BucketExists(context.Background(), bucketName)
	if err != nil {
		return minio.UploadInfo{}, err
	}
	if !exists {
		_ = m.CreateBucket(bucketName)
	}
	return m.Client.PutObject(context.Background(), bucketName, objectName, reader, objectSize, minio.PutObjectOptions{ContentType: contentType})
}

// 获取文件
func (m *Minio) GetFile(bucketName, objectName string) (*minio.Object, error) {
	return m.Client.GetObject(context.Background(), bucketName, objectName, minio.GetObjectOptions{})
}

// 删除文件
func (m *Minio) DeleteFile(bucketName, objectName string) error {
	return m.Client.RemoveObject(context.Background(), bucketName, objectName, minio.RemoveObjectOptions{})
}

// GetFileURL 获取文件的URL，永久访问
func (m *Minio) GetFileURL(bucketName, objectName string) (string, error) {
	return m.GetPublicURL(objectName), nil
}

// ImageInfo 图片信息结构体
type ImageInfo struct {
	BucketName   string `json:"bucket_name"`   // 桶名
	ObjectName   string `json:"object_name"`   // 对象名
	FileName     string `json:"file_name"`     // 文件名
	Width        int    `json:"width"`         // 图片宽度
	Height       int    `json:"height"`        // 图片高度
	URL          string `json:"url"`           // 图片URL
	ThumbnailURL string `json:"thumbnail_url"` // 缩略图URL
	Size         int64  `json:"size"`          // 文件大小
	ContentType  string `json:"content_type"`  // 文件类型
}

// GetImageInfo 获取图片详细信息
func (m *Minio) GetImageInfo(bucketName, objectName string) (*ImageInfo, error) {
	// 获取图片对象
	obj, err := m.GetFile(bucketName, objectName)
	if err != nil {
		return nil, err
	}
	defer func(obj *minio.Object) {
		_ = obj.Close()
	}(obj)

	// 获取图片信息
	stat, err := obj.Stat()
	if err != nil {
		return nil, err
	}

	// 解码图片获取宽高
	img, _, err := image.DecodeConfig(obj)
	if err != nil {
		return nil, err
	}

	// 获取原图URL
	originalURL, err := m.GetFileURL(bucketName, objectName)
	if err != nil {
		return nil, err
	}

	// 构建缩略图对象名称
	thumbnailName := "thumb_" + objectName
	thumbnailURL, _ := m.GetFileURL(bucketName, thumbnailName)

	return &ImageInfo{
		BucketName:   bucketName,
		ObjectName:   objectName,
		FileName:     filepath.Base(objectName),
		Width:        img.Width,
		Height:       img.Height,
		URL:          originalURL,
		ThumbnailURL: thumbnailURL,
		Size:         stat.Size,
		ContentType:  stat.ContentType,
	}, nil
}

// UploadImageWithThumbnail 上传图片并生成缩略图
func (m *Minio) UploadImageWithThumbnail(bucketName, objectName string, reader io.Reader, objectSize int64, contentType string, thumbnailWidth uint) (ImageInfo, error) {
	exists, err := m.Client.BucketExists(context.Background(), bucketName)
	if err != nil {
		return ImageInfo{}, err
	}
	if !exists {
		err := m.CreateBucket(bucketName)
		if err != nil {
			return ImageInfo{}, err
		}
	}
	// 读取原始图片数据
	imageData, err := io.ReadAll(reader)
	if err != nil {
		return ImageInfo{}, err
	}

	// 上传原图
	originalReader := bytes.NewReader(imageData)
	_, err = m.UploadFile(bucketName, objectName, originalReader, objectSize, contentType)
	if err != nil {
		return ImageInfo{}, err
	}

	// 生成缩略图
	thumbnailData, err := m.generateThumbnail(imageData, thumbnailWidth, contentType)
	if err != nil {
		return ImageInfo{}, err
	}

	// 上传缩略图
	thumbnailName := "thumb_" + objectName
	thumbnailReader := bytes.NewReader(thumbnailData)
	_, err = m.UploadFile(bucketName, thumbnailName, thumbnailReader, int64(len(thumbnailData)), contentType)
	if err != nil {
		return ImageInfo{}, err
	}

	// 获取URL
	originalURL, err := m.GetFileURL(bucketName, objectName)
	if err != nil {
		return ImageInfo{}, err
	}

	thumbnailURL, err := m.GetFileURL(bucketName, "thumb_"+objectName)
	if err != nil {
		return ImageInfo{}, err
	}

	return ImageInfo{
		BucketName:   bucketName,
		ObjectName:   objectName,
		FileName:     filepath.Base(objectName),
		Width:        0,
		Height:       0,
		URL:          originalURL,
		ThumbnailURL: thumbnailURL,
		Size:         objectSize,
		ContentType:  contentType,
	}, nil
}

// generateThumbnail 生成缩略图
func (m *Minio) generateThumbnail(imageData []byte, width uint, contentType string) ([]byte, error) {
	// 解码原始图片
	originalImg, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, err
	}

	// 等比例缩放
	thumbnail := resize.Resize(width, 0, originalImg, resize.Lanczos3)

	// 编码缩略图
	var buf bytes.Buffer
	switch contentType {
	case "image/jpeg":
		err = jpeg.Encode(&buf, thumbnail, &jpeg.Options{Quality: 85})
	case "image/png":
		err = png.Encode(&buf, thumbnail)
	case "image/gif":
		// 对于GIF，我将其转换为静态图片（取第一帧）
		err = gif.Encode(&buf, thumbnail, &gif.Options{
			NumColors: 256,
		})
	default:
		return nil, fmt.Errorf("不支持的图片格式: %s", contentType)
	}

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// FileInfo 文件信息结构体
type FileInfo struct {
	BucketName  string `json:"bucket_name"`  // 桶名
	ObjectName  string `json:"object_name"`  // 对象名
	FileName    string `json:"file_name"`    // 文件名
	URL         string `json:"url"`          // 文件URL
	Size        int64  `json:"size"`         // 文件大小
	ContentType string `json:"content_type"` // 文件类型
}

// UploadFileWithInfo 上传文件并返回文件信息
func (m *Minio) UploadFileWithInfo(bucketName, objectName string, reader io.Reader, objectSize int64, contentType string) (FileInfo, error) {
	// 上传文件
	_, err := m.UploadFile(bucketName, objectName, reader, objectSize, contentType)
	if err != nil {
		return FileInfo{}, err
	}

	// 获取文件URL
	fileURL, err := m.GetFileURL(bucketName, objectName)
	if err != nil {
		return FileInfo{}, err
	}

	return FileInfo{
		BucketName:  bucketName,
		ObjectName:  objectName,
		FileName:    filepath.Base(objectName),
		URL:         fileURL,
		Size:        objectSize,
		ContentType: contentType,
	}, nil
}

// GetPublicURL 获取永久公开访问的URL
func (m *Minio) GetPublicURL(objectName string) string {
	return fmt.Sprintf("http://%s/%s/%s", "localhost:9000", m.Bucket, objectName)
}
