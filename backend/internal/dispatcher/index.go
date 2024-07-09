package dispatcher

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"sync"
	"time"
	"zt/backend/internal/request"
	"zt/backend/pkg/extractor"
	"zt/backend/pkg/utils"
)

const (
	downloadFileExt = ".download"
)

// Aria2RPCData Aria2 RPC data structure
type Aria2RPCData struct {
	JsonRpc string         `json:"jsonrpc"`
	ID      string         `json:"id"`
	Method  string         `json:"method"`
	Params  [3]interface{} `json:"params"` // secret, uris, options
}

// Aria2Input is options for `aria2.addUri`
// https://aria2.github.io/manual/en/html/aria2c.html#id3
type Aria2Input struct {
	// The file name of the downloaded file
	Out string `json:"out"`
	// For a simple download, only add headers
	Header []string `json:"header"`
}

// FileMetaInfo define file meta info
type FileMetaInfo struct {
	Index float32
	Start int64
	End   int64
	Cur   int64
}

type Options struct {
	Silent         bool
	AudioOnly      bool
	Stream         string
	Refer          string
	OutputPath     string
	OutputName     string
	FileNameLength int
	Caption        bool

	MultiThread  bool
	ThreadNumber int
	RetryTimes   int
	ChunkSizeMB  int
	// Aria2
	UseAria2RPC bool
	Aria2Token  string
	Aria2Method string
	Aria2Addr   string
}

type Downloader struct {
	options *Options
}

func NewDispatcher(options *Options) *Downloader {
	return &Downloader{
		options: options,
	}
}

func (receiver *Downloader) sortStreams(streams map[string]*extractor.Stream) []*extractor.Stream {
	sortedStreams := make([]*extractor.Stream, 0, len(streams))
	for _, data := range streams {
		sortedStreams = append(sortedStreams, data)
	}
	if len(sortedStreams) > 1 {
		sort.SliceStable(sortedStreams, func(i, j int) bool {
			return sortedStreams[i].Size > sortedStreams[j].Size
		})
	}
	return sortedStreams
}

func (receiver *Downloader) Loader(data *extractor.Data) error {
	if len(data.Streams) == 0 {
		return errors.Errorf("no streams in outputName %s", data.Title)
	}

	var (
		outputName = receiver.options.OutputName // 输出文件名
		streamName = receiver.options.Stream     // 文件流名称
	)

	sortedStreams := receiver.sortStreams(data.Streams)
	if outputName == "" {
		outputName = data.Title
	}
	outputName = utils.FileName(outputName, "", receiver.options.FileNameLength)

	if streamName == "" {
		streamName = sortedStreams[0].ID
	}
	stream, ok := data.Streams[streamName]
	if !ok {
		return errors.Errorf("no stream named %s", streamName)
	}

	// only download audio
	if receiver.options.AudioOnly {
		var isFound bool
		reg, err := regexp.Compile("audio+")
		if err != nil {
			return errors.Wrapf(err, "regexp.Compile failed")
		}

		for _, s := range sortedStreams {
			// Looking for the best quality
			if reg.MatchString(s.Quality) {
				isFound = true
				stream = data.Streams[s.ID]
				break
			}
		}
		if !isFound {
			return errors.Errorf("No audio stream found")
		}
	}

	// download caption
	entry := slog.With(
		slog.String("title", outputName),
	)
	if receiver.options.Caption && data.Captions != nil {
		entry.Info("download caption")
		for k, v := range data.Captions {
			if v != nil {
				entry.Info("downloading [%s] caption...", k)
				if err := receiver.caption(v.URL, outputName, v.Ext, v.Transform); err != nil {
					entry.With("error", err).Error("download caption failed")
				}
			}
		}
	}

	// Use aria2 rpc to download
	if receiver.options.UseAria2RPC {
		return receiver.aria2(outputName, stream)
	}

	// Skip the complete file that has been merged
	mergedFilePath, err := utils.FilePath(
		outputName,
		stream.Ext,
		receiver.options.FileNameLength,
		receiver.options.OutputPath,
		false,
	)
	if err != nil {
		return errors.Wrapf(err, "utils.FilePath failed")
	}
	_, exists, err := utils.FileSize(mergedFilePath)
	if err != nil {
		return errors.Wrapf(err, "utils.FileSize failed")
	}
	// 合并之后文件大小已经改变，所以不在匹配文件大小
	if exists {
		slog.With(
			slog.String("mergedFilePath", mergedFilePath),
		).Info("file already exists, skipping")
		return nil
	}

	if len(stream.Parts) == 1 {
		return receiver.singlePart(stream, data, outputName)
	}

	return receiver.multiPart(
		stream,
		data,
		outputName,
		mergedFilePath,
	)
}

