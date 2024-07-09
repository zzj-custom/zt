package app

import (
	"context"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"log/slog"
	"zt/backend/cmd/bootstrapper"
	_ "zt/backend/pkg/extractor/bilibili"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
// NewApp 创建一个新的 App 应用程序
func NewApp() *App {
	return &App{}
}

// Startup is called at application Startup
//
//	在应用程序启动时调用
func (a *App) Startup(ctx context.Context) {
	// Perform your setup here
	// 在这里执行初始化设置

	// 初始化配置
	bootstrapper.Bootstrap("build/config.toml", func(in fsnotify.Event) {})

	// 监听文件上传选择时间

	a.ctx = ctx
}

// DomReady is called after the front-end dom has been loaded
// DomReady 在前端Dom加载完毕后调用
func (a *App) DomReady(ctx context.Context) {
	// Add your action here
	// 在这里添加你的操作
	fmt.Println("DomReady")
}

// BeforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue,
// false will continue shutdown as normal.
// beforeClose在单击窗口关闭按钮或调用runtime.Quit即将退出应用程序时被调用.
// 返回 true 将导致应用程序继续，false 将继续正常关闭。
func (a *App) BeforeClose(ctx context.Context) (prevent bool) {
	return false
}

// Shutdown is called at application termination
// 在应用程序终止时被调用
func (a *App) Shutdown(ctx context.Context) {
	// Perform your teardown here
	// 在此处做一些资源释放的操作
	bootstrapper.Release()
}

func (a *App) eventSelectFile(ctx context.Context) {
	runtime.EventsOn(ctx, "selectFile", func(optionalData ...interface{}) {
		file, err := runtime.OpenFileDialog(ctx, runtime.OpenDialogOptions{
			Title:           "请选择需要转换的文件",
			ShowHiddenFiles: true,
			Filters: []runtime.FileFilter{
				{
					DisplayName: "All Files",
					Pattern:     "*.ncm",
				},
			},
		})
		if err != nil {
			slog.With("error", err).Error("selectFile error")
			return
		}

		if file != "" {
			// 给前端通过事件发送选中的文件
			runtime.EventsEmit(ctx, "selectedFile")
		}
	})
}
