package proxy

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/ttmars/goproxy"
	"golang.org/x/sys/windows/registry"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var DefaultProxy = NewNepalProxy()

type NepalProxy struct {
	ProxyPort string
	CaCert []byte
	CaKey []byte
	CachePath string
	DownloadPath string
	Proxy *goproxy.ProxyHttpServer
}

func NewNepalProxy() *NepalProxy {
	curDir,_ := os.Getwd()
	downloadPath := curDir + "\\download"
	cachePath := os.TempDir() + "\\nepal"
	if _,err := os.Stat(cachePath);err != nil {
		os.MkdirAll(cachePath, 0755)
	}
	if _,err := os.Stat(downloadPath);err != nil {
		os.MkdirAll(downloadPath, 0755)
	}
	return &NepalProxy{
		ProxyPort:    "7777",
		CaCert:       CaCert,
		CaKey:        CaKey,
		CachePath:    cachePath,
		DownloadPath: downloadPath,
		Proxy:        goproxy.NewProxyHttpServer(),
	}
}

func (NP *NepalProxy)Run() {
	NP.Proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm) 				// 设置https请求截取
	NP.Proxy.CertStore = NewCertStorage()                  				//设置storage
	NP.Proxy.Verbose = false

	NP.SetCA()											   				// 设置自签证书
	NP.HandleResp("", "all", "all", "all")		// 初始化过滤条件

	log.Fatal(http.ListenAndServe(":"+NP.ProxyPort, NP.Proxy))
}

func (NP *NepalProxy)SetCA() error {
	goproxyCa, err := tls.X509KeyPair(NP.CaCert, NP.CaKey)
	if err != nil {
		return err
	}
	if goproxyCa.Leaf, err = x509.ParseCertificate(goproxyCa.Certificate[0]); err != nil {
		return err
	}
	goproxy.GoproxyCa = goproxyCa
	goproxy.OkConnect = &goproxy.ConnectAction{Action: goproxy.ConnectAccept, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	goproxy.MitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectMitm, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	goproxy.HTTPMitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectHTTPMitm, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	goproxy.RejectConnect = &goproxy.ConnectAction{Action: goproxy.ConnectReject, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	return nil
}

// ClearProxyData 清理缓存并还原系统代理
func (NP *NepalProxy)ClearProxyData() {
	_ = os.RemoveAll(NP.CachePath)
	NP.UnSetWinProxy()
}

func (NP *NepalProxy)HandleResp(host, typ, method, code string)  {
	NP.Proxy.ClearRespHandlers()		// 清空条件调用链，可以动态设置筛选条件;源码中未实现这个方法，自行fork添加！！！
	NP.Proxy.OnResponse(respCond(host,typ, method, code)).DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		contentType := resp.Header.Get("Content-Type")
		var suffix string
		if contentType != "" {
			suffix = "." + strings.Split(strings.Split(contentType, ";")[0], "/")[1]
		}

		fileName := GetRandomString2(6) + suffix
		b,err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("readall err")
		}
		resp.Body = io.NopCloser(bytes.NewBuffer(b))
		err = os.WriteFile(NP.CachePath+"\\"+fileName, b, 0755)
		if err != nil {
			fmt.Println("writefile err")
		}

		m := resp.Request.Method

		url := resp.Request.URL.String()
		if len(url) >= 75 {
			url = url[:75]
		}
		if len(contentType) >= 10 {
			contentType = contentType[:10]
		}

		var size string
		lb  := len(b)
		if lb < 1024 {
			size = fmt.Sprintf("%dB", lb)
		}else if lb < 1048576{
			size = fmt.Sprintf("%.2fK", float64(lb)/1024)
		}else {
			size = fmt.Sprintf("%.2fM", float64(lb)/1048576)
		}
		//Data = append(Data, Item{URI: url, ContentType: contentType, Method: method, Size: size, SizeInt: lb, StatusCode: strconv.Itoa(resp.StatusCode), CacheName: fileName})
		Data = append([]Item{{URI: url, ContentType: contentType, Method: m, Size: size, SizeInt: lb, StatusCode: strconv.Itoa(resp.StatusCode), CacheName: fileName}}, Data...)
		List.Refresh()
		return resp
	})
}

func respCond(host string, typ string, method string, code string) goproxy.RespCondition {
	return goproxy.RespConditionFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) bool {
		var hostCond, typCond, methodCond, codeCond bool
		s := resp.StatusCode

		if host==""||strings.Contains(resp.Request.URL.String(), host){
			hostCond=true
		}
		if typ == "all" || strings.Contains(resp.Header.Get("Content-Type"), typ) {
			typCond = true
		}
		if method == "all" || resp.Request.Method == method {
			methodCond = true
		}
		if (code=="all")||(code=="1xx"&&s>=100&&s<200)||(code=="2xx"&&s>=200&&s<300)||(code=="3xx"&&s>=300&&s<400)||(code=="4xx"&&s>=400&&s<500)||(code=="5xx"&&s>=500&&s<600){
			codeCond=true
		}

		if hostCond && typCond && methodCond && codeCond{
			return true
		}
		return false
	})
}

func (NP *NepalProxy)SetWinProxy()  {
	if !editReg("1", "localhost:"+NP.ProxyPort) {
		log.Println("设置代理失败！")
	}
}

func (NP *NepalProxy)UnSetWinProxy() {
	if !editReg("0", "") {
		log.Println("取消代理失败！")
	}
}

func editReg(enable, proxy string) bool {
	key, exists, err := registry.CreateKey(registry.CURRENT_USER, "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Internet Settings", registry.ALL_ACCESS)
	if err != nil {
		log.Fatal(err)
	}
	defer key.Close()

	if !exists {
		return false
	}

	err = key.SetStringValue("ProxyEnable", enable)
	if err != nil {
		return false
	}

	err = key.SetStringValue("ProxyServer", proxy)
	if err != nil {
		return false
	}

	return true
}

func GetRandomString2(n int) string {
	rand.Seed(time.Now().UnixNano())
	randBytes := make([]byte, n/2)
	rand.Read(randBytes)
	return fmt.Sprintf("%x", randBytes)
}