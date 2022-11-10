package pkg

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"net/url"
	"os"
	"strings"
)

// MakeMyMenu 菜单组件
func MakeMyMenu(a fyne.App, w fyne.Window) *fyne.MainMenu {
	saveMenuItem := fyne.NewMenuItem("设置", func() {
		cw := a.NewWindow("设置")
		cw.Resize(fyne.NewSize(600,400))
		cw.CenterOnScreen()

		savePath := widget.NewEntry()
		savePath.SetText(a.Preferences().String("downloadPath"))

		form := &widget.Form{
			SubmitText: "确定",
			CancelText: "取消",
			Items: []*widget.FormItem{
				{Text: "下载路径", Widget: savePath, HintText: "文件保存路径"},
			},
			OnSubmit: func() {
				if _,err := os.Stat(savePath.Text);err != nil {
					return
				}
				a.Preferences().SetString("downloadPath", strings.TrimRight(savePath.Text, "\\"))
				DownloadPath = a.Preferences().String("downloadPath")
				os.MkdirAll(DownloadPath, 0755)
				cw.Close()
			},
			OnCancel: func() {
				cw.Close()
			},
		}

		cw.SetContent(form)
		cw.Show()
	})

	helpMenuItem := fyne.NewMenuItem("doc", func() {
		u, _ := url.Parse("https://github.com/ttmars/nepal")
		_ = a.OpenURL(u)
	})

	// a quit item will be appended to our first (File) menu
	setting := fyne.NewMenu("菜单", saveMenuItem)
	help := fyne.NewMenu("帮助", helpMenuItem)
	mainMenu := fyne.NewMainMenu(
		setting,
		help,
	)
	return mainMenu
}
