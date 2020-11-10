package download

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// SFTP 结构体
type SFTP struct {
}

// Download SFTP 下载文件
func (s *SFTP) Download(urlStr, host, port, user, password string, filePath, targetPath string) (fileName string, err error) {
	sClient, err := SftpCreate(user, password, host, port)
	if err != nil {
		return "", err
	}
	defer sClient.Close()
	srcFile, err := sClient.Open(filePath)
	if err != nil {
		return "", err
	}
	defer srcFile.Close()

	var localFileName = path.Base(filePath)
	dstFile, err := os.Create(path.Join(targetPath, localFileName))
	if err != nil {
		return "", err
	}
	defer dstFile.Close()
	if _, err = srcFile.WriteTo(dstFile); err != nil {
		return "", err
	}

	return localFileName, err
}

// SftpCreate 创建sftp 客户端
func SftpCreate(user string, passwd string, host string, port string) (*sftp.Client, error) {
	auth := make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(passwd))
	config := ssh.ClientConfig{
		User:            user,
		Auth:            auth,
		Timeout:         10 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	addr := fmt.Sprintf("%s:%s", host, port)
	conn, err := ssh.Dial("tcp", addr, &config)
	if err != nil {
		return nil, err
	}

	//fmt.Println("connet sftp server ok.")
	c, err := sftp.NewClient(conn)
	if err != nil {
		_ = conn.Close()
		return nil, err
	}
	//fmt.Println("create sftp client ok.")
	return c, nil
}

// SftpReadAll 获取sftp sPath下所有文件列表
func SftpReadAll(user string, passwd string, host string, port string, sPath string) ([]byte, error) {
	sClient, err := SftpCreate(user, passwd, host, port)
	if err != nil {
		fmt.Printf("create ftp client failed, err:%+v.", err)
		return nil, err
	}
	defer sClient.Close()
	fp, err := sClient.Open(sPath)
	if err != nil {
		fmt.Printf("open remote file :%s, err:%+v failed.", sPath, err)
		return nil, err
	}
	defer fp.Close()
	bytes, err := ioutil.ReadAll(fp)
	return bytes, err
}
