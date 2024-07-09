package iMysql

import (
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func PageHelper(r *http.Request) (page int, pageSize int) {
	err := r.ParseForm()
	if err != nil {
		return 1, 20
	}
	var (
		getPage     string
		getPageSize string
	)
	if r.Method == http.MethodGet {
		getPage = r.Form.Get("page")
		getPageSize = r.Form.Get("pageSize")
	} else {
		getPage = r.PostForm.Get("page")
		getPageSize = r.PostForm.Get("pageSize")
	}

	if getPage == "" {
		page = 1
	} else {
		page, err = strconv.Atoi(getPage)
		if err != nil {
			page = 1
		}
	}
	pageSize, err = strconv.Atoi(getPageSize)
	if err != nil {
		pageSize = 20
	}
	return
}

func Paginate(page int, pageSize int) func(db *gorm.DB) *gorm.DB {
	if page <= 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset).Limit(pageSize)
	}
}

type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
	Total    int `json:"total"`
	Pages    int `json:"pages"`
}

func (p *Pagination) Scope() func(db *gorm.DB) *gorm.DB {
	return Paginate(p.Page, p.PageSize)
}

func NewPagination(total int, pageSize int, page int) Pagination {
	pages := total / pageSize
	if total%pageSize > 0 {
		pages += 1
	}
	if page <= 1 {
		page = 1
	}
	if page >= pages {
		page = pages
	}
	return Pagination{
		Page:     page,
		PageSize: pageSize,
		Total:    total,
		Pages:    pages,
	}
}
