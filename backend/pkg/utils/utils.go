package utils

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"slices"
	"strings"
	"zt/backend/internal/request"
)

const (
	Ellipses = "..."
)

// MatchOneOf match one of the patterns
func MatchOneOf(text string, patterns ...string) []string {
	var (
		re    *regexp.Regexp
		value []string
	)
	for _, pattern := range patterns {
		// (?flags): set flags within current group; non-capturing
		// s: let . match \n (default false)
		// https://github.com/google/re2/wiki/Syntax
		re = regexp.MustCompile(pattern)
		value = re.FindStringSubmatch(text)
		if len(value) > 0 {
			return value
		}
	}
	return nil
}

// MatchAll return all matching results
func MatchAll(text, pattern string) [][]string {
	re := regexp.MustCompile(pattern)
	value := re.FindAllStringSubmatch(text, -1)
	return value
}

// FileSize return the file size of the specified path file
func FileSize(filePath string) (int64, bool, error) {
	file, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, false, nil
		}
		return 0, false, err
	}
	return file.Size(), true, nil
}

// Domain get the domain of given URL
func Domain(url string) string {
	domainPattern := `([a-z0-9][-a-z0-9]{0,62})\.` +
		`(com\.cn|com\.hk|` +
		`cn|com|net|edu|gov|biz|org|info|pro|name|xxx|xyz|be|` +
		`me|top|cc|tv|tt|vn)`
	domain := MatchOneOf(url, domainPattern)
	if domain != nil {
		return domain[1]
	}
	return ""
}

// LimitLength Handle overly long strings
func LimitLength(s string, length int) string {
	if length == 0 || len(s) <= length {
		return s
	}
	str := []rune(s)
	return string(str[:length-len(Ellipses)]) + Ellipses
}

// FileName Converts a string to a valid filename
func FileName(name, ext string, length int) string {
	rep := strings.NewReplacer("\n", " ", "/", " ", "|", "-", ": ", "：", ":", "：", "'", "’")
	if runtime.GOOS == "windows" {
		rep = strings.NewReplacer("\"", " ", "?", " ", "*", " ", "\\", " ", "<", " ", ">", " ")
	}
	name = rep.Replace(name)
	name = LimitLength(name, length)
	if ext == "" {
		return name
	}
	return fmt.Sprintf("%s.%s", name, ext)
}

// FilePath gen valid file path
func FilePath(name, ext string, length int, outputPath string, escape bool) (string, error) {
	if outputPath != "" {
		if _, err := os.Stat(outputPath); err != nil {
			return "", errors.Wrapf(err, "invalid output path: %s", outputPath)
		}
	}
	var fileName string
	if escape {
		fileName = FileName(name, ext, length)
	} else {
		fileName = fmt.Sprintf("%s.%s", name, ext)
	}
	return filepath.Join(outputPath, fileName), nil
}

// FileLineCounter Counts line in file
func FileLineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

// ParseInputFile Parses input file into args
func ParseInputFile(r io.Reader, items string, itemStart, itemEnd int) []string {
	scanner := bufio.NewScanner(r)

	temp := make([]string, 0)
	totalLines := 0
	for scanner.Scan() {
		totalLines++
		universalURL := strings.TrimSpace(scanner.Text())
		temp = append(temp, universalURL)
	}

	wantedItems := NeedDownloadList(items, itemStart, itemEnd, totalLines)

	itemList := make([]string, 0, len(wantedItems))
	for i, item := range temp {
		if slices.Contains(wantedItems, i+1) {
			itemList = append(itemList, item)
		}
	}

	return itemList
}

// GetNameAndExt return the name and ext of the URL
// https://img9.bcyimg.com/drawer/15294/post/1799t/1f5a87801a0711e898b12b640777720f.jpg ->
// 1f5a87801a0711e898b12b640777720f, jpg
func GetNameAndExt(uri string) (string, string, error) {
	u, err := url.ParseRequestURI(uri)
	if err != nil {
		return "", "", err
	}
	s := strings.Split(u.Path, "/")
	filename := strings.Split(s[len(s)-1], ".")
	if len(filename) > 1 {
		return filename[0], filename[1], nil
	}
	// Image url like this
	// https://img9.bcyimg.com/drawer/15294/post/1799t/1f5a87801a0711e898b12b640777720f.jpg/w650
	// has no suffix
	contentType, err := request.DefaultRequest().ContentType(uri, uri)
	if err != nil {
		return "", "", err
	}
	return filename[0], strings.Split(contentType, "/")[1], nil
}

// Md5 md5 hash
func Md5(text string) string {
	sign := md5.New()
	sign.Write([]byte(text)) // nolint
	return fmt.Sprintf("%x", sign.Sum(nil))
}

// M3u8URLs get all urls from m3u8 url
func M3u8URLs(uri string) ([]string, error) {
	if len(uri) == 0 {
		return nil, errors.New("url is null")
	}

	html, err := request.DefaultRequest().Get(uri, "", nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	lines := strings.Split(html, "\n")
	var urls []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			if strings.HasPrefix(line, "http") {
				urls = append(urls, line)
			} else {
				base, err := url.Parse(uri)
				if err != nil {
					continue
				}
				u, err := url.Parse(line)
				if err != nil {
					continue
				}
				urls = append(urls, base.ResolveReference(u).String())
			}
		}
	}
	return urls, nil
}

// Reverse Reverse a string
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// Range 按范围生成数字序列
func Range(min, max int) []int {
	items := make([]int, max-min+1)
	for index := range items {
		items[index] = min + index
	}
	return items
}
