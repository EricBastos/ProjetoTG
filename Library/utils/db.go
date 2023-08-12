package utils

import (
	"gorm.io/gorm"
)

func Paginate(p, s int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if p == 0 {
			p = 1
		}

		switch {
		case s > 100:
			s = 100
		case s <= 0:
			s = 10
		}

		offset := (p - 1) * s
		return db.Offset(offset).Limit(s)
	}
}
