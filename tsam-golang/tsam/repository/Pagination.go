package repository

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type Paginator interface {
	Paginate() QueryProcessor
}

// Pagination is struct for DTO only.
type Pagination struct {
	Limit      int
	Offset     int
	TotalCount int
}

// TenantPagination needs to be discussed.
type TenantPagination struct {
	TenantID uuid.UUID
	Pagination
}

// Paginate will restrict the output of query with limit and offset & fill totalCount with total records.
func (p *Pagination) Paginate() QueryProcessor {
	return func(db *gorm.DB, out interface{}) (*gorm.DB, error) {
		if out != nil {
			if err := db.Model(out).Count(&p.TotalCount).Error; err != nil {
				return db, err
			}
		}

		if p.Limit != -1 {
			db = db.Limit(p.Limit)
		}

		if p.Offset > 0 {
			db = db.Offset(p.Limit * p.Offset)
		}
		return db, nil
	}
}

// `json:"-" gorm:"-"`
// `json:"-" gorm:"-"`
// `json:"-" gorm:"-"`

// Limit gets limit.
// func (p *Pagination) Limit() int {
// 	return p.limit
// }

// // SetLimit sets limit.
// func (p *Pagination) SetLimit(limit int) {
// 	p.limit = limit
// }

// // Offset gets offset.
// func (p *Pagination) Offset() int {
// 	return p.offset
// }

// // SetOffset sets offset.
// func (p *Pagination) SetOffset(offset int) {
// 	p.offset = offset
// }

// // TotalCount gets totalCount.
// func (p *Pagination) TotalCount() *int {
// 	return p.totalCount
// }

// // SetTotalCount sets totalCount.
// func (p *Pagination) SetTotalCount(totalCount int) {
// 	p.totalCount = &totalCount
// }
