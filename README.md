### 简介
FyneProxy是一个http/https代理工具，可用于抓包分析、资源下载等等
实现功能：
- 多种方式过滤请求，如URL、请求方法、资源类型、响应状态码等
- 一键开启/关闭系统代理，退出后自动关闭
- 资源预览、下载保存
- 对资源大小进行排序
- 简单易用，初次使用需安装证书

### 编译环境
- golang
- gcc
- fyne打包工具

### 打包
fyne package -os windows -icon logo.png

### 下载
[https://github.com/ttmars/nepal/releases](https://github.com/ttmars/nepal/releases)

### 效果图
![image](https://raw.githubusercontent.com/ttmars/image/master/github/fyneProxy.png)

# 如何对安卓进行抓包？
以夜神模拟器为例，参考：[https://blog.csdn.net/qq_43278826/article/details/124291040](https://blog.csdn.net/qq_43278826/article/details/124291040)

### 下载adb工具

[https://androidstudio.io/downloads/tools/download-the-latest-version-of-adb.exe.html](https://androidstudio.io/downloads/tools/download-the-latest-version-of-adb.exe.html)

### 开启USB调试

[https://support.yeshen.com/zh-CN/often/kfz](https://support.yeshen.com/zh-CN/often/kfz)

1. 进入设置，多次点击版本号解锁开发者选项
2. 进入开发者选项，启动USB调试

### 证书制作

#### 生成证书

Go SDK中提供了证书制作工具，执行后会在当前路径下生成cert和key文件

[https://zhuanlan.zhihu.com/p/514004767](https://zhuanlan.zhihu.com/p/514004767)

```shell
go run $GOROOT/src/crypto/tls/generate_cert.go --host localhost
```

#### 计算hash值

计算证书的hash值，并重命名

```shell
# 计算hash
[root@iZ8vb8qjajxkxytobaq81hZ test]# openssl x509 -inform PEM -subject_hash_old -in xx.cert 
ea1cd156
-----BEGIN CERTIFICATE-----
...

# 重命名证书
mv xx.cert ea1cd156.0
```

最终得到一个符合Android的系统证书文件：ea1cd156.0，因为Android7版本的用户证书不受信任！

### 导入证书

```shell
# 连接Android模拟器
C:\Users\lee>adb connect 127.0.0.1:62001
already connected to 127.0.0.1:62001

# 列出已连接的设备
C:\Users\lee>adb devices
List of devices attached
127.0.0.1:62001 device

# 进入模拟器中的shell环境
C:\Users\lee>adb shell
beyond1q:/ #

# 修改证书目录的权限
beyond1q:/ # cd /system/etc/security/
beyond1q:/system/etc/security # ls
cacerts  mac_permissions.xml  otacerts.zip
beyond1q:/system/etc/security # chmod 777 cacerts/

# 另起一个终端，导入上一步制作好的证书
C:\Users\lee>adb push C:\lee\Downloads\ea1cd156.0 /system/etc/security/cacerts
[100%] /system/etc/security/cacerts/ea1cd156.0

# 修改证书的权限
beyond1q:/system/etc/security/cacerts # ll ea1cd156.0
-rw-r--r-- 1 root root 1110 2022-11-10 19:09 ea1cd156.0
beyond1q:/system/etc/security/cacerts # chmod 777 ea1cd156.0
beyond1q:/system/etc/security/cacerts # ll ea1cd156.0
-rwxrwxrwx 1 root root 1110 2022-11-10 19:09 ea1cd156.0
```

### 手动设置代理

进入WiFi设置，手动将代理设置为：宿主机IP+端口，注意不要设置成了localhost+端口。

到这里就可以愉快的抓包了！