package proxy

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"io"
	"nepal/myTheme"
	"net/http"
	"os"
	"sort"
	"strings"
)

var List *widget.List
var SortFlag = true
var typSelect,methodSelect,codeSelect *widget.Select
var App fyne.App
var Window fyne.Window

var Data []Item
type Item struct {
	URI string
	ContentType string
	Method string
	Size string
	SizeInt int
	StatusCode string
	CacheName string
}

func RunApp()  {
	App = app.NewWithID("nepal")           			// 创建APP
	Window = App.NewWindow("FyneProxy") 			// 创建窗口
	//myApp.SetIcon(myTheme.ResourceLogoJpg)			// 设置logo
	App.SetIcon(theme.FyneLogo())               		// 默认logo
	App.Settings().SetTheme(&myTheme.MyTheme{}) 		// 设置APP主题，嵌入字体，解决乱码
	Window.Resize(fyne.NewSize(1200,800))       	// 设置窗口大小
	Window.CenterOnScreen()                     		// 窗口居中显示
	Window.SetMaster()                          		// 设置为主窗口

	Window.SetMainMenu(MakeMyMenu(App, Window))			// 加载菜单
	Window.SetContent(MakeApp())                		// 加载主界面

	defer DefaultProxy.ClearProxyData()					// 注册清理
	InitSetting()						    			// 初始化设置

	Window.ShowAndRun()
}

func InitSetting()  {
	downloadPath := App.Preferences().String("downloadPath")
	if _,err := os.Stat(downloadPath); err == nil {
		DefaultProxy.DownloadPath = downloadPath
	}
}

func MakeApp() fyne.CanvasObject {
	return container.NewBorder(container.NewVBox(MakeOperate(), MakeListLabel()),nil,nil,nil, MakeList())
}

func MakeOperate() fyne.CanvasObject {
	hostLabel := widget.NewLabel("过滤链接")
	hostEntry := widget.NewEntry()
	c1 := container.NewBorder(nil,nil,hostLabel,nil,hostEntry)

	typLabel := widget.NewLabel("过滤类型")
	typSelect = widget.NewSelect([]string{"all", "video", "audio", "image", "text", "application"}, func(value string) {
	})
	typSelect.SetSelected("all")
	c2 := container.NewBorder(nil,nil,typLabel,nil, typSelect)

	switchProxy := widget.NewCheck("开启代理", func(b bool) {
		if b {
			DefaultProxy.SetWinProxy()
		}else{
			DefaultProxy.UnSetWinProxy()
		}
	})
	switchProxy.SetChecked(true)
	certInstall := widget.NewHyperlink("安装证书", nil)
	os.WriteFile(DefaultProxy.CachePath + "\\" + "ca.crt", CaCert, 0755)
	certInstall.SetURLFromString("file:///" + DefaultProxy.CachePath + "\\" + "ca.crt")


	c3 := container.NewGridWithColumns(3, c1,c2,container.NewGridWithColumns(2,switchProxy,certInstall))

	methodLabel := widget.NewLabel("过滤方法")
	methodSelect = widget.NewSelect([]string{"all", http.MethodGet, http.MethodPost, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodConnect, http.MethodOptions, http.MethodTrace}, func(value string) {
	})
	methodSelect.SetSelected("all")
	c4 := container.NewBorder(nil,nil,methodLabel,nil, methodSelect)

	codeLabel := widget.NewLabel("过滤code")
	codeSelect = widget.NewSelect([]string{"all", "1xx", "2xx", "3xx", "4xx", "5xx"}, func(value string) {
	})
	codeSelect.SetSelected("all")
	c5 := container.NewBorder(nil,nil,codeLabel,nil, codeSelect)

	setCondButton := widget.NewButton("清空列表", func() {
		Data = Data[0:0]
		List.Refresh()
		DefaultProxy.HandleResp(hostEntry.Text, typSelect.Selected, methodSelect.Selected, codeSelect.Selected)
	})

	// 动态设置过滤条件
	typSelect.OnChanged = func(value string) {
		DefaultProxy.HandleResp(hostEntry.Text, value, methodSelect.Selected, codeSelect.Selected)
	}
	methodSelect.OnChanged = func(value string) {
		DefaultProxy.HandleResp(hostEntry.Text, typSelect.Selected, value, codeSelect.Selected)
	}
	codeSelect.OnChanged = func(value string) {
		DefaultProxy.HandleResp(hostEntry.Text, typSelect.Selected, methodSelect.Selected, value)
	}

	c6 := container.NewGridWithColumns(3, c4,c5,setCondButton)

	c := container.NewVBox(c3,c6)
	return c
}

