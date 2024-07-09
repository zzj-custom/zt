package app

import (
	"github.com/pkg/errors"
	"github.com/sagikazarmark/slog-shim"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"zt/backend/internal/repository"
	"zt/backend/internal/response"
	"zt/backend/internal/utils"
	"zt/backend/pkg/ncm"
)

const (
	savePath = "/Users/Apple/Application/github/zt/file"
	dstExt   = ".ncm"
)

type ConvertResponse struct {
	Name          string          `json:"name"`          // 文件名称
	Flag          string          `json:"flag"`          // 文件标识
	MusicID       int             `json:"musicId"`       // 音乐id
	MusicName     string          `json:"musicName"`     // 音乐名称
	Artist        [][]interface{} `json:"artist"`        // [[string,int],] 艺术家
	AlbumID       int             `json:"albumId"`       // 专辑id
	Album         string          `json:"album"`         // 专辑
	AlbumPicDocID interface{}     `json:"albumPicDocId"` // string or int 专辑图片id
	AlbumPic      string          `json:"albumPic"`      // 专辑图片
	BitRate       int             `json:"bitrate"`       // 音乐比特率
	Mp3DocID      string          `json:"mp3DocId"`      // 音乐文件id
	Duration      int             `json:"duration"`      // 音乐时长
	MvID          int             `json:"mvId"`          // 视频id
	Alias         []string        `json:"alias"`         // 音乐别名
	TransNames    []interface{}   `json:"transNames"`    // string or int 翻译名
	Format        string          `json:"format"`        // 音乐格式
	Status        int             `json:"status"`        // 状态
}

type ProcessRequest struct {
	Flag    string `json:"flag"`
	PType   int    `json:"pType"`
	OutPath string `json:"outPath"`
}

type ProcessResponse struct {
	Flag   string `json:"flag"`
	Status int    `json:"status"`
}

func (a *App) ChooseFile() *response.Reply {
	file, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title:           "请选择需要转换的文件",
		ShowHiddenFiles: true,
		Filters: []runtime.FileFilter{
			{
				DisplayName: "All Files",
				Pattern:     "*" + dstExt,
			},
		},
	})
	if err != nil {
		return response.FailReply(response.ChooseFileFail)
	}

	// 上传文件到本地
	sf, err := a.saveFile(file)
	if err != nil {
		slog.With(slog.String("file", file)).With("err", err).Error("保存文件失败")
		return response.FailReply(response.ChooseFileFail)
	}

	// 解析文件
	info, err := a.parseInfo(sf)
	if err != nil {
		slog.With(slog.String("savaFile", sf)).With("err", err).Error("解析文件失败")
		return response.FailReply(response.ChooseFileFail)
	}

	// 文件信息入库
	musicRepo := repository.NewMusicRepository()
	if err = musicRepo.CreateByModel(&repository.Music{
		Flag:       info.Flag,
		MusicId:    info.MusicID,
		MusicName:  info.MusicName,
		Artist:     info.Artist,
		AlbumPic:   info.AlbumPic,
		MusicDocId: info.Mp3DocID,
		Duration:   info.Duration,
		MvId:       info.MvID,
		Format:     info.Format,
		SrcPath:    file,
		DstPath:    sf,
	}); err != nil {
		return response.FailReply(response.ChooseFileFail)
	}
	slog.Info("文件信息入库成功")
	return response.OkReply(info)
}

func (a *App) saveFile(req string) (string, error) {
	if req == "" {
		return "", errors.New("文件路径为空")
	}

	f, err := os.Open(req)
	if err != nil {
		return "", errors.Wrapf(err, "打开文件[%s]失败", req)
	}
	defer func() { _ = f.Close() }()

	sp := filepath.Join(savePath, filepath.Base(req))
	dst, err := os.Create(sp)
	if err != nil {
		return "", errors.Wrapf(err, "创建文件[%s]失败", sp)
	}

	_, err = io.Copy(dst, f)
	if err != nil {
		return "", errors.Wrapf(err, "拷贝文件[%s]到[%s]失败", req, sp)
	}
	return sp, nil
}

func (a *App) ChooseFolder() *response.Reply {
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

	// 获取所有的目标文件
	fs := a.getDstFile(dir)
	if len(fs) == 0 {
		return response.FailReply(response.FolderNotFile)
	}

	insertData := make([]*repository.Music, 0, len(fs))
	resp := make([]*ConvertResponse, 0, len(fs))
	for _, item := range fs {
		// 上传文件到本地
		sf, err := a.saveFile(item)
		if err != nil {
			slog.With(slog.String("file", item)).With("err", err).Error("保存文件失败")
			return response.FailReply(response.ChooseFolderFail)
		}

		// 解析文件
		info, err := a.parseInfo(sf)
		if err != nil {
			slog.With(slog.String("savaFile", sf)).With("err", err).Error("解析文件失败")
			return response.FailReply(response.ChooseFolderFail)
		}

		resp = append(resp, info)

		insertData = append(insertData, &repository.Music{
			Flag:       info.Flag,
			MusicId:    info.MusicID,
			MusicName:  info.MusicName,
			Artist:     info.Artist,
			AlbumPic:   info.AlbumPic,
			MusicDocId: info.Mp3DocID,
			Duration:   info.Duration,
			MvId:       info.MvID,
			Format:     info.Format,
			SrcPath:    item,
			DstPath:    sf,
		})
	}

	// 批量入库
	musicRepo := repository.NewMusicRepository()
	if err = musicRepo.BulkCreate(insertData); err != nil {
		return response.FailReply(response.ChooseFolderFail)
	}
	return response.OkReply(resp)
}

