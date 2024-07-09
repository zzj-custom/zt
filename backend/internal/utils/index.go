package utils

import (
	"github.com/oklog/ulid/v2"
	"github.com/pkg/errors"
	"math/rand"
	"os"
	"strconv"
	"time"
)

// GenerateRandomNumber 生成随机数字
func GenerateRandomNumber(length int) int {
	rand.NewSource(time.Now().UnixNano())
	var str string
	for i := 0; i < length; i++ {
		str += strconv.Itoa(rand.Intn(9) + 1)
	}

	i, _ := strconv.Atoi(str)
	return i
}

func DirExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// CreateDirIfNotExists Create a directory if it does not exist
func CreateDirIfNotExists(path string) error {
	if !DirExists(path) {
		return os.MkdirAll(path, os.ModePerm)
	}
	return nil
}

func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, errors.Wrapf(err, "Failed to check file exists")
}

func MustULid() string {
	// uLid生成的默认是26位的
	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
	ms := ulid.Timestamp(time.Now())
	return ulid.MustNew(ms, entropy).String()
}