func (receiver *Downloader) singlePart(stream *extractor.Stream, data *extractor.Data, outputName string) error {
	// 如果只有一个分片，直接下载
	var err error
	if receiver.options.MultiThread {
		err = receiver.multiThreadSave(stream.Parts[0], data.URL, outputName)
	} else {
		err = receiver.save(stream.Parts[0], data.URL, outputName)
	}

	if err != nil {
		return errors.Wrapf(err, "failed to save %s", outputName)
	}
	return nil
}

func (receiver *Downloader) multiPart(
	stream *extractor.Stream,
	data *extractor.Data,
	outputName string,
	mergedFilePath string,
) error {
	wg := &sync.WaitGroup{}
	errs := make([]error, 0)
	lock := sync.Mutex{}
	parts := make([]string, len(stream.Parts))
	for index, part := range stream.Parts {
		partFileName := fmt.Sprintf("%s[%d]", outputName, index)
		partFilePath, err := utils.FilePath(
			partFileName,
			part.Ext,
			receiver.options.FileNameLength,
			receiver.options.OutputPath,
			false,
		)
		if err != nil {
			return errors.Wrapf(err, "failed to get file path of %s", partFileName)
		}
		parts[index] = partFilePath

		wg.Add(1)
		go func(part *extractor.Part, fileName string) {
			defer wg.Done()
			if receiver.options.MultiThread {
				err = receiver.multiThreadSave(part, data.URL, fileName)
			} else {
				err = receiver.save(part, data.URL, fileName)
			}
			if err != nil {
				lock.Lock()
				errs = append(errs, err)
				lock.Unlock()
			}
		}(part, partFileName)
	}
	wg.Wait()
	if len(errs) > 0 {
		return receiver.ErrorMerge(errs)
	}

	if data.Type != extractor.DataTypeVideo {
		return nil
	}

	if !receiver.options.Silent {
		fmt.Printf("Merging video parts into %s\n", mergedFilePath)
	}
	if stream.Ext != "mp4" || stream.NeedMux {
		return utils.MergeFilesWithSameExtension(parts, mergedFilePath)
	}
	return utils.MergeToMP4(parts, mergedFilePath, outputName)
}

func (receiver *Downloader) caption(url, fileName, ext string, transform func([]byte) ([]byte, error)) error {
	refer := receiver.options.Refer
	if refer == "" {
		refer = url
	}
	body, err := request.DefaultRequest().GetByte(url, refer, nil)
	if err != nil {
		return errors.Wrapf(err, "failed to download caption")
	}

	if transform != nil {
		body, err = transform(body)
		if err != nil {
			return errors.Wrapf(err, "failed to transform caption")
		}
	}

	captionPath, err := utils.FilePath(
		fileName,
		ext,
		receiver.options.FileNameLength,
		receiver.options.OutputPath,
		true,
	)
	if err != nil {
		return errors.Wrapf(err, "failed to get caption file path")
	}
	file, err := os.Create(captionPath)
	if err != nil {
		return errors.Wrapf(err, "failed to create file %s", captionPath)
	}
	defer func() { _ = file.Close() }()

	if _, err = file.Write(body); err != nil {
		return errors.Wrapf(err, "failed to write file %s", captionPath)
	}
	return nil
}

func (receiver *Downloader) aria2(title string, stream *extractor.Stream) error {
	rpcData := Aria2RPCData{
		JsonRpc: "2.0",
		ID:      "zzjlovetl", // 可以修改
		Method:  "aria2.addUri",
	}
	rpcData.Params[0] = "token:" + receiver.options.Aria2Token

	urls := make([]string, 0, len(stream.Parts))
	for _, p := range stream.Parts {
		urls = append(urls, p.URL)
	}

	var inputs Aria2Input
	inputs.Header = append(inputs.Header, "Referer: "+receiver.options.Refer)
	for i := range urls {
		rpcData.Params[1] = urls[i : i+1]
		inputs.Out = fmt.Sprintf("%s[%d].%s", title, i, stream.Parts[0].Ext)
		rpcData.Params[2] = &inputs
		jsonData, err := json.Marshal(rpcData)
		if err != nil {
			return err
		}
		reqURL := fmt.Sprintf("%s://%s/jsonrpc", receiver.options.Aria2Method, receiver.options.Aria2Addr)
		req, err := http.NewRequest(http.MethodPost, reqURL, bytes.NewBuffer(jsonData))
		if err != nil {
			return errors.Wrapf(err, "failed to create aria2 request")
		}
		req.Header.Set("Content-Type", "application/json")

		var client = http.Client{Timeout: 30 * time.Second}
		res, err := client.Do(req)
		if err != nil {
			return errors.Wrapf(err, "failed to send aria2 request")
		}
		// The http Client and Transport guarantee that Body is always
		// non-nil, even on responses without a body or responses with
		// a zero-length body.
		_ = res.Body.Close()
	}
	return nil
}

