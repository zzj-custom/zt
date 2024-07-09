package app

import (
	"zt/backend/internal/repository"
	"zt/backend/internal/response"
)

type BingListResponse struct {
	Id            uint   `json:"id"`
	Name          string `json:"name"`
	Copyright     string `json:"copyright"`
	CopyrightLink string `json:"copyrightLink"`
	Url           string `json:"url"`
	Start         string `json:"start"`
	End           string `json:"end"`
	Location      string `json:"location"`
	ClickCount    int    `json:"clickCount"`
	DownloadCount int    `json:"downloadCount"`
}

func (a *App) Images() *response.Reply {
	repo := repository.NewBingImagesRepository()
	result, err := repo.GetAll()
	if err != nil {
		return response.FailReply(response.AcquiredImagesFailed)
	}

	if len(result) == 0 {
		return response.OkReply(make([]*BingListResponse, 0))
	}

	resp := make([]*BingListResponse, 0, len(result))
	for _, item := range result {
		resp = append(resp, &BingListResponse{
			Id:            item.Id,
			Name:          item.Name,
			Copyright:     item.Copyright,
			CopyrightLink: item.CopyrightLink,
			Url:           item.Url,
			Start:         item.Start,
			End:           item.End,
			Location:      item.Location,
			ClickCount:    item.ClickCount,
			DownloadCount: item.DownloadCount,
		})
	}

	return response.OkReply(result)
}
