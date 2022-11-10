package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"nepal/myTheme"
	"nepal/pkg"
	"os"
)

func main()  {
	myApp := app.NewWithID("go!")						// 创建APP
	myWindow := myApp.NewWindow("HttpHijack")			// 创建窗口
	pkg.W = myWindow
	defer os.RemoveAll(pkg.CachePath)						// 退出程序后，清理缓存
	defer pkg.UnSetWinProxy()								// 退出程序后，取消代理

	myApp.SetIcon(myTheme.ResourceLogoJpg)			    	// 设置logo
	myApp.Settings().SetTheme(&myTheme.MyTheme{})			// 设置APP主题，嵌入字体，解决乱码
	myWindow.Resize(fyne.NewSize(1200,800))			// 设置窗口大小
	myWindow.CenterOnScreen()								// 窗口居中显示
	myWindow.SetMaster()									// 设置为主窗口
	//myWindow.SetCloseIntercept(func() {myWindow.Hide()})	// 设置窗口托盘显示
	//if desk, ok := myApp.(desktop.App); ok {
	//	m := fyne.NewMenu("MyApp",
	//		fyne.NewMenuItem("Show", func() {
	//			myWindow.Show()
	//		}))
	//	desk.SetSystemTrayMenu(m)
	//}

	myWindow.SetContent(pkg.MakeApp())

	go pkg.OpenProxy()

	myWindow.ShowAndRun()			// 事件循环
}