func (receiver *Downloader) multiThreadSave(dataPart *extractor.Part, refer, fileName string) error {
	filePath, err := utils.FilePath(
		fileName,
		dataPart.Ext,
		receiver.options.FileNameLength,
		receiver.options.OutputPath,
		false,
	)
	if err != nil {
		return errors.Wrapf(err, "multiThreadSave failed to get file path")
	}
	fileSize, exists, err := utils.FileSize(filePath)
	if err != nil {
		return errors.Wrapf(err, "multiThreadSave failed to get file size")
	}

	// Skip segment file
	// TODO: Live video URLs will not return the size
	if exists && fileSize == dataPart.Size {
		slog.With(
			slog.String("filePath", filePath),
			slog.Int("fileSize", int(fileSize)),
		).Info("Skip segment file")
		return nil
	}
	tmpFilePath := filePath + downloadFileExt
	tmpFileSize, tmpExists, err := utils.FileSize(tmpFilePath)
	if err != nil {
		return errors.Wrapf(err, "multiThreadSave failed to get tmp file size")
	}

	// 如果临时文件存在，并且大小一致，那么将临时文件改为正式文件名称
	if tmpExists {
		if tmpFileSize == dataPart.Size {
			return os.Rename(tmpFilePath, filePath)
		}

		if err = os.Remove(tmpFilePath); err != nil {
			return errors.Wrapf(err, "multiThreadSave failed to remove tmp file")
		}
	}

	// 扫描所有片段
	parts, err := receiver.readDirAllFilePart(filePath, fileName, dataPart.Ext)
	if err != nil {
		return errors.Wrapf(err, "multiThreadSave failed to read dir all file part")
	}

	var unfinishedPart []*FileMetaInfo
	savedSize := int64(0)
	if len(parts) > 0 {
		end := int64(-1)
		for i, part := range parts {
			// If some parts are lost, re-insert one part.
			if part.Start-end != 1 {
				newPart := &FileMetaInfo{
					Index: part.Index - 0.000001,
					Start: end + 1,
					End:   part.Start - 1,
					Cur:   end + 1,
				}
				tmp := append([]*FileMetaInfo{}, parts[:i]...)
				tmp = append(tmp, newPart)
				parts = append(tmp, parts[i:]...)
				unfinishedPart = append(unfinishedPart, newPart)
			}
			// When the part has been downloaded in whole, part.Cur is equal to part.End + 1
			if part.Cur <= part.End+1 {
				savedSize += part.Cur - part.Start
				if part.Cur < part.End+1 {
					unfinishedPart = append(unfinishedPart, part)
				}
			} else {
				// The size of this part has been saved greater than the part size, delete it transparently and re-download.
				err = os.Remove(receiver.filePartPath(filePath, part))
				if err != nil {
					return err
				}
				part.Cur = part.Start
				unfinishedPart = append(unfinishedPart, part)
			}
			end = part.End
		}
		if end != dataPart.Size-1 {
			newPart := &FileMetaInfo{
				Index: parts[len(parts)-1].Index + 1,
				Start: end + 1,
				End:   dataPart.Size - 1,
				Cur:   end + 1,
			}
			parts = append(parts, newPart)
			unfinishedPart = append(unfinishedPart, newPart)
		}
	} else {
		var start, end, partSize int64
		var i float32
		partSize = dataPart.Size / int64(receiver.options.ThreadNumber)
		i = 0
		for start < dataPart.Size {
			end = start + partSize - 1
			if end > dataPart.Size {
				end = dataPart.Size - 1
			} else if int(i+1) == receiver.options.ThreadNumber && end < dataPart.Size {
				end = dataPart.Size - 1
			}
			part := &FileMetaInfo{
				Index: i,
				Start: start,
				End:   end,
				Cur:   start,
			}
			parts = append(parts, part)
			unfinishedPart = append(unfinishedPart, part)
			start = end + 1
			i++
		}
	}
	if savedSize > 0 {
		if savedSize == dataPart.Size {
			return receiver.mergeMultiPart(filePath, parts)
		}
	}

	wgp := &sync.WaitGroup{}
	var errs []error
	var mu sync.Mutex
	for _, part := range unfinishedPart {
		wgp.Add(1)
		go func(part *FileMetaInfo) {
			file, err := os.OpenFile(receiver.filePartPath(filePath, part), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
				return
			}
			defer func() {
				_ = file.Close() // nolint
				wgp.Done()
			}()

			var end, chunkSize int64
			headers := map[string]string{
				"Referer": refer,
			}
			if receiver.options.ChunkSizeMB <= 0 {
				chunkSize = part.End - part.Start + 1
			} else {
				chunkSize = int64(receiver.options.ChunkSizeMB) * 1024 * 1024
			}
			remainingSize := part.End - part.Cur + 1
			if part.Cur == part.Start {
				// Only write part to new file.
				err = receiver.writeFilePartMeta(file, part)
				if err != nil {
					mu.Lock()
					errs = append(errs, err)
					mu.Unlock()
					return
				}
			}
			for remainingSize > 0 {
				end = receiver.computeEnd(part.Cur, chunkSize, part.End)
				headers["Range"] = fmt.Sprintf("bytes=%d-%d", part.Cur, end)
				temp := part.Cur
				for i := 0; ; i++ {
					written, err := receiver.writeFile(dataPart.URL, file, headers)
					if err == nil {
						remainingSize -= chunkSize
						break
					} else if i+1 >= receiver.options.RetryTimes {
						mu.Lock()
						errs = append(errs, err)
						mu.Unlock()
						return
					}
					temp += written
					headers["Range"] = fmt.Sprintf("bytes=%d-%d", temp, end)
				}
				part.Cur = end + 1
			}
		}(part)
	}
	wgp.Wait()
	if len(errs) > 0 {
		return receiver.ErrorMerge(errs)
	}
	return receiver.mergeMultiPart(filePath, parts)
}

