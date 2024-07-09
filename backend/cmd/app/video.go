package app

import (
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"log/slog"
	"zt/backend/internal/dispatcher"
	"zt/backend/internal/response"
	"zt/backend/pkg/extractor"
)

type DownloadRequest struct {
	U        string `json:"u" validate:"required,url"`         // 下载地址
	OutPath  string `json:"outPath" validate:"omitempty,gt=0"` // 保存路径
	PlayList bool   `json:"playList"`                          // 是否是播放列表
	Items    string `json:"items"`                             // 下载范围，定义例如：1,5,6,8-10
}

type ListResponse struct {
	URL      string           `json:"url"`
	Site     string           `json:"site"`
	Title    string           `json:"title"`
	Type     string           `json:"type"`
	Streams  []*Stream        `json:"streams"`
	Captions map[string]*Part `json:"caption"`
}

type Stream struct {
	ID      string  `json:"id"`
	Quality string  `json:"quality"`
	Parts   []*Part `json:"parts"`
	Size    int64   `json:"size"`
	Ext     string  `json:"ext"`
	NeedMux bool
}

type Part struct {
	URL  string `json:"url"`
	Size int64  `json:"size"`
	Ext  string `json:"ext"`
}

func (a *App) List(u string) *response.Reply {
	if err := validator.New().Var(u, "required,url"); err != nil {
		return response.FailReply(response.QueryParamsError)
	}

	videos, err := extractor.Dispatch(u, extractor.Options{})
	if err != nil {
		logrus.WithField("url", u).WithError(err).Error("get video list error")
		return response.FailReply(response.AcquiredVideoList)
	}

	resp := make([]*ListResponse, 0, len(videos))
	for _, video := range videos {
		streams := make([]*Stream, 0, len(video.Streams))
		if video.Streams != nil {
			for _, stream := range video.Streams {
				parts := make([]*Part, 0, len(stream.Parts))
				if stream.Parts != nil {
					for _, part := range stream.Parts {
						parts = append(parts, &Part{
							URL:  part.URL,
							Size: part.Size,
							Ext:  part.Ext,
						})
					}
				}

				streams = append(streams, &Stream{
					ID:      stream.ID,
					Quality: stream.Quality,
					Parts:   parts,
					Size:    stream.Size,
					Ext:     stream.Ext,
					NeedMux: stream.NeedMux,
				})
			}
		}

		captions := make(map[string]*Part)
		if video.Captions != nil {
			for key, caption := range video.Captions {
				if caption == nil {
					continue
				}
				captions[key] = &Part{
					URL:  caption.URL,
					Size: caption.Size,
					Ext:  caption.Ext,
				}
			}
		}

		resp = append(resp, &ListResponse{
			URL:      video.URL,
			Site:     video.Site,
			Title:    video.Title,
			Type:     string(video.Type),
			Streams:  streams,
			Captions: captions,
		})
	}

	return response.OkReply(resp)
}

func (a *App) Download(req *DownloadRequest) *response.Reply {
	if err := validator.New().Struct(req); err != nil {
		slog.With("req", req, "err", err).Error("请求参数错误")
		return response.FailReply(response.QueryParamsError)
	}

	videos, err := extractor.Dispatch(req.U, extractor.Options{
		Playlist: req.PlayList,
		Items:    req.Items,
	})
	if err != nil || len(videos) == 0 {
		logrus.WithField("url", req.U).WithError(err).Error("获取视频列表失败")
		return response.FailReply(response.AcquiredVideoList)
	}

	dispatch := dispatcher.NewDispatcher(&dispatcher.Options{
		OutputPath: req.OutPath,
	})
	for _, item := range videos {
		if dispatch.Loader(item) != nil {
			logrus.WithField("url", req.U).WithError(err).Error("下载视频失败")
			return response.FailReply(response.DownloadError)
		}
	}
	return response.OkReply(videos)
}
