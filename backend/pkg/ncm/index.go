package ncm

import (
	"crypto/aes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"github.com/bogem/id3v2"
	"github.com/go-flac/flacpicture"
	"github.com/go-flac/flacvorbis"
	"github.com/go-flac/go-flac"
	"github.com/pkg/errors"
	"io"
	"log/slog"
	"mime/multipart"
	"os"
	"path/filepath"
	"reflect"
	"zt/backend/pkg/request"
)

var (
	aesCoreKey   = []byte{0x68, 0x7A, 0x48, 0x52, 0x41, 0x6D, 0x73, 0x6F, 0x35, 0x6B, 0x49, 0x6E, 0x62, 0x61, 0x78, 0x57}
	aesModifyKey = []byte{0x23, 0x31, 0x34, 0x6C, 0x6A, 0x6B, 0x5F, 0x21, 0x5C, 0x5D, 0x26, 0x30, 0x55, 0x3C, 0x27, 0x28}
	mp3Type      = "mp3"
	flacType     = "flac"
	jpegMime     = "image/jpeg"
	pngMime      = "image/png"
)

type MetaInfo struct {
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
}

type Ncm struct {
	Path    string `json:"path"`
	OutPath string `json:"outPath"`
}

func NewNcm(path string, outPath string) *Ncm {
	return &Ncm{
		Path:    path,
		OutPath: outPath,
	}
}

func (receiver *Ncm) validate(fp *os.File, rBuf []byte) error {
	uLen := receiver.readUint32(rBuf, fp)

	if uLen != 0x4e455443 {
		return errors.New("isn't netEase cloud music copyright file!")
	}

	uLen = receiver.readUint32(rBuf, fp)
	if uLen != 0x4d414446 {
		return errors.New("isn't netEase cloud music copyright file!")
	}
	return nil
}

func (receiver *Ncm) parseInfo(fp *os.File, rBuf []byte) (MetaInfo, []byte, error) {
	var meta MetaInfo
	err := receiver.validate(fp, rBuf)
	if err != nil {
		return meta, nil, errors.Wrapf(err, "validate fail")
	}

	// whence, 0 - 表示相对于文件开头的位置。1 - 表示相对于当前文件指针的位置。2 - 表示相对于文件末尾的位置。
	// offset是相较于whence的偏移量
	_, err = fp.Seek(2, 1)
	if err != nil {
		return meta, nil, errors.Wrap(err, "seek fail")
	}
	uLen := receiver.readUint32(rBuf, fp)

	var keyData = make([]byte, uLen)
	_, err = fp.Read(keyData)
	if err != nil {
		return meta, nil, errors.Wrap(err, "read key data fail")
	}

	for i := range keyData {
		keyData[i] ^= 0x64
	}

	deKeyData, err := receiver.decryptAes128Ecb(aesCoreKey, receiver.fixBlockSize(keyData))
	if err != nil {
		return meta, nil, errors.Wrap(err, "decrypt aes128 data fail")
	}

	// 去掉前缀（neteasecloudmusic）
	deKeyData = deKeyData[17:]

	uLen = receiver.readUint32(rBuf, fp)
	var modifyData = make([]byte, uLen)
	_, err = fp.Read(modifyData)
	if err != nil {
		return meta, nil, errors.Wrap(err, "read modify data fail")
	}

	for i := range modifyData {
		modifyData[i] ^= 0x63
	}
	deModifyData := make([]byte, base64.StdEncoding.DecodedLen(len(modifyData)-22))
	_, err = base64.StdEncoding.Decode(deModifyData, modifyData[22:])
	if err != nil {
		return meta, nil, errors.Wrap(err, "base64 decode fail")
	}

	deData, err := receiver.decryptAes128Ecb(aesModifyKey, receiver.fixBlockSize(deModifyData))
	if err != nil {
		return meta, nil, errors.Wrap(err, "decrypt aes128 data fail")
	}

	// 去掉前缀（music:）
	deData = deData[6:]
	err = json.Unmarshal(deData, &meta)
	if err != nil {
		return meta, nil, errors.Wrap(err, "json unmarshal music fail")
	}
	return meta, deKeyData, nil
}

func (receiver *Ncm) ParseMateInfo() (*MetaInfo, error) {
	fp, err := os.Open(receiver.Path)
	if err != nil {
		return nil, errors.Wrap(err, "open file fail")
	}
	defer func() { _ = fp.Close() }()

	var rBuf = make([]byte, 4)

	meta, _, err := receiver.parseInfo(fp, rBuf)
	if err != nil {
		return nil, errors.Wrap(err, "parse info fail")
	}
	return &meta, nil
}