func (receiver *Downloader) ErrorMerge(err []error) error {
	if len(err) == 1 {
		return err[0]
	}
	var e error
	for _, item := range err {
		e = errors.WithStack(item)
	}
	return e
}

func (receiver *Downloader) readDirAllFilePart(filePath, filename, extname string) ([]*FileMetaInfo, error) {
	dirPath := filepath.Dir(filePath)
	dir, err := os.Open(dirPath)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer func() { _ = dir.Close() }()
	fns, err := dir.Readdir(0)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var metas []*FileMetaInfo
	reg := regexp.MustCompile(fmt.Sprintf("%s.%s.part.+", regexp.QuoteMeta(filename), extname))
	for _, fn := range fns {
		if reg.MatchString(fn.Name()) {
			meta, err := receiver.parseFilePartMeta(path.Join(dirPath, fn.Name()), fn.Size())
			if err != nil {
				return nil, errors.WithStack(err)
			}
			metas = append(metas, meta)
		}
	}
	sort.SliceStable(metas, func(i, j int) bool {
		return metas[i].Index < metas[j].Index
	})
	return metas, nil
}

func (receiver *Downloader) parseFilePartMeta(filepath string, fileSize int64) (*FileMetaInfo, error) {
	meta := new(FileMetaInfo)
	size := binary.Size(*meta)
	file, err := os.OpenFile(filepath, os.O_RDWR, 0666)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer func() { _ = file.Close() }()
	var buf [512]byte
	readSize, err := file.ReadAt(buf[0:size], 0)
	if err != nil && err != io.EOF {
		return nil, errors.WithStack(err)
	}
	if readSize < size {
		return nil, errors.Errorf("the file has been broken, please delete all part files and re-download")
	}
	err = binary.Read(bytes.NewBuffer(buf[:size]), binary.LittleEndian, meta)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	savedSize := fileSize - int64(binary.Size(meta))
	meta.Cur = meta.Start + savedSize
	return meta, nil
}

