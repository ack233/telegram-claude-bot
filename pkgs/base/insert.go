package base

import "gorm.io/gorm"

func InsertFileContentToDB(db *gorm.DB, _s string, tablefunc func(line string) interface{}) error {
	if _s == "" {
		return nil
	}
	text := tablefunc(_s)
	return db.Create(text).Error
}