func MakeListLabel() fyne.CanvasObject {
	listLabel1 := widget.NewLabelWithStyle("链接", fyne.TextAlignLeading, fyne.TextStyle{Bold:true})
	listLabel2 := widget.NewLabelWithStyle("类型", fyne.TextAlignLeading, fyne.TextStyle{Bold:true})
	listLabel7 := widget.NewLabelWithStyle("方法", fyne.TextAlignLeading, fyne.TextStyle{Bold:true})
	listLabel3 := widget.NewLabelWithStyle("状态码", fyne.TextAlignLeading, fyne.TextStyle{Bold:true})
	listLabel4 := widget.NewHyperlink("大小排序", nil)
	listLabel5 := widget.NewLabel("预览")
	listLabel6 := widget.NewLabel("下载")
	listLabel4.OnTapped = func() {
		if SortFlag {
			sort.Slice(Data, func(i, j int) bool {
				if Data[i].SizeInt < Data[j].SizeInt {
					return false
				}
				return true
			})
			List.Refresh()
			SortFlag = false
		}else{
			sort.Slice(Data, func(i, j int) bool {
				if Data[i].SizeInt < Data[j].SizeInt {
					return true
				}
				return false
			})
			List.Refresh()
			SortFlag = true
		}
	}

	c1 := container.NewGridWithColumns(6, listLabel2, listLabel7, listLabel3,listLabel4, listLabel5, listLabel6)
	c := container.NewGridWithColumns(2, listLabel1, c1)
	return c
}

func MakeList() fyne.CanvasObject {
	List = widget.NewList(
		func() int {
			return len(Data)
		},
		func() fyne.CanvasObject {
			urlLabel := widget.NewLabel("hello")

			ContentTypeLabel := widget.NewLabel("hello")
			StatusCodeLabel := widget.NewLabel("hello")
			SizeLabel := widget.NewLabel("hello")
			MethodLabel := widget.NewLabel("hello")
			previewLabel := widget.NewHyperlink("预览", nil)
			downloadLink := widget.NewHyperlink("下载", nil)

			c1 := container.NewGridWithColumns(6, ContentTypeLabel,MethodLabel, StatusCodeLabel,SizeLabel, previewLabel, downloadLink)
			c := container.NewGridWithColumns(2, urlLabel, c1)
			return c
		},
		func(id widget.ListItemID, Item fyne.CanvasObject) {
			if id >= len(Data) {
				return
			}
			d := Data[id]
			Item.(*fyne.Container).Objects[0].(*widget.Label).SetText(d.URI)

			c2 := Item.(*fyne.Container).Objects[1].(*fyne.Container)
			c2.Objects[0].(*widget.Label).SetText(d.ContentType)
			c2.Objects[1].(*widget.Label).SetText(d.Method)
			c2.Objects[2].(*widget.Label).SetText(d.StatusCode)
			c2.Objects[3].(*widget.Label).SetText(d.Size)


			var p,prefix string
			typ := d.ContentType
			if typ != "" {
				prefix = strings.Split(strings.Split(d.ContentType, ";")[0], "/")[0]
			}
			p = d.URI
			if prefix == "video" || prefix == "audio" {
				p = fmt.Sprintf("file:///%s\\%s", DefaultProxy.CachePath,d.CacheName)
			}
			c2.Objects[4].(*widget.Hyperlink).SetURLFromString(p)
			c2.Objects[5].(*widget.Hyperlink).OnTapped = func() {
				srcPath := DefaultProxy.CachePath + "\\" + d.CacheName
				dstPath := DefaultProxy.DownloadPath + "\\" + d.CacheName

				src,err := os.OpenFile(srcPath, os.O_RDONLY, 0755)
				if err != nil {
					dialog.ShowInformation("保存失败!", dstPath, Window)
					return
				}
				dst,err := os.OpenFile(dstPath, os.O_RDWR|os.O_CREATE, 0755)
				if err != nil {
					dialog.ShowInformation("保存失败!", dstPath, Window)
					return
				}
				io.Copy(dst, src)
				dialog.ShowInformation("保存成功!", dstPath, Window)
			}
		},
	)
	//List.OnSelected = func(id widget.ListItemID) {
	//	label.SetText(d)
	//	icon.SetResource(theme.DocumentIcon())
	//}
	//List.OnUnselected = func(id widget.ListItemID) {
	//	label.SetText("Select An Item From The List")
	//	icon.SetResource(nil)
	//}
	//List.Select(12)
	
	return List
}