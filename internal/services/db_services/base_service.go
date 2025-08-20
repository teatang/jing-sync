package db_services

import (
	"jing-sync/internal/utils"

	"gorm.io/gorm"
)

type BaseService[T any] struct {
	db *gorm.DB
}

func (bs *BaseService[T]) GetPageList(page, size int) (*utils.PageList[T], error) {
	var infos []T
	offset := (page - 1) * size
	r := bs.db.Where("status = 1").Offset(offset).Limit(size).Find(&infos)

	// 计算总数
	var count int64
	var m T
	bs.db.Model(&m).Where("status = 1").Count(&count)
	return &utils.PageList[T]{
		List:       infos,
		Pagination: utils.PageInfo{Page: page, Size: size, Total: count},
	}, r.Error
}

func (bs *BaseService[T]) GetByID(id string) (*T, error) {
	var info T
	err := bs.db.First(&info, id).Error
	return &info, err
}

func (bs *BaseService[T]) Update(info *T) error {
	return bs.db.Save(info).Error
}

func (bs *BaseService[T]) Create(info *T) error {
	return bs.db.Create(info).Error
}

func (us *BaseService[T]) Delete(id string) error {
	// delete 实际操作为 update status = 0
	var m T
	return us.db.Model(&m).Where("id = ?", id).Update("status", 0).Error
}
