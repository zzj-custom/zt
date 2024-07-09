package utils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/pkg/errors"
)

func findFFmpegExecutable() string {
	ffmpegFileName := "ffmpeg"
	if runtime.GOOS == "windows" {
		ffmpegFileName = "ffmpeg.exe" // 在Windows上添加.exe擴展名
	}
	// 嘗試在當前目錄中查找ffmpeg
	matches, err := filepath.Glob("./" + ffmpegFileName)
	if err == nil && len(matches) > 0 {
		// 如果在當前目錄找到了ffmpeg，直接返回這個路徑
		return "./" + ffmpegFileName
	}

	// 返回從PATH中找到的ffmpeg路徑
	return ffmpegFileName
}

func runMergeCmd(cmd *exec.Cmd, paths []string, mergeFilePath string) error {
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return errors.Errorf("%s\n%s", err, stderr.String())
	}

	if mergeFilePath != "" {
		_ = os.Remove(mergeFilePath)
	}
	// remove parts
	for _, path := range paths {
		_ = os.Remove(path)
	}
	return nil
}

// MergeFilesWithSameExtension merges files that have the same extension into one.
// Can also handle merging audio and video.
func MergeFilesWithSameExtension(paths []string, mergedFilePath string) error {
	command := []string{
		"-y",
	}
	for _, path := range paths {
		command = append(command, "-i", path)
	}
	command = append(command, "-c:v", "copy", "-c:a", "copy", mergedFilePath)

	return runMergeCmd(exec.Command(findFFmpegExecutable(), command...), paths, "")
}

// MergeToMP4 merges video parts to an MP4 file.
func MergeToMP4(paths []string, mergedFilePath string, filename string) error {
	mergeFilePath := filename + ".txt" // merge list file should be in the current directory

	// write ffmpeg input file list
	mergeFile, _ := os.Create(mergeFilePath)
	for _, path := range paths {
		_, _ = mergeFile.Write([]byte(fmt.Sprintf("file '%s'\n", path)))
	}
	_ = mergeFile.Close()

	cmd := exec.Command(
		findFFmpegExecutable(), "-y", "-f", "concat", "-safe", "0",
		"-i", mergeFilePath, "-c", "copy", "-bsf:a", "aac_adtstoasc", mergedFilePath,
	)
	return runMergeCmd(cmd, paths, mergeFilePath)
}
