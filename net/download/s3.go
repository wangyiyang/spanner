package download

import (
	"io"
	"os"

	"github.com/minio/minio-go/v6"
)

//s3参数
var (
	URL       string
	AccessKey string
	SecretKey string
	SSL       bool
)

func newClient() (*minio.Client, error) {
	return minio.New(URL, AccessKey, SecretKey, SSL)
}

func download(bucket, fileKey string) (object *minio.Object, err error) {
	client, err := newClient()
	if err != nil {
		return nil, err
	}
	return client.GetObject(bucket, fileKey, minio.GetObjectOptions{})
}

// S3Download s3下载文件
func S3Download(filePath, bucket, fileKey string) (string, error) {
	object, err := download(bucket, fileKey)
	if err != nil {
		return "", err
	}
	targetPath := filePath + "/" + fileKey
	f, err := os.Create(targetPath)
	if err != nil {
		return "", err
	}
	if _, err = io.Copy(f, object); err != nil {
		return "", err
	}
	return fileKey, err
}
