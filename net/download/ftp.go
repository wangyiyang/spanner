package download

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/jlaffaye/ftp"
)

// FTP 默认参数
var (
	FTPDefaultPort     = "21"
	FTPDefaultUsername = "anonymous"
	FTPDefaultPassword = "anonymous"
)

// FTP 结构体
type FTP struct {
}

// Download FTP 下载
func (f *FTP) Download(urlStr, host, port, user, password string, filePath, targetPath string) (string, error) {
	c, err := ftpConnect(host, port, user, password)
	defer c.Quit()
	if err != nil {
		return "", err
	}
	_, fileName := filepath.Split(filePath)
	r, err := c.Retr(filePath)
	if err != nil {
		return "", err
	}
	fmt.Println(targetPath + "/" + fileName)
	file, err := os.Create(targetPath + "/" + fileName)
	if err != nil {
		return "", err
	}
	_, er2 := io.Copy(file, r)
	if er2 != nil {
		return "", er2
	}
	_ = file.Close()
	return fileName, nil
}

func ftpConnect(host, port, user, password string) (*ftp.ServerConn, error) {
	if port == "" {
		port = FTPDefaultPort
	}
	linkStr := fmt.Sprintf("%s:%s", host, port)
	c, err := ftp.Dial(linkStr, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		log.Fatal(err)
	}
	if user == "" {
		user = FTPDefaultUsername
		password = FTPDefaultPassword
	}
	err = c.Login(user, password)
	if err != nil {
		log.Fatal(err)
	}
	return c, err
}
