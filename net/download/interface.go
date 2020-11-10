package download

// Tool 下载接口结构体
type Tool interface {
	Download(urlStr, host, port, user, password string, filePath, targetPath string) (string, error)
}
