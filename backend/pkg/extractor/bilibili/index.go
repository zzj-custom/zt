package bilibili

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"zt/backend/internal/request"
	"zt/backend/pkg/extractor"
	"zt/backend/pkg/utils"
)

const (
	bAPI       = "https://api.bilibili.com/x/player/playurl?"
	bangumiAPI = "https://api.bilibili.com/pgc/player/web/playurl?"
	bTokenAPI  = "https://api.bilibili.com/x/player/playurl/token?"
	referer    = "https://www.bilibili.com"
)

var (
	bName = "bilibili"
	b23   = "b23"
)

func init() {
	extractor.NewManagerDownloader().Register(bName, &Bili{})
	extractor.NewManagerDownloader().Register(b23, &Bili{})
}

type Bili struct{}

func (b Bili) Handler(u string, option extractor.Options) ([]*extractor.Data, error) {
	html, err := request.DefaultRequest().Get(u, referer, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// 手动将线程数设置为1以避免http 412错误
	option.ThreadNumber = 1

	if strings.Contains(u, "bangumi") {
		// handle bangumi
		return extractBangumi(u, html, option)
	} else if strings.Contains(u, "festival") {
		return extractFestival(u, html, option)
	} else {
		// handle normal video
		return extractNormalVideo(u, html, option)
	}
}

func extractBangumi(url, html string, extractOption extractor.Options) ([]*extractor.Data, error) {
	dataString := utils.MatchOneOf(html, `<script\s+id="__NEXT_DATA__"\s+type="application/json"\s*>(.*?)</script\s*>`)[1]
	epArrayString := utils.MatchOneOf(dataString, `"episodes"\s*:\s*(.+?)\s*,\s*"user_status"`)[1]
	fullVideoIdString := utils.MatchOneOf(dataString, `"videoId"\s*:\s*"(ep|ss)(\d+)"`)
	epSsString := fullVideoIdString[1] // "ep" or "ss"
	videoIdString := fullVideoIdString[2]

	var epArray []json.RawMessage
	err := json.Unmarshal([]byte(epArrayString), &epArray)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var data bangumiData
	for _, jsonByte := range epArray {
		var epInfo bangumiEpData
		err := json.Unmarshal(jsonByte, &epInfo)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		videoId, err := strconv.ParseInt(videoIdString, 10, 0)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if epInfo.ID == int(videoId) || (epSsString == "ss" && epInfo.TitleFormat == "第1话") {
			data.EpInfo = epInfo
		}
		data.EpList = append(data.EpList, epInfo)
	}

	sort.Slice(data.EpList, func(i, j int) bool {
		return data.EpList[i].EpID < data.EpList[j].EpID
	})

	if !extractOption.Playlist {
		aid := data.EpInfo.Aid
		cid := data.EpInfo.Cid
		bvid := data.EpInfo.BVid
		titleFormat := data.EpInfo.TitleFormat
		longTitle := data.EpInfo.LongTitle
		if aid <= 0 || cid <= 0 || bvid == "" {
			aid = data.EpList[0].Aid
			cid = data.EpList[0].Cid
			bvid = data.EpList[0].BVid
			titleFormat = data.EpList[0].TitleFormat
			longTitle = data.EpList[0].LongTitle
		}
		options := bOptions{
			url:     url,
			html:    html,
			bangumi: true,
			aid:     aid,
			cid:     cid,
			bvid:    bvid,

			subtitle: fmt.Sprintf("%s %s", titleFormat, longTitle),
		}
		return []*extractor.Data{download(options, extractOption)}, nil
	}

	// handle bangumi playlist
	needDownloadItems := utils.NeedDownloadList(extractOption.Items, extractOption.ItemStart, extractOption.ItemEnd, len(data.EpList))
	extractedData := make([]*extractor.Data, len(needDownloadItems))
	wg := &sync.WaitGroup{}
	dataIndex := 0
	for index, u := range data.EpList {
		if !slices.Contains(needDownloadItems, index+1) {
			continue
		}
		wg.Add(1)
		id := u.EpID
		if id == 0 {
			id = u.ID
		}
		// html content can't be reused here
		options := bOptions{
			url:     fmt.Sprintf("https://www.bilibili.com/bangumi/play/ep%d", id),
			bangumi: true,
			aid:     u.Aid,
			cid:     u.Cid,
			bvid:    u.BVid,

			subtitle: fmt.Sprintf("%s %s", u.TitleFormat, u.LongTitle),
		}
		go func(index int, options bOptions, extractedData []*extractor.Data) {
			defer wg.Done()
			extractedData[index] = download(options, extractOption)
		}(dataIndex, options, extractedData)
		dataIndex++
	}
	wg.Wait()
	return extractedData, nil
}

func extractFestival(url, html string, extractOption extractor.Options) ([]*extractor.Data, error) {
	matches := utils.MatchAll(html, "<\\s*script[^>]*>\\s*window\\.__INITIAL_STATE__=([\\s\\S]*?);\\s?\\(function[\\s\\S]*?<\\/\\s*script\\s*>")
	if len(matches) < 1 {
		return nil, errors.Wrapf(errors.New("could not find video in page"), "html: %s", html)
	}
	if len(matches[0]) < 2 {
		return nil, errors.New("could not find video in page")
	}

	var festivalData festival
	err := json.Unmarshal([]byte(matches[0][1]), &festivalData)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	options := bOptions{
		url:  url,
		html: html,
		aid:  festivalData.VideoInfo.Aid,
		bvid: festivalData.VideoInfo.BVid,
		cid:  festivalData.VideoInfo.Cid,
		page: 0,
	}

	return []*extractor.Data{download(options, extractOption)}, nil
}

func getMultiPageData(html string) (*multiPage, error) {
	var data multiPage
	multiPageDataString := utils.MatchOneOf(
		html, `window.__INITIAL_STATE__=(.+?);\(function`,
	)
	if multiPageDataString == nil {
		return &data, errors.New("this page has no playlist")
	}
	err := json.Unmarshal([]byte(multiPageDataString[1]), &data)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &data, nil
}

func extractNormalVideo(url, html string, extractOption extractor.Options) ([]*extractor.Data, error) {
	pageData, err := getMultiPageData(html)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if !extractOption.Playlist {
		pageString := utils.MatchOneOf(url, `\?p=(\d+)`)
		var p int
		if pageString == nil {
			// https://www.bilibili.com/video/av20827366/
			p = 1
		} else {
			// https://www.bilibili.com/video/av20827366/?p=2
			p, _ = strconv.Atoi(pageString[1])
		}

		if len(pageData.VideoData.Pages) < p || p < 1 {
			return nil, errors.Wrapf(errors.New("could not find video in page"), "html: %s", html)
		}

		page := pageData.VideoData.Pages[p-1]
		options := bOptions{
			url:  url,
			html: html,
			aid:  pageData.Aid,
			bvid: pageData.BVid,
			cid:  page.Cid,
			page: p,
		}
		// "part":"" or "part":"Untitled"
		if page.Part == "Untitled" || len(pageData.VideoData.Pages) == 1 {
			options.subtitle = ""
		} else {
			options.subtitle = page.Part
		}
		return []*extractor.Data{
			download(options, extractOption),
		}, nil
	}

	// handle normal video playlist
	if len(pageData.Sections) == 0 {
		// https://www.bilibili.com/video/av20827366/?p=* each video in playlist has different p=?
		return multiPageDownload(url, html, extractOption, pageData)
	}
	// handle another kind of playlist
	// https://www.bilibili.com/video/av*** each video in playlist has different av/bv id
	return multiEpisodeDownload(url, html, extractOption, pageData)
}

// download is the download function for a single URL
func download(options bOptions, extractOption extractor.Options) *extractor.Data {
	var (
		err  error
		html string
	)
	if options.html != "" {
		// reuse html string, but this can't be reused in case of playlist
		html = options.html
	} else {
		html, err = request.DefaultRequest().Get(options.url, referer, nil)
		if err != nil {
			return EmptyData(options.url, err)
		}
	}

	// Get "accept_quality" and "accept_description"
	// "accept_description":["超高清 8K","超清 4K","高清 1080P+","高清 1080P","高清 720P","清晰 480P","流畅 360P"],
	// "accept_quality":[127，120,112,80,48,32,16],
	api, err := genAPI(
		options.aid,
		options.cid,
		127,
		options.bvid,
		options.bangumi,
		extractOption.Cookie,
	)
	if err != nil {
		return EmptyData(options.url, err)
	}
	jsonString, err := request.DefaultRequest().Get(api, referer, nil)
	if err != nil {
		return EmptyData(options.url, err)
	}

	var data dash
	err = json.Unmarshal([]byte(jsonString), &data)
	if err != nil {
		return EmptyData(options.url, err)
	}
	var dashData dashInfo
	if data.Data.Description == nil {
		dashData = data.Result
	} else {
		dashData = data.Data
	}

	var audioPart *extractor.Part
	if dashData.Streams.Audio != nil {
		// Get audio part
		var audioID int
		audios := map[int]string{}
		bandwidth := 0
		for _, stream := range dashData.Streams.Audio {
			if stream.Bandwidth > bandwidth {
				audioID = stream.ID
				bandwidth = stream.Bandwidth
			}
			audios[stream.ID] = stream.BaseURL
		}
		s, err := request.DefaultRequest().Size(audios[audioID], referer)
		if err != nil {
			return EmptyData(options.url, err)
		}
		audioPart = &extractor.Part{
			URL:  audios[audioID],
			Size: s,
			Ext:  "m4a",
		}
	}

	streams := make(map[string]*extractor.Stream, len(dashData.Quality))
	for _, stream := range dashData.Streams.Video {
		s, err := request.DefaultRequest().Size(stream.BaseURL, referer)
		if err != nil {
			return EmptyData(options.url, err)
		}
		parts := make([]*extractor.Part, 0, 2)
		parts = append(parts, &extractor.Part{
			URL:  stream.BaseURL,
			Size: s,
			Ext:  getExtFromMimeType(stream.MimeType),
		})
		if audioPart != nil {
			parts = append(parts, audioPart)
		}
		var size int64
		for _, part := range parts {
			size += part.Size
		}
		id := fmt.Sprintf("%d-%d", stream.ID, stream.Codecid)
		streams[id] = &extractor.Stream{
			Parts:   parts,
			Size:    size,
			Quality: fmt.Sprintf("%s %s", qualityString[stream.ID], stream.Codecs),
		}
		if audioPart != nil {
			streams[id].NeedMux = true
		}
	}

	for _, durl := range dashData.DURLs {
		var ext string
		switch dashData.DURLFormat {
		case "flv", "flv480":
			ext = "flv"
		case "mp4", "hdmp4": // nolint
			ext = "mp4"
		}

		parts := make([]*extractor.Part, 0, 1)
		parts = append(parts, &extractor.Part{
			URL:  durl.URL,
			Size: durl.Size,
			Ext:  ext,
		})

		streams[strconv.Itoa(dashData.CurQuality)] = &extractor.Stream{
			Parts:   parts,
			Size:    durl.Size,
			Quality: qualityString[dashData.CurQuality],
		}
	}

	// get the title
	doc, err := utils.GetDoc(html)
	if err != nil {
		return EmptyData(options.url, err)
	}
	title := utils.Title(doc)
	if options.subtitle != "" {
		pageString := ""
		if options.page > 0 {
			pageString = fmt.Sprintf("P%d ", options.page)
		}
		if extractOption.EpisodeTitleOnly {
			title = fmt.Sprintf("%s%s", pageString, options.subtitle)
		} else {
			title = fmt.Sprintf("%s %s%s", title, pageString, options.subtitle)
		}
	}

	return &extractor.Data{
		Site:    "哔哩哔哩 bilibili.com",
		Title:   title,
		Type:    extractor.DataTypeVideo,
		Streams: streams,
		Captions: map[string]*extractor.CaptionPart{
			"danmaku": {
				Part: extractor.Part{
					URL: fmt.Sprintf("https://comment.bilibili.com/%d.xml", options.cid),
					Ext: "xml",
				},
			},
			"subtitle": getSubTitleCaptionPart(options.aid, options.cid),
		},
		URL: options.url,
	}
}

// handle multi page download
func multiPageDownload(url, html string, extractOption extractor.Options, pageData *multiPage) ([]*extractor.Data, error) {
	needDownloadItems := utils.NeedDownloadList(
		extractOption.Items,
		extractOption.ItemStart,
		extractOption.ItemEnd,
		len(pageData.VideoData.Pages),
	)
	extractedData := make([]*extractor.Data, len(needDownloadItems))
	wg := &sync.WaitGroup{}
	dataIndex := 0
	for index, u := range pageData.VideoData.Pages {
		if !slices.Contains(needDownloadItems, index+1) {
			continue
		}
		wg.Add(1)
		options := bOptions{
			url:      url,
			html:     html,
			aid:      pageData.Aid,
			bvid:     pageData.BVid,
			cid:      u.Cid,
			subtitle: u.Part,
			page:     u.Page,
		}
		go func(index int, options bOptions, extractedData []*extractor.Data) {
			defer wg.Done()
			extractedData[index] = download(options, extractOption)
		}(dataIndex, options, extractedData)
		dataIndex++
	}
	wg.Wait()
	return extractedData, nil
}

// handle multi episode download
func multiEpisodeDownload(url, html string, extractOption extractor.Options, pageData *multiPage) ([]*extractor.Data, error) {
	needDownloadItems := utils.NeedDownloadList(extractOption.Items, extractOption.ItemStart, extractOption.ItemEnd, len(pageData.Sections[0].Episodes))
	extractedData := make([]*extractor.Data, len(needDownloadItems))
	wg := &sync.WaitGroup{}
	dataIndex := 0
	for index, u := range pageData.Sections[0].Episodes {
		if !slices.Contains(needDownloadItems, index+1) {
			continue
		}
		wg.Add(1)
		options := bOptions{
			url:      url,
			html:     html,
			aid:      u.Aid,
			bvid:     u.BVid,
			cid:      u.Cid,
			subtitle: fmt.Sprintf("%s P%d", u.Title, index+1),
		}
		go func(index int, options bOptions, extractedData []*extractor.Data) {
			defer wg.Done()
			extractedData[index] = download(options, extractOption)
		}(dataIndex, options, extractedData)
		dataIndex++
	}
	wg.Wait()
	return extractedData, nil
}

// EmptyData returns an "empty" Data object with the given URL and error.
func EmptyData(url string, err error) *extractor.Data {
	return &extractor.Data{
		URL: url,
		Err: err,
	}
}

func getExtFromMimeType(mimeType string) string {
	exts := strings.Split(mimeType, "/")
	if len(exts) == 2 {
		return exts[1]
	}
	return "mp4"
}

func getSubTitleCaptionPart(aid int, cid int) *extractor.CaptionPart {
	jsonString, err := request.NewRequest(&request.ReqOptions{}).Get(
		fmt.Sprintf("http://api.bilibili.com/x/player/wbi/v2?aid=%d&cid=%d", aid, cid), referer, nil,
	)
	if err != nil {
		return nil
	}
	stu := bilibiliWebInterface{}
	err = json.Unmarshal([]byte(jsonString), &stu)
	if err != nil || len(stu.Data.SubtitleInfo.SubtitleList) == 0 {
		return nil
	}
	return &extractor.CaptionPart{
		Part: extractor.Part{
			URL: fmt.Sprintf("https:%s", stu.Data.SubtitleInfo.SubtitleList[0].SubtitleUrl),
			Ext: "srt",
		},
		Transform: subtitleTransform,
	}
}

func subtitleTransform(body []byte) ([]byte, error) {
	bytes := ""
	captionData := bilibiliSubtitleFormat{}
	err := json.Unmarshal(body, &captionData)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for i := 0; i < len(captionData.Body); i++ {
		bytes += fmt.Sprintf("%d\n%s --> %s\n%s\n\n",
			i,
			time.Unix(0, int64(captionData.Body[i].From*1000)*int64(time.Millisecond)).UTC().Format("15:04:05.000"),
			time.Unix(0, int64(captionData.Body[i].To*1000)*int64(time.Millisecond)).UTC().Format("15:04:05.000"),
			captionData.Body[i].Content,
		)
	}
	return []byte(bytes), nil
}

func genAPI(aid, cid, quality int, bvid string, bangumi bool, cookie string) (string, error) {
	var (
		err        error
		baseAPIURL string
		params     string
		utoken     string
	)
	if cookie != "" && utoken == "" {
		utoken, err = request.NewRequest(&request.ReqOptions{}).Get(
			fmt.Sprintf("%said=%d&cid=%d", bTokenAPI, aid, cid),
			referer,
			nil,
		)
		if err != nil {
			return "", err
		}
		var t token
		err = json.Unmarshal([]byte(utoken), &t)
		if err != nil {
			return "", err
		}
		if t.Code != 0 {
			return "", errors.Errorf("cookie error: %s", t.Message)
		}
		utoken = t.Data.Token
	}
	var api string
	if bangumi {
		// The parameters need to be sorted by name
		// qn=0 flag makes the CDN address different every time
		// quality=120(4k) is the highest quality so far
		params = fmt.Sprintf(
			"cid=%d&bvid=%s&qn=%d&type=&otype=json&fourk=1&fnver=0&fnval=16",
			cid, bvid, quality,
		)
		baseAPIURL = bangumiAPI
	} else {
		params = fmt.Sprintf(
			"avid=%d&cid=%d&bvid=%s&qn=%d&type=&otype=json&fourk=1&fnver=0&fnval=2000",
			aid, cid, bvid, quality,
		)
		baseAPIURL = bAPI
	}
	api = baseAPIURL + params
	// bangumi utoken also need to put in params to sign, but the ordinary video doesn't need
	if !bangumi && utoken != "" {
		api = fmt.Sprintf("%s&utoken=%s", api, utoken)
	}
	return api, nil
}
