package pkg

import (
	"bytes"
	"fmt"
	"fyne.io/fyne/v2"
	"github.com/ttmars/goproxy"
	"golang.org/x/sys/windows/registry"
	"io"
	"math/rand"
	"time"

	//"github.com/elazarl/goproxy"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var proxyPort = "7777"
var Proxy *goproxy.ProxyHttpServer
var CachePath string
var DownloadPath string

// InitPreferences 初始化Preferences变量
func InitPreferences(a fyne.App)  {
	if _,err := os.Stat(a.Preferences().String("downloadPath"));err != nil {
		curDir,_ := os.Getwd()
		a.Preferences().SetString("downloadPath", curDir + "\\download")
	}
	DownloadPath = a.Preferences().String("downloadPath")
	CachePath = os.TempDir() + "\\nepal"
	if _,err := os.Stat(CachePath);err != nil {
		os.MkdirAll(CachePath, 0755)
	}
	if _,err := os.Stat(DownloadPath);err != nil {
		os.MkdirAll(DownloadPath, 0755)
	}
}

func OpenProxy() {
	Proxy = goproxy.NewProxyHttpServer()
	Proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm) // 设置https请求截取
	Proxy.CertStore = NewCertStorage() //设置storage
	Proxy.Verbose = false

	SetCA(CaCert, CaKey)				// 设置自签证书
	SetResp("", "all", "all", "all")		// 初始化过滤条件

	log.Fatal(http.ListenAndServe(":"+proxyPort, Proxy))
}

func SetResp(host string, typ string, method string, code string)  {
	Proxy.ClearRespHandlers()		// 清空条件调用链，可以动态设置筛选条件;源码中未实现这个方法，自行fork添加！！！
	Proxy.OnResponse(RespCond(host,typ, method, code)).DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		return HandleResp(resp)
	})
}

func GetRandomString2(n int) string {
	rand.Seed(time.Now().UnixNano())
	randBytes := make([]byte, n/2)
	rand.Read(randBytes)
	return fmt.Sprintf("%x", randBytes)
}

func HandleResp(resp *http.Response) *http.Response  {
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
	err = os.WriteFile(CachePath+"\\"+fileName, b, 0755)
	if err != nil {
		fmt.Println("writefile err")
	}

	method := resp.Request.Method

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
	Data = append(Data, Item{URI: url, ContentType: contentType, Method: method, Size: size, SizeInt: lb, StatusCode: strconv.Itoa(resp.StatusCode), CacheName: fileName})
	List.Refresh()
	return resp
}

func RespCond(host string, typ string, method string, code string) goproxy.RespCondition {
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

func SetWinProxy()  {
	if !editReg("1", "localhost:"+proxyPort) {
		log.Println("设置代理失败！")
	}
}

func UnSetWinProxy()  {
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