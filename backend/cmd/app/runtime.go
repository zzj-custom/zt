package app

import (
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"zt/backend/internal/response"
)

func (a *App) ChooseDirectory() *response.Reply {
	// 使用Wails的runtime包来打开目录选择对话框
	dir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title:                      "请选择文件夹", // 对话框的标题
		CanCreateDirectories:       true,     // 是否允许创建文件夹
		ResolvesAliases:            true,     // 是否解析别名
		TreatPackagesAsDirectories: true,     // 是否将包视为文件夹
	})
	if err != nil {
		return response.FailReply(response.ChooseDirectoryFail)
	}
	return response.OkReply(dir)
}
