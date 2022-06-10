package pb

import (
	"bytes"
	context "context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	sync "sync"
	"time"

	"github.com/h2non/filetype"
	svg "github.com/h2non/go-is-svg"
	"golang.org/x/net/http/httpproxy"
)

type ctxKey string

// envProxyOnce System proxies load only one
var envProxyOnce sync.Once

// envProxyFuncValue System proxies get funcation
var envProxyFuncValue func(*url.URL) (*url.URL, error)

type Downloader struct {
	Dir    string
	Client *http.Client
}

// proxyFunc http.Transport.Proxy return proxy
func proxyFunc(req *http.Request) (*url.URL, error) {
	// 从上下文管理器中获取代理配置，实现代理和请求的一对一配置关系
	value := req.Context().Value(ctxKey("key")).(map[string]interface{})
	proxy, ok := value["proxy"]
	if !ok {
		return nil, nil
	}
	p := proxy.(*Proxy)
	// If there is no proxy set, use default proxy from environment.
	// This mitigates expensive lookups on some platforms (e.g. Windows).
	envProxyOnce.Do(func() {
		// get proxies from system env
		envProxyFuncValue = httpproxy.FromEnvironment().ProxyFunc()
	})
	if p != nil && p.ProxyUrl != "" {
		proxyURL, err := urlParse(p.ProxyUrl)
		if err != nil {
			err := fmt.Sprint(err.Error())
			return nil, errors.New(err)
		}
		return proxyURL, nil
	}
	return envProxyFuncValue(req.URL)
}

// redirectFunc redirect handle funcation
// limit max redirect times
func redirectFunc(req *http.Request, via []*http.Request) error {
	value := req.Context().Value(ctxKey("key")).(map[string]interface{})
	redirectNum := value["redirectNum"].(int)
	if len(via) > redirectNum {
		err := fmt.Errorf("stopped after %d redirects", redirectNum)
		return err
	}
	return nil
}
func urlParse(URL string) (*url.URL, error) {
	return url.Parse(URL)
}

// SpiderDownloader get a new spider downloader
func NewDownloader(dir string) *Downloader {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false,
		},
		Proxy: proxyFunc,
		DialContext: (&net.Dialer{
			Timeout:   2 * 60 * time.Second,
			KeepAlive: 2 * 60 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     false,
		MaxIdleConns:          1024,
		IdleConnTimeout:       2 * 60 * time.Second,
		TLSHandshakeTimeout:   60 * time.Second,
		ExpectContinueTimeout: 60 * time.Second,
		MaxIdleConnsPerHost:   1024,
		MaxConnsPerHost:       1024,
	}
	client := http.Client{
		Transport:     transport,
		CheckRedirect: redirectFunc,
	}
	return &Downloader{
		Dir:    dir,
		Client: &client,
	}

}
func (d *Downloader) IsFileExist(filename string, filesize int64) (bool, int64) {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false, 0
	}
	return true, info.Size()
}
func (d *Downloader) GetFileContentType(localFile *os.File) (string, error) {

	// 只需要前 512 个字节就可以了
	head := make([]byte, 1024)
	localFile.Seek(0, 0)
	localFile.Read(head)
	defer localFile.Seek(0, 0)
	kind, _ := filetype.Match(head)
	if kind == filetype.Unknown {
		localFile.Seek(0, 0)
		lo, err := ioutil.ReadAll(localFile)
		if err != nil {
			return "", err
		}
		if svg.Is(lo) {
			return "svg", nil
		} else {
			return "", errors.New("unknown file type")

		}
	}

	return kind.Extension, nil
}
func (d *Downloader) Download(request *Request, filename string) (string, error) {
	var (
		fsize int64
	)
	tmpFilePath := path.Join(d.Dir, filename+".download")
	exist, fsize := d.IsFileExist(tmpFilePath, fsize)
	var file *os.File
	var err error
	if exist {
		file, err = os.OpenFile(tmpFilePath, os.O_RDWR|os.O_APPEND, 0666)
	} else {
		file, err = os.Create(tmpFilePath)
	}
	defer func() {
		if file != nil {
			file.Close()
		}
	}()

	if err != nil {
		return "", err
	}
	ctxValue := map[string]interface{}{}
	if request.Proxy != nil {

		ctxValue["proxy"] = request.Proxy

	}
	ctxValue["redirectNum"] = 3
	var asCtxKey ctxKey = "key"
	valCtx := context.WithValue(context.TODO(), asCtxKey, ctxValue)

	req, err := http.NewRequestWithContext(valCtx, request.Method, request.Url, bytes.NewReader(request.Body))
	if err != nil {
		return "", err
	}
	fmt.Printf("downloading RANGE %d", fsize)
	if exist {
		req.Header.Set("Range", "bytes="+strconv.FormatInt(fsize, 10)+"-")
	}
	for _, header := range request.Header {
		req.Header.Set(header.Key, header.Value)
	}
	resp, err := d.Client.Do(req)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if err != nil {
		return "", err
	}
	if resp.Body == nil {
		return "", errors.New("body is null")
	}
	defer resp.Body.Close()
	if resp.StatusCode >= http.StatusBadRequest {
		return "", errors.New("status code is " + strconv.Itoa(resp.StatusCode))
	}
	if exist {
		file.Seek(0, os.SEEK_END)
	}
	_, err = io.Copy(file, resp.Body)
	if err == nil {
		localFile := path.Join(d.Dir, filename)
		err = os.Rename(tmpFilePath, localFile)
		return localFile, err
	}
	return "", err
}