func (receiver *Ncm) Process() error {
	fp, err := os.Open(receiver.Path)
	if err != nil {
		return errors.Wrap(err, "open file fail")
	}
	defer func() { _ = fp.Close() }()

	var rBuf = make([]byte, 4)

	meta, deKeyData, err := receiver.parseInfo(fp, rBuf)
	if err != nil {
		return errors.Wrap(err, "parse info fail")
	}

	// crc32 check
	_, err = fp.Seek(9, 1)
	if err != nil {
		return errors.Wrap(err, "seek fail")
	}

	imgLen := receiver.readUint32(rBuf, fp)
	imgData := func() []byte {
		if imgLen > 0 {
			data := make([]byte, imgLen)
			_, err = fp.Read(data)
			if err != nil {
				return nil
			}
			return data
		}
		return nil
	}()

	if err = receiver.save(meta, fp, deKeyData); err != nil {
		return errors.Wrap(err, "save music fail")
	}
	slog.With(
		slog.String("musicName", meta.MusicName),
	).Info("save music success")

	switch meta.Format {
	case mp3Type:
		return receiver.addMP3Tag(imgData, &meta, fp)
	case flacType:
		return receiver.addFLACTag(imgData, &meta, fp)
	}
	return nil
}

func (receiver *Ncm) outPathName(meta *MetaInfo, name string) string {
	if name == "" {
		return filepath.Join(receiver.OutPath, meta.MusicName+"."+meta.Format)
	}

	return filepath.Join(receiver.OutPath, receiver.filename(name)+"."+meta.Format)
}

func (receiver *Ncm) save(meta MetaInfo, fp *os.File, deKeyData []byte) error {
	outputName := receiver.outPathName(&meta, fp.Name())
	fpOut, err := os.OpenFile(outputName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return errors.Wrap(err, "open outfile fail")
	}
	defer func() { _ = fpOut.Close() }()

	box := receiver.buildKeyBox(deKeyData)
	n := 0x8000

	var tb = make([]byte, n)
	for {
		_, err = fp.Read(tb)
		if err == io.EOF { // read EOF
			break
		}

		if err != nil {
			return errors.Wrap(err, "read origin file fail")
		}

		for i := 0; i < n; i++ {
			j := byte((i + 1) & 0xff)
			tb[i] ^= box[(box[j]+box[(box[j]+j)&0xff])&0xff]
		}
		_, err = fpOut.Write(tb)
		if err != nil {
			return errors.Wrap(err, "write outfile fail")
		}
	}
	return nil
}

func (receiver *Ncm) readUint32(rBuf []byte, fp multipart.File) uint32 {
	_, err := fp.Read(rBuf)
	if err != nil {
		slog.With(
			slog.String("rBuf", string(rBuf)),
		).With("err", err).Error("read uint32 fail")
		return 0
	}
	return binary.LittleEndian.Uint32(rBuf)
}

func (receiver *Ncm) decryptAes128Ecb(key, data []byte) ([]byte, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	dataLen := len(data)
	decrypted := make([]byte, dataLen)
	bs := block.BlockSize()
	for i := 0; i <= dataLen-bs; i += bs {
		block.Decrypt(decrypted[i:i+bs], data[i:i+bs])
	}
	return receiver.pkcs7UnPadding(decrypted), nil
}

func (receiver *Ncm) pkcs7UnPadding(src []byte) []byte {
	length := len(src)
	return src[:(length - int(src[length-1]))]
}

func (receiver *Ncm) fixBlockSize(src []byte) []byte {
	return src[:len(src)/aes.BlockSize*aes.BlockSize]
}

func (receiver *Ncm) buildKeyBox(key []byte) []byte {
	box := make([]byte, 256)
	for i := 0; i < 256; i++ {
		box[i] = byte(i)
	}
	keyLen := byte(len(key))
	var c, lastByte, keyOffset byte
	for i := 0; i < 256; i++ {
		c = (box[i] + lastByte + key[keyOffset]) & 0xff
		keyOffset++
		if keyOffset >= keyLen {
			keyOffset = 0
		}
		box[i], box[c] = box[c], box[i]
		lastByte = c
	}
	return box
}

func (receiver *Ncm) containPNGHeader(data []byte) bool {
	if len(data) < 8 {
		return false
	}
	return string(data[:8]) == string([]byte{137, 80, 78, 71, 13, 10, 26, 10})
}

