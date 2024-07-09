package extractor

import "slices"

type Options struct {
	// Playlist indicates if we need to extract the whole playlist rather than the single video.
	Playlist bool
	// Items defines wanted items from a playlist. Separated by commas like: 1,5,6,8-10.
	Items string
	// ItemStart defines the starting item of a playlist.
	ItemStart int
	// ItemEnd defines the ending item of a playlist.
	ItemEnd int

	// ThreadNumber defines how many threads will use in the extraction, only works when Playlist is true.
	ThreadNumber int
	Cookie       string

	// EpisodeTitleOnly indicates file name of each bilibili episode doesn't include the playlist title
	EpisodeTitleOnly bool

	YKCode     string // 优酷code
	YKKey      string // 优酷key
	YKPassword string // 优酷密码
}

// Part is the data structure for a single part of the video stream information.
type Part struct {
	URL  string `json:"url"`
	Size int64  `json:"size"`
	Ext  string `json:"ext"`
}

type CaptionPart struct {
	Part
	Transform func([]byte) ([]byte, error) `json:"-"`
}

// Stream is the data structure for each video stream, eg: 720P, 1080P.
type Stream struct {
	// eg: "1080"
	ID string `json:"id"`
	// eg: "1080P xxx"
	Quality string `json:"quality"`
	// [Part: {URL, Size, Ext}, ...]
	// Some video stream have multiple parts,
	// and can also be used to download multiple image files at once
	Parts []*Part `json:"parts"`
	// total size of all urls
	Size int64 `json:"size"`
	// the file extension after video parts merged
	Ext string `json:"ext"`
	// if the parts need mux
	NeedMux bool
}

// DataType indicates the type of extracted data, eg: video or image.
type DataType string

const (
	// DataTypeVideo indicates the type of extracted data is the video.
	DataTypeVideo DataType = "video"
	// DataTypeImage indicates the type of extracted data is the image.
	DataTypeImage DataType = "image"
	// DataTypeAudio indicates the type of extracted data is the audio.
	DataTypeAudio DataType = "audio"
)

// Data is the main data structure for the whole video data.
type Data struct {
	// URL is used to record the address of this download
	URL   string   `json:"url"`
	Site  string   `json:"site"`
	Title string   `json:"title"`
	Type  DataType `json:"type"`
	// each stream has its own Parts and Quality
	Streams map[string]*Stream `json:"streams"`
	// danmaku(弹幕), subtitles, etc
	Captions map[string]*CaptionPart `json:"caption"`
	// Err is used to record whether an error occurred when extracting the list data
	Err error `json:"err"`
}

func (d *Data) FillStreamsData() {
	for id, stream := range d.Streams {
		stream.ID = id
		if stream.Quality == "" {
			stream.Quality = id
		}

		// 生成文件扩展名
		if d.Type == DataTypeVideo && stream.Ext == "" {
			ext := stream.Parts[0].Ext
			if slices.Contains([]string{"ts", "flv", "f4v"}, ext) {
				ext = "mp4"
			}
			stream.Ext = ext
		}

		// 计算总大小
		if stream.Size > 0 {
			continue
		}
		var size int64
		for _, part := range stream.Parts {
			size += part.Size
		}
		stream.Size = size
	}
}
