package services

import (
	"gorm.io/gorm"
)

type PageInfo struct {
	Page  int   `json:"page"`
	Size  int   `json:"size"`
	Total int64 `json:"total"`
}

type PageList[T any] struct {
	List       []T      `json:"list"`
	Pagination PageInfo `json:"pagination"`
}

type BaseService[T any] struct {
	db *gorm.DB
}

func (bs *BaseService[T]) GetPageList(page, size int) (*PageList[T], error) {
	var infos []T
	offset := (page - 1) * size
	r := bs.db.Offset(offset).Limit(size).Find(&infos)
	return &PageList[T]{
		List:       infos,
		Pagination: PageInfo{Page: page, Size: size, Total: r.RowsAffected},
	}, r.Error
}
