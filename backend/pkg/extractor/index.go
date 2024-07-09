package extractor

import (
	"github.com/pkg/errors"
	"net/url"
	"strings"
	"sync"
	"zt/backend/pkg/utils"
)

var (
	bURL = map[string]string{
		"av": "https://www.bilibili.com/video/",
		"BV": "https://www.bilibili.com/video/",
		"ep": "https://www.bilibili.com/bangumi/play/",
	}
	bDomain   = "bilibili"
	domainMap = map[string]string{
		"haokan.baidu.com": "haokan",
		"xhslink.com":      "xiaohongshu",
	}
)

type Scheduler interface {
	Handler(u string, opts Options) ([]*Data, error)
}

type ManagerDownloader struct {
	md  map[string]Scheduler
	mux sync.RWMutex
}

var (
	mdr     *ManagerDownloader
	mdrOnce sync.Once
)

func NewManagerDownloader() *ManagerDownloader {
	mdrOnce.Do(func() {
		mdr = &ManagerDownloader{
			md:  make(map[string]Scheduler),
			mux: sync.RWMutex{},
		}
	})
	return mdr
}

func (receiver *ManagerDownloader) Register(f string, s Scheduler) {
	receiver.mux.Lock()
	defer receiver.mux.Unlock()
	if _, ok := receiver.md[f]; ok {
		return
	}
	receiver.md[f] = s
}

func (receiver *ManagerDownloader) Handler(f string) Scheduler {
	receiver.mux.RLock()
	defer receiver.mux.RUnlock()
	if _, ok := receiver.md[f]; !ok {
		return nil
	}
	return receiver.md[f]
}

func Dispatch(u string, option Options) ([]*Data, error) {
	u = strings.TrimSpace(u)
	var domain string

	bShortLink := utils.MatchOneOf(u, `^(av|BV|ep)\w+`)
	if len(bShortLink) > 1 {
		domain = bDomain
		u = bURL[bShortLink[1]] + u
	} else {
		pru, err := url.ParseRequestURI(u)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		dm, ok := domainMap[pru.Host]
		if ok {
			domain = dm
		} else {
			domain = utils.Domain(pru.Host)
		}
	}
	handler := NewManagerDownloader().Handler(domain)
	if handler == nil {
		return nil, errors.WithStack(errors.New("handler not found"))
	}
	videos, err := handler.Handler(u, option)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	for _, v := range videos {
		v.FillStreamsData()
	}
	return videos, nil
}