func (a *App) getDstFile(dir string) []string {
	resp := make([]string, 0)
	if dir == "" || !utils.DirExists(dir) {
		return resp
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return resp
	}
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), dstExt) {
			resp = append(resp, filepath.Join(dir, f.Name()))
		}
	}
	return resp
}

func (a *App) parseInfo(req string) (*ConvertResponse, error) {
	n := ncm.NewNcm(req, "")
	mi, err := n.ParseMateInfo()
	if err != nil {
		return nil, errors.Wrapf(err, "解析文件[%s]失败", req)
	}
	return &ConvertResponse{
		Flag:          utils.MustULid(),
		Name:          filepath.Base(req),
		MusicID:       mi.MusicID,
		MusicName:     mi.MusicName,
		Artist:        mi.Artist,
		AlbumID:       mi.AlbumID,
		Album:         mi.Album,
		AlbumPicDocID: mi.AlbumPic,
		AlbumPic:      mi.AlbumPic,
		BitRate:       mi.BitRate,
		Mp3DocID:      mi.Mp3DocID,
		Duration:      mi.Duration,
		MvID:          mi.MvID,
		Alias:         mi.Alias,
		TransNames:    mi.TransNames,
		Format:        mi.Format,
		Status:        repository.MusicStatusIncomplete,
	}, nil
}

func (a *App) Process(req []*ProcessRequest) *response.Reply {
	if len(req) == 0 {
		return response.FailReply(response.ChooseFile)
	}

	var (
		resp = make([]*ProcessResponse, 0, len(req))
		wg   sync.WaitGroup
	)

	for _, item := range req {
		wg.Add(1)
		musicRepo := repository.NewMusicRepository()
		m, err := musicRepo.StatByFlag(item.Flag)
		if err != nil {
			slog.With(
				slog.String("flag", item.Flag),
				slog.String("outPath", item.OutPath),
			).With("err", err).Error("查询文件失败")
			return response.FailReply(response.ProcessFail)
		}

		op := item.OutPath
		if item.PType == repository.MusicPTypeSrc {
			op = filepath.Dir(m.SrcPath)
		}

		go func(item *ProcessRequest) {
			defer func() { wg.Done() }()
			if err = ncm.NewNcm(m.DstPath, op).Process(); err != nil {
				slog.With(
					slog.String("flag", item.Flag),
					slog.String("outPath", op),
				).With("err", err).Error("处理文件失败")
				return
			}

			// 更新status,p_type,parse_path
			if err = musicRepo.UpdateByModel(m.Id, &repository.Music{
				Status:    repository.MusicStatusComplete,
				PType:     item.PType,
				ParsePath: op,
			}); err != nil {
				slog.With(
					slog.String("flag", item.Flag),
					slog.String("parsePath", op),
					slog.Int("pType", item.PType),
				).With("err", err).Error("更新数据失败")
				return
			}
			resp = append(resp, &ProcessResponse{
				Flag:   item.Flag,
				Status: repository.MusicStatusComplete,
			})
		}(item)
	}

	wg.Wait()
	slog.Info("处理完毕")
	return response.OkReply(resp)
}

func (a *App) View(req string) *response.Reply {
	slog.With(slog.String("req", req)).Info("接收到view的参数")
	if req == "" {
		return response.FailReply(response.QueryParamsError)
	}

	musicRepo := repository.NewMusicRepository()
	data, err := musicRepo.StatByFlag(req)
	if err != nil {
		slog.With(slog.String("flag", req)).With("err", err).Error("查询文件失败")
		return response.FailReply(response.MusicNotExists)
	}

	if data.Status == repository.MusicStatusIncomplete {
		slog.With(slog.Int("status", data.Status)).Info("状态不符合")
		return response.FailReply(response.MusicIncomplete)
	}

	slog.With(slog.String("parsePath", data.ParsePath)).Info("打开的文件地址")

	// 使用Wails的runtime包来打开目录选择对话框
	dir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		DefaultDirectory: data.ParsePath,
		Title:            "查看文件", // 对话框的标题
	})
	if err != nil {
		return response.FailReply(response.ChooseDirectoryFail)
	}
	return response.OkReply(dir)
}
