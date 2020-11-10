package download

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// HTTP 结构体
type HTTP struct {
}

// Download HTTP 下载方法
func (h *HTTP) Download(urlStr, host, port, user, password string, filePath, targetPath string) (string, error) {
	_, fileName := filepath.Split(filePath)
	client := &http.Client{}
	req, err := http.NewRequest("GET", urlStr, nil)
	fmt.Println(err)
	if err != nil {
		return "", err
	}
	if user != "" {
		req.SetBasicAuth(user, password)
	}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	f, err := os.Create(targetPath + "/" + fileName)
	if err != nil {
		return "", err
	}
	_, err = io.Copy(f, res.Body)
	if err != nil {
		return "", err
	}
	return fileName, nil
}

// EncodeJSON  结构体转Json
func EncodeJSON(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

// SendRequest 发送请求
func SendRequest(method, baseURL, path string, body interface{}, header map[string]string, cookie map[string]string, RegistryUserName, RegistryPwd string) ([]byte, int, error) {
	client := http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(60 * time.Second)
				c, err := net.DialTimeout(netw, addr, time.Second*50)
				if err != nil {
					return nil, err
				}
				_ = c.SetDeadline(deadline)
				return c, nil
			},
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	encodeBody, err := EncodeJSON(body)
	if err != nil {
		return nil, 0, err
	}

	url := fmt.Sprintf("%s%s", baseURL, path)
	r, err := http.NewRequest(method, url, bytes.NewBuffer(encodeBody))
	if err != nil {
		return nil, 0, err
	}

	r.Header.Set("Content-Type", "application/json")
	for k, v := range header {
		r.Header.Set(k, v)
	}

	for k, v := range cookie {
		cookieStr := fmt.Sprintf("%s=%s", k, v)
		r.Header.Set("Cookie", cookieStr)
	}
	if RegistryUserName != "" {
		r.SetBasicAuth(RegistryUserName, RegistryPwd)
	}
	w, err := client.Do(r)
	if err != nil {
		return nil, 0, err
	}
	defer w.Body.Close()

	data, err := ioutil.ReadAll(w.Body)
	if err != nil {
		return nil, 0, err
	}

	return data, w.StatusCode, nil
}
