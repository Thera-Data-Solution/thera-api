package utils

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"os"
	"strings"
	"time"

	_ "golang.org/x/image/webp"

	"github.com/chai2010/webp"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOUploader struct {
	Client     *minio.Client
	BucketName string
	Endpoint   string
	Secure     bool
}

func NewMinIOUploader() (*MinIOUploader, error) {
	endpoint := os.Getenv("S3_ENDPOINT")
	accessKey := os.Getenv("S3_ACCESS_KEY")
	secretKey := os.Getenv("S3_SECRET_KEY")
	bucket := os.Getenv("S3_BUCKET")

	if accessKey == "" || secretKey == "" || bucket == "" || endpoint == "" {
		return nil, fmt.Errorf("salah satu MinIO environment variables (ENDPOINT, ACCESS_KEY, SECRET_KEY, BUCKET) hilang")
	}

	var host string
	secure := false
	if strings.HasPrefix(endpoint, "https://") {
		host = strings.TrimPrefix(endpoint, "https://")
		secure = true
	} else if strings.HasPrefix(endpoint, "http://") {
		host = strings.TrimPrefix(endpoint, "http://")
	} else {
		host = endpoint
	}
	host = strings.TrimSuffix(host, "/")

	client, err := minio.New(host, &minio.Options{
		Creds:  credentials.NewStatic(accessKey, secretKey, "", credentials.SignatureV4),
		Secure: secure,
	})
	if err != nil {
		return nil, fmt.Errorf("gagal membuat MinIO client: %w", err)
	}

	return &MinIOUploader{
		Client:     client,
		BucketName: bucket,
		Endpoint:   endpoint,
		Secure:     secure,
	}, nil
}

func (u *MinIOUploader) UploadFile(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	defer file.Close()

	// Pastikan pointer file di awal
	if _, err := file.Seek(0, 0); err != nil {
		return "", fmt.Errorf("gagal reset pointer file: %w", err)
	}

	// Baca semua file ke buffer (agar bisa di-decode berulang)
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("gagal membaca file: %w", err)
	}

	imgReader := bytes.NewReader(fileBytes)
	img, _, err := image.Decode(imgReader)
	if err != nil {
		return "", fmt.Errorf("gagal decode gambar: %w", err)
	}

	// Kompres ke WebP
	var buf bytes.Buffer
	quality := float32(80)
	for {
		buf.Reset()
		err = webp.Encode(&buf, img, &webp.Options{Quality: quality})
		if err != nil {
			return "", fmt.Errorf("gagal encode ke webp: %w", err)
		}

		if buf.Len() < 1_500_000 || quality <= 30 {
			break
		}
		quality -= 5
	}

	// Gunakan timestamp sebagai nama file
	timestamp := time.Now().UnixMilli()
	objectName := fmt.Sprintf("logos/%d.webp", timestamp)

	info, err := u.Client.PutObject(ctx, u.BucketName, objectName, &buf, int64(buf.Len()), minio.PutObjectOptions{
		ContentType: "image/webp",
		UserMetadata: map[string]string{
			"x-amz-acl": "public-read",
		},
	})
	if err != nil {
		return "", fmt.Errorf("operasi PutObject MinIO gagal: %w", err)
	}

	fmt.Printf("✅ File berhasil diupload: %s (%.2f KB)\n", info.Key, float64(info.Size)/1024)
	finalURL := fmt.Sprintf("%s/%s/%s", strings.TrimRight(u.Endpoint, "/"), u.BucketName, objectName)
	return finalURL, nil
}
