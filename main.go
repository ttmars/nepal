package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
	"nepal/myTheme"
	"nepal/pkg"
	"os"
)

func main()  {
	myApp := app.NewWithID("go")						// 创建APP
	myWindow := myApp.NewWindow("HttpHijack")			// 创建窗口
	//myApp.SetIcon(myTheme.ResourceLogoJpg)			    	// 设置logo
	myApp.SetIcon(theme.FyneLogo())							// 默认logo
	myApp.Settings().SetTheme(&myTheme.MyTheme{})			// 设置APP主题，嵌入字体，解决乱码
	myWindow.Resize(fyne.NewSize(1200,800))			// 设置窗口大小
	myWindow.CenterOnScreen()								// 窗口居中显示
	myWindow.SetMaster()									// 设置为主窗口

	pkg.W = myWindow
	pkg.InitPreferences(myApp)
	defer os.RemoveAll(pkg.CachePath)						// 退出程序后，清理缓存目录
	defer pkg.UnSetWinProxy()								// 退出程序后，取消系统代理

	myWindow.SetMainMenu(pkg.MakeMyMenu(myApp, myWindow))	// 加载菜单
	myWindow.SetContent(pkg.MakeApp())						// 加载主界面
	go pkg.OpenProxy()

	myWindow.ShowAndRun()
}