func (receiver *Ncm) addFLACTag(imgData []byte, meta *MetaInfo, fp *os.File) error {
	outPathName := receiver.outPathName(meta, fp.Name())
	f, err := flac.ParseFile(outPathName)
	if err != nil {
		return errors.Wrap(err, "parse flac file fail")
	}

	imgData = receiver.imageInfo(imgData, meta)

	f = receiver.addFLACMeta(f, imgData, meta)

	cmtmeta := receiver.cmtMeta(f)

	cmts, err := receiver.cmts(f)
	if err != nil {
		return errors.Wrap(err, "get cmt meta fail")
	}

	cmts = receiver.FLACMusicName(cmts, meta)
	cmts = receiver.FLACAlbum(cmts, meta)
	cmts = receiver.FLACArticle(cmts, meta)

	res := cmts.Marshal()
	if cmtmeta != nil {
		*cmtmeta = res
	} else {
		f.Meta = append(f.Meta, &res)
	}

	return f.Save(outPathName)
}

func (receiver *Ncm) FLACArticle(
	cmts *flacvorbis.MetaDataBlockVorbisComment,
	meta *MetaInfo,
) *flacvorbis.MetaDataBlockVorbisComment {
	artists, err := cmts.Get(flacvorbis.FIELD_ARTIST)
	if err != nil {
		slog.With("err", err).Error("get artist fail")
		return cmts
	}
	if len(artists) > 0 {
		return cmts
	}

	if len(meta.Artist) < 1 {
		return cmts
	}
	for _, artist := range meta.Artist {
		if astr, ok := artist[0].(string); ok {
			if err := cmts.Add(flacvorbis.FIELD_ARTIST, astr); err != nil {
				slog.With("err", err).Error("add artist fail")
				return cmts
			}
		}
	}

	return cmts
}

func (receiver *Ncm) FLACAlbum(
	cmts *flacvorbis.MetaDataBlockVorbisComment,
	meta *MetaInfo,
) *flacvorbis.MetaDataBlockVorbisComment {
	albums, err := cmts.Get(flacvorbis.FIELD_ALBUM)
	if err != nil {
		slog.With("err", err).Error("get album fail")
		return cmts
	}

	if len(albums) == 0 && meta.Album != "" {
		if err = cmts.Add(flacvorbis.FIELD_ALBUM, meta.Album); err != nil {
			slog.With("err", err).Error("add album fail")
			return cmts
		}
	}

	return cmts
}

func (receiver *Ncm) FLACMusicName(
	cmts *flacvorbis.MetaDataBlockVorbisComment,
	meta *MetaInfo,
) *flacvorbis.MetaDataBlockVorbisComment {
	titles, err := cmts.Get(flacvorbis.FIELD_TITLE)
	if err != nil {
		slog.With("err", err).Error("get title fail")
		return cmts
	}

	if len(titles) == 0 && meta.MusicName != "" {
		if err = cmts.Add(flacvorbis.FIELD_TITLE, meta.MusicName); err != nil {
			slog.With("err", err).Error("add title fail")
			return cmts
		}
	}

	return cmts
}

func (receiver *Ncm) cmts(f *flac.File) (*flacvorbis.MetaDataBlockVorbisComment, error) {
	var (
		mdvc *flacvorbis.MetaDataBlockVorbisComment
		err  error
	)

	m := receiver.cmtMeta(f)
	if m != nil {
		mdvc, err = flacvorbis.ParseFromMetaDataBlock(*m)
		if err != nil {
			return nil, errors.Wrap(err, "parse metadata block fail")
		}
	} else {
		mdvc = flacvorbis.New()
	}
	return mdvc, nil
}

func (receiver *Ncm) cmtMeta(f *flac.File) *flac.MetaDataBlock {
	for _, m := range f.Meta {
		if m.Type == flac.VorbisComment {
			return m
		}
	}
	return nil
}

func (receiver *Ncm) addFLACMeta(f *flac.File, imgData []byte, meta *MetaInfo) *flac.File {
	if imgData != nil {
		picture, err := flacpicture.NewFromImageData(
			flacpicture.PictureTypeFrontCover,
			"Front cover", imgData,
			receiver.mime(imgData),
		)
		if err != nil {
			slog.With("err", err).Error("create picture fail")
			return f
		}
		pictureMeta := picture.Marshal()
		f.Meta = append(f.Meta, &pictureMeta)
		return f
	}

	if meta.AlbumPic != "" {
		picture := &flacpicture.MetadataBlockPicture{
			PictureType: flacpicture.PictureTypeFrontCover,
			MIME:        "-->",
			Description: "Front cover",
			ImageData:   []byte(meta.AlbumPic),
		}
		pictureMeta := picture.Marshal()
		f.Meta = append(f.Meta, &pictureMeta)
	}
	return f
}

