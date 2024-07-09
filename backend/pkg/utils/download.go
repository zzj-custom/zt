package utils

import (
	"strconv"
	"strings"
)

// NeedDownloadList 返回需要下载的播放列表的索引
func NeedDownloadList(items string, itemStart, itemEnd, length int) []int {
	if items == "" {
		return emptyDownloadList(itemStart, itemEnd, length)
	}
	return notEmptyDownloadList(items)
}

func notEmptyDownloadList(items string) []int {
	var (
		itemList         []int
		selStart, selEnd int
	)

	temp := strings.Split(items, ",")
	for _, i := range temp {
		selection := strings.Split(i, "-")
		selStart, _ = strconv.Atoi(strings.TrimSpace(selection[0]))

		if len(selection) >= 2 {
			selEnd, _ = strconv.Atoi(strings.TrimSpace(selection[1]))
		} else {
			selEnd = selStart
		}

		for item := selStart; item <= selEnd; item++ {
			itemList = append(itemList, item)
		}
	}
	return itemList
}

func emptyDownloadList(itemStart, itemEnd, length int) []int {
	if itemStart < 1 {
		itemStart = 1
	}
	if itemEnd == 0 {
		itemEnd = length
	}
	if itemEnd < itemStart {
		itemStart, itemEnd = itemEnd, itemStart
	}
	return Range(itemStart, itemEnd)
}