func (receiver *Downloader) save(part *extractor.Part, refer, fileName string) error {
	filePath, err := utils.FilePath(fileName, part.Ext, receiver.options.FileNameLength, receiver.options.OutputPath, false)
	if err != nil {
		return err
	}
	fileSize, exists, err := utils.FileSize(filePath)
	if err != nil {
		return err
	}
	// Skip segment file
	// TODO: Live video URLs will not return the size
	if exists && fileSize == part.Size {
		return nil
	}

	tempFilePath := filePath + downloadFileExt
	tempFileSize, _, err := utils.FileSize(tempFilePath)
	if err != nil {
		return err
	}
	headers := map[string]string{
		"Referer": refer,
	}
	var (
		file      *os.File
		fileError error
	)
	if tempFileSize > 0 {
		// range start from 0, 0-1023 means the first 1024 bytes of the file
		headers["Range"] = fmt.Sprintf("bytes=%d-", tempFileSize)
		file, fileError = os.OpenFile(tempFilePath, os.O_APPEND|os.O_WRONLY, 0644)
	} else {
		file, fileError = os.Create(tempFilePath)
	}
	if fileError != nil {
		return fileError
	}

	// close and rename temp file at the end of this function
	defer func() {
		// must close the file before rename or it will cause
		// `The process cannot access the file because it is being used by another process.` error.
		if err = file.Close(); err == nil {
			if err = os.Rename(tempFilePath, filePath); err != nil {
				slog.With(
					slog.String("path", filePath),
				).With("err", err).Error("failed to rename file")
			}
		}
	}()

	if receiver.options.ChunkSizeMB > 0 {
		var start, end, chunkSize int64
		chunkSize = int64(receiver.options.ChunkSizeMB) * 1024 * 1024
		remainingSize := part.Size
		if tempFileSize > 0 {
			start = tempFileSize
			remainingSize -= tempFileSize
		}
		chunk := remainingSize / chunkSize
		if remainingSize%chunkSize != 0 {
			chunk++
		}
		var i int64 = 1
		for ; i <= chunk; i++ {
			end = start + chunkSize - 1
			headers["Range"] = fmt.Sprintf("bytes=%d-%d", start, end)
			temp := start
			for i := 0; ; i++ {
				written, err := receiver.writeFile(part.URL, file, headers)
				if err == nil {
					break
				} else if i+1 >= receiver.options.RetryTimes {
					return err
				}
				temp += written
				headers["Range"] = fmt.Sprintf("bytes=%d-%d", temp, end)
				time.Sleep(1 * time.Second)
			}
			start = end + 1
		}
	} else {
		temp := tempFileSize
		for i := 0; ; i++ {
			written, err := receiver.writeFile(part.URL, file, headers)
			if err == nil {
				break
			} else if i+1 >= receiver.options.RetryTimes {
				return err
			}
			temp += written
			headers["Range"] = fmt.Sprintf("bytes=%d-", temp)
			time.Sleep(1 * time.Second)
		}
	}

	return nil
}

func (receiver *Downloader) writeFile(url string, file *os.File, headers map[string]string) (int64, error) {
	res, err := request.DefaultRequest().Do(http.MethodGet, url, nil, headers)
	if err != nil {
		return 0, errors.Wrapf(err, "request error: %s", url)
	}
	defer func() { _ = res.Body.Close() }()

	barWriter := bufio.NewWriter(file)
	// Note that io.Copy reads 32kb(maximum) from input and writes them to output, then repeats.
	// So don't worry about memory.
	written, copyErr := io.Copy(barWriter, res.Body)
	if copyErr != nil && copyErr != io.EOF {
		return written, errors.Errorf("file copy error: %s", copyErr)
	}
	return written, nil
}

func (receiver *Downloader) filePartPath(filepath string, part *FileMetaInfo) string {
	return fmt.Sprintf("%s.part%f", filepath, part.Index)
}

func (receiver *Downloader) mergeMultiPart(filepath string, parts []*FileMetaInfo) error {
	tempFilePath := filepath + downloadFileExt
	tempFile, err := os.OpenFile(tempFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	var partFiles []*os.File
	defer func() {
		for _, f := range partFiles {
			_ = f.Close()
			if os.Remove(f.Name()) != nil {
				slog.With(
					slog.String("path", f.Name()),
				).With("err", err).Error("failed to remove file")
			}
		}
	}()
	for _, part := range parts {
		file, err := os.Open(receiver.filePartPath(filepath, part))
		if err != nil {
			return err
		}
		partFiles = append(partFiles, file)
		_, err = file.Seek(int64(binary.Size(part)), 0)
		if err != nil {
			return err
		}
		_, err = io.Copy(tempFile, file)
		if err != nil {
			return err
		}
	}
	_ = tempFile.Close()
	err = os.Rename(tempFilePath, filepath)
	return err
}

func (receiver *Downloader) writeFilePartMeta(file *os.File, meta *FileMetaInfo) error {
	return binary.Write(file, binary.LittleEndian, meta)
}

func (receiver *Downloader) computeEnd(s, chunkSize, max int64) int64 {
	var end int64
	end = s + chunkSize - 1
	if end > max {
		end = max
	}
	return end
}