func (receiver *Ncm) addMP3Tag(imgData []byte, meta *MetaInfo, fp *os.File) error {
	tag, err := id3v2.Open(receiver.outPathName(meta, fp.Name()), id3v2.Options{Parse: true})
	if err != nil {
		return errors.Wrapf(err, "Failed to open file %s", receiver.outPathName(meta, fp.Name()))
	}
	defer func() { _ = tag.Close() }()

	imgData = receiver.imageInfo(imgData, meta)

	receiver.attachedPicture(imgData, tag, meta)
	receiver.MusicNameFrame(tag, meta)
	receiver.AlbumFrame(tag, meta)
	receiver.ArtistFrame(tag, meta)

	return tag.Save()
}

func (receiver *Ncm) attachedPicture(ida []byte, tag *id3v2.Tag, meta *MetaInfo) {
	var pic id3v2.PictureFrame
	if ida != nil {
		pic = id3v2.PictureFrame{
			Encoding:    id3v2.EncodingISO,
			MimeType:    receiver.mime(ida),
			PictureType: id3v2.PTFrontCover,
			Description: "Front cover",
			Picture:     ida,
		}
	} else if meta.AlbumPic != "" {
		pic = id3v2.PictureFrame{
			Encoding:    id3v2.EncodingISO,
			MimeType:    "-->",
			PictureType: id3v2.PTFrontCover,
			Description: "Front cover",
			Picture:     []byte(meta.AlbumPic),
		}
	}
	tag.AddAttachedPicture(pic)
}

func (receiver *Ncm) MusicNameFrame(tag *id3v2.Tag, meta *MetaInfo) {
	if tag.GetTextFrame("TIT2").Text == "" {
		if meta.MusicName != "" {
			slog.Info("Adding music name")
			tag.AddTextFrame("TIT2", id3v2.EncodingUTF8, meta.MusicName)
		}
	}
}

func (receiver *Ncm) AlbumFrame(tag *id3v2.Tag, meta *MetaInfo) {
	if tag.GetTextFrame("TALB").Text == "" {
		if meta.Album != "" {
			slog.Info("Adding album name")
			tag.AddTextFrame("TALB", id3v2.EncodingUTF8, meta.Album)
		}
	}
}

func (receiver *Ncm) ArtistFrame(tag *id3v2.Tag, meta *MetaInfo) {
	if len(meta.Artist) < 1 {
		slog.Info("No artist")
		return
	}
	for _, artist := range meta.Artist {
		slog.Info("Adding artist")
		if aster, ok := artist[0].(string); ok {
			tag.AddTextFrame("TPE1", id3v2.EncodingUTF8, aster)
		}
	}
}

func (receiver *Ncm) mime(ida []byte) string {
	if receiver.containPNGHeader(ida) {
		return pngMime
	}
	return jpegMime
}

func (receiver *Ncm) imageInfo(ida []byte, meta *MetaInfo) []byte {
	if ida == nil {
		return receiver.imgInfo(meta.AlbumPic)
	}
	return ida
}

func (receiver *Ncm) imgInfo(url string) []byte {
	var resp ImageResponse
	if err := request.NewClient(&request.Config{}).Do(&ImageRequest{
		request.Request{
			Url: url,
		},
	}, &resp); err != nil {
		slog.With("err", err).Error("Failed to fetch image")
		return nil
	}
	return resp.Data
}

func (receiver *Ncm) defaultValue(o any, d any) any {
	if receiver.isZero(o) {
		return d
	}
	return o
}

func (receiver *Ncm) defaultValueFn(o any, d any, fn func() bool) any {
	if fn() {
		return d
	}
	return o
}

func (receiver *Ncm) isZero(v any) bool {
	rv := reflect.ValueOf(v)
	switch v.(type) {
	case nil:
		return true
	default:
		return rv.IsZero()
	}
}

func (receiver *Ncm) filename(filePath string) string {
	base := filepath.Base(filePath)
	// 使用 filepath.Ext 获取文件的扩展名
	ext := filepath.Ext(base)
	// 从文件名中排除扩展名
	fileNameWithoutExtension := base[:len(base)-len(ext)]
	return string(fileNameWithoutExtension)
}